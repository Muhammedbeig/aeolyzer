package interop

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"aeolyzer/internal/runtime"
)

var (
	// PERFORMANCE: Precompile regular expressions during init to avoid redundant AST evaluations at runtime.
	titlePattern        = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	metaPattern         = regexp.MustCompile(`(?is)<meta\s+[^>]*(?:name|property)\s*=\s*["']([^"']+)["'][^>]*content\s*=\s*["']([^"']*)["'][^>]*>`)
	reversedMetaPattern = regexp.MustCompile(`(?is)<meta\s+[^>]*content\s*=\s*["']([^"']*)["'][^>]*(?:name|property)\s*=\s*["']([^"']+)["'][^>]*>`)
	linkPattern         = regexp.MustCompile(`(?is)<link\s+[^>]*rel\s*=\s*["'][^"']*(?:icon|shortcut icon)[^"']*["'][^>]*href\s*=\s*["']([^"']+)["'][^>]*>`)
	hrefPattern         = regexp.MustCompile(`(?is)<a\s+[^>]*href\s*=\s*["'](https?://[^"'#\s]+)["'][^>]*>(.*?)</a>`)
	tagPattern          = regexp.MustCompile(`<[^>]+>`)
	spacePattern        = regexp.MustCompile(`\s+`)
)

var ignoredCompetitorHosts = map[string]struct{}{
	"facebook.com": {}, "instagram.com": {}, "linkedin.com": {}, "tiktok.com": {},
	"twitter.com": {}, "x.com": {}, "youtube.com": {}, "github.com": {},
	"pinterest.com": {}, "schema.org": {}, "w3.org": {}, "iana.org": {},
}

var competitorTerms = []string{
	"alternative", "compare", "comparison", "competitor", "versus", " vs ",
}

type SiteClient struct {
	client *http.Client
}

func NewSiteClient(timeout time.Duration) *SiteClient {
	return &SiteClient{
		client: &http.Client{
			// STATE MANAGEMENT: Apply strict timeout to bound thread block durations.
			Timeout: timeout,
			// SECURITY INVARIANT: Disable redirects to prevent silent SSRF attacks chaining through 30x codes.
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return errors.New("redirects are disabled")
			},
		},
	}
}

func (c *SiteClient) Inspect(ctx context.Context, targetURL string, maxBytes int64) (runtime.ExecutionResult, error) {
	if c == nil || c.client == nil {
		return runtime.ExecutionResult{}, errors.New("site connector is not configured")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return runtime.ExecutionResult{}, fmt.Errorf("build site request: %w", err)
	}
	request.Header.Set("Accept", "text/html,application/xhtml+xml")
	request.Header.Set("User-Agent", "AEOlyzer-Site-Inspector/1.0")

	response, err := c.client.Do(request)
	if err != nil {
		return runtime.ExecutionResult{}, fmt.Errorf("inspect site: %w", err)
	}
	// EDGE CASE: Ensure body is closed to prevent socket exhaustion and FD leaks.
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return runtime.ExecutionResult{}, fmt.Errorf("inspect site: unexpected status %d", response.StatusCode)
	}
	mediaType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil || (mediaType != "text/html" && mediaType != "application/xhtml+xml") {
		return runtime.ExecutionResult{}, errors.New("inspect site: unsupported content type")
	}

	// PERFORMANCE & SECURITY INVARIANT: Cap memory buffering via LimitReader to maxBytes + 1, rejecting payloads exceeding the boundary without reading further.
	body, err := io.ReadAll(io.LimitReader(response.Body, maxBytes+1))
	if err != nil {
		return runtime.ExecutionResult{}, fmt.Errorf("read site: %w", err)
	}
	if int64(len(body)) > maxBytes {
		return runtime.ExecutionResult{}, errors.New("inspect site: response exceeds limit")
	}

	return extractSiteMetadata(targetURL, string(body)), nil
}

func extractSiteMetadata(targetURL, document string) runtime.ExecutionResult {
	meta := collectMeta(document)
	title := firstNonEmpty(meta["og:site_name"], meta["application-name"], submatch(titlePattern, document))
	description := firstNonEmpty(meta["description"], meta["og:description"])
	category := firstNonEmpty(meta["article:section"], meta["category"])
	iconURL := resolveReference(targetURL, submatch(linkPattern, document))
	if iconURL == "" {
		iconURL = resolveReference(targetURL, "/favicon.ico")
	}

	return runtime.ExecutionResult{
		// STATE MANAGEMENT: Sanitize outputs via cleanHTMLText to strip rogue markup before boundary egress.
		Title:                cleanHTMLText(title),
		Description:          cleanHTMLText(description),
		IconURL:              iconURL,
		Category:             cleanHTMLText(category),
		CandidateCompetitors: collectExternalHosts(targetURL, document),
	}
}

func collectMeta(document string) map[string]string {
	result := make(map[string]string)
	for _, match := range metaPattern.FindAllStringSubmatch(document, -1) {
		result[strings.ToLower(strings.TrimSpace(match[1]))] = match[2]
	}
	for _, match := range reversedMetaPattern.FindAllStringSubmatch(document, -1) {
		result[strings.ToLower(strings.TrimSpace(match[2]))] = match[1]
	}
	return result
}

func collectExternalHosts(targetURL, document string) []string {
	base, err := url.Parse(targetURL)
	if err != nil {
		return nil
	}
	ownHost := normalizeHost(base.Hostname())
	hosts := make(map[string]struct{})
	for _, match := range hrefPattern.FindAllStringSubmatch(document, -1) {
		if len(match) < 3 || !containsCompetitorTerm(cleanHTMLText(match[2])) {
			continue
		}
		parsed, err := url.Parse(match[1])
		if err != nil {
			continue
		}
		host := normalizeHost(parsed.Hostname())
		if host == "" || host == ownHost || strings.HasSuffix(host, "."+ownHost) {
			continue
		}
		if _, ignored := ignoredCompetitorHosts[host]; ignored {
			continue
		}
		hosts[host] = struct{}{}
	}

	result := make([]string, 0, len(hosts))
	for host := range hosts {
		result = append(result, host)
	}
	sort.Strings(result)
	if len(result) > 8 {
		result = result[:8]
	}
	return result
}

func containsCompetitorTerm(value string) bool {
	value = " " + strings.ToLower(value) + " "
	for _, term := range competitorTerms {
		if strings.Contains(value, term) {
			return true
		}
	}
	return false
}

func resolveReference(targetURL, reference string) string {
	base, err := url.Parse(targetURL)
	if err != nil {
		return ""
	}
	ref, err := url.Parse(strings.TrimSpace(reference))
	if err != nil {
		return ""
	}
	resolved := base.ResolveReference(ref)
	if resolved.Scheme != "http" && resolved.Scheme != "https" {
		return ""
	}
	return resolved.String()
}

func submatch(pattern *regexp.Regexp, value string) string {
	match := pattern.FindStringSubmatch(value)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func cleanHTMLText(value string) string {
	value = tagPattern.ReplaceAllString(value, " ")
	value = strings.NewReplacer(
		"&amp;", "&",
		"&quot;", `"`,
		"&#39;", "'",
		"&lt;", "<",
		"&gt;", ">",
	).Replace(value)
	return strings.TrimSpace(spacePattern.ReplaceAllString(value, " "))
}

func normalizeHost(host string) string {
	return strings.TrimPrefix(strings.ToLower(host), "www.")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
