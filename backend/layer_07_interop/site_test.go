package interop

import "testing"

func TestExtractSiteMetadata(t *testing.T) {
	t.Parallel()

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
	if got.Title != "Example Brand" {
		t.Fatalf("extractSiteMetadata() title = %q", got.Title)
	}
	if got.IconURL != "https://example.com/mark.svg" {
		t.Fatalf("extractSiteMetadata() icon = %q", got.IconURL)
	}
	if len(got.CandidateCompetitors) != 1 || got.CandidateCompetitors[0] != "competitor.example" {
		t.Fatalf("extractSiteMetadata() competitors = %#v", got.CandidateCompetitors)
	}
}
