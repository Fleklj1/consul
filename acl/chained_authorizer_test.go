package acl

import (
	"testing"
)

type testAuthorizer EnforcementDecision

func (authz testAuthorizer) ACLRead(*EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) ACLWrite(*EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) AgentRead(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) AgentWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) EventRead(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) EventWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) IntentionDefaultAllow(*EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) IntentionRead(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) IntentionWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) KeyList(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) KeyRead(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) KeyWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) KeyWritePrefix(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) KeyringRead(*EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) KeyringWrite(*EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) NodeRead(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) NodeWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) OperatorRead(*EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) OperatorWrite(*EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) PreparedQueryRead(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) PreparedQueryWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) ServiceRead(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) ServiceWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) SessionRead(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) SessionWrite(string, *EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}
func (authz testAuthorizer) Snapshot(*EnterpriseAuthorizerContext) EnforcementDecision {
	return EnforcementDecision(authz)
}

func TestChainedAuthorizer(t *testing.T) {
	t.Parallel()

	t.Run("No Authorizers", func(t *testing.T) {
		t.Parallel()

		authz := NewChainedAuthorizer([]Authorizer{})
		checkDenyACLRead(t, authz, "foo", nil)
		checkDenyACLWrite(t, authz, "foo", nil)
		checkDenyAgentRead(t, authz, "foo", nil)
		checkDenyAgentWrite(t, authz, "foo", nil)
		checkDenyEventRead(t, authz, "foo", nil)
		checkDenyEventWrite(t, authz, "foo", nil)
		checkDenyIntentionDefaultAllow(t, authz, "foo", nil)
		checkDenyIntentionRead(t, authz, "foo", nil)
		checkDenyIntentionWrite(t, authz, "foo", nil)
		checkDenyKeyRead(t, authz, "foo", nil)
		checkDenyKeyList(t, authz, "foo", nil)
		checkDenyKeyringRead(t, authz, "foo", nil)
		checkDenyKeyringWrite(t, authz, "foo", nil)
		checkDenyKeyWrite(t, authz, "foo", nil)
		checkDenyKeyWritePrefix(t, authz, "foo", nil)
		checkDenyNodeRead(t, authz, "foo", nil)
		checkDenyNodeWrite(t, authz, "foo", nil)
		checkDenyOperatorRead(t, authz, "foo", nil)
		checkDenyOperatorWrite(t, authz, "foo", nil)
		checkDenyPreparedQueryRead(t, authz, "foo", nil)
		checkDenyPreparedQueryWrite(t, authz, "foo", nil)
		checkDenyServiceRead(t, authz, "foo", nil)
		checkDenyServiceWrite(t, authz, "foo", nil)
		checkDenySessionRead(t, authz, "foo", nil)
		checkDenySessionWrite(t, authz, "foo", nil)
		checkDenySnapshot(t, authz, "foo", nil)
	})

	t.Run("Authorizer Defaults", func(t *testing.T) {
		t.Parallel()

		authz := NewChainedAuthorizer([]Authorizer{testAuthorizer(Default)})
		checkDenyACLRead(t, authz, "foo", nil)
		checkDenyACLWrite(t, authz, "foo", nil)
		checkDenyAgentRead(t, authz, "foo", nil)
		checkDenyAgentWrite(t, authz, "foo", nil)
		checkDenyEventRead(t, authz, "foo", nil)
		checkDenyEventWrite(t, authz, "foo", nil)
		checkDenyIntentionDefaultAllow(t, authz, "foo", nil)
		checkDenyIntentionRead(t, authz, "foo", nil)
		checkDenyIntentionWrite(t, authz, "foo", nil)
		checkDenyKeyRead(t, authz, "foo", nil)
		checkDenyKeyList(t, authz, "foo", nil)
		checkDenyKeyringRead(t, authz, "foo", nil)
		checkDenyKeyringWrite(t, authz, "foo", nil)
		checkDenyKeyWrite(t, authz, "foo", nil)
		checkDenyKeyWritePrefix(t, authz, "foo", nil)
		checkDenyNodeRead(t, authz, "foo", nil)
		checkDenyNodeWrite(t, authz, "foo", nil)
		checkDenyOperatorRead(t, authz, "foo", nil)
		checkDenyOperatorWrite(t, authz, "foo", nil)
		checkDenyPreparedQueryRead(t, authz, "foo", nil)
		checkDenyPreparedQueryWrite(t, authz, "foo", nil)
		checkDenyServiceRead(t, authz, "foo", nil)
		checkDenyServiceWrite(t, authz, "foo", nil)
		checkDenySessionRead(t, authz, "foo", nil)
		checkDenySessionWrite(t, authz, "foo", nil)
		checkDenySnapshot(t, authz, "foo", nil)
	})

	t.Run("Authorizer No Defaults", func(t *testing.T) {
		t.Parallel()

		authz := NewChainedAuthorizer([]Authorizer{testAuthorizer(Allow)})
		checkAllowACLRead(t, authz, "foo", nil)
		checkAllowACLWrite(t, authz, "foo", nil)
		checkAllowAgentRead(t, authz, "foo", nil)
		checkAllowAgentWrite(t, authz, "foo", nil)
		checkAllowEventRead(t, authz, "foo", nil)
		checkAllowEventWrite(t, authz, "foo", nil)
		checkAllowIntentionDefaultAllow(t, authz, "foo", nil)
		checkAllowIntentionRead(t, authz, "foo", nil)
		checkAllowIntentionWrite(t, authz, "foo", nil)
		checkAllowKeyRead(t, authz, "foo", nil)
		checkAllowKeyList(t, authz, "foo", nil)
		checkAllowKeyringRead(t, authz, "foo", nil)
		checkAllowKeyringWrite(t, authz, "foo", nil)
		checkAllowKeyWrite(t, authz, "foo", nil)
		checkAllowKeyWritePrefix(t, authz, "foo", nil)
		checkAllowNodeRead(t, authz, "foo", nil)
		checkAllowNodeWrite(t, authz, "foo", nil)
		checkAllowOperatorRead(t, authz, "foo", nil)
		checkAllowOperatorWrite(t, authz, "foo", nil)
		checkAllowPreparedQueryRead(t, authz, "foo", nil)
		checkAllowPreparedQueryWrite(t, authz, "foo", nil)
		checkAllowServiceRead(t, authz, "foo", nil)
		checkAllowServiceWrite(t, authz, "foo", nil)
		checkAllowSessionRead(t, authz, "foo", nil)
		checkAllowSessionWrite(t, authz, "foo", nil)
		checkAllowSnapshot(t, authz, "foo", nil)
	})

	t.Run("First Found", func(t *testing.T) {
		t.Parallel()

		authz := NewChainedAuthorizer([]Authorizer{testAuthorizer(Deny), testAuthorizer(Allow)})
		checkDenyACLRead(t, authz, "foo", nil)
		checkDenyACLWrite(t, authz, "foo", nil)
		checkDenyAgentRead(t, authz, "foo", nil)
		checkDenyAgentWrite(t, authz, "foo", nil)
		checkDenyEventRead(t, authz, "foo", nil)
		checkDenyEventWrite(t, authz, "foo", nil)
		checkDenyIntentionDefaultAllow(t, authz, "foo", nil)
		checkDenyIntentionRead(t, authz, "foo", nil)
		checkDenyIntentionWrite(t, authz, "foo", nil)
		checkDenyKeyRead(t, authz, "foo", nil)
		checkDenyKeyList(t, authz, "foo", nil)
		checkDenyKeyringRead(t, authz, "foo", nil)
		checkDenyKeyringWrite(t, authz, "foo", nil)
		checkDenyKeyWrite(t, authz, "foo", nil)
		checkDenyKeyWritePrefix(t, authz, "foo", nil)
		checkDenyNodeRead(t, authz, "foo", nil)
		checkDenyNodeWrite(t, authz, "foo", nil)
		checkDenyOperatorRead(t, authz, "foo", nil)
		checkDenyOperatorWrite(t, authz, "foo", nil)
		checkDenyPreparedQueryRead(t, authz, "foo", nil)
		checkDenyPreparedQueryWrite(t, authz, "foo", nil)
		checkDenyServiceRead(t, authz, "foo", nil)
		checkDenyServiceWrite(t, authz, "foo", nil)
		checkDenySessionRead(t, authz, "foo", nil)
		checkDenySessionWrite(t, authz, "foo", nil)
		checkDenySnapshot(t, authz, "foo", nil)

		authz = NewChainedAuthorizer([]Authorizer{testAuthorizer(Default), testAuthorizer(Allow)})
		checkAllowACLRead(t, authz, "foo", nil)
		checkAllowACLWrite(t, authz, "foo", nil)
		checkAllowAgentRead(t, authz, "foo", nil)
		checkAllowAgentWrite(t, authz, "foo", nil)
		checkAllowEventRead(t, authz, "foo", nil)
		checkAllowEventWrite(t, authz, "foo", nil)
		checkAllowIntentionDefaultAllow(t, authz, "foo", nil)
		checkAllowIntentionRead(t, authz, "foo", nil)
		checkAllowIntentionWrite(t, authz, "foo", nil)
		checkAllowKeyRead(t, authz, "foo", nil)
		checkAllowKeyList(t, authz, "foo", nil)
		checkAllowKeyringRead(t, authz, "foo", nil)
		checkAllowKeyringWrite(t, authz, "foo", nil)
		checkAllowKeyWrite(t, authz, "foo", nil)
		checkAllowKeyWritePrefix(t, authz, "foo", nil)
		checkAllowNodeRead(t, authz, "foo", nil)
		checkAllowNodeWrite(t, authz, "foo", nil)
		checkAllowOperatorRead(t, authz, "foo", nil)
		checkAllowOperatorWrite(t, authz, "foo", nil)
		checkAllowPreparedQueryRead(t, authz, "foo", nil)
		checkAllowPreparedQueryWrite(t, authz, "foo", nil)
		checkAllowServiceRead(t, authz, "foo", nil)
		checkAllowServiceWrite(t, authz, "foo", nil)
		checkAllowSessionRead(t, authz, "foo", nil)
		checkAllowSessionWrite(t, authz, "foo", nil)
		checkAllowSnapshot(t, authz, "foo", nil)
	})

}
