package xds

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/mitchellh/copystructure"

	"github.com/stretchr/testify/require"

	envoy "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyendpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/hashicorp/consul/agent/proxycfg"
	"github.com/hashicorp/consul/agent/structs"
	testinf "github.com/mitchellh/go-testing-interface"
)

func Test_makeLoadAssignment(t *testing.T) {

	testCheckServiceNodes := structs.CheckServiceNodes{
		structs.CheckServiceNode{
			Node: &structs.Node{
				ID:         "node1-id",
				Node:       "node1",
				Address:    "10.10.10.10",
				Datacenter: "dc1",
			},
			Service: &structs.NodeService{
				Service: "web",
				Port:    1234,
			},
			Checks: structs.HealthChecks{
				&structs.HealthCheck{
					Node:    "node1",
					CheckID: "serfHealth",
					Status:  "passing",
				},
				&structs.HealthCheck{
					Node:      "node1",
					ServiceID: "web",
					CheckID:   "web:check",
					Status:    "passing",
				},
			},
		},
		structs.CheckServiceNode{
			Node: &structs.Node{
				ID:         "node2-id",
				Node:       "node2",
				Address:    "10.10.10.20",
				Datacenter: "dc1",
			},
			Service: &structs.NodeService{
				Service: "web",
				Port:    1234,
			},
			Checks: structs.HealthChecks{
				&structs.HealthCheck{
					Node:    "node2",
					CheckID: "serfHealth",
					Status:  "passing",
				},
				&structs.HealthCheck{
					Node:      "node2",
					ServiceID: "web",
					CheckID:   "web:check",
					Status:    "passing",
				},
			},
		},
	}

	testWeightedCheckServiceNodesRaw, err := copystructure.Copy(testCheckServiceNodes)
	require.NoError(t, err)
	testWeightedCheckServiceNodes := testWeightedCheckServiceNodesRaw.(structs.CheckServiceNodes)

	testWeightedCheckServiceNodes[0].Service.Weights = &structs.Weights{
		Passing: 10,
		Warning: 1,
	}
	testWeightedCheckServiceNodes[1].Service.Weights = &structs.Weights{
		Passing: 5,
		Warning: 0,
	}

	testWarningCheckServiceNodesRaw, err := copystructure.Copy(testWeightedCheckServiceNodes)
	require.NoError(t, err)
	testWarningCheckServiceNodes := testWarningCheckServiceNodesRaw.(structs.CheckServiceNodes)

	testWarningCheckServiceNodes[0].Checks[0].Status = "warning"
	testWarningCheckServiceNodes[1].Checks[0].Status = "warning"

	tests := []struct {
		name                   string
		clusterName            string
		overprovisioningFactor int
		endpoints              []structs.CheckServiceNodes
		want                   *envoy.ClusterLoadAssignment
	}{
		{
			name:        "no instances",
			clusterName: "service:test",
			endpoints: []structs.CheckServiceNodes{
				{},
			},
			want: &envoy.ClusterLoadAssignment{
				ClusterName: "service:test",
				Endpoints: []envoyendpoint.LocalityLbEndpoints{{
					LbEndpoints: []envoyendpoint.LbEndpoint{},
				}},
			},
		},
		{
			name:        "instances, no weights",
			clusterName: "service:test",
			endpoints: []structs.CheckServiceNodes{
				testCheckServiceNodes,
			},
			want: &envoy.ClusterLoadAssignment{
				ClusterName: "service:test",
				Endpoints: []envoyendpoint.LocalityLbEndpoints{{
					LbEndpoints: []envoyendpoint.LbEndpoint{
						envoyendpoint.LbEndpoint{
							HostIdentifier: &envoyendpoint.LbEndpoint_Endpoint{
								Endpoint: &envoyendpoint.Endpoint{
									Address: makeAddressPtr("10.10.10.10", 1234),
								}},
							HealthStatus:        core.HealthStatus_HEALTHY,
							LoadBalancingWeight: makeUint32Value(1),
						},
						envoyendpoint.LbEndpoint{
							HostIdentifier: &envoyendpoint.LbEndpoint_Endpoint{
								Endpoint: &envoyendpoint.Endpoint{
									Address: makeAddressPtr("10.10.10.20", 1234),
								}},
							HealthStatus:        core.HealthStatus_HEALTHY,
							LoadBalancingWeight: makeUint32Value(1),
						},
					},
				}},
			},
		},
		{
			name:        "instances, healthy weights",
			clusterName: "service:test",
			endpoints: []structs.CheckServiceNodes{
				testWeightedCheckServiceNodes,
			},
			want: &envoy.ClusterLoadAssignment{
				ClusterName: "service:test",
				Endpoints: []envoyendpoint.LocalityLbEndpoints{{
					LbEndpoints: []envoyendpoint.LbEndpoint{
						envoyendpoint.LbEndpoint{
							HostIdentifier: &envoyendpoint.LbEndpoint_Endpoint{
								Endpoint: &envoyendpoint.Endpoint{
									Address: makeAddressPtr("10.10.10.10", 1234),
								}},
							HealthStatus:        core.HealthStatus_HEALTHY,
							LoadBalancingWeight: makeUint32Value(10),
						},
						envoyendpoint.LbEndpoint{
							HostIdentifier: &envoyendpoint.LbEndpoint_Endpoint{
								Endpoint: &envoyendpoint.Endpoint{
									Address: makeAddressPtr("10.10.10.20", 1234),
								}},
							HealthStatus:        core.HealthStatus_HEALTHY,
							LoadBalancingWeight: makeUint32Value(5),
						},
					},
				}},
			},
		},
		{
			name:        "instances, warning weights",
			clusterName: "service:test",
			endpoints: []structs.CheckServiceNodes{
				testWarningCheckServiceNodes,
			},
			want: &envoy.ClusterLoadAssignment{
				ClusterName: "service:test",
				Endpoints: []envoyendpoint.LocalityLbEndpoints{{
					LbEndpoints: []envoyendpoint.LbEndpoint{
						envoyendpoint.LbEndpoint{
							HostIdentifier: &envoyendpoint.LbEndpoint_Endpoint{
								Endpoint: &envoyendpoint.Endpoint{
									Address: makeAddressPtr("10.10.10.10", 1234),
								}},
							HealthStatus:        core.HealthStatus_HEALTHY,
							LoadBalancingWeight: makeUint32Value(1),
						},
						envoyendpoint.LbEndpoint{
							HostIdentifier: &envoyendpoint.LbEndpoint_Endpoint{
								Endpoint: &envoyendpoint.Endpoint{
									Address: makeAddressPtr("10.10.10.20", 1234),
								}},
							HealthStatus:        core.HealthStatus_UNHEALTHY,
							LoadBalancingWeight: makeUint32Value(1),
						},
					},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeLoadAssignment(tt.clusterName, tt.overprovisioningFactor, tt.endpoints, "dc1")
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_endpointsFromSnapshot(t *testing.T) {

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
			name:   "defaults",
			create: proxycfg.TestConfigSnapshot,
			setup:  nil, // Default snapshot
		},
		{
			name:   "mesh-gateway",
			create: proxycfg.TestConfigSnapshotMeshGateway,
			setup:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			// Sanity check default with no overrides first
			snap := tt.create(t)

			// We need to replace the TLS certs with deterministic ones to make golden
			// files workable. Note we don't update these otherwise they'd change
			// golden files for every test case and so not be any use!
			if snap.ConnectProxy.Leaf != nil {
				snap.ConnectProxy.Leaf.CertPEM = golden(t, "test-leaf-cert", "")
				snap.ConnectProxy.Leaf.PrivateKeyPEM = golden(t, "test-leaf-key", "")
			}
			if snap.Roots != nil {
				snap.Roots.Roots[0].RootCert = golden(t, "test-root-cert", "")
			}

			if tt.setup != nil {
				tt.setup(snap)
			}

			// Need server just for logger dependency
			s := Server{Logger: log.New(os.Stderr, "", log.LstdFlags)}

			endpoints, err := s.endpointsFromSnapshot(snap, "my-token")
			require.NoError(err)
			r, err := createResponse(EndpointType, "00000001", "00000001", endpoints)
			require.NoError(err)

			gotJSON := responseToJSON(t, r)

			gName := tt.name
			if tt.overrideGoldenName != "" {
				gName = tt.overrideGoldenName
			}

			require.JSONEq(golden(t, path.Join("endpoints", gName), gotJSON), gotJSON)
		})
	}
}
