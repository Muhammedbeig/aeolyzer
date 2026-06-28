package datasecuritymesh

import (
	"crypto/x509"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// MTLSIdentity contains verified, non-secret peer identity.
type MTLSIdentity struct {
	SPIFFEID     string    `json:"spiffe_id"`
	SerialNumber string    `json:"serial_number"`
	NotAfter     time.Time `json:"not_after"`
}

// VerifyMTLSIdentity verifies the full certificate chain and an exact SPIFFE
// URI SAN for server authentication.
func VerifyMTLSIdentity(
	peerCertificates []*x509.Certificate,
	roots *x509.CertPool,
	expectedSPIFFEID string,
	now time.Time,
) (MTLSIdentity, error) {
	if len(peerCertificates) == 0 ||
		roots == nil ||
		expectedSPIFFEID == "" ||
		now.IsZero() {
		return MTLSIdentity{}, errors.New("mtls verification context is incomplete")
	}
	expectedURI, err := url.Parse(expectedSPIFFEID)
	if err != nil || expectedURI.Scheme != "spiffe" || expectedURI.Host == "" {
		return MTLSIdentity{}, errors.New("expected spiffe identity is invalid")
	}

	intermediates := x509.NewCertPool()
	for _, certificate := range peerCertificates[1:] {
		intermediates.AddCert(certificate)
	}
	leaf := peerCertificates[0]
	if _, err := leaf.Verify(x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediates,
		CurrentTime:   now,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}); err != nil {
		return MTLSIdentity{}, fmt.Errorf("verify mtls certificate chain: %w", err)
	}

	matches := 0
	for _, identity := range leaf.URIs {
		if identity.String() == expectedURI.String() {
			matches++
		}
	}
	if matches != 1 {
		return MTLSIdentity{}, errors.New("mtls spiffe identity mismatch")
	}
	return MTLSIdentity{
		SPIFFEID:     expectedURI.String(),
		SerialNumber: leaf.SerialNumber.Text(16),
		NotAfter:     leaf.NotAfter.UTC(),
	}, nil
}
