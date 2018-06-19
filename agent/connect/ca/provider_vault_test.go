package ca

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/hashicorp/consul/agent/connect"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func testVaultCluster(t *testing.T) (*VaultProvider, *vault.Core, net.Listener) {
	if err := vault.AddTestLogicalBackend("pki", pki.Factory); err != nil {
		t.Fatal(err)
	}
	core, _, token := vault.TestCoreUnsealedRaw(t)

	ln, addr := vaulthttp.TestServer(t, core)

	provider, err := NewVaultProvider(map[string]interface{}{
		"Address":             addr,
		"Token":               token,
		"RootPKIPath":         "pki-root/",
		"IntermediatePKIPath": "pki-intermediate/",
	}, "asdf")
	if err != nil {
		t.Fatal(err)
	}

	return provider, core, ln
}

func TestVaultCAProvider_Bootstrap(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	provider, core, listener := testVaultCluster(t)
	defer core.Shutdown()
	defer listener.Close()
	client, err := vaultapi.NewClient(&vaultapi.Config{
		Address: "http://" + listener.Addr().String(),
	})
	require.NoError(err)
	client.SetToken(provider.config.Token)

	cases := []struct {
		certFunc    func() (string, error)
		backendPath string
	}{
		{
			certFunc:    provider.ActiveRoot,
			backendPath: "pki-root/",
		},
		{
			certFunc:    provider.ActiveIntermediate,
			backendPath: "pki-intermediate/",
		},
	}

	// Verify the root and intermediate certs match the ones in the vault backends
	for _, tc := range cases {
		cert, err := tc.certFunc()
		require.NoError(err)
		req := client.NewRequest("GET", "/v1/"+tc.backendPath+"ca/pem")
		resp, err := client.RawRequest(req)
		require.NoError(err)
		bytes, err := ioutil.ReadAll(resp.Body)
		require.NoError(err)
		require.Equal(cert, string(bytes))

		// Should be a valid CA cert
		parsed, err := connect.ParseCert(cert)
		require.NoError(err)
		require.True(parsed.IsCA)
		require.Len(parsed.URIs, 1)
		require.Equal(parsed.URIs[0].String(), fmt.Sprintf("spiffe://%s.consul", provider.clusterId))
	}
}

func TestVaultCAProvider_SignLeaf(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	provider, core, listener := testVaultCluster(t)
	defer core.Shutdown()
	defer listener.Close()
	client, err := vaultapi.NewClient(&vaultapi.Config{
		Address: "http://" + listener.Addr().String(),
	})
	require.NoError(err)
	client.SetToken(provider.config.Token)

	spiffeService := &connect.SpiffeIDService{
		Host:       "node1",
		Namespace:  "default",
		Datacenter: "dc1",
		Service:    "foo",
	}

	// Generate a leaf cert for the service.
	var firstSerial uint64
	{
		raw, _ := connect.TestCSR(t, spiffeService)

		csr, err := connect.ParseCSR(raw)
		require.NoError(err)

		cert, err := provider.Sign(csr)
		require.NoError(err)

		parsed, err := connect.ParseCert(cert)
		require.NoError(err)
		require.Equal(parsed.URIs[0], spiffeService.URI())
		require.Equal(parsed.Subject.CommonName, "foo")
		firstSerial = parsed.SerialNumber.Uint64()

		// Ensure the cert is valid now and expires within the correct limit.
		require.True(parsed.NotAfter.Sub(time.Now()) < 3*24*time.Hour)
		require.True(parsed.NotBefore.Before(time.Now()))
	}

	// Generate a new cert for another service and make sure
	// the serial number is unique.
	spiffeService.Service = "bar"
	{
		raw, _ := connect.TestCSR(t, spiffeService)

		csr, err := connect.ParseCSR(raw)
		require.NoError(err)

		cert, err := provider.Sign(csr)
		require.NoError(err)

		parsed, err := connect.ParseCert(cert)
		require.NoError(err)
		require.Equal(parsed.URIs[0], spiffeService.URI())
		require.Equal(parsed.Subject.CommonName, "bar")
		require.NotEqual(firstSerial, parsed.SerialNumber.Uint64())

		// Ensure the cert is valid now and expires within the correct limit.
		require.True(parsed.NotAfter.Sub(time.Now()) < 3*24*time.Hour)
		require.True(parsed.NotBefore.Before(time.Now()))
	}
}

func TestVaultCAProvider_CrossSignCA(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	provider1, core1, listener1 := testVaultCluster(t)
	defer core1.Shutdown()
	defer listener1.Close()

	provider2, core2, listener2 := testVaultCluster(t)
	defer core2.Shutdown()
	defer listener2.Close()

	// Have provider2 generate a cross-signing CSR
	csr, err := provider2.GetCrossSigningCSR()
	require.NoError(err)
	oldSubject := csr.Subject.CommonName
	_, err = provider2.GenerateIntermediate()
	require.NoError(err)

	// Have provider1 cross sign our new CA cert.
	xcPEM, err := provider1.CrossSignCA(csr)
	require.NoError(err)
	xc, err := connect.ParseCert(xcPEM)
	require.NoError(err)

	rootPEM, err := provider1.ActiveRoot()
	require.NoError(err)
	root, err := connect.ParseCert(rootPEM)
	require.NoError(err)

	// AuthorityKeyID should be the signing root's, SubjectKeyId should be different.
	require.Equal(root.AuthorityKeyId, xc.AuthorityKeyId)
	require.NotEqual(root.SubjectKeyId, xc.SubjectKeyId)

	// Subject name should not have changed.
	require.NotEqual(root.Subject.CommonName, xc.Subject.CommonName)
	require.Equal(oldSubject, xc.Subject.CommonName)

	// Issuer should be the signing root.
	require.Equal(root.Issuer.CommonName, xc.Issuer.CommonName)

	// Get a leaf cert so we can verify against the cross-signed cert.
	spiffeService := &connect.SpiffeIDService{
		Host:       "node1",
		Namespace:  "default",
		Datacenter: "dc1",
		Service:    "foo",
	}
	raw, _ := connect.TestCSR(t, spiffeService)

	leafCsr, err := connect.ParseCSR(raw)
	require.NoError(err)

	leafPEM, err := provider2.Sign(leafCsr)
	require.NoError(err)

	cert, err := connect.ParseCert(leafPEM)
	require.NoError(err)

	intermediatePool := x509.NewCertPool()
	intermediatePool.AddCert(xc)

	rootPool := x509.NewCertPool()
	rootPool.AddCert(root)

	// Check that the leaf signed by the new cert can be verified by the
	// chain of cross-signed cert + old root (as would be the case on any
	// proxies that haven't received the new root yet) for backwards compatibility.
	_, err = cert.Verify(x509.VerifyOptions{
		Intermediates: intermediatePool,
		Roots:         rootPool,
	})
	require.NoError(err)
}
