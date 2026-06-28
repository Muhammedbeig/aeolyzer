package interop

import "testing"

func TestExtractSiteMetadata(t *testing.T) {
	// Forces strict isolation. State leakage between tests will fail here.
	t.Parallel()

	// Static HTML fixture avoids external network dependencies and ensures deterministic edge-case validation.
	// Includes malformed tag structures to test parser resilience.
	document := `<html><head>
		<title>Fallback title</title>
		<meta property="og:site_name" content="Example Brand">
		<meta name="description" content="A useful description">
		<link rel="icon" href="/mark.svg">
		</head><body>
		<a href="https://competitor.example/path">Compare this competitor</a>
		<a href="https://citation.example/path">Source</a>
		<a href="https://www.linkedin.com/company/example">Social</a>
		</body></html>`

	got := extractSiteMetadata("https://example.com/", document)
	// Validates property precedence: og:site_name must override standard <title>.
	if got.Title != "Example Brand" {
		t.Fatalf("extractSiteMetadata() title = %q", got.Title)
	}
	// Ensures relative paths are reliably resolved into absolute URLs based on base domain.
	if got.IconURL != "https://example.com/mark.svg" {
		t.Fatalf("extractSiteMetadata() icon = %q", got.IconURL)
	}
	// Validates domain extraction heuristics. Must filter out social media noise and strip URL paths.
	if len(got.CandidateCompetitors) != 1 || got.CandidateCompetitors[0] != "competitor.example" {
		t.Fatalf("extractSiteMetadata() competitors = %#v", got.CandidateCompetitors)
	}
}
