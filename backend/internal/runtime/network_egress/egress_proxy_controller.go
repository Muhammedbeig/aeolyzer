// Package networkegress enforces domain, DNS, protocol, header, and bandwidth
// policy before any Layer 6 outbound request.
package networkegress

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrEgressDenied indicates that an outbound request violates policy.
	ErrEgressDenied = errors.New("network egress denied")
	// ErrResponseLimit indicates that a response exceeded the approved budget.
	ErrResponseLimit = errors.New("network response exceeds approved limit")
)

// Resolver resolves every IP for policy validation.
type Resolver interface {
	LookupIPAddr(context.Context, string) ([]net.IPAddr, error)
}

// Policy is an immutable request allowlist.
type Policy struct {
	AllowedDomains        []string
	AllowedMethods        []string
	AllowedRequestHeaders []string
	AllowHTTP             bool
	MaxRequestBytes       int64
	MaxResponseBytes      int64
	Timeout               time.Duration
}

// Request is one bounded egress operation.
type Request struct {
	URL     string
	Method  string
	Headers http.Header
	Body    []byte
}

// Response contains a bounded response body and safe metadata.
type Response struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

type validatedTarget struct {
	url       *url.URL
	addresses []netip.Addr
}

type sender func(context.Context, validatedTarget, Request, Policy) (Response, error)

// Controller validates and executes approved egress requests.
type Controller struct {
	resolver Resolver
	policy   Policy
	send     sender
}

// NewController builds an egress controller with a DNS-pinned HTTP transport.
func NewController(resolver Resolver, policy Policy) (*Controller, error) {
	if resolver == nil {
		return nil, errors.New("egress resolver is required")
	}
	normalized, err := normalizePolicy(policy)
	if err != nil {
		return nil, err
	}
	return &Controller{
		resolver: resolver,
		policy:   normalized,
		send:     sendPinnedHTTP,
	}, nil
}

// Do validates policy, resolves all addresses, rejects private/special ranges,
// pins the connection to an approved address, and enforces byte limits.
func (c *Controller) Do(ctx context.Context, request Request) (Response, error) {
	if c == nil || c.resolver == nil || c.send == nil {
		return Response{}, errors.New("egress controller is not configured")
	}
	if err := ctx.Err(); err != nil {
		return Response{}, fmt.Errorf("execute egress request: %w", err)
	}
	target, err := c.validateTarget(ctx, request)
	if err != nil {
		return Response{}, err
	}
	return c.send(ctx, target, request, c.policy)
}

func (c *Controller) validateTarget(
	ctx context.Context,
	request Request,
) (validatedTarget, error) {
	if int64(len(request.Body)) > c.policy.MaxRequestBytes {
		return validatedTarget{}, ErrEgressDenied
	}
	method := strings.ToUpper(strings.TrimSpace(request.Method))
	if method == "" {
		method = http.MethodGet
	}
	if !contains(c.policy.AllowedMethods, method) {
		return validatedTarget{}, ErrEgressDenied
	}
	for key := range request.Headers {
		if !containsFold(c.policy.AllowedRequestHeaders, key) {
			return validatedTarget{}, ErrEgressDenied
		}
	}

	parsed, err := url.Parse(request.URL)
	if err != nil ||
		parsed.Hostname() == "" ||
		parsed.User != nil ||
		parsed.Fragment != "" {
		return validatedTarget{}, ErrEgressDenied
	}
	if parsed.Scheme != "https" && !(c.policy.AllowHTTP && parsed.Scheme == "http") {
		return validatedTarget{}, ErrEgressDenied
	}
	host := strings.ToLower(strings.TrimSuffix(parsed.Hostname(), "."))
	if !domainAllowed(host, c.policy.AllowedDomains) {
		return validatedTarget{}, ErrEgressDenied
	}

	resolved, err := c.resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return validatedTarget{}, fmt.Errorf("resolve egress host: %w", err)
	}
	if len(resolved) == 0 {
		return validatedTarget{}, ErrEgressDenied
	}
	addresses := make([]netip.Addr, 0, len(resolved))
	for _, item := range resolved {
		address, ok := netip.AddrFromSlice(item.IP)
		if !ok {
			return validatedTarget{}, ErrEgressDenied
		}
		address = address.Unmap()
		if !allowedPublicAddress(address) {
			return validatedTarget{}, ErrEgressDenied
		}
		addresses = append(addresses, address)
	}
	return validatedTarget{url: parsed, addresses: addresses}, nil
}

func sendPinnedHTTP(
	ctx context.Context,
	target validatedTarget,
	request Request,
	policy Policy,
) (Response, error) {
	method := strings.ToUpper(strings.TrimSpace(request.Method))
	if method == "" {
		method = http.MethodGet
	}
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		method,
		target.url.String(),
		bytes.NewReader(request.Body),
	)
	if err != nil {
		return Response{}, fmt.Errorf("build egress request: %w", err)
	}
	httpRequest.Header = request.Headers.Clone()

	port := target.url.Port()
	if port == "" {
		if target.url.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	addressIndex := 0
	dialer := &net.Dialer{Timeout: policy.Timeout}
	transport := &http.Transport{
		Proxy:                 nil,
		DisableKeepAlives:     true,
		MaxIdleConns:          0,
		ResponseHeaderTimeout: policy.Timeout,
		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			if addressIndex >= len(target.addresses) {
				return nil, ErrEgressDenied
			}
			address := net.JoinHostPort(target.addresses[addressIndex].String(), port)
			addressIndex++
			return dialer.DialContext(ctx, network, address)
		},
	}
	defer transport.CloseIdleConnections()
	client := &http.Client{
		Transport: transport,
		Timeout:   policy.Timeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return ErrEgressDenied
		},
	}
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return Response{}, fmt.Errorf("send egress request: %w", err)
	}
	defer httpResponse.Body.Close()

	body, err := io.ReadAll(io.LimitReader(httpResponse.Body, policy.MaxResponseBytes+1))
	if err != nil {
		return Response{}, fmt.Errorf("read egress response: %w", err)
	}
	if int64(len(body)) > policy.MaxResponseBytes {
		return Response{}, ErrResponseLimit
	}
	return Response{
		StatusCode:  httpResponse.StatusCode,
		ContentType: httpResponse.Header.Get("Content-Type"),
		Body:        body,
	}, nil
}

func normalizePolicy(policy Policy) (Policy, error) {
	if len(policy.AllowedDomains) == 0 ||
		len(policy.AllowedMethods) == 0 ||
		policy.MaxRequestBytes < 0 ||
		policy.MaxResponseBytes < 1 ||
		policy.MaxRequestBytes > 4<<20 ||
		policy.MaxResponseBytes > 32<<20 ||
		policy.Timeout < time.Second ||
		policy.Timeout > 2*time.Minute {
		return Policy{}, errors.New("egress policy is invalid")
	}
	result := policy
	result.AllowedDomains = make([]string, 0, len(policy.AllowedDomains))
	for _, domain := range policy.AllowedDomains {
		normalized := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))
		if normalized == "" ||
			strings.Contains(normalized, "://") ||
			strings.ContainsAny(normalized, "/?#@") {
			return Policy{}, errors.New("egress domain allowlist is invalid")
		}
		result.AllowedDomains = append(result.AllowedDomains, normalized)
	}
	result.AllowedMethods = make([]string, 0, len(policy.AllowedMethods))
	for _, method := range policy.AllowedMethods {
		normalized := strings.ToUpper(strings.TrimSpace(method))
		switch normalized {
		case http.MethodGet, http.MethodHead, http.MethodPost:
			result.AllowedMethods = append(result.AllowedMethods, normalized)
		default:
			return Policy{}, errors.New("egress method allowlist is invalid")
		}
	}
	result.AllowedRequestHeaders = append([]string(nil), policy.AllowedRequestHeaders...)
	return result, nil
}

func domainAllowed(host string, allowed []string) bool {
	for _, pattern := range allowed {
		if strings.HasPrefix(pattern, "*.") {
			suffix := strings.TrimPrefix(pattern, "*")
			if strings.HasSuffix(host, suffix) && host != strings.TrimPrefix(suffix, ".") {
				return true
			}
			continue
		}
		if host == pattern {
			return true
		}
	}
	return false
}

func allowedPublicAddress(address netip.Addr) bool {
	if !address.IsValid() ||
		!address.IsGlobalUnicast() ||
		address.IsPrivate() ||
		address.IsLoopback() ||
		address.IsLinkLocalUnicast() ||
		address.IsLinkLocalMulticast() ||
		address.IsMulticast() ||
		address.IsUnspecified() {
		return false
	}
	for _, prefix := range []string{
		"100.64.0.0/10",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"2001:db8::/32",
	} {
		parsed := netip.MustParsePrefix(prefix)
		if parsed.Contains(address) {
			return false
		}
	}
	return true
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func containsFold(values []string, expected string) bool {
	for _, value := range values {
		if strings.EqualFold(value, expected) {
			return true
		}
	}
	return false
}

func hostPort(target *url.URL) string {
	port := target.Port()
	if port == "" {
		if target.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	return net.JoinHostPort(target.Hostname(), port)
}

func parsePort(target *url.URL) (int, error) {
	port := target.Port()
	if port == "" {
		if target.Scheme == "https" {
			return 443, nil
		}
		return 80, nil
	}
	return strconv.Atoi(port)
}
