// Code generated by mog. DO NOT EDIT.

package pbcommon

import "github.com/hashicorp/consul/agent/structs"

func QueryMetaToStructs(s *QueryMeta, t *structs.QueryMeta) {
	if s == nil {
		return
	}
	t.Index = s.Index
	t.LastContact = structs.DurationFromProto(s.LastContact)
	t.KnownLeader = s.KnownLeader
	t.ConsistencyLevel = s.ConsistencyLevel
	t.ResultsFilteredByACLs = s.ResultsFilteredByACLs
}
func QueryMetaFromStructs(t *structs.QueryMeta, s *QueryMeta) {
	if s == nil {
		return
	}
	s.Index = t.Index
	s.LastContact = structs.DurationToProto(t.LastContact)
	s.KnownLeader = t.KnownLeader
	s.ConsistencyLevel = t.ConsistencyLevel
	s.ResultsFilteredByACLs = t.ResultsFilteredByACLs
}
func QueryOptionsToStructs(s *QueryOptions, t *structs.QueryOptions) {
	if s == nil {
		return
	}
	t.Token = s.Token
	t.MinQueryIndex = s.MinQueryIndex
	t.MaxQueryTime = structs.DurationFromProto(s.MaxQueryTime)
	t.AllowStale = s.AllowStale
	t.RequireConsistent = s.RequireConsistent
	t.UseCache = s.UseCache
	t.MaxStaleDuration = structs.DurationFromProto(s.MaxStaleDuration)
	t.MaxAge = structs.DurationFromProto(s.MaxAge)
	t.MustRevalidate = s.MustRevalidate
	t.Filter = s.Filter
}
func QueryOptionsFromStructs(t *structs.QueryOptions, s *QueryOptions) {
	if s == nil {
		return
	}
	s.Token = t.Token
	s.MinQueryIndex = t.MinQueryIndex
	s.MaxQueryTime = structs.DurationToProto(t.MaxQueryTime)
	s.AllowStale = t.AllowStale
	s.RequireConsistent = t.RequireConsistent
	s.UseCache = t.UseCache
	s.MaxStaleDuration = structs.DurationToProto(t.MaxStaleDuration)
	s.MaxAge = structs.DurationToProto(t.MaxAge)
	s.MustRevalidate = t.MustRevalidate
	s.Filter = t.Filter
}
