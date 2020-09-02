package xds

import (
	"path/filepath"
	"sort"
	"testing"
	"time"

	envoy "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/consul/agent/connect"
	"github.com/hashicorp/consul/agent/consul/discoverychain"
	"github.com/hashicorp/consul/agent/proxycfg"
	"github.com/hashicorp/consul/agent/structs"
	"github.com/hashicorp/consul/agent/xds/proxysupport"
	testinf "github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

func TestRoutesFromSnapshot(t *testing.T) {

	tests := []struct {
		name   string
		create func(t testinf.T) *proxycfg.ConfigSnapshot
		// Setup is called before the test starts. It is passed the snapshot from
		// create func and is allowed to modify it in any way to setup the
		// test input.
		setup              func(snap *proxycfg.ConfigSnapshot)
		overrideGoldenName string
	}{
		{
			name:   "defaults-no-chain",
			create: proxycfg.TestConfigSnapshot,
			setup:  nil, // Default snapshot
		},
		{
			name:   "connect-proxy-with-chain",
			create: proxycfg.TestConfigSnapshotDiscoveryChain,
			setup:  nil,
		},
		{
			name:   "connect-proxy-with-chain-external-sni",
			create: proxycfg.TestConfigSnapshotDiscoveryChainExternalSNI,
			setup:  nil,
		},
		{
			name:   "connect-proxy-with-chain-and-overrides",
			create: proxycfg.TestConfigSnapshotDiscoveryChainWithOverrides,
			setup:  nil,
		},
		{
			name:   "splitter-with-resolver-redirect",
			create: proxycfg.TestConfigSnapshotDiscoveryChain_SplitterWithResolverRedirectMultiDC,
			setup:  nil,
		},
		{
			name:   "connect-proxy-with-chain-and-splitter",
			create: proxycfg.TestConfigSnapshotDiscoveryChainWithSplitter,
			setup:  nil,
		},
		{
			name:   "connect-proxy-with-grpc-router",
			create: proxycfg.TestConfigSnapshotDiscoveryChainWithGRPCRouter,
			setup:  nil,
		},
		{
			name:   "connect-proxy-with-chain-and-router",
			create: proxycfg.TestConfigSnapshotDiscoveryChainWithRouter,
			setup:  nil,
		},
		{
			name:   "connect-proxy-lb-in-resolver",
			create: proxycfg.TestConfigSnapshotDiscoveryChainWithLB,
			setup:  nil,
		},
		// TODO(rb): test match stanza skipped for grpc
		// Start ingress gateway test cases
		{
			name:   "ingress-defaults-no-chain",
			create: proxycfg.TestConfigSnapshotIngressGateway,
			setup:  nil, // Default snapshot
		},
		{
			name:   "ingress-with-chain",
			create: proxycfg.TestConfigSnapshotIngress,
			setup:  nil,
		},
		{
			name:   "ingress-with-chain-external-sni",
			create: proxycfg.TestConfigSnapshotIngressExternalSNI,
			setup:  nil,
		},
		{
			name:   "ingress-with-chain-and-overrides",
			create: proxycfg.TestConfigSnapshotIngressWithOverrides,
			setup:  nil,
		},
		{
			name:   "ingress-splitter-with-resolver-redirect",
			create: proxycfg.TestConfigSnapshotIngress_SplitterWithResolverRedirectMultiDC,
			setup:  nil,
		},
		{
			name:   "ingress-with-chain-and-splitter",
			create: proxycfg.TestConfigSnapshotIngressWithSplitter,
			setup:  nil,
		},
		{
			name:   "ingress-with-grpc-router",
			create: proxycfg.TestConfigSnapshotIngressWithGRPCRouter,
			setup:  nil,
		},
		{
			name:   "ingress-with-chain-and-router",
			create: proxycfg.TestConfigSnapshotIngressWithRouter,
			setup:  nil,
		},
		{
			name:   "ingress-lb-in-resolver",
			create: proxycfg.TestConfigSnapshotIngressWithLB,
			setup:  nil,
		},
		{
			name:   "ingress-http-multiple-services",
			create: proxycfg.TestConfigSnapshotIngress_HTTPMultipleServices,
			setup: func(snap *proxycfg.ConfigSnapshot) {
				snap.IngressGateway.Upstreams = map[proxycfg.IngressListenerKey]structs.Upstreams{
					{Protocol: "http", Port: 8080}: {
						{
							DestinationName: "foo",
							LocalBindPort:   8080,
							IngressHosts: []string{
								"test1.example.com",
								"test2.example.com",
								"test2.example.com:8080",
							},
						},
						{
							DestinationName: "bar",
							LocalBindPort:   8080,
						},
					},
					{Protocol: "http", Port: 443}: {
						{
							DestinationName: "baz",
							LocalBindPort:   443,
						},
						{
							DestinationName: "qux",
							LocalBindPort:   443,
						},
					},
				}

				// We do not add baz/qux here so that we test the chain.IsDefault() case
				entries := []structs.ConfigEntry{
					&structs.ProxyConfigEntry{
						Kind: structs.ProxyDefaults,
						Name: structs.ProxyConfigGlobal,
						Config: map[string]interface{}{
							"protocol": "http",
						},
					},
					&structs.ServiceResolverConfigEntry{
						Kind:           structs.ServiceResolver,
						Name:           "foo",
						ConnectTimeout: 22 * time.Second,
					},
					&structs.ServiceResolverConfigEntry{
						Kind:           structs.ServiceResolver,
						Name:           "bar",
						ConnectTimeout: 22 * time.Second,
					},
				}
				fooChain := discoverychain.TestCompileConfigEntries(t, "foo", "default", "dc1", connect.TestClusterID+".consul", "dc1", nil, entries...)
				barChain := discoverychain.TestCompileConfigEntries(t, "bar", "default", "dc1", connect.TestClusterID+".consul", "dc1", nil, entries...)
				bazChain := discoverychain.TestCompileConfigEntries(t, "baz", "default", "dc1", connect.TestClusterID+".consul", "dc1", nil, entries...)
				quxChain := discoverychain.TestCompileConfigEntries(t, "qux", "default", "dc1", connect.TestClusterID+".consul", "dc1", nil, entries...)

				snap.IngressGateway.DiscoveryChain = map[string]*structs.CompiledDiscoveryChain{
					"foo": fooChain,
					"bar": barChain,
					"baz": bazChain,
					"qux": quxChain,
				}
			},
		},
		{
			name:   "terminating-gateway-lb-config",
			create: proxycfg.TestConfigSnapshotTerminatingGateway,
			setup: func(snap *proxycfg.ConfigSnapshot) {
				snap.TerminatingGateway.ServiceResolvers = map[structs.ServiceName]*structs.ServiceResolverConfigEntry{
					structs.NewServiceName("web", nil): {
						Kind:          structs.ServiceResolver,
						Name:          "web",
						DefaultSubset: "v2",
						Subsets: map[string]structs.ServiceResolverSubset{
							"v1": {
								Filter: "Service.Meta.Version == 1",
							},
							"v2": {
								Filter:      "Service.Meta.Version == 2",
								OnlyPassing: true,
							},
						},
						LoadBalancer: &structs.LoadBalancer{
							EnvoyConfig: &structs.EnvoyLBConfig{
								Policy: "ring_hash",
								RingHashConfig: &structs.RingHashConfig{
									MinimumRingSize: 20,
									MaximumRingSize: 50,
								},
								HashPolicies: []structs.HashPolicy{
									{
										Field:      structs.HashPolicyCookie,
										FieldValue: "chocolate-chip",
										Terminal:   true,
									},
									{
										Field:      structs.HashPolicyHeader,
										FieldValue: "x-user-id",
									},
									{
										SourceIP: true,
										Terminal: true,
									},
								},
							},
						},
					},
				}
			},
		},
	}

	for _, envoyVersion := range proxysupport.EnvoyVersions {
		sf, err := determineSupportedProxyFeaturesFromString(envoyVersion)
		require.NoError(t, err)
		t.Run("envoy-"+envoyVersion, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					require := require.New(t)

					// Sanity check default with no overrides first
					snap := tt.create(t)

					// We need to replace the TLS certs with deterministic ones to make golden
					// files workable. Note we don't update these otherwise they'd change
					// golden files for every test case and so not be any use!
					setupTLSRootsAndLeaf(t, snap)

					if tt.setup != nil {
						tt.setup(snap)
					}

					cInfo := connectionInfo{
						Token:         "my-token",
						ProxyFeatures: sf,
					}
					routes, err := routesFromSnapshot(cInfo, snap)
					require.NoError(err)
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].(*envoy.RouteConfiguration).Name < routes[j].(*envoy.RouteConfiguration).Name
					})
					r, err := createResponse(RouteType, "00000001", "00000001", routes)
					require.NoError(err)

					gotJSON := responseToJSON(t, r)

					gName := tt.name
					if tt.overrideGoldenName != "" {
						gName = tt.overrideGoldenName
					}

					require.JSONEq(goldenEnvoy(t, filepath.Join("routes", gName), envoyVersion, gotJSON), gotJSON)
				})
			}
		})
	}
}

func TestEnvoyLBConfig_InjectToRouteAction(t *testing.T) {
	var tests = []struct {
		name     string
		lb       *structs.EnvoyLBConfig
		expected envoyroute.RouteAction
	}{
		{
			name: "empty",
			lb: &structs.EnvoyLBConfig{
				Policy: "",
			},
			// we only modify route actions for hash-based LB policies
			expected: envoyroute.RouteAction{},
		},
		{
			name: "least request",
			lb: &structs.EnvoyLBConfig{
				Policy: structs.LBPolicyLeastRequest,
				LeastRequestConfig: &structs.LeastRequestConfig{
					ChoiceCount: 3,
				},
			},
			// we only modify route actions for hash-based LB policies
			expected: envoyroute.RouteAction{},
		},
		{
			name: "headers",
			lb: &structs.EnvoyLBConfig{
				Policy: "ring_hash",
				RingHashConfig: &structs.RingHashConfig{
					MinimumRingSize: 3,
					MaximumRingSize: 7,
				},
				HashPolicies: []structs.HashPolicy{
					{
						Field:      structs.HashPolicyHeader,
						FieldValue: "x-route-key",
						Terminal:   true,
					},
				},
			},
			expected: envoyroute.RouteAction{
				HashPolicy: []*envoyroute.RouteAction_HashPolicy{
					{
						PolicySpecifier: &envoyroute.RouteAction_HashPolicy_Header_{
							Header: &envoyroute.RouteAction_HashPolicy_Header{
								HeaderName: "x-route-key",
							},
						},
						Terminal: true,
					},
				},
			},
		},
		{
			name: "cookies",
			lb: &structs.EnvoyLBConfig{
				Policy: structs.LBPolicyMaglev,
				HashPolicies: []structs.HashPolicy{
					{
						Field:      structs.HashPolicyCookie,
						FieldValue: "red-velvet",
						Terminal:   true,
					},
					{
						Field:      structs.HashPolicyCookie,
						FieldValue: "oatmeal",
					},
				},
			},
			expected: envoyroute.RouteAction{
				HashPolicy: []*envoyroute.RouteAction_HashPolicy{
					{
						PolicySpecifier: &envoyroute.RouteAction_HashPolicy_Cookie_{
							Cookie: &envoyroute.RouteAction_HashPolicy_Cookie{
								Name: "red-velvet",
							},
						},
						Terminal: true,
					},
					{
						PolicySpecifier: &envoyroute.RouteAction_HashPolicy_Cookie_{
							Cookie: &envoyroute.RouteAction_HashPolicy_Cookie{
								Name: "oatmeal",
							},
						},
					},
				},
			},
		},
		{
			name: "source addr",
			lb: &structs.EnvoyLBConfig{
				Policy: structs.LBPolicyMaglev,
				HashPolicies: []structs.HashPolicy{
					{
						SourceIP: true,
						Terminal: true,
					},
				},
			},
			expected: envoyroute.RouteAction{
				HashPolicy: []*envoyroute.RouteAction_HashPolicy{
					{
						PolicySpecifier: &envoyroute.RouteAction_HashPolicy_ConnectionProperties_{
							ConnectionProperties: &envoyroute.RouteAction_HashPolicy_ConnectionProperties{
								SourceIp: true,
							},
						},
						Terminal: true,
					},
				},
			},
		},
		{
			name: "kitchen sink",
			lb: &structs.EnvoyLBConfig{
				Policy: structs.LBPolicyMaglev,
				HashPolicies: []structs.HashPolicy{
					{
						SourceIP: true,
						Terminal: true,
					},
					{
						Field:      structs.HashPolicyCookie,
						FieldValue: "oatmeal",
						CookieConfig: &structs.CookieConfig{
							TTL:  10 * time.Second,
							Path: "/oven",
						},
					},
					{
						Field:      structs.HashPolicyHeader,
						FieldValue: "special-header",
						Terminal:   true,
					},
				},
			},
			expected: envoyroute.RouteAction{
				HashPolicy: []*envoyroute.RouteAction_HashPolicy{
					{
						PolicySpecifier: &envoyroute.RouteAction_HashPolicy_ConnectionProperties_{
							ConnectionProperties: &envoyroute.RouteAction_HashPolicy_ConnectionProperties{
								SourceIp: true,
							},
						},
						Terminal: true,
					},
					{
						PolicySpecifier: &envoyroute.RouteAction_HashPolicy_Cookie_{
							Cookie: &envoyroute.RouteAction_HashPolicy_Cookie{
								Name: "oatmeal",
								Ttl:  ptypes.DurationProto(10 * time.Second),
								Path: "/oven",
							},
						},
					},
					{
						PolicySpecifier: &envoyroute.RouteAction_HashPolicy_Header_{
							Header: &envoyroute.RouteAction_HashPolicy_Header{
								HeaderName: "special-header",
							},
						},
						Terminal: true,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var ra envoyroute.RouteAction
			err := injectLBToRouteAction(tc.lb, &ra)
			require.NoError(t, err)

			require.Equal(t, &tc.expected, &ra)
		})
	}
}
