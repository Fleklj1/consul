package connect

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var invalidDNSNameChars = regexp.MustCompile(`[^a-z0-9]`)

// ServiceCN returns the common name for a service's certificate. We can't use
// SPIFFE URIs because some CAs require valid FQDN format. We can't use SNI
// values because they are often too long than the 64 bytes allowed by
// CommonNames. We could attempt to encode more information into this to make
// identifying which instance/node it was issued to in a management tool easier
// but that just introduces more complications around length. It's also strange
// that the Common Name would encode more information than the actual
// identifying URI we use to assert anything does and my lead to bad assumptions
// that the common name is in some way "secure" or verified - there is nothing
// inherently provable here except that the requestor had ACLs for that service
// name in that DC.
//
// Format is:
//   <sanitized_service_name>.svc.<trust_domain_first_8>.consul
//
//   service name is sanitized by removing any chars that are not legal in a DNS
//   name and lower casing. It is truncated to the first X chars to keep the
//   total at 64.
//
//   trust domain is truncated to keep the whole name short
func ServiceCN(serviceName, trustDomain string) string {
	svc := invalidDNSNameChars.ReplaceAllString(strings.ToLower(serviceName), "")
	// 20 = 7 bytes for ".consul", 8 bytes for trust domain, 5 bytes for ".svc."
	return fmt.Sprintf("%s.svc.%s.consul",
		truncateTo(svc, 64-20), truncateTo(trustDomain, 8))
}

// AgentCN returns the common name for an agent certificate. See ServiceCN for
// more details on rationale.
//
// Format is:
//   <sanitized_node_name>.agnt.<trust_domain_first_8>.consul
//
//   node name is sanitized by removing any chars that are not legal in a DNS
//   name and lower casing. It is truncated to the first X chars to keep the
//   total at 64.
//
//   trust domain is truncated to keep the whole name short
func AgentCN(node, trustDomain string) string {
	nodeSan := invalidDNSNameChars.ReplaceAllString(strings.ToLower(node), "")
	// 21 = 7 bytes for ".consul", 8 bytes for trust domain, 6 bytes for ".agnt."
	return fmt.Sprintf("%s.agnt.%s.consul",
		truncateTo(nodeSan, 64-21), truncateTo(trustDomain, 8))
}

// CompactUID returns a crypto random Unique Identifier string consiting of 8
// characters of base36 encoded random value. This has roughly 41 bits of
// entropy so is suitable for infrequently occuring events with low probability
// of collision. It is not suitable for UUIDs for very frequent events. It's
// main purpose is to assign unique values to CA certificate Common Names which
// need to be unique in some providers - see CACN - but without using up large
// amounts of the limited 64 character Common Name. It also makes the values
// more easily digestable by humans considering there are likely to be few of
// them ever in use.
func CompactUID() (string, error) {
	// 48 bits (6 bytes) is enough to fill 8 bytes in base36 but it's simpler to
	// have a whole uint8 to convert from.
	var raw [8]byte
	_, err := rand.Read(raw[:])
	if err != nil {
		return "", err
	}

	i := binary.LittleEndian.Uint64(raw[:])
	return truncateTo(strconv.FormatInt(int64(i), 36), 8), nil
}

// CACN returns the common name for a CA certificate. See ServiceCN for more
// details on rationale. A uniqueID is requires because some providers (e.g.
// Vault) cache by subject and so produce incorrect results - for example they
// won't cross-sign an older CA certificate with the same common name since they
// think they already have a valid cert for that CN and just return the current
// root.
//
// This can be generated by any means but will be truncated to 8 chars and
// sanitised to DNS-safe chars. CompactUID generates suitable UIDs for this
// specific purpose.
//
// Format is:
//   {provider}-{uniqueID_first8}.{pri|sec}.ca.<trust_domain_first_8>.consul
//
//   trust domain is truncated to keep the whole name short
func CACN(provider, uniqueID, trustDomain string, primaryDC bool) string {
	providerSan := invalidDNSNameChars.ReplaceAllString(strings.ToLower(provider), "")
	typ := "pri"
	if !primaryDC {
		typ = "sec"
	}
	// 32 = 7 bytes for ".consul", 8 bytes for trust domain, 8 bytes for
	// ".pri.ca.", 9 bytes for "-{uniqueID-8-b36}"
	uidSAN := invalidDNSNameChars.ReplaceAllString(strings.ToLower(uniqueID), "")
	return fmt.Sprintf("%s-%s.%s.ca.%s.consul", typ, truncateTo(uidSAN, 8),
		truncateTo(providerSan, 64-32), truncateTo(trustDomain, 8))
}

func truncateTo(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

// CNForCertURI returns the correct common name for a given cert URI type. It
// doesn't work for CA Signing IDs since more context is needed and CA Providers
// always know their CN from their own context.
func CNForCertURI(uri CertURI) (string, error) {
	// Even though leafs should be from our own CSRs which should have the same CN
	// logic as here, override anyway to account for older version clients that
	// didn't include the Common Name in the CSR.
	switch id := uri.(type) {
	case *SpiffeIDService:
		return ServiceCN(id.Service, id.Host), nil
	case *SpiffeIDAgent:
		return AgentCN(id.Agent, id.Host), nil
	case *SpiffeIDSigning:
		return "", fmt.Errorf("CertURI is a SpiffeIDSigning, not enough context to generate Common Name")
	default:
		return "", fmt.Errorf("CertURI type not recognized")
	}
}
