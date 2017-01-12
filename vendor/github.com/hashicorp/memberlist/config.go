package memberlist

import (
	"io"
	"log"
	"os"
	"time"
)

type Config struct {
	// The name of this node. This must be unique in the cluster.
	Name string

	// Configuration related to what address to bind to and ports to
	// listen on. The port is used for both UDP and TCP gossip.
	// It is assumed other nodes are running on this port, but they
	// do not need to.
	BindAddr string
	BindPort int

	// Configuration related to what address to advertise to other
	// cluster members. Used for nat traversal.
	AdvertiseAddr string
	AdvertisePort int

	// ProtocolVersion is the configured protocol version that we
	// will _speak_. This must be between ProtocolVersionMin and
	// ProtocolVersionMax.
	ProtocolVersion uint8

	// TCPTimeout is the timeout for establishing a TCP connection with
	// a remote node for a full state sync.
	TCPTimeout time.Duration

	// IndirectChecks is the number of nodes that will be asked to perform
	// an indirect probe of a node in the case a direct probe fails. Memberlist
	// waits for an ack from any single indirect node, so increasing this
	// number will increase the likelihood that an indirect probe will succeed
	// at the expense of bandwidth.
	IndirectChecks int

	// RetransmitMult is the multiplier for the number of retransmissions
	// that are attempted for messages broadcasted over gossip. The actual
	// count of retransmissions is calculated using the formula:
	//
	//   Retransmits = RetransmitMult * log(N+1)
	//
	// This allows the retransmits to scale properly with cluster size. The
	// higher the multiplier, the more likely a failed broadcast is to converge
	// at the expense of increased bandwidth.
	RetransmitMult int

	// SuspicionMult is the multiplier for determining the time an
	// inaccessible node is considered suspect before declaring it dead.
	// The actual timeout is calculated using the formula:
	//
	//   SuspicionTimeout = SuspicionMult * log(N+1) * ProbeInterval
	//
	// This allows the timeout to scale properly with expected propagation
	// delay with a larger cluster size. The higher the multiplier, the longer
	// an inaccessible node is considered part of the cluster before declaring
	// it dead, giving that suspect node more time to refute if it is indeed
	// still alive.
	SuspicionMult int

	// SuspicionMaxTimeoutMult is the multiplier applied to the
	// SuspicionTimeout used as an upper bound on detection time. This max
	// timeout is calculated using the formula:
	//
	// SuspicionMaxTimeout = SuspicionMaxTimeoutMult * SuspicionTimeout
	//
	// If everything is working properly, confirmations from other nodes will
	// accelerate suspicion timers in a manner which will cause the timeout
	// to reach the base SuspicionTimeout before that elapses, so this value
	// will typically only come into play if a node is experiencing issues
	// communicating with other nodes. It should be set to a something fairly
	// large so that a node having problems will have a lot of chances to
	// recover before falsely declaring other nodes as failed, but short
	// enough for a legitimately isolated node to still make progress marking
	// nodes failed in a reasonable amount of time.
	SuspicionMaxTimeoutMult int

	// PushPullInterval is the interval between complete state syncs.
	// Complete state syncs are done with a single node over TCP and are
	// quite expensive relative to standard gossiped messages. Setting this
	// to zero will disable state push/pull syncs completely.
	//
	// Setting this interval lower (more frequent) will increase convergence
	// speeds across larger clusters at the expense of increased bandwidth
	// usage.
	PushPullInterval time.Duration

	// ProbeInterval and ProbeTimeout are used to configure probing
	// behavior for memberlist.
	//
	// ProbeInterval is the interval between random node probes. Setting
	// this lower (more frequent) will cause the memberlist cluster to detect
	// failed nodes more quickly at the expense of increased bandwidth usage.
	//
	// ProbeTimeout is the timeout to wait for an ack from a probed node
	// before assuming it is unhealthy. This should be set to 99-percentile
	// of RTT (round-trip time) on your network.
	ProbeInterval time.Duration
	ProbeTimeout  time.Duration

	// DisableTcpPings will turn off the fallback TCP pings that are attempted
	// if the direct UDP ping fails. These get pipelined along with the
	// indirect UDP pings.
	DisableTcpPings bool

	// AwarenessMaxMultiplier will increase the probe interval if the node
	// becomes aware that it might be degraded and not meeting the soft real
	// time requirements to reliably probe other nodes.
	AwarenessMaxMultiplier int

	// GossipInterval and GossipNodes are used to configure the gossip
	// behavior of memberlist.
	//
	// GossipInterval is the interval between sending messages that need
	// to be gossiped that haven't been able to piggyback on probing messages.
	// If this is set to zero, non-piggyback gossip is disabled. By lowering
	// this value (more frequent) gossip messages are propagated across
	// the cluster more quickly at the expense of increased bandwidth.
	//
	// GossipNodes is the number of random nodes to send gossip messages to
	// per GossipInterval. Increasing this number causes the gossip messages
	// to propagate across the cluster more quickly at the expense of
	// increased bandwidth.
	//
	// GossipToTheDeadTime is the interval after which a node has died that
	// we will still try to gossip to it. This gives it a chance to refute.
	GossipInterval      time.Duration
	GossipNodes         int
	GossipToTheDeadTime time.Duration

	// EnableCompression is used to control message compression. This can
	// be used to reduce bandwidth usage at the cost of slightly more CPU
	// utilization. This is only available starting at protocol version 1.
	EnableCompression bool

	// SecretKey is used to initialize the primary encryption key in a keyring.
	// The primary encryption key is the only key used to encrypt messages and
	// the first key used while attempting to decrypt messages. Providing a
	// value for this primary key will enable message-level encryption and
	// verification, and automatically install the key onto the keyring.
	// The value should be either 16, 24, or 32 bytes to select AES-128,
	// AES-192, or AES-256.
	SecretKey []byte

	// The keyring holds all of the encryption keys used internally. It is
	// automatically initialized using the SecretKey and SecretKeys values.
	Keyring *Keyring

	// Delegate and Events are delegates for receiving and providing
	// data to memberlist via callback mechanisms. For Delegate, see
	// the Delegate interface. For Events, see the EventDelegate interface.
	//
	// The DelegateProtocolMin/Max are used to guarantee protocol-compatibility
	// for any custom messages that the delegate might do (broadcasts,
	// local/remote state, etc.). If you don't set these, then the protocol
	// versions will just be zero, and version compliance won't be done.
	Delegate                Delegate
	DelegateProtocolVersion uint8
	DelegateProtocolMin     uint8
	DelegateProtocolMax     uint8
	Events                  EventDelegate
	Conflict                ConflictDelegate
	Merge                   MergeDelegate
	Ping                    PingDelegate
	Alive                   AliveDelegate

	// DNSConfigPath points to the system's DNS config file, usually located
	// at /etc/resolv.conf. It can be overridden via config for easier testing.
	DNSConfigPath string

	// LogOutput is the writer where logs should be sent. If this is not
	// set, logging will go to stderr by default. You cannot specify both LogOutput
	// and Logger at the same time.
	LogOutput io.Writer

	// Logger is a custom logger which you provide. If Logger is set, it will use
	// this for the internal logger. If Logger is not set, it will fall back to the
	// behavior for using LogOutput. You cannot specify both LogOutput and Logger
	// at the same time.
	Logger *log.Logger

	// Size of Memberlist's internal channel which handles UDP messages. The
	// size of this determines the size of the queue which Memberlist will keep
	// while UDP messages are handled.
	HandoffQueueDepth int
}

// DefaultLANConfig returns a sane set of configurations for Memberlist.
// It uses the hostname as the node name, and otherwise sets very conservative
// values that are sane for most LAN environments. The default configuration
// errs on the side of caution, choosing values that are optimized
// for higher convergence at the cost of higher bandwidth usage. Regardless,
// these values are a good starting point when getting started with memberlist.
func DefaultLANConfig() *Config {
	hostname, _ := os.Hostname()
	return &Config{
		Name:                    hostname,
		BindAddr:                "0.0.0.0",
		BindPort:                7946,
		AdvertiseAddr:           "",
		AdvertisePort:           7946,
		ProtocolVersion:         ProtocolVersion2Compatible,
		TCPTimeout:              10 * time.Second,       // Timeout after 10 seconds
		IndirectChecks:          3,                      // Use 3 nodes for the indirect ping
		RetransmitMult:          4,                      // Retransmit a message 4 * log(N+1) nodes
		SuspicionMult:           5,                      // Suspect a node for 5 * log(N+1) * Interval
		SuspicionMaxTimeoutMult: 6,                      // For 10k nodes this will give a max timeout of 120 seconds
		PushPullInterval:        30 * time.Second,       // Low frequency
		ProbeTimeout:            500 * time.Millisecond, // Reasonable RTT time for LAN
		ProbeInterval:           1 * time.Second,        // Failure check every second
		DisableTcpPings:         false,                  // TCP pings are safe, even with mixed versions
		AwarenessMaxMultiplier:  8,                      // Probe interval backs off to 8 seconds

		GossipNodes:         3,                      // Gossip to 3 nodes
		GossipInterval:      200 * time.Millisecond, // Gossip more rapidly
		GossipToTheDeadTime: 30 * time.Second,       // Same as push/pull

		EnableCompression: true, // Enable compression by default

		SecretKey: nil,
		Keyring:   nil,

		DNSConfigPath: "/etc/resolv.conf",

		HandoffQueueDepth: 1024,
	}
}

// DefaultWANConfig works like DefaultConfig, however it returns a configuration
// that is optimized for most WAN environments. The default configuration is
// still very conservative and errs on the side of caution.
func DefaultWANConfig() *Config {
	conf := DefaultLANConfig()
	conf.TCPTimeout = 30 * time.Second
	conf.SuspicionMult = 6
	conf.PushPullInterval = 60 * time.Second
	conf.ProbeTimeout = 3 * time.Second
	conf.ProbeInterval = 5 * time.Second
	conf.GossipNodes = 4 // Gossip less frequently, but to an additional node
	conf.GossipInterval = 500 * time.Millisecond
	conf.GossipToTheDeadTime = 60 * time.Second
	return conf
}

// DefaultLocalConfig works like DefaultConfig, however it returns a configuration
// that is optimized for a local loopback environments. The default configuration is
// still very conservative and errs on the side of caution.
func DefaultLocalConfig() *Config {
	conf := DefaultLANConfig()
	conf.TCPTimeout = time.Second
	conf.IndirectChecks = 1
	conf.RetransmitMult = 2
	conf.SuspicionMult = 3
	conf.PushPullInterval = 15 * time.Second
	conf.ProbeTimeout = 200 * time.Millisecond
	conf.ProbeInterval = time.Second
	conf.GossipInterval = 100 * time.Millisecond
	conf.GossipToTheDeadTime = 15 * time.Second
	return conf
}

// Returns whether or not encryption is enabled
func (c *Config) EncryptionEnabled() bool {
	return c.Keyring != nil && len(c.Keyring.GetKeys()) > 0
}
