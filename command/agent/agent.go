package agent

import (
	"fmt"
	"github.com/hashicorp/consul/consul"
	"io"
	"log"
	"os"
	"sync"
)

/*
 The agent is the long running process that is run on every machine.
 It exposes an RPC interface that is used by the CLI to control the
 agent. The agent runs the query interfaces like HTTP, DNS, and RPC.
 However, it can run in either a client, or server mode. In server
 mode, it runs a full Consul server. In client-only mode, it only forwards
 requests to other Consul servers.
*/
type Agent struct {
	config *Config

	// Used for writing our logs
	logger *log.Logger

	// Output sink for logs
	logOutput io.Writer

	// We have one of a client or a server, depending
	// on our configuration
	server *consul.Server
	client *consul.Client

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// Create is used to create a new Agent. Returns
// the agent or potentially an error.
func Create(config *Config, logOutput io.Writer) (*Agent, error) {
	// Ensure we have a log sink
	if logOutput == nil {
		logOutput = os.Stderr
	}

	// Validate the config
	if config.Datacenter == "" {
		return nil, fmt.Errorf("Must configure a Datacenter")
	}
	if config.DataDir == "" {
		return nil, fmt.Errorf("Must configure a DataDir")
	}

	agent := &Agent{
		config:     config,
		logger:     log.New(logOutput, "", log.LstdFlags),
		logOutput:  logOutput,
		shutdownCh: make(chan struct{}),
	}

	// Setup either the client or the server
	var err error
	if config.Server {
		err = agent.setupServer()
	} else {
		err = agent.setupClient()
	}
	if err != nil {
		return nil, err
	}

	return agent, nil
}

// consulConfig is used to return a consul configuration
func (a *Agent) consulConfig() *consul.Config {
	// Start with the provided config or default config
	var base *consul.Config
	if a.config.ConsulConfig != nil {
		base = a.config.ConsulConfig
	} else {
		base = consul.DefaultConfig()
	}

	// Override with our config
	if a.config.Datacenter != "" {
		base.Datacenter = a.config.Datacenter
	}
	if a.config.DataDir != "" {
		base.DataDir = a.config.DataDir
	}
	if a.config.EncryptKey != "" {
		key, _ := a.config.EncryptBytes()
		base.SerfLANConfig.MemberlistConfig.SecretKey = key
		base.SerfWANConfig.MemberlistConfig.SecretKey = key
	}
	if a.config.NodeName != "" {
		base.NodeName = a.config.NodeName
	}
	if a.config.SerfBindAddr != "" {
		base.SerfLANConfig.MemberlistConfig.BindAddr = a.config.SerfBindAddr
		base.SerfWANConfig.MemberlistConfig.BindAddr = a.config.SerfBindAddr
	}
	if a.config.SerfLanPort != 0 {
		base.SerfLANConfig.MemberlistConfig.BindPort = a.config.SerfLanPort
	}
	if a.config.SerfWanPort != 0 {
		base.SerfWANConfig.MemberlistConfig.BindPort = a.config.SerfWanPort
	}
	if a.config.ServerAddr != "" {
		base.RPCAddr = a.config.ServerAddr
	}
	if a.config.Bootstrap {
		base.Bootstrap = true
	}

	// Setup the loggers
	base.LogOutput = a.logOutput
	return base
}

// setupServer is used to initialize the Consul server
func (a *Agent) setupServer() error {
	server, err := consul.NewServer(a.consulConfig())
	if err != nil {
		return fmt.Errorf("Failed to start Consul server: %v", err)
	}
	a.server = server
	return nil
}

// setupClient is used to initialize the Consul client
func (a *Agent) setupClient() error {
	client, err := consul.NewClient(a.consulConfig())
	if err != nil {
		return fmt.Errorf("Failed to start Consul client: %v", err)
	}
	a.client = client
	return nil
}

// RPC is used to make an RPC call to the Consul servers
// This allows the agent to implement the Consul.Interface
func (a *Agent) RPC(method string, args interface{}, reply interface{}) error {
	if a.server != nil {
		return a.server.RPC(method, args, reply)
	}
	return a.client.RPC(method, args, reply)
}

// Leave prepares the agent for a graceful shutdown
func (a *Agent) Leave() error {
	if a.server != nil {
		return a.server.Leave()
	} else {
		return a.client.Leave()
	}
}

// Shutdown is used to hard stop the agent. Should be preceeded
// by a call to Leave to do it gracefully.
func (a *Agent) Shutdown() error {
	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()

	if a.shutdown {
		return nil
	}

	a.logger.Println("[INFO] agent: requesting shutdown")
	var err error
	if a.server != nil {
		err = a.server.Shutdown()
	} else {
		err = a.client.Shutdown()
	}

	a.logger.Println("[INFO] agent: shutdown complete")
	a.shutdown = true
	close(a.shutdownCh)
	return err
}

// ShutdownCh returns a channel that can be selected to wait
// for the agent to perform a shutdown.
func (a *Agent) ShutdownCh() <-chan struct{} {
	return a.shutdownCh
}
