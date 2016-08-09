package consul

import (
	"net"
	"strings"

	"github.com/hashicorp/consul/consul/agent"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
)

const (
	// StatusReap is used to update the status of a node if we
	// are handling a EventMemberReap
	StatusReap = serf.MemberStatus(-1)

	// userEventPrefix is pre-pended to a user event to distinguish it
	userEventPrefix = "consul:event:"
)

// userEventName computes the name of a user event
func userEventName(name string) string {
	return userEventPrefix + name
}

// isUserEvent checks if a serf event is a user event
func isUserEvent(name string) bool {
	return strings.HasPrefix(name, userEventPrefix)
}

// rawUserEventName is used to get the raw user event name
func rawUserEventName(name string) string {
	return strings.TrimPrefix(name, userEventPrefix)
}

// lanEventHandler is used to handle events from the lan Serf cluster
func (s *Server) lanEventHandler() {
	for {
		select {
		case e := <-s.eventChLAN:
			switch e.EventType() {
			case serf.EventMemberJoin:
				s.lanNodeJoin(e.(serf.MemberEvent))
				s.localMemberEvent(e.(serf.MemberEvent))

			case serf.EventMemberLeave, serf.EventMemberFailed:
				s.lanNodeFailed(e.(serf.MemberEvent))
				s.localMemberEvent(e.(serf.MemberEvent))

			case serf.EventMemberReap:
				s.localMemberEvent(e.(serf.MemberEvent))
			case serf.EventUser:
				s.localEvent(e.(serf.UserEvent))
			case serf.EventMemberUpdate: // Ignore
			case serf.EventQuery: // Ignore
			default:
				s.logger.Printf("[WARN] consul: Unhandled LAN Serf Event: %#v", e)
			}

		case <-s.shutdownCh:
			return
		}
	}
}

// wanEventHandler is used to handle events from the wan Serf cluster
func (s *Server) wanEventHandler() {
	for {
		select {
		case e := <-s.eventChWAN:
			switch e.EventType() {
			case serf.EventMemberJoin:
				s.wanNodeJoin(e.(serf.MemberEvent))
			case serf.EventMemberLeave, serf.EventMemberFailed:
				s.wanNodeFailed(e.(serf.MemberEvent))
			case serf.EventMemberUpdate: // Ignore
			case serf.EventMemberReap: // Ignore
			case serf.EventUser:
			case serf.EventQuery: // Ignore
			default:
				s.logger.Printf("[WARN] consul: Unhandled WAN Serf Event: %#v", e)
			}

		case <-s.shutdownCh:
			return
		}
	}
}

// localMemberEvent is used to reconcile Serf events with the strongly
// consistent store if we are the current leader
func (s *Server) localMemberEvent(me serf.MemberEvent) {
	// Do nothing if we are not the leader
	if !s.IsLeader() {
		return
	}

	// Check if this is a reap event
	isReap := me.EventType() == serf.EventMemberReap

	// Queue the members for reconciliation
	for _, m := range me.Members {
		// Change the status if this is a reap event
		if isReap {
			m.Status = StatusReap
		}
		select {
		case s.reconcileCh <- m:
		default:
		}
	}
}

// localEvent is called when we receive an event on the local Serf
func (s *Server) localEvent(event serf.UserEvent) {
	// Handle only consul events
	if !strings.HasPrefix(event.Name, "consul:") {
		return
	}

	switch name := event.Name; {
	case name == newLeaderEvent:
		s.logger.Printf("[INFO] consul: New leader elected: %s", event.Payload)

		// Trigger the callback
		if s.config.ServerUp != nil {
			s.config.ServerUp()
		}
	case isUserEvent(name):
		event.Name = rawUserEventName(name)
		s.logger.Printf("[DEBUG] consul: User event: %s", event.Name)

		// Trigger the callback
		if s.config.UserEventHandler != nil {
			s.config.UserEventHandler(event)
		}
	default:
		s.logger.Printf("[WARN] consul: Unhandled local event: %v", event)
	}
}

// lanNodeJoin is used to handle join events on the LAN pool.
func (s *Server) lanNodeJoin(me serf.MemberEvent) {
	for _, m := range me.Members {
		ok, parts := agent.IsConsulServer(m)
		if !ok {
			continue
		}
		s.logger.Printf("[INFO] consul: Adding LAN server %s", parts)

		// See if it's configured as part of our DC.
		if parts.Datacenter == s.config.Datacenter {
			s.localLock.Lock()
			s.localConsuls[raft.ServerAddress(parts.Addr.String())] = parts
			s.localLock.Unlock()
		}

		// If we still expecting to bootstrap, may need to handle this.
		if s.config.BootstrapExpect != 0 {
			s.maybeBootstrap()
		}
	}
}

// wanNodeJoin is used to handle join events on the WAN pool.
func (s *Server) wanNodeJoin(me serf.MemberEvent) {
	for _, m := range me.Members {
		ok, parts := agent.IsConsulServer(m)
		if !ok {
			s.logger.Printf("[WARN] consul: Non-server in WAN pool: %s", m.Name)
			continue
		}
		s.logger.Printf("[INFO] consul: Adding WAN server %s", parts)

		// Search for this node in our existing remotes.
		found := false
		s.remoteLock.Lock()
		existing := s.remoteConsuls[parts.Datacenter]
		for idx, e := range existing {
			if e.Name == parts.Name {
				existing[idx] = parts
				found = true
				break
			}
		}

		// Add to the list if not known.
		if !found {
			s.remoteConsuls[parts.Datacenter] = append(existing, parts)
		}
		s.remoteLock.Unlock()
	}
}

// maybeBootsrap is used to handle bootstrapping when a new consul server joins
func (s *Server) maybeBootstrap() {
	// Bootstrap can only be done if there are no committed logs, remove our
	// expectations of bootstrapping. This is slightly cheaper than the full
	// check that BootstrapCluster will do, so this is a good pre-filter.
	index, err := s.raftStore.LastIndex()
	if err != nil {
		s.logger.Printf("[ERR] consul: Failed to read last raft index: %v", err)
		return
	}
	if index != 0 {
		s.config.BootstrapExpect = 0
		return
	}

	// Scan for all the known servers.
	members := s.serfLAN.Members()
	addrs := make([]string, 0)
	for _, member := range members {
		valid, p := agent.IsConsulServer(member)
		if !valid {
			continue
		}
		if p.Datacenter != s.config.Datacenter {
			s.logger.Printf("[ERR] consul: Member %v has a conflicting datacenter, ignoring", member)
			continue
		}
		if p.Expect != 0 && p.Expect != s.config.BootstrapExpect {
			s.logger.Printf("[ERR] consul: Member %v has a conflicting expect value. All nodes should expect the same number.", member)
			return
		}
		if p.Bootstrap {
			s.logger.Printf("[ERR] consul: Member %v has bootstrap mode. Expect disabled.", member)
			return
		}
		addr := &net.TCPAddr{IP: member.Addr, Port: p.Port}
		addrs = append(addrs, addr.String())
	}

	// Skip if we haven't met the minimum expect count.
	if len(addrs) < s.config.BootstrapExpect {
		return
	}

	// Attempt a live bootstrap!
	var configuration raft.Configuration
	for _, addr := range addrs {
		// TODO (slackpad) - This will need to be updated once we support
		// node IDs.
		server := raft.Server{
			ID:      raft.ServerID(addr),
			Address: raft.ServerAddress(addr),
		}
		configuration.Servers = append(configuration.Servers, server)
	}
	s.logger.Printf("[INFO] consul: Found expected number of peers (%s), attempting to bootstrap cluster...",
		strings.Join(addrs, ","))
	future := s.raft.BootstrapCluster(configuration)
	if err := future.Error(); err != nil {
		s.logger.Printf("[ERR] consul: Failed to bootstrap cluster: %v", err)
	}

	// Bootstrapping complete, don't enter this again.
	s.config.BootstrapExpect = 0
}

// lanNodeFailed is used to handle fail events on the LAN pool.
func (s *Server) lanNodeFailed(me serf.MemberEvent) {
	for _, m := range me.Members {
		ok, parts := agent.IsConsulServer(m)
		if !ok {
			continue
		}
		s.logger.Printf("[INFO] consul: Removing LAN server %s", parts)

		s.localLock.Lock()
		delete(s.localConsuls, raft.ServerAddress(parts.Addr.String()))
		s.localLock.Unlock()
	}
}

// wanNodeFailed is used to handle fail events on the WAN pool.
func (s *Server) wanNodeFailed(me serf.MemberEvent) {
	for _, m := range me.Members {
		ok, parts := agent.IsConsulServer(m)
		if !ok {
			continue
		}
		s.logger.Printf("[INFO] consul: Removing WAN server %s", parts)

		// Remove the server if known
		s.remoteLock.Lock()
		existing := s.remoteConsuls[parts.Datacenter]
		n := len(existing)
		for i := 0; i < n; i++ {
			if existing[i].Name == parts.Name {
				existing[i], existing[n-1] = existing[n-1], nil
				existing = existing[:n-1]
				n--
				break
			}
		}

		// Trim the list if all known consuls are dead
		if n == 0 {
			delete(s.remoteConsuls, parts.Datacenter)
		} else {
			s.remoteConsuls[parts.Datacenter] = existing
		}
		s.remoteLock.Unlock()
	}
}
