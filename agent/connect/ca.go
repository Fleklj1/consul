package connect

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
)

// ParseCert parses the x509 certificate from a PEM-encoded value.
func ParseCert(pemValue string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(pemValue))
	if block == nil {
		return nil, fmt.Errorf("no PEM-encoded data found")
	}

	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("first PEM-block should be CERTIFICATE type")
	}

	return x509.ParseCertificate(block.Bytes)
}

// ParseSigner parses a crypto.Signer from a PEM-encoded key. The private key
// is expected to be the first block in the PEM value.
func ParseSigner(pemValue string) (crypto.Signer, error) {
	block, _ := pem.Decode([]byte(pemValue))
	if block == nil {
		return nil, fmt.Errorf("no PEM-encoded data found")
	}

	switch block.Type {
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)

	default:
		return nil, fmt.Errorf("unknown PEM block type for signing key: %s", block.Type)
	}
}

// ParseCSR parses a CSR from a PEM-encoded value. The certificate request
// must be the the first block in the PEM value.
func ParseCSR(pemValue string) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode([]byte(pemValue))
	if block == nil {
		return nil, fmt.Errorf("no PEM-encoded data found")
	}

	if block.Type != "CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("first PEM-block should be CERTIFICATE REQUEST type")
	}

	return x509.ParseCertificateRequest(block.Bytes)
}

// KeyId returns a x509 KeyId from the given signing key. The key must be
// an *ecdsa.PublicKey, but is an interface type to support crypto.Signer.
func KeyId(raw interface{}) ([]byte, error) {
	pub, ok := raw.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid key type: %T", raw)
	}

	// This is not standard; RFC allows any unique identifier as long as they
	// match in subject/authority chains but suggests specific hashing of DER
	// bytes of public key including DER tags. I can't be bothered to do esp.
	// since ECDSA keys don't have a handy way to marshal the publick key alone.
	h := sha256.New()
	h.Write(pub.X.Bytes())
	h.Write(pub.Y.Bytes())
	return h.Sum([]byte{}), nil
}

// HexString returns a standard colon-separated hex value for the input
// byte slice. This should be used with cert serial numbers and so on.
func HexString(input []byte) string {
	return strings.Replace(fmt.Sprintf("% x", input), " ", ":", -1)
}
