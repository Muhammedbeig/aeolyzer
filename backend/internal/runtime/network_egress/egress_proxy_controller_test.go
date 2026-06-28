package networkegress

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"
)

type resolverStub map[string][]net.IPAddr

func (r resolverStub) LookupIPAddr(_ context.Context, host string) ([]net.IPAddr, error) {
	return r[host], nil
}

func TestControllerRejectsSSRFAndPolicyViolations(t *testing.T) {
	controller := newTestController(t, resolverStub{
		"allowed.example": {{IP: net.ParseIP("93.184.216.34")}},
		"private.example": {{IP: net.ParseIP("10.0.0.4")}},
		"mixed.example": {
			{IP: net.ParseIP("93.184.216.34")},
			{IP: net.ParseIP("169.254.169.254")},
		},
	})
	tests := map[string]Request{
		"unlisted domain": {
			URL:    "https://attacker.invalid",
			Method: http.MethodGet,
		},
		"private dns": {
			URL:    "https://private.example",
			Method: http.MethodGet,
		},
		"mixed dns rebinding": {
			URL:    "https://mixed.example",
			Method: http.MethodGet,
		},
		"metadata literal": {
			URL:    "https://169.254.169.254/latest/meta-data",
			Method: http.MethodGet,
		},
		"userinfo": {
			URL:    "https://user:pass@allowed.example",
			Method: http.MethodGet,
		},
		"method": {
			URL:    "https://allowed.example",
			Method: http.MethodDelete,
		},
		"credential header": {
			URL:     "https://allowed.example",
			Method:  http.MethodGet,
			Headers: http.Header{"Authorization": {"Bearer secret"}},
		},
	}
	for name, request := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := controller.Do(context.Background(), request); !errors.Is(err, ErrEgressDenied) {
				t.Fatalf("Controller.Do() error = %v, want %v", err, ErrEgressDenied)
			}
		})
	}
}

func TestControllerPassesOnlyValidatedTargetToSender(t *testing.T) {
	controller := newTestController(t, resolverStub{
		"allowed.example": {{IP: net.ParseIP("93.184.216.34")}},
	})
	called := false
	controller.send = func(
		_ context.Context,
		target validatedTarget,
		request Request,
		_ Policy,
	) (Response, error) {
		called = true
		if target.url.Hostname() != "allowed.example" ||
			len(target.addresses) != 1 ||
			target.addresses[0].String() != "93.184.216.34" {
			t.Fatalf("sender target = %+v, want validated pinned target", target)
		}
		if request.Method != http.MethodGet {
			t.Fatalf("sender method = %q, want GET", request.Method)
		}
		return Response{StatusCode: http.StatusOK}, nil
	}
	response, err := controller.Do(context.Background(), Request{
		URL:    "https://allowed.example/path",
		Method: http.MethodGet,
	})
	if err != nil {
		t.Fatalf("Controller.Do() failed: %v", err)
	}
	if !called || response.StatusCode != http.StatusOK {
		t.Fatal("Controller.Do() did not call validated sender")
	}
}

func TestWildcardDomainDoesNotMatchApexOrSuffixConfusion(t *testing.T) {
	if domainAllowed("example.com", []string{"*.example.com"}) {
		t.Fatal("wildcard unexpectedly matched apex")
	}
	if domainAllowed("example.com.attacker.invalid", []string{"*.example.com"}) {
		t.Fatal("wildcard matched suffix-confusion domain")
	}
	if !domainAllowed("api.example.com", []string{"*.example.com"}) {
		t.Fatal("wildcard failed to match subdomain")
	}
}

func newTestController(t *testing.T, resolver Resolver) *Controller {
	t.Helper()
	controller, err := NewController(resolver, Policy{
		AllowedDomains:        []string{"allowed.example", "private.example", "mixed.example"},
		AllowedMethods:        []string{http.MethodGet},
		AllowedRequestHeaders: []string{"Accept"},
		MaxRequestBytes:       1024,
		MaxResponseBytes:      4096,
		Timeout:               5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewController() failed: %v", err)
	}
	controller.send = func(
		context.Context,
		validatedTarget,
		Request,
		Policy,
	) (Response, error) {
		return Response{StatusCode: http.StatusOK}, nil
	}
	return controller
}
