package ca

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/consul/agent/connect"
	"github.com/hashicorp/consul/agent/consul/state"
	"github.com/hashicorp/consul/agent/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type consulCAMockDelegate struct {
	state *state.Store
}

func (c *consulCAMockDelegate) State() *state.Store {
	return c.state
}

func (c *consulCAMockDelegate) ApplyCARequest(req *structs.CARequest) error {
	idx, _, err := c.state.CAConfig()
	if err != nil {
		return err
	}

	switch req.Op {
	case structs.CAOpSetProviderState:
		_, err := c.state.CASetProviderState(idx+1, req.ProviderState)
		if err != nil {
			return err
		}

		return nil
	case structs.CAOpDeleteProviderState:
		if err := c.state.CADeleteProviderState(req.ProviderState.ID); err != nil {
			return err
		}

		return nil
	default:
		return fmt.Errorf("Invalid CA operation '%s'", req.Op)
	}
}

func newMockDelegate(t *testing.T, conf *structs.CAConfiguration) *consulCAMockDelegate {
	s, err := state.NewStateStore(nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if s == nil {
		t.Fatalf("missing state store")
	}
	if err := s.CASetConfig(conf.RaftIndex.CreateIndex, conf); err != nil {
		t.Fatalf("err: %s", err)
	}

	return &consulCAMockDelegate{s}
}

func testConsulCAConfig() *structs.CAConfiguration {
	return &structs.CAConfiguration{
		ClusterID: "asdf",
		Provider:  "consul",
		Config:    map[string]interface{}{},
	}
}

func TestConsulCAProvider_Bootstrap(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	conf := testConsulCAConfig()
	delegate := newMockDelegate(t, conf)

	provider, err := NewConsulProvider(conf.Config, delegate)
	assert.NoError(err)

	root, err := provider.ActiveRoot()
	assert.NoError(err)

	// Intermediate should be the same cert.
	inter, err := provider.ActiveIntermediate()
	assert.NoError(err)
	assert.Equal(root, inter)

	// Should be a valid cert
	parsed, err := connect.ParseCert(root)
	assert.NoError(err)
	assert.Equal(parsed.URIs[0].String(), fmt.Sprintf("spiffe://%s.consul", conf.ClusterID))
}

func TestConsulCAProvider_Bootstrap_WithCert(t *testing.T) {
	t.Parallel()

	// Make sure setting a custom private key/root cert works.
	assert := assert.New(t)
	rootCA := connect.TestCA(t, nil)
	conf := testConsulCAConfig()
	conf.Config = map[string]interface{}{
		"PrivateKey": rootCA.SigningKey,
		"RootCert":   rootCA.RootCert,
	}
	delegate := newMockDelegate(t, conf)

	provider, err := NewConsulProvider(conf.Config, delegate)
	assert.NoError(err)

	root, err := provider.ActiveRoot()
	assert.NoError(err)
	assert.Equal(root, rootCA.RootCert)
}

func TestConsulCAProvider_SignLeaf(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	conf := testConsulCAConfig()
	delegate := newMockDelegate(t, conf)

	provider, err := NewConsulProvider(conf.Config, delegate)
	assert.NoError(err)

	spiffeService := &connect.SpiffeIDService{
		Host:       "node1",
		Namespace:  "default",
		Datacenter: "dc1",
		Service:    "foo",
	}

	// Generate a leaf cert for the service.
	{
		raw, _ := connect.TestCSR(t, spiffeService)

		csr, err := connect.ParseCSR(raw)
		assert.NoError(err)

		cert, err := provider.Sign(csr)
		assert.NoError(err)

		parsed, err := connect.ParseCert(cert)
		assert.NoError(err)
		assert.Equal(parsed.URIs[0], spiffeService.URI())
		assert.Equal(parsed.Subject.CommonName, "foo")
		assert.Equal(uint64(2), parsed.SerialNumber.Uint64())

		// Ensure the cert is valid now and expires within the correct limit.
		assert.True(parsed.NotAfter.Sub(time.Now()) < 3*24*time.Hour)
		assert.True(parsed.NotBefore.Before(time.Now()))
	}

	// Generate a new cert for another service and make sure
	// the serial number is incremented.
	spiffeService.Service = "bar"
	{
		raw, _ := connect.TestCSR(t, spiffeService)

		csr, err := connect.ParseCSR(raw)
		assert.NoError(err)

		cert, err := provider.Sign(csr)
		assert.NoError(err)

		parsed, err := connect.ParseCert(cert)
		assert.NoError(err)
		assert.Equal(parsed.URIs[0], spiffeService.URI())
		assert.Equal(parsed.Subject.CommonName, "bar")
		assert.Equal(parsed.SerialNumber.Uint64(), uint64(2))

		// Ensure the cert is valid now and expires within the correct limit.
		assert.True(parsed.NotAfter.Sub(time.Now()) < 3*24*time.Hour)
		assert.True(parsed.NotBefore.Before(time.Now()))
	}
}

func TestConsulCAProvider_CrossSignCA(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	conf1 := testConsulCAConfig()
	delegate1 := newMockDelegate(t, conf1)
	provider1, err := NewConsulProvider(conf1.Config, delegate1)

	conf2 := testConsulCAConfig()
	conf2.CreateIndex = 10
	delegate2 := newMockDelegate(t, conf2)
	provider2, err := NewConsulProvider(conf2.Config, delegate2)
	require.NoError(err)

	require.NoError(err)

	// Have provider2 generate a cross-signing CSR
	csr, err := provider2.GetCrossSigningCSR()
	require.NoError(err)
	oldSubject := csr.Subject.CommonName

	// Have the provider cross sign our new CA cert.
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
}
