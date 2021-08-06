package consul

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	msgpackrpc "github.com/hashicorp/net-rpc-msgpackrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/consul/acl"
	"github.com/hashicorp/consul/agent/consul/state"
	"github.com/hashicorp/consul/agent/pool"
	"github.com/hashicorp/consul/agent/structs"
	tokenStore "github.com/hashicorp/consul/agent/token"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/sdk/testutil/retry"
	"github.com/hashicorp/consul/testrpc"
)

func TestRPC_NoLeader_Fail(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()
	dir1, s1 := testServerWithConfig(t, func(c *Config) {
		c.RPCHoldTimeout = 1 * time.Millisecond
	})
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	codec := rpcClient(t, s1)
	defer codec.Close()

	arg := structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
	}
	var out struct{}

	// Make sure we eventually fail with a no leader error, which we should
	// see given the short timeout.
	err := msgpackrpc.CallWithCodec(codec, "Catalog.Register", &arg, &out)
	if err == nil || err.Error() != structs.ErrNoLeader.Error() {
		t.Fatalf("bad: %v", err)
	}

	// Now make sure it goes through.
	testrpc.WaitForTestAgent(t, s1.RPC, "dc1")
	err = msgpackrpc.CallWithCodec(codec, "Catalog.Register", &arg, &out)
	if err != nil {
		t.Fatalf("bad: %v", err)
	}
}

func TestRPC_NoLeader_Fail_on_stale_read(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()
	dir1, s1 := testServerWithConfig(t, func(c *Config) {
		c.RPCHoldTimeout = 1 * time.Millisecond
	})
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	codec := rpcClient(t, s1)
	defer codec.Close()

	arg := structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
	}
	var out struct{}

	// Make sure we eventually fail with a no leader error, which we should
	// see given the short timeout.
	err := msgpackrpc.CallWithCodec(codec, "Catalog.Register", &arg, &out)
	if err == nil || err.Error() != structs.ErrNoLeader.Error() {
		t.Fatalf("bad: %v", err)
	}

	// Until leader has never been known, stale should fail
	getKeysReq := structs.KeyListRequest{
		Datacenter:   "dc1",
		Prefix:       "",
		Seperator:    "/",
		QueryOptions: structs.QueryOptions{AllowStale: true},
	}
	var keyList structs.IndexedKeyList
	if err := msgpackrpc.CallWithCodec(codec, "KVS.ListKeys", &getKeysReq, &keyList); err.Error() != structs.ErrNoLeader.Error() {
		t.Fatalf("expected %v but got err: %v", structs.ErrNoLeader, err)
	}

	testrpc.WaitForTestAgent(t, s1.RPC, "dc1")
	if err := msgpackrpc.CallWithCodec(codec, "KVS.ListKeys", &getKeysReq, &keyList); err != nil {
		t.Fatalf("Did not expect any error but got err: %v", err)
	}
}

func TestRPC_NoLeader_Retry(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()
	dir1, s1 := testServerWithConfig(t, func(c *Config) {
		c.RPCHoldTimeout = 10 * time.Second
	})
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	codec := rpcClient(t, s1)
	defer codec.Close()

	arg := structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
	}
	var out struct{}

	// This isn't sure-fire but tries to check that we don't have a
	// leader going into the RPC, so we exercise the retry logic.
	if ok, _, _ := s1.getLeader(); ok {
		t.Fatalf("should not have a leader yet")
	}

	// The timeout is long enough to ride out any reasonable leader
	// election.
	err := msgpackrpc.CallWithCodec(codec, "Catalog.Register", &arg, &out)
	if err != nil {
		t.Fatalf("bad: %v", err)
	}
}

func TestRPC_getLeader_ErrLeaderNotTracked(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	cluster := newTestCluster(t, &testClusterConfig{
		Datacenter: "dc1",
		Servers:    3,
		ServerWait: func(t *testing.T, srv *Server) {
			// The test cluster waits for a leader to be established
			// but not for all the RPC tracking of all servers to be updated
			// so we also want to wait for that here
			retry.Run(t, func(r *retry.R) {
				if !srv.IsLeader() {
					_, _, err := srv.getLeader()
					require.NoError(r, err)
				}
			})

		},
	})

	// At this point we know we have a cluster with a leader and all followers are tracking that
	// leader in the serverLookup struct. We need to find a follower to hack its server lookup
	// to force the error we desire

	var follower *Server
	for _, srv := range cluster.Servers {
		if !srv.IsLeader() {
			follower = srv
			break
		}
	}

	_, leaderMeta, err := follower.getLeader()
	require.NoError(t, err)

	// now do some behind the scenes trickery on the followers server lookup
	// to remove the leader from it so that we can force a ErrLeaderNotTracked error
	follower.serverLookup.RemoveServer(leaderMeta)

	isLeader, meta, err := follower.getLeader()
	require.Error(t, err)
	require.True(t, errors.Is(err, structs.ErrLeaderNotTracked))
	require.Nil(t, meta)
	require.False(t, isLeader)
}

type MockSink struct {
	*bytes.Buffer
	cancel bool
}

func (m *MockSink) ID() string {
	return "Mock"
}

func (m *MockSink) Cancel() error {
	m.cancel = true
	return nil
}

func (m *MockSink) Close() error {
	return nil
}

func TestRPC_blockingQuery(t *testing.T) {
	t.Parallel()
	dir, s := testServer(t)
	defer os.RemoveAll(dir)
	defer s.Shutdown()

	require := require.New(t)
	assert := assert.New(t)

	// Perform a non-blocking query. Note that it's significant that the meta has
	// a zero index in response - the implied opts.MinQueryIndex is also zero but
	// this should not block still.
	{
		var opts structs.QueryOptions
		var meta structs.QueryMeta
		var calls int
		fn := func(_ memdb.WatchSet, _ *state.Store) error {
			calls++
			return nil
		}
		if err := s.blockingQuery(&opts, &meta, fn); err != nil {
			t.Fatalf("err: %v", err)
		}
		if calls != 1 {
			t.Fatalf("bad: %d", calls)
		}
	}

	// Perform a blocking query that gets woken up and loops around once.
	{
		opts := structs.QueryOptions{
			MinQueryIndex: 3,
		}
		var meta structs.QueryMeta
		var calls int
		fn := func(ws memdb.WatchSet, _ *state.Store) error {
			if calls == 0 {
				meta.Index = 3

				fakeCh := make(chan struct{})
				close(fakeCh)
				ws.Add(fakeCh)
			} else {
				meta.Index = 4
			}
			calls++
			return nil
		}
		if err := s.blockingQuery(&opts, &meta, fn); err != nil {
			t.Fatalf("err: %v", err)
		}
		if calls != 2 {
			t.Fatalf("bad: %d", calls)
		}
	}

	// Perform a blocking query that returns a zero index from blocking func (e.g.
	// no state yet). This should still return an empty response immediately, but
	// with index of 1 and then block on the next attempt. In one sense zero index
	// is not really a valid response from a state method that is not an error but
	// in practice a lot of state store operations do return it unless they
	// explicitly special checks to turn 0 into 1. Often this is not caught or
	// covered by tests but eventually when hit in the wild causes blocking
	// clients to busy loop and burn CPU. This test ensure that blockingQuery
	// systematically does the right thing to prevent future bugs like that.
	{
		opts := structs.QueryOptions{
			MinQueryIndex: 0,
		}
		var meta structs.QueryMeta
		var calls int
		fn := func(ws memdb.WatchSet, _ *state.Store) error {
			if opts.MinQueryIndex > 0 {
				// If client requested blocking, block forever. This is simulating
				// waiting for the watched resource to be initialized/written to giving
				// it a non-zero index. Note the timeout on the query options is relied
				// on to stop the test taking forever.
				fakeCh := make(chan struct{})
				ws.Add(fakeCh)
			}
			meta.Index = 0
			calls++
			return nil
		}
		require.NoError(s.blockingQuery(&opts, &meta, fn))
		assert.Equal(1, calls)
		assert.Equal(uint64(1), meta.Index,
			"expect fake index of 1 to force client to block on next update")

		// Simulate client making next request
		opts.MinQueryIndex = 1
		opts.MaxQueryTime = 20 * time.Millisecond // Don't wait too long

		// This time we should block even though the func returns index 0 still
		t0 := time.Now()
		require.NoError(s.blockingQuery(&opts, &meta, fn))
		t1 := time.Now()
		assert.Equal(2, calls)
		assert.Equal(uint64(1), meta.Index,
			"expect fake index of 1 to force client to block on next update")
		assert.True(t1.Sub(t0) > 20*time.Millisecond,
			"should have actually blocked waiting for timeout")

	}

	// Perform a query that blocks and gets interrupted when the state store
	// is abandoned.
	{
		opts := structs.QueryOptions{
			MinQueryIndex: 3,
		}
		var meta structs.QueryMeta
		var calls int
		fn := func(_ memdb.WatchSet, _ *state.Store) error {
			if calls == 0 {
				meta.Index = 3

				snap, err := s.fsm.Snapshot()
				if err != nil {
					t.Fatalf("err: %v", err)
				}
				defer snap.Release()

				buf := bytes.NewBuffer(nil)
				sink := &MockSink{buf, false}
				if err := snap.Persist(sink); err != nil {
					t.Fatalf("err: %v", err)
				}

				if err := s.fsm.Restore(sink); err != nil {
					t.Fatalf("err: %v", err)
				}
			}
			calls++
			return nil
		}
		if err := s.blockingQuery(&opts, &meta, fn); err != nil {
			t.Fatalf("err: %v", err)
		}
		if calls != 1 {
			t.Fatalf("bad: %d", calls)
		}
	}
}

func TestRPC_ReadyForConsistentReads(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()
	dir, s := testServerWithConfig(t, func(c *Config) {
		c.RPCHoldTimeout = 2 * time.Millisecond
	})
	defer os.RemoveAll(dir)
	defer s.Shutdown()

	testrpc.WaitForLeader(t, s.RPC, "dc1")

	if !s.isReadyForConsistentReads() {
		t.Fatal("Server should be ready for consistent reads")
	}

	s.resetConsistentReadReady()
	err := s.consistentRead()
	if err.Error() != "Not ready to serve consistent reads" {
		t.Fatal("Server should NOT be ready for consistent reads")
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		s.setConsistentReadReady()
	}()

	retry.Run(t, func(r *retry.R) {
		if err := s.consistentRead(); err != nil {
			r.Fatalf("Expected server to be ready for consistent reads, got error %v", err)
		}
	})
}

func TestRPC_MagicByteTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()
	dir1, s1 := testServerWithConfig(t, func(c *Config) {
		c.RPCHandshakeTimeout = 10 * time.Millisecond
	})
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	// Connect to the server with bare TCP to simulate a malicious client trying
	// to hold open resources.
	addr := s1.config.RPCAdvertise
	conn, err := net.DialTimeout("tcp", addr.String(), time.Second)
	require.NoError(t, err)
	defer conn.Close()

	// Wait for more than the timeout. This is timing dependent so could fail if
	// the CPU is super overloaded so the handler goroutine so I'm using a retry
	// loop below to be sure but this feels like a pretty generous margin for
	// error (10x the timeout and 100ms of scheduling time).
	time.Sleep(100 * time.Millisecond)

	// Set a read deadline on the Conn in case the timeout is not working we don't
	// want the read below to block forever. Needs to be much longer than what we
	// expect and the error should be different too.
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	retry.Run(t, func(r *retry.R) {
		// Sanity check the conn was closed by attempting to read from it (a write
		// might not detect the close).
		buf := make([]byte, 10)
		_, err = conn.Read(buf)
		require.Error(r, err)
		require.Contains(r, err.Error(), "EOF")
	})
}

func TestRPC_TLSHandshakeTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()

	dir1, s1 := testServerWithConfig(t, func(c *Config) {
		c.RPCHandshakeTimeout = 10 * time.Millisecond
		c.TLSConfig.CAFile = "../../test/hostname/CertAuth.crt"
		c.TLSConfig.CertFile = "../../test/hostname/Alice.crt"
		c.TLSConfig.KeyFile = "../../test/hostname/Alice.key"
		c.TLSConfig.VerifyServerHostname = true
		c.TLSConfig.VerifyOutgoing = true
		c.TLSConfig.VerifyIncoming = true
	})
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	// Connect to the server with TLS magic byte delivered on time
	addr := s1.config.RPCAdvertise
	conn, err := net.DialTimeout("tcp", addr.String(), time.Second)
	require.NoError(t, err)
	defer conn.Close()

	// Write TLS byte to avoid being closed by either the (outer) first byte
	// timeout or the fact that server requires TLS
	_, err = conn.Write([]byte{pool.RPCTLS})
	require.NoError(t, err)

	// Wait for more than the timeout before we start a TLS handshake. This is
	// timing dependent so could fail if the CPU is super overloaded so the
	// handler goroutine so I'm using a retry loop below to be sure but this feels
	// like a pretty generous margin for error (10x the timeout and 100ms of
	// scheduling time).
	time.Sleep(100 * time.Millisecond)

	// Set a read deadline on the Conn in case the timeout is not working we don't
	// want the read below to block forever. Needs to be much longer than what we
	// expect and the error should be different too.
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	retry.Run(t, func(r *retry.R) {
		// Sanity check the conn was closed by attempting to read from it (a write
		// might not detect the close).
		buf := make([]byte, 10)
		_, err = conn.Read(buf)
		require.Error(r, err)
		require.Contains(r, err.Error(), "EOF")
	})
}

func TestRPC_PreventsTLSNesting(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()

	cases := []struct {
		name      string
		outerByte pool.RPCType
		innerByte pool.RPCType
		wantClose bool
	}{
		{
			// Base case, sanity check normal RPC in TLS works
			name:      "RPC in TLS",
			outerByte: pool.RPCTLS,
			innerByte: pool.RPCConsul,
			wantClose: false,
		},
		{
			// Nested TLS-in-TLS
			name:      "TLS in TLS",
			outerByte: pool.RPCTLS,
			innerByte: pool.RPCTLS,
			wantClose: true,
		},
		{
			// Nested TLS-in-TLS
			name:      "TLS in Insecure TLS",
			outerByte: pool.RPCTLSInsecure,
			innerByte: pool.RPCTLS,
			wantClose: true,
		},
		{
			// Nested TLS-in-TLS
			name:      "Insecure TLS in TLS",
			outerByte: pool.RPCTLS,
			innerByte: pool.RPCTLSInsecure,
			wantClose: true,
		},
		{
			// Nested TLS-in-TLS
			name:      "Insecure TLS in Insecure TLS",
			outerByte: pool.RPCTLSInsecure,
			innerByte: pool.RPCTLSInsecure,
			wantClose: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir1, s1 := testServerWithConfig(t, func(c *Config) {
				c.TLSConfig.CAFile = "../../test/hostname/CertAuth.crt"
				c.TLSConfig.CertFile = "../../test/hostname/Alice.crt"
				c.TLSConfig.KeyFile = "../../test/hostname/Alice.key"
				c.TLSConfig.VerifyServerHostname = true
				c.TLSConfig.VerifyOutgoing = true
				c.TLSConfig.VerifyIncoming = false // saves us getting client cert setup
				c.TLSConfig.Domain = "consul"
			})
			defer os.RemoveAll(dir1)
			defer s1.Shutdown()

			// Connect to the server with TLS magic byte delivered on time
			addr := s1.config.RPCAdvertise
			conn, err := net.DialTimeout("tcp", addr.String(), time.Second)
			require.NoError(t, err)
			defer conn.Close()

			// Write Outer magic byte
			_, err = conn.Write([]byte{byte(tc.outerByte)})
			require.NoError(t, err)

			// Start tls client
			tlsWrap := s1.tlsConfigurator.OutgoingRPCWrapper()
			tlsConn, err := tlsWrap("dc1", conn)
			require.NoError(t, err)

			// Write Inner magic byte
			_, err = tlsConn.Write([]byte{byte(tc.innerByte)})
			require.NoError(t, err)

			if tc.wantClose {
				// Allow up to a second for a read failure to indicate conn was closed by
				// server.
				conn.SetReadDeadline(time.Now().Add(1 * time.Second))

				retry.Run(t, func(r *retry.R) {
					// Sanity check the conn was closed by attempting to read from it (a
					// write might not detect the close).
					buf := make([]byte, 10)
					_, err = tlsConn.Read(buf)
					require.Error(r, err)
					require.Contains(r, err.Error(), "EOF")
				})
			} else {
				// Set a shorter read deadline that should typically be enough to detect
				// immediate close but will also not make test hang forever. This
				// positive case is mostly just a sanity check that the test code here
				// is actually not failing just due to some other error in the way we
				// setup TLS. It also sanity checks that we still allow valid TLS conns
				// but if it produces possible false-positives in CI sometimes that's
				// not such a huge deal - CI won't be brittle and it will have done it's
				// job as a sanity check most of the time.
				conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
				buf := make([]byte, 10)
				_, err = tlsConn.Read(buf)
				require.Error(t, err)
				require.Contains(t, err.Error(), "i/o timeout")
			}
		})
	}
}

func connectClient(t *testing.T, s1 *Server, mb pool.RPCType, useTLS, wantOpen bool, message string) net.Conn {
	t.Helper()

	addr := s1.config.RPCAdvertise
	tlsWrap := s1.tlsConfigurator.OutgoingRPCWrapper()

	conn, err := net.DialTimeout("tcp", addr.String(), time.Second)
	require.NoError(t, err)

	// Write magic byte so we aren't timed out
	outerByte := mb
	if useTLS {
		outerByte = pool.RPCTLS
	}
	_, err = conn.Write([]byte{byte(outerByte)})
	require.NoError(t, err)

	if useTLS {
		tlsConn, err := tlsWrap(s1.config.Datacenter, conn)
		// Subtly, tlsWrap will NOT actually do a handshake in this case - it only
		// does so for some configs, so even if the server closed the conn before
		// handshake this won't fail and it's only when we attempt to read or write
		// that we'll see the broken pipe.
		require.NoError(t, err, "%s: wanted open conn, failed TLS handshake: %s",
			message, err)
		conn = tlsConn

		// Write Inner magic byte
		_, err = conn.Write([]byte{byte(mb)})
		if !wantOpen {
			// TLS Handshake will be done on this attempt to write and should fail
			require.Error(t, err, "%s: wanted closed conn, TLS Handshake succeeded", message)
		} else {
			require.NoError(t, err, "%s: wanted open conn, failed writing inner magic byte: %s",
				message, err)
		}
	}

	// Check if the conn is in the state we want.
	retry.Run(t, func(r *retry.R) {
		// Don't wait around as server won't be sending data but the read will fail
		// immediately if the conn is closed.
		conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
		buf := make([]byte, 10)
		_, err := conn.Read(buf)
		require.Error(r, err)
		if wantOpen {
			require.Contains(r, err.Error(), "i/o timeout",
				"%s: wanted an open conn (read timeout)", message)
		} else {
			if useTLS {
				require.Error(r, err)
				// TLS may fail during either read or write of the handshake so there
				// are a few different errors that come up.
				if !strings.Contains(err.Error(), "read: connection reset by peer") &&
					!strings.Contains(err.Error(), "write: connection reset by peer") &&
					!strings.Contains(err.Error(), "write: broken pipe") {
					r.Fatalf("%s: wanted closed conn got err: %s", message, err)
				}
			} else {
				require.Contains(r, err.Error(), "EOF", "%s: wanted a closed conn",
					message)
			}
		}
	})

	return conn
}

func TestRPC_RPCMaxConnsPerClient(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()

	cases := []struct {
		name       string
		magicByte  pool.RPCType
		tlsEnabled bool
	}{
		{"RPC", pool.RPCMultiplexV2, false},
		{"RPC TLS", pool.RPCMultiplexV2, true},
		{"Raft", pool.RPCRaft, false},
		{"Raft TLS", pool.RPCRaft, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dir1, s1 := testServerWithConfig(t, func(c *Config) {
				c.RPCMaxConnsPerClient = 2
				if tc.tlsEnabled {
					c.TLSConfig.CAFile = "../../test/hostname/CertAuth.crt"
					c.TLSConfig.CertFile = "../../test/hostname/Alice.crt"
					c.TLSConfig.KeyFile = "../../test/hostname/Alice.key"
					c.TLSConfig.VerifyServerHostname = true
					c.TLSConfig.VerifyOutgoing = true
					c.TLSConfig.VerifyIncoming = false // saves us getting client cert setup
					c.TLSConfig.Domain = "consul"
				}
			})
			defer os.RemoveAll(dir1)
			defer s1.Shutdown()

			// Connect to the server with bare TCP
			conn1 := connectClient(t, s1, tc.magicByte, tc.tlsEnabled, true, "conn1")
			defer conn1.Close()

			// Two conns should succeed
			conn2 := connectClient(t, s1, tc.magicByte, tc.tlsEnabled, true, "conn2")
			defer conn2.Close()

			// Third should be closed byt the limiter
			conn3 := connectClient(t, s1, tc.magicByte, tc.tlsEnabled, false, "conn3")
			defer conn3.Close()

			// If we close one of the earlier ones, we should be able to open another
			addr := conn1.RemoteAddr()
			conn1.Close()
			retry.Run(t, func(r *retry.R) {
				if n := s1.rpcConnLimiter.NumOpen(addr); n >= 2 {
					r.Fatal("waiting for open conns to drop")
				}
			})
			conn4 := connectClient(t, s1, tc.magicByte, tc.tlsEnabled, true, "conn4")
			defer conn4.Close()

			// Reload config with higher limit
			rc := ReloadableConfig{
				RPCRateLimit:         s1.config.RPCRateLimit,
				RPCMaxBurst:          s1.config.RPCMaxBurst,
				RPCMaxConnsPerClient: 10,
			}
			require.NoError(t, s1.ReloadConfig(rc))

			// Now another conn should be allowed
			conn5 := connectClient(t, s1, tc.magicByte, tc.tlsEnabled, true, "conn5")
			defer conn5.Close()
		})
	}
}

func TestRPC_readUint32(t *testing.T) {
	cases := []struct {
		name    string
		writeFn func(net.Conn)
		readFn  func(*testing.T, net.Conn)
	}{
		{
			name: "timeouts irrelevant",
			writeFn: func(conn net.Conn) {
				_ = binary.Write(conn, binary.BigEndian, uint32(42))
				_ = binary.Write(conn, binary.BigEndian, uint32(math.MaxUint32))
				_ = binary.Write(conn, binary.BigEndian, uint32(1))
			},
			readFn: func(t *testing.T, conn net.Conn) {
				t.Helper()
				v, err := readUint32(conn, 5*time.Second)
				require.NoError(t, err)
				require.Equal(t, uint32(42), v)

				v, err = readUint32(conn, 5*time.Second)
				require.NoError(t, err)
				require.Equal(t, uint32(math.MaxUint32), v)

				v, err = readUint32(conn, 5*time.Second)
				require.NoError(t, err)
				require.Equal(t, uint32(1), v)
			},
		},
		{
			name: "triggers timeout on last read",
			writeFn: func(conn net.Conn) {
				_ = binary.Write(conn, binary.BigEndian, uint32(42))
				_ = binary.Write(conn, binary.BigEndian, uint32(math.MaxUint32))
				_ = binary.Write(conn, binary.BigEndian, uint16(1)) // half as many bytes as expected
			},
			readFn: func(t *testing.T, conn net.Conn) {
				t.Helper()
				v, err := readUint32(conn, 5*time.Second)
				require.NoError(t, err)
				require.Equal(t, uint32(42), v)

				v, err = readUint32(conn, 5*time.Second)
				require.NoError(t, err)
				require.Equal(t, uint32(math.MaxUint32), v)

				_, err = readUint32(conn, 50*time.Millisecond)
				require.Error(t, err)
				nerr, ok := err.(net.Error)
				require.True(t, ok)
				require.True(t, nerr.Timeout())
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var doneWg sync.WaitGroup
			defer doneWg.Wait()

			client, server := net.Pipe()
			defer client.Close()
			defer server.Close()

			// Client pushes some data.
			doneWg.Add(1)
			go func() {
				doneWg.Done()
				tc.writeFn(client)
			}()

			// The server tests the function for us.
			tc.readFn(t, server)
		})
	}
}

func TestRPC_LocalTokenStrippedOnForward(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()
	dir1, s1 := testServerWithConfig(t, func(c *Config) {
		c.PrimaryDatacenter = "dc1"
		c.ACLsEnabled = true
		c.ACLResolverSettings.ACLDefaultPolicy = "deny"
		c.ACLMasterToken = "root"
	})
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	testrpc.WaitForLeader(t, s1.RPC, "dc1")
	codec := rpcClient(t, s1)
	defer codec.Close()

	dir2, s2 := testServerWithConfig(t, func(c *Config) {
		c.Datacenter = "dc2"
		c.PrimaryDatacenter = "dc1"
		c.ACLsEnabled = true
		c.ACLResolverSettings.ACLDefaultPolicy = "deny"
		c.ACLTokenReplication = true
		c.ACLReplicationRate = 100
		c.ACLReplicationBurst = 100
		c.ACLReplicationApplyLimit = 1000000
	})
	s2.tokens.UpdateReplicationToken("root", tokenStore.TokenSourceConfig)
	testrpc.WaitForLeader(t, s2.RPC, "dc2")
	defer os.RemoveAll(dir2)
	defer s2.Shutdown()
	codec2 := rpcClient(t, s2)
	defer codec2.Close()

	// Try to join.
	joinWAN(t, s2, s1)
	testrpc.WaitForLeader(t, s1.RPC, "dc1")
	testrpc.WaitForLeader(t, s1.RPC, "dc2")

	// Wait for legacy acls to be disabled so we are clear that
	// legacy replication isn't meddling.
	waitForNewACLs(t, s1)
	waitForNewACLs(t, s2)
	waitForNewACLReplication(t, s2, structs.ACLReplicateTokens, 1, 1, 0)

	// create simple kv policy
	kvPolicy, err := upsertTestPolicyWithRules(codec, "root", "dc1", `
	key_prefix "" { policy = "write" }
	`)
	require.NoError(t, err)

	// Wait for it to replicate
	retry.Run(t, func(r *retry.R) {
		_, p, err := s2.fsm.State().ACLPolicyGetByID(nil, kvPolicy.ID, &structs.EnterpriseMeta{})
		require.Nil(r, err)
		require.NotNil(r, p)
	})

	// create local token that only works in DC2
	localToken2, err := upsertTestToken(codec, "root", "dc2", func(token *structs.ACLToken) {
		token.Local = true
		token.Policies = []structs.ACLTokenPolicyLink{
			{ID: kvPolicy.ID},
		}
	})
	require.NoError(t, err)

	// Try to use it locally (it should work)
	arg := structs.KVSRequest{
		Datacenter: "dc2",
		Op:         api.KVSet,
		DirEnt: structs.DirEntry{
			Key:   "foo",
			Value: []byte("bar"),
		},
		WriteRequest: structs.WriteRequest{Token: localToken2.SecretID},
	}
	var out bool
	err = msgpackrpc.CallWithCodec(codec2, "KVS.Apply", &arg, &out)
	require.NoError(t, err)
	require.Equal(t, localToken2.SecretID, arg.WriteRequest.Token, "token should not be stripped")

	// Try to use it remotely
	arg = structs.KVSRequest{
		Datacenter: "dc1",
		Op:         api.KVSet,
		DirEnt: structs.DirEntry{
			Key:   "foo",
			Value: []byte("bar"),
		},
		WriteRequest: structs.WriteRequest{Token: localToken2.SecretID},
	}
	err = msgpackrpc.CallWithCodec(codec2, "KVS.Apply", &arg, &out)
	if !acl.IsErrPermissionDenied(err) {
		t.Fatalf("err: %v", err)
	}

	// Update the anon token to also be able to write to kv
	{
		tokenUpsertReq := structs.ACLTokenSetRequest{
			Datacenter: "dc1",
			ACLToken: structs.ACLToken{
				AccessorID: structs.ACLTokenAnonymousID,
				Policies: []structs.ACLTokenPolicyLink{
					{
						ID: kvPolicy.ID,
					},
				},
			},
			WriteRequest: structs.WriteRequest{Token: "root"},
		}
		token := structs.ACLToken{}
		err = msgpackrpc.CallWithCodec(codec, "ACL.TokenSet", &tokenUpsertReq, &token)
		require.NoError(t, err)
		require.NotEmpty(t, token.SecretID)
	}

	// Try to use it remotely again, but this time it should fallback to anon
	arg = structs.KVSRequest{
		Datacenter: "dc1",
		Op:         api.KVSet,
		DirEnt: structs.DirEntry{
			Key:   "foo",
			Value: []byte("bar"),
		},
		WriteRequest: structs.WriteRequest{Token: localToken2.SecretID},
	}
	err = msgpackrpc.CallWithCodec(codec2, "KVS.Apply", &arg, &out)
	require.NoError(t, err)
	require.Equal(t, localToken2.SecretID, arg.WriteRequest.Token, "token should not be stripped")
}

func TestCanRetry(t *testing.T) {
	type testCase struct {
		name     string
		req      structs.RPCInfo
		err      error
		expected bool
		timeout  time.Time
	}
	config := DefaultConfig()
	now := time.Now()
	config.RPCHoldTimeout = 7 * time.Second
	run := func(t *testing.T, tc testCase) {
		timeOutValue := tc.timeout
		if timeOutValue.IsZero() {
			timeOutValue = now
		}
		require.Equal(t, tc.expected, canRetry(tc.req, tc.err, timeOutValue, config))
	}

	var testCases = []testCase{
		{
			name:     "unexpected error",
			err:      fmt.Errorf("some arbitrary error"),
			expected: false,
		},
		{
			name:     "checking error",
			err:      fmt.Errorf("some wrapping :%w", ErrChunkingResubmit),
			expected: true,
		},
		{
			name:     "no leader error",
			err:      fmt.Errorf("some wrapping: %w", structs.ErrNoLeader),
			expected: true,
		},
		{
			name:     "EOF on read request",
			req:      isReadRequest{},
			err:      io.EOF,
			expected: true,
		},
		{
			name:     "EOF error",
			req:      &structs.DCSpecificRequest{},
			err:      io.EOF,
			expected: true,
		},
		{
			name:     "HasTimedOut implementation with no error",
			req:      &structs.DCSpecificRequest{},
			err:      nil,
			expected: false,
		},
		{
			name:     "HasTimedOut implementation timedOut with no error",
			req:      &structs.DCSpecificRequest{},
			err:      nil,
			expected: false,
			timeout:  now.Add(-(config.RPCHoldTimeout + time.Second)),
		},
		{
			name:     "HasTimedOut implementation timedOut (with EOF error)",
			req:      &structs.DCSpecificRequest{},
			err:      io.EOF,
			expected: false,
			timeout:  now.Add(-(config.RPCHoldTimeout + time.Second)),
		},
		{
			name:     "HasTimedOut implementation timedOut blocking call",
			req:      &structs.DCSpecificRequest{QueryOptions: structs.QueryOptions{MaxQueryTime: 300, MinQueryIndex: 1}},
			err:      nil,
			expected: false,
			timeout:  now.Add(-(config.RPCHoldTimeout + config.MaxQueryTime + time.Second)),
		},
		{
			name:     "HasTimedOut implementation timedOut blocking call (MaxQueryTime not set)",
			req:      &structs.DCSpecificRequest{QueryOptions: structs.QueryOptions{MinQueryIndex: 1}},
			err:      nil,
			expected: false,
			timeout:  now.Add(-(config.RPCHoldTimeout + config.MaxQueryTime + time.Second)),
		},
		{
			name:     "EOF on write request",
			err:      io.EOF,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

type isReadRequest struct {
	structs.RPCInfo
}

func (r isReadRequest) IsRead() bool {
	return true
}

func (r isReadRequest) HasTimedOut(since time.Time, rpcHoldTimeout, maxQueryTime, defaultQueryTime time.Duration) bool {
	return false
}
