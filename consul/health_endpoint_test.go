package consul

import (
	"github.com/hashicorp/consul/consul/structs"
	"os"
	"testing"
	"time"
)

func TestHealth_ChecksInState(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	client := rpcClient(t, s1)
	defer client.Close()

	// Wait for leader
	time.Sleep(100 * time.Millisecond)

	arg := structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
		Check: &structs.HealthCheck{
			Name:   "memory utilization",
			Status: structs.HealthPassing,
		},
	}
	var out struct{}
	if err := client.Call("Catalog.Register", &arg, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	var checks structs.HealthChecks
	inState := structs.ChecksInStateRequest{
		Datacenter: "dc1",
		State:      structs.HealthPassing,
	}
	if err := client.Call("Health.ChecksInState", &inState, &checks); err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(checks) != 1 {
		t.Fatalf("Bad: %v", checks)
	}
	if checks[0].Name != "memory utilization" {
		t.Fatalf("Bad: %v", checks)
	}
}

func TestHealth_NodeChecks(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	client := rpcClient(t, s1)
	defer client.Close()

	// Wait for leader
	time.Sleep(100 * time.Millisecond)

	arg := structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
		Check: &structs.HealthCheck{
			Name:   "memory utilization",
			Status: structs.HealthPassing,
		},
	}
	var out struct{}
	if err := client.Call("Catalog.Register", &arg, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	var checks structs.HealthChecks
	node := structs.NodeSpecificRequest{
		Datacenter: "dc1",
		Node:       "foo",
	}
	if err := client.Call("Health.NodeChecks", &node, &checks); err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(checks) != 1 {
		t.Fatalf("Bad: %v", checks)
	}
	if checks[0].Name != "memory utilization" {
		t.Fatalf("Bad: %v", checks)
	}
}

func TestHealth_ServiceChecks(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	client := rpcClient(t, s1)
	defer client.Close()

	// Wait for leader
	time.Sleep(100 * time.Millisecond)

	arg := structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
		Service: &structs.NodeService{
			ID:      "db",
			Service: "db",
		},
		Check: &structs.HealthCheck{
			Name:      "db connect",
			Status:    structs.HealthPassing,
			ServiceID: "db",
		},
	}
	var out struct{}
	if err := client.Call("Catalog.Register", &arg, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	var checks structs.HealthChecks
	node := structs.ServiceSpecificRequest{
		Datacenter:  "dc1",
		ServiceName: "db",
	}
	if err := client.Call("Health.ServiceChecks", &node, &checks); err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(checks) != 1 {
		t.Fatalf("Bad: %v", checks)
	}
	if checks[0].Name != "db connect" {
		t.Fatalf("Bad: %v", checks)
	}
}

func TestHealth_ServiceNodes(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	client := rpcClient(t, s1)
	defer client.Close()

	// Wait for leader
	time.Sleep(100 * time.Millisecond)

	arg := structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
		Service: &structs.NodeService{
			ID:      "db",
			Service: "db",
			Tag:     "master",
		},
		Check: &structs.HealthCheck{
			Name:      "db connect",
			Status:    structs.HealthPassing,
			ServiceID: "db",
		},
	}
	var out struct{}
	if err := client.Call("Catalog.Register", &arg, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	arg = structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "bar",
		Address:    "127.0.0.2",
		Service: &structs.NodeService{
			ID:      "db",
			Service: "db",
			Tag:     "slave",
		},
		Check: &structs.HealthCheck{
			Name:      "db connect",
			Status:    structs.HealthWarning,
			ServiceID: "db",
		},
	}
	if err := client.Call("Catalog.Register", &arg, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	var nodes structs.CheckServiceNodes
	req := structs.ServiceSpecificRequest{
		Datacenter:  "dc1",
		ServiceName: "db",
		ServiceTag:  "master",
		TagFilter:   false,
	}
	if err := client.Call("Health.ServiceNodes", &req, &nodes); err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(nodes) != 2 {
		t.Fatalf("Bad: %v", nodes)
	}
	if nodes[0].Node.Node != "foo" {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if nodes[1].Node.Node != "bar" {
		t.Fatalf("Bad: %v", nodes[1])
	}
	if nodes[0].Service.Tag != "master" {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if nodes[1].Service.Tag != "slave" {
		t.Fatalf("Bad: %v", nodes[1])
	}
	if nodes[0].Checks[0].Status != structs.HealthPassing {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if nodes[1].Checks[0].Status != structs.HealthWarning {
		t.Fatalf("Bad: %v", nodes[1])
	}
}
