package token

import (
	"sync"

	"crypto/subtle"
)

type TokenSource bool

const (
	TokenSourceConfig TokenSource = false
	TokenSourceAPI    TokenSource = true
)

// Store is used to hold the special ACL tokens used by Consul agents. It is
// designed to update the tokens on the fly, so the token store itself should be
// plumbed around and used to get tokens at runtime, don't save the resulting
// tokens.
type Store struct {
	// l synchronizes access to the token store.
	l sync.RWMutex

	// userToken is passed along for requests when the user didn't supply a
	// token, and may be left blank to use the anonymous token. This will
	// also be used for agent operations if the agent token isn't set.
	userToken string

	// userTokenSource indicates where this token originated from
	userTokenSource TokenSource

	// agentToken is used for internal agent operations like self-registering
	// with the catalog and anti-entropy, but should never be used for
	// user-initiated operations.
	agentToken string

	// agentTokenSource indicates where this token originated from
	agentTokenSource TokenSource

	// agentMasterToken is a special token that's only used locally for
	// access to the /v1/agent utility operations if the servers aren't
	// available.
	agentMasterToken string

	// agentMasterTokenSource indicates where this token originated from
	agentMasterTokenSource TokenSource

	// replicationToken is a special token that's used by servers to
	// replicate data from the primary datacenter.
	replicationToken string

	// replicationTokenSource indicates where this token originated from
	replicationTokenSource TokenSource

	// enterpriseTokens contains tokens only used in consul-enterprise
	enterpriseTokens
}

// UpdateUserToken replaces the current user token in the store.
// Returns true if it was changed.
func (t *Store) UpdateUserToken(token string, source TokenSource) bool {
	t.l.Lock()
	changed := (t.userToken != token || t.userTokenSource != source)
	t.userToken = token
	t.userTokenSource = source
	t.l.Unlock()
	return changed
}

// UpdateAgentToken replaces the current agent token in the store.
// Returns true if it was changed.
func (t *Store) UpdateAgentToken(token string, source TokenSource) bool {
	t.l.Lock()
	changed := (t.agentToken != token || t.agentTokenSource != source)
	t.agentToken = token
	t.agentTokenSource = source
	t.l.Unlock()
	return changed
}

// UpdateAgentMasterToken replaces the current agent master token in the store.
// Returns true if it was changed.
func (t *Store) UpdateAgentMasterToken(token string, source TokenSource) bool {
	t.l.Lock()
	changed := (t.agentMasterToken != token || t.agentMasterTokenSource != source)
	t.agentMasterToken = token
	t.agentMasterTokenSource = source
	t.l.Unlock()
	return changed
}

// UpdateReplicationToken replaces the current replication token in the store.
// Returns true if it was changed.
func (t *Store) UpdateReplicationToken(token string, source TokenSource) bool {
	t.l.Lock()
	changed := (t.replicationToken != token || t.replicationTokenSource != source)
	t.replicationToken = token
	t.replicationTokenSource = source
	t.l.Unlock()
	return changed
}

// UserToken returns the best token to use for user operations.
func (t *Store) UserToken() string {
	t.l.RLock()
	defer t.l.RUnlock()

	return t.userToken
}

// AgentToken returns the best token to use for internal agent operations.
func (t *Store) AgentToken() string {
	t.l.RLock()
	defer t.l.RUnlock()

	if tok := t.enterpriseAgentToken(); tok != "" {
		return tok
	}

	if t.agentToken != "" {
		return t.agentToken
	}
	return t.userToken
}

func (t *Store) AgentMasterToken() string {
	t.l.RLock()
	defer t.l.RUnlock()

	return t.agentMasterToken
}

// ReplicationToken returns the replication token.
func (t *Store) ReplicationToken() string {
	t.l.RLock()
	defer t.l.RUnlock()

	return t.replicationToken
}

// UserToken returns the best token to use for user operations.
func (t *Store) UserTokenAndSource() (string, TokenSource) {
	t.l.RLock()
	defer t.l.RUnlock()

	return t.userToken, t.userTokenSource
}

// AgentToken returns the best token to use for internal agent operations.
func (t *Store) AgentTokenAndSource() (string, TokenSource) {
	t.l.RLock()
	defer t.l.RUnlock()

	return t.agentToken, t.agentTokenSource
}

func (t *Store) AgentMasterTokenAndSource() (string, TokenSource) {
	t.l.RLock()
	defer t.l.RUnlock()

	return t.agentMasterToken, t.agentMasterTokenSource
}

// ReplicationToken returns the replication token.
func (t *Store) ReplicationTokenAndSource() (string, TokenSource) {
	t.l.RLock()
	defer t.l.RUnlock()

	return t.replicationToken, t.replicationTokenSource
}

// IsAgentMasterToken checks to see if a given token is the agent master token.
// This will never match an empty token for safety.
func (t *Store) IsAgentMasterToken(token string) bool {
	t.l.RLock()
	defer t.l.RUnlock()

	return (token != "") && (subtle.ConstantTimeCompare([]byte(token), []byte(t.agentMasterToken)) == 1)
}
