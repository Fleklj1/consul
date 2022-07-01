package pbconfigentry

import (
	"fmt"

	"github.com/hashicorp/consul/acl"
	"github.com/hashicorp/consul/agent/structs"
	"github.com/hashicorp/consul/proto/pbcommon"
	"github.com/hashicorp/consul/types"
)

func ConfigEntryToStructs(s *ConfigEntry) structs.ConfigEntry {
	switch s.Kind {
	case Kind_KindMeshConfig:
		var target structs.MeshConfigEntry
		MeshConfigToStructs(s.GetMeshConfig(), &target)
		pbcommon.RaftIndexToStructs(s.RaftIndex, &target.RaftIndex)
		pbcommon.EnterpriseMetaToStructs(s.EnterpriseMeta, &target.EnterpriseMeta)
		return &target
	case Kind_KindServiceResolver:
		var target structs.ServiceResolverConfigEntry
		target.Name = s.Name

		ServiceResolverToStructs(s.GetServiceResolver(), &target)
		pbcommon.RaftIndexToStructs(s.RaftIndex, &target.RaftIndex)
		pbcommon.EnterpriseMetaToStructs(s.EnterpriseMeta, &target.EnterpriseMeta)
		return &target
	case Kind_KindIngressGateway:
		var target structs.IngressGatewayConfigEntry
		target.Name = s.Name

		IngressGatewayToStructs(s.GetIngressGateway(), &target)
		pbcommon.RaftIndexToStructs(s.RaftIndex, &target.RaftIndex)
		pbcommon.EnterpriseMetaToStructs(s.EnterpriseMeta, &target.EnterpriseMeta)
		return &target
	default:
		panic(fmt.Sprintf("unable to convert ConfigEntry of kind %s to structs", s.Kind))
	}
}

func ConfigEntryFromStructs(s structs.ConfigEntry) *ConfigEntry {
	configEntry := &ConfigEntry{
		Name:           s.GetName(),
		EnterpriseMeta: pbcommon.NewEnterpriseMetaFromStructs(*s.GetEnterpriseMeta()),
	}

	var raftIndex pbcommon.RaftIndex
	pbcommon.RaftIndexFromStructs(s.GetRaftIndex(), &raftIndex)
	configEntry.RaftIndex = &raftIndex

	switch v := s.(type) {
	case *structs.MeshConfigEntry:
		var meshConfig MeshConfig
		MeshConfigFromStructs(v, &meshConfig)

		configEntry.Kind = Kind_KindMeshConfig
		configEntry.Entry = &ConfigEntry_MeshConfig{
			MeshConfig: &meshConfig,
		}
	case *structs.ServiceResolverConfigEntry:
		var serviceResolver ServiceResolver
		ServiceResolverFromStructs(v, &serviceResolver)

		configEntry.Kind = Kind_KindServiceResolver
		configEntry.Entry = &ConfigEntry_ServiceResolver{
			ServiceResolver: &serviceResolver,
		}
	case *structs.IngressGatewayConfigEntry:
		var ingressGateway IngressGateway
		IngressGatewayFromStructs(v, &ingressGateway)

		configEntry.Kind = Kind_KindIngressGateway
		configEntry.Entry = &ConfigEntry_IngressGateway{
			IngressGateway: &ingressGateway,
		}
	default:
		panic(fmt.Sprintf("unable to convert %T to proto", s))
	}

	return configEntry
}

func tlsVersionToStructs(s string) types.TLSVersion {
	return types.TLSVersion(s)
}

func tlsVersionFromStructs(t types.TLSVersion) string {
	return t.String()
}

func cipherSuitesToStructs(cs []string) []types.TLSCipherSuite {
	cipherSuites := make([]types.TLSCipherSuite, len(cs))
	for idx, suite := range cs {
		cipherSuites[idx] = types.TLSCipherSuite(suite)
	}
	return cipherSuites
}

func cipherSuitesFromStructs(cs []types.TLSCipherSuite) []string {
	cipherSuites := make([]string, len(cs))
	for idx, suite := range cs {
		cipherSuites[idx] = suite.String()
	}
	return cipherSuites
}

func enterpriseMetaToStructs(m *pbcommon.EnterpriseMeta) acl.EnterpriseMeta {
	var entMeta acl.EnterpriseMeta
	pbcommon.EnterpriseMetaToStructs(m, &entMeta)
	return entMeta
}

func enterpriseMetaFromStructs(m acl.EnterpriseMeta) *pbcommon.EnterpriseMeta {
	return pbcommon.NewEnterpriseMetaFromStructs(m)
}
