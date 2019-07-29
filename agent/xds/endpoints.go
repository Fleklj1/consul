package xds

import (
	"errors"
	"fmt"

	envoy "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyendpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/gogo/protobuf/proto"

	"github.com/hashicorp/consul/agent/proxycfg"
	"github.com/hashicorp/consul/agent/structs"
	"github.com/hashicorp/consul/api"

	bexpr "github.com/hashicorp/go-bexpr"
)

// endpointsFromSnapshot returns the xDS API representation of the "endpoints"
func (s *Server) endpointsFromSnapshot(cfgSnap *proxycfg.ConfigSnapshot, token string) ([]proto.Message, error) {
	if cfgSnap == nil {
		return nil, errors.New("nil config given")
	}

	switch cfgSnap.Kind {
	case structs.ServiceKindConnectProxy:
		return s.endpointsFromSnapshotConnectProxy(cfgSnap, token)
	case structs.ServiceKindMeshGateway:
		return s.endpointsFromSnapshotMeshGateway(cfgSnap, token)
	default:
		return nil, fmt.Errorf("Invalid service kind: %v", cfgSnap.Kind)
	}
}

// endpointsFromSnapshotConnectProxy returns the xDS API representation of the "endpoints"
// (upstream instances) in the snapshot.
func (s *Server) endpointsFromSnapshotConnectProxy(cfgSnap *proxycfg.ConfigSnapshot, token string) ([]proto.Message, error) {
	// TODO(rb): this sizing is a low bound.
	resources := make([]proto.Message, 0, len(cfgSnap.ConnectProxy.UpstreamEndpoints))

	// TODO(rb): should naming from 1.5 -> 1.6 for clusters remain unchanged?

	for _, u := range cfgSnap.Proxy.Upstreams {
		id := u.Identifier()

		var chain *structs.CompiledDiscoveryChain
		if u.DestinationType != structs.UpstreamDestTypePreparedQuery {
			chain = cfgSnap.ConnectProxy.DiscoveryChain[id]
		}

		if chain == nil {
			// We ONLY want this branch for prepared queries.

			sni := UpstreamSNI(&u, "", cfgSnap)
			endpoints, ok := cfgSnap.ConnectProxy.UpstreamEndpoints[id]
			if ok {
				la := makeLoadAssignment(
					sni,
					0,
					[]loadAssignmentEndpointGroup{
						{Endpoints: endpoints},
					},
					cfgSnap.Datacenter,
				)
				resources = append(resources, la)
			}

		} else {
			// Newfangled discovery chain plumbing.

			chainEndpointMap, ok := cfgSnap.ConnectProxy.WatchedUpstreamEndpoints[id]
			if !ok {
				continue // skip the upstream (should not happen)
			}

			for target, node := range chain.GroupResolverNodes {
				groupResolver := node.GroupResolver
				failover := groupResolver.Failover

				endpoints, ok := chainEndpointMap[target]
				if !ok {
					continue // skip the cluster (should not happen)
				}

				var (
					endpointGroups         []loadAssignmentEndpointGroup
					overprovisioningFactor int
				)

				primaryGroup := loadAssignmentEndpointGroup{
					Endpoints:   endpoints,
					OnlyPassing: chain.SubsetDefinitionForTarget(target).OnlyPassing,
				}

				if failover != nil && len(failover.Targets) > 0 {
					endpointGroups = make([]loadAssignmentEndpointGroup, 0, len(failover.Targets)+1)

					endpointGroups = append(endpointGroups, primaryGroup)

					if failover.Definition.OverprovisioningFactor > 0 {
						overprovisioningFactor = failover.Definition.OverprovisioningFactor
					}
					if overprovisioningFactor <= 0 {
						// We choose such a large value here that the failover math should
						// in effect not happen until zero instances are healthy.
						overprovisioningFactor = 100000
					}

					for _, failTarget := range failover.Targets {
						failEndpoints, ok := chainEndpointMap[failTarget]
						if !ok {
							continue // skip the failover target (should not happen)
						}

						endpointGroups = append(endpointGroups, loadAssignmentEndpointGroup{
							Endpoints:   failEndpoints,
							OnlyPassing: chain.SubsetDefinitionForTarget(failTarget).OnlyPassing,
						})
					}
				} else {
					endpointGroups = append(endpointGroups, primaryGroup)
				}

				sni := TargetSNI(target, cfgSnap)

				la := makeLoadAssignment(
					sni,
					overprovisioningFactor,
					endpointGroups,
					cfgSnap.Datacenter,
				)
				resources = append(resources, la)
			}
		}
	}

	return resources, nil
}

func (s *Server) endpointsFromSnapshotMeshGateway(cfgSnap *proxycfg.ConfigSnapshot, token string) ([]proto.Message, error) {
	resources := make([]proto.Message, 0, len(cfgSnap.MeshGateway.GatewayGroups)+len(cfgSnap.MeshGateway.ServiceGroups))

	// generate the endpoints for the gateways in the remote datacenters
	for dc, endpoints := range cfgSnap.MeshGateway.GatewayGroups {
		clusterName := DatacenterSNI(dc, cfgSnap)
		la := makeLoadAssignment(
			clusterName,
			0,
			[]loadAssignmentEndpointGroup{
				{Endpoints: endpoints},
			},
			cfgSnap.Datacenter,
		)
		resources = append(resources, la)
	}

	// generate the endpoints for the local service groups
	for svc, endpoints := range cfgSnap.MeshGateway.ServiceGroups {
		clusterName := ServiceSNI(svc, "", "default", cfgSnap.Datacenter, cfgSnap)
		la := makeLoadAssignment(
			clusterName,
			0,
			[]loadAssignmentEndpointGroup{
				{Endpoints: endpoints},
			},
			cfgSnap.Datacenter,
		)
		resources = append(resources, la)
	}

	// generate the endpoints for the service subsets
	for svc, resolver := range cfgSnap.MeshGateway.ServiceResolvers {
		for subsetName, subset := range resolver.Subsets {
			clusterName := ServiceSNI(svc, subsetName, "default", cfgSnap.Datacenter, cfgSnap)

			endpoints := cfgSnap.MeshGateway.ServiceGroups[svc]

			// locally execute the subsets filter
			if subset.Filter != "" {
				filter, err := bexpr.CreateFilter(subset.Filter, nil, endpoints)
				if err != nil {
					return nil, err
				}

				raw, err := filter.Execute(endpoints)
				if err != nil {
					return nil, err
				}
				endpoints = raw.(structs.CheckServiceNodes)
			}

			la := makeLoadAssignment(
				clusterName,
				0,
				[]loadAssignmentEndpointGroup{
					{
						Endpoints:   endpoints,
						OnlyPassing: subset.OnlyPassing,
					},
				},
				cfgSnap.Datacenter,
			)
			resources = append(resources, la)
		}
	}

	return resources, nil
}

func makeEndpoint(clusterName, host string, port int) envoyendpoint.LbEndpoint {
	return envoyendpoint.LbEndpoint{
		HostIdentifier: &envoyendpoint.LbEndpoint_Endpoint{
			Endpoint: &envoyendpoint.Endpoint{
				Address: makeAddressPtr(host, port),
			},
		},
	}
}

type loadAssignmentEndpointGroup struct {
	Endpoints   structs.CheckServiceNodes
	OnlyPassing bool
}

func makeLoadAssignment(
	clusterName string,
	overprovisioningFactor int,
	endpointGroups []loadAssignmentEndpointGroup,
	localDatacenter string,
) *envoy.ClusterLoadAssignment {
	cla := &envoy.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints:   make([]envoyendpoint.LocalityLbEndpoints, 0, len(endpointGroups)),
	}
	if overprovisioningFactor > 0 {
		cla.Policy = &envoy.ClusterLoadAssignment_Policy{
			OverprovisioningFactor: makeUint32Value(overprovisioningFactor),
		}
	}

	for priority, endpointGroup := range endpointGroups {
		endpoints := endpointGroup.Endpoints
		es := make([]envoyendpoint.LbEndpoint, 0, len(endpoints))

		for _, ep := range endpoints {
			// TODO (mesh-gateway) - should we respect the translate_wan_addrs configuration here or just always use the wan for cross-dc?
			addr, port := ep.BestAddress(localDatacenter != ep.Node.Datacenter)
			healthStatus := envoycore.HealthStatus_HEALTHY
			weight := 1
			if ep.Service.Weights != nil {
				weight = ep.Service.Weights.Passing
			}

			for _, chk := range ep.Checks {
				if chk.Status == api.HealthCritical {
					healthStatus = envoycore.HealthStatus_UNHEALTHY
				}
				if endpointGroup.OnlyPassing && chk.Status != api.HealthPassing {
					healthStatus = envoycore.HealthStatus_UNHEALTHY
				}
				if chk.Status == api.HealthWarning && ep.Service.Weights != nil {
					weight = ep.Service.Weights.Warning
				}
			}
			// Make weights fit Envoy's limits. A zero weight means that either Warning
			// (likely) or Passing (weirdly) weight has been set to 0 effectively making
			// this instance unhealthy and should not be sent traffic.
			if weight < 1 {
				healthStatus = envoycore.HealthStatus_UNHEALTHY
				weight = 1
			}
			if weight > 128 {
				weight = 128
			}
			es = append(es, envoyendpoint.LbEndpoint{
				HostIdentifier: &envoyendpoint.LbEndpoint_Endpoint{
					Endpoint: &envoyendpoint.Endpoint{
						Address: makeAddressPtr(addr, port),
					},
				},
				HealthStatus:        healthStatus,
				LoadBalancingWeight: makeUint32Value(weight),
			})
		}

		cla.Endpoints = append(cla.Endpoints, envoyendpoint.LocalityLbEndpoints{
			Priority:    uint32(priority),
			LbEndpoints: es,
		})
	}

	return cla
}
