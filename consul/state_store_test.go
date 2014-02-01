package consul

import (
	"github.com/hashicorp/consul/consul/structs"
	"reflect"
	"sort"
	"testing"
)

func TestEnsureNode(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v")
	}

	found, addr := store.GetNode("foo")
	if !found || addr != "127.0.0.1" {
		t.Fatalf("Bad: %v %v", found, addr)
	}

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.2"}); err != nil {
		t.Fatalf("err: %v")
	}

	found, addr = store.GetNode("foo")
	if !found || addr != "127.0.0.2" {
		t.Fatalf("Bad: %v %v", found, addr)
	}
}

func TestGetNodes(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureNode(structs.Node{"bar", "127.0.0.2"}); err != nil {
		t.Fatalf("err: %v")
	}

	nodes := store.Nodes()
	if len(nodes) != 2 {
		t.Fatalf("Bad: %v", nodes)
	}
	if nodes[1].Node != "foo" && nodes[0].Node != "bar" {
		t.Fatalf("Bad: %v", nodes)
	}
}

func BenchmarkGetNodes(b *testing.B) {
	store, err := NewStateStore()
	if err != nil {
		b.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		b.Fatalf("err: %v")
	}

	if err := store.EnsureNode(structs.Node{"bar", "127.0.0.2"}); err != nil {
		b.Fatalf("err: %v")
	}

	for i := 0; i < b.N; i++ {
		store.Nodes()
	}
}

func TestEnsureService(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := store.EnsureService("foo", "api", "api", "", 5000); err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := store.EnsureService("foo", "api", "api", "", 5001); err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := store.EnsureService("foo", "db", "db", "master", 8000); err != nil {
		t.Fatalf("err: %v", err)
	}

	services := store.NodeServices("foo")

	entry, ok := services.Services["api"]
	if !ok {
		t.Fatalf("missing api: %#v", services)
	}
	if entry.Tag != "" || entry.Port != 5001 {
		t.Fatalf("Bad entry: %#v", entry)
	}

	entry, ok = services.Services["db"]
	if !ok {
		t.Fatalf("missing db: %#v", services)
	}
	if entry.Tag != "master" || entry.Port != 8000 {
		t.Fatalf("Bad entry: %#v", entry)
	}
}

func TestEnsureService_DuplicateNode(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := store.EnsureService("foo", "api1", "api", "", 5000); err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := store.EnsureService("foo", "api2", "api", "", 5001); err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := store.EnsureService("foo", "api3", "api", "", 5002); err != nil {
		t.Fatalf("err: %v", err)
	}

	services := store.NodeServices("foo")

	entry, ok := services.Services["api1"]
	if !ok {
		t.Fatalf("missing api: %#v", services)
	}
	if entry.Tag != "" || entry.Port != 5000 {
		t.Fatalf("Bad entry: %#v", entry)
	}

	entry, ok = services.Services["api2"]
	if !ok {
		t.Fatalf("missing api: %#v", services)
	}
	if entry.Tag != "" || entry.Port != 5001 {
		t.Fatalf("Bad entry: %#v", entry)
	}

	entry, ok = services.Services["api3"]
	if !ok {
		t.Fatalf("missing api: %#v", services)
	}
	if entry.Tag != "" || entry.Port != 5002 {
		t.Fatalf("Bad entry: %#v", entry)
	}
}

func TestDeleteNodeService(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := store.EnsureService("foo", "api", "api", "", 5000); err != nil {
		t.Fatalf("err: %v", err)
	}

	check := &structs.HealthCheck{
		Node:      "foo",
		CheckID:   "api",
		Name:      "Can connect",
		Status:    structs.HealthPassing,
		ServiceID: "api",
	}
	if err := store.EnsureCheck(check); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.DeleteNodeService("foo", "api"); err != nil {
		t.Fatalf("err: %v", err)
	}

	services := store.NodeServices("foo")
	_, ok := services.Services["api"]
	if ok {
		t.Fatalf("has api: %#v", services)
	}

	checks := store.NodeChecks("foo")
	if len(checks) != 0 {
		t.Fatalf("has check: %#v", checks)
	}
}

func TestDeleteNodeService_One(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := store.EnsureService("foo", "api", "api", "", 5000); err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := store.EnsureService("foo", "api2", "api", "", 5001); err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := store.DeleteNodeService("foo", "api"); err != nil {
		t.Fatalf("err: %v", err)
	}

	services := store.NodeServices("foo")
	_, ok := services.Services["api"]
	if ok {
		t.Fatalf("has api: %#v", services)
	}
	_, ok = services.Services["api2"]
	if !ok {
		t.Fatalf("does not have api2: %#v", services)
	}
}

func TestDeleteNode(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "api", "api", "", 5000); err != nil {
		t.Fatalf("err: %v")
	}

	check := &structs.HealthCheck{
		Node:      "foo",
		CheckID:   "db",
		Name:      "Can connect",
		Status:    structs.HealthPassing,
		ServiceID: "api",
	}
	if err := store.EnsureCheck(check); err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := store.DeleteNode("foo"); err != nil {
		t.Fatalf("err: %v", err)
	}

	services := store.NodeServices("foo")
	_, ok := services.Services["api"]
	if ok {
		t.Fatalf("has api: %#v", services)
	}

	checks := store.NodeChecks("foo")
	if len(checks) > 0 {
		t.Fatalf("has checks: %v", checks)
	}

	found, _ := store.GetNode("foo")
	if found {
		t.Fatalf("found node")
	}
}

func TestGetServices(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureNode(structs.Node{"bar", "127.0.0.2"}); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "api", "api", "", 5000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "db", "db", "master", 8000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("bar", "db", "db", "slave", 8000); err != nil {
		t.Fatalf("err: %v")
	}

	services := store.Services()

	tags, ok := services["api"]
	if !ok {
		t.Fatalf("missing api: %#v", services)
	}
	if len(tags) != 1 || tags[0] != "" {
		t.Fatalf("Bad entry: %#v", tags)
	}

	tags, ok = services["db"]
	sort.Strings(tags)
	if !ok {
		t.Fatalf("missing db: %#v", services)
	}
	if len(tags) != 2 || tags[0] != "master" || tags[1] != "slave" {
		t.Fatalf("Bad entry: %#v", tags)
	}
}

func TestServiceNodes(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureNode(structs.Node{"bar", "127.0.0.2"}); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "api", "api", "", 5000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("bar", "api", "api", "", 5000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "db", "db", "master", 8000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("bar", "db", "db", "slave", 8000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("bar", "db2", "db", "slave", 8001); err != nil {
		t.Fatalf("err: %v")
	}

	nodes := store.ServiceNodes("db")
	if len(nodes) != 3 {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[0].Node != "foo" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[0].Address != "127.0.0.1" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[0].ServiceID != "db" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[0].ServiceTag != "master" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[0].ServicePort != 8000 {
		t.Fatalf("bad: %v", nodes)
	}

	if nodes[1].Node != "bar" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[1].Address != "127.0.0.2" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[1].ServiceID != "db" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[1].ServiceTag != "slave" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[1].ServicePort != 8000 {
		t.Fatalf("bad: %v", nodes)
	}

	if nodes[2].Node != "bar" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[2].Address != "127.0.0.2" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[2].ServiceID != "db2" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[2].ServiceTag != "slave" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[2].ServicePort != 8001 {
		t.Fatalf("bad: %v", nodes)
	}
}

func TestServiceTagNodes(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureNode(structs.Node{"bar", "127.0.0.2"}); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "db", "db", "master", 8000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "db2", "db", "slave", 8001); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("bar", "db", "db", "slave", 8000); err != nil {
		t.Fatalf("err: %v")
	}

	nodes := store.ServiceTagNodes("db", "master")
	if len(nodes) != 1 {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[0].Node != "foo" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[0].Address != "127.0.0.1" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[0].ServiceTag != "master" {
		t.Fatalf("bad: %v", nodes)
	}
	if nodes[0].ServicePort != 8000 {
		t.Fatalf("bad: %v", nodes)
	}
}

func TestStoreSnapshot(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureNode(structs.Node{"bar", "127.0.0.2"}); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "db", "db", "master", 8000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "db2", "db", "slave", 8001); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("bar", "db", "db", "slave", 8000); err != nil {
		t.Fatalf("err: %v")
	}

	check := &structs.HealthCheck{
		Node:      "foo",
		CheckID:   "db",
		Name:      "Can connect",
		Status:    structs.HealthPassing,
		ServiceID: "db",
	}
	if err := store.EnsureCheck(check); err != nil {
		t.Fatalf("err: %v")
	}

	// Take a snapshot
	snap, err := store.Snapshot()
	if err != nil {
		t.Fatalf("err: %v")
	}
	defer snap.Close()

	// Check snapshot has old values
	nodes := snap.Nodes()
	if len(nodes) != 2 {
		t.Fatalf("bad: %v", nodes)
	}

	// Ensure we get the service entries
	services := snap.NodeServices("foo")
	if services.Services["db"].Tag != "master" {
		t.Fatalf("bad: %v", services)
	}
	if services.Services["db2"].Tag != "slave" {
		t.Fatalf("bad: %v", services)
	}

	services = snap.NodeServices("bar")
	if services.Services["db"].Tag != "slave" {
		t.Fatalf("bad: %v", services)
	}

	// Ensure we get the checks
	checks := snap.NodeChecks("foo")
	if len(checks) != 1 {
		t.Fatalf("bad: %v", checks)
	}
	if !reflect.DeepEqual(checks[0], check) {
		t.Fatalf("bad: %v", checks[0])
	}

	// Make some changes!
	if err := store.EnsureService("foo", "db", "db", "slave", 8000); err != nil {
		t.Fatalf("err: %v", err)
	}
	if err := store.EnsureService("bar", "db", "db", "master", 8000); err != nil {
		t.Fatalf("err: %v", err)
	}
	if err := store.EnsureNode(structs.Node{"baz", "127.0.0.3"}); err != nil {
		t.Fatalf("err: %v", err)
	}
	checkAfter := &structs.HealthCheck{
		Node:      "foo",
		CheckID:   "db",
		Name:      "Can connect",
		Status:    structs.HealthCritical,
		ServiceID: "db",
	}
	if err := store.EnsureCheck(checkAfter); err != nil {
		t.Fatalf("err: %v")
	}

	// Check snapshot has old values
	nodes = snap.Nodes()
	if len(nodes) != 2 {
		t.Fatalf("bad: %v", nodes)
	}

	// Ensure old service entries
	services = snap.NodeServices("foo")
	if services.Services["db"].Tag != "master" {
		t.Fatalf("bad: %v", services)
	}
	if services.Services["db2"].Tag != "slave" {
		t.Fatalf("bad: %v", services)
	}

	services = snap.NodeServices("bar")
	if services.Services["db"].Tag != "slave" {
		t.Fatalf("bad: %v", services)
	}

	checks = snap.NodeChecks("foo")
	if len(checks) != 1 {
		t.Fatalf("bad: %v", checks)
	}
	if !reflect.DeepEqual(checks[0], check) {
		t.Fatalf("bad: %v", checks[0])
	}
}

func TestEnsureCheck(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if err := store.EnsureService("foo", "db1", "db", "master", 8000); err != nil {
		t.Fatalf("err: %v")
	}
	check := &structs.HealthCheck{
		Node:      "foo",
		CheckID:   "db",
		Name:      "Can connect",
		Status:    structs.HealthPassing,
		ServiceID: "db1",
	}
	if err := store.EnsureCheck(check); err != nil {
		t.Fatalf("err: %v")
	}

	check2 := &structs.HealthCheck{
		Node:    "foo",
		CheckID: "memory",
		Name:    "memory utilization",
		Status:  structs.HealthWarning,
	}
	if err := store.EnsureCheck(check2); err != nil {
		t.Fatalf("err: %v")
	}

	checks := store.NodeChecks("foo")
	if len(checks) != 2 {
		t.Fatalf("bad: %v", checks)
	}
	if !reflect.DeepEqual(checks[0], check) {
		t.Fatalf("bad: %v", checks[0])
	}
	if !reflect.DeepEqual(checks[1], check2) {
		t.Fatalf("bad: %v", checks[1])
	}

	checks = store.ServiceChecks("db")
	if len(checks) != 1 {
		t.Fatalf("bad: %v", checks)
	}
	if !reflect.DeepEqual(checks[0], check) {
		t.Fatalf("bad: %v", checks[0])
	}

	checks = store.ChecksInState(structs.HealthPassing)
	if len(checks) != 1 {
		t.Fatalf("bad: %v", checks)
	}
	if !reflect.DeepEqual(checks[0], check) {
		t.Fatalf("bad: %v", checks[0])
	}

	checks = store.ChecksInState(structs.HealthWarning)
	if len(checks) != 1 {
		t.Fatalf("bad: %v", checks)
	}
	if !reflect.DeepEqual(checks[0], check2) {
		t.Fatalf("bad: %v", checks[0])
	}
}

func TestDeleteNodeCheck(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if err := store.EnsureService("foo", "db1", "db", "master", 8000); err != nil {
		t.Fatalf("err: %v")
	}
	check := &structs.HealthCheck{
		Node:      "foo",
		CheckID:   "db",
		Name:      "Can connect",
		Status:    structs.HealthPassing,
		ServiceID: "db1",
	}
	if err := store.EnsureCheck(check); err != nil {
		t.Fatalf("err: %v")
	}

	check2 := &structs.HealthCheck{
		Node:    "foo",
		CheckID: "memory",
		Name:    "memory utilization",
		Status:  structs.HealthWarning,
	}
	if err := store.EnsureCheck(check2); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.DeleteNodeCheck("foo", "db"); err != nil {
		t.Fatalf("err: %v", err)
	}

	checks := store.NodeChecks("foo")
	if len(checks) != 1 {
		t.Fatalf("bad: %v", checks)
	}
	if !reflect.DeepEqual(checks[0], check2) {
		t.Fatalf("bad: %v", checks[0])
	}
}

func TestCheckServiceNodes(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if err := store.EnsureService("foo", "db1", "db", "master", 8000); err != nil {
		t.Fatalf("err: %v")
	}
	check := &structs.HealthCheck{
		Node:      "foo",
		CheckID:   "db",
		Name:      "Can connect",
		Status:    structs.HealthPassing,
		ServiceID: "db1",
	}
	if err := store.EnsureCheck(check); err != nil {
		t.Fatalf("err: %v")
	}
	check = &structs.HealthCheck{
		Node:    "foo",
		CheckID: SerfCheckID,
		Name:    SerfCheckName,
		Status:  structs.HealthPassing,
	}
	if err := store.EnsureCheck(check); err != nil {
		t.Fatalf("err: %v")
	}

	nodes := store.CheckServiceNodes("db")
	if len(nodes) != 1 {
		t.Fatalf("Bad: %v", nodes)
	}

	if nodes[0].Node.Node != "foo" {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if nodes[0].Service.ID != "db1" {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if len(nodes[0].Checks) != 2 {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if nodes[0].Checks[0].CheckID != "db" {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if nodes[0].Checks[1].CheckID != SerfCheckID {
		t.Fatalf("Bad: %v", nodes[0])
	}

	nodes = store.CheckServiceTagNodes("db", "master")
	if len(nodes) != 1 {
		t.Fatalf("Bad: %v", nodes)
	}

	if nodes[0].Node.Node != "foo" {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if nodes[0].Service.ID != "db1" {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if len(nodes[0].Checks) != 2 {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if nodes[0].Checks[0].CheckID != "db" {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if nodes[0].Checks[1].CheckID != SerfCheckID {
		t.Fatalf("Bad: %v", nodes[0])
	}
}
func BenchmarkCheckServiceNodes(t *testing.B) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode(structs.Node{"foo", "127.0.0.1"}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if err := store.EnsureService("foo", "db1", "db", "master", 8000); err != nil {
		t.Fatalf("err: %v")
	}
	check := &structs.HealthCheck{
		Node:      "foo",
		CheckID:   "db",
		Name:      "Can connect",
		Status:    structs.HealthPassing,
		ServiceID: "db1",
	}
	if err := store.EnsureCheck(check); err != nil {
		t.Fatalf("err: %v")
	}
	check = &structs.HealthCheck{
		Node:    "foo",
		CheckID: SerfCheckID,
		Name:    SerfCheckName,
		Status:  structs.HealthPassing,
	}
	if err := store.EnsureCheck(check); err != nil {
		t.Fatalf("err: %v")
	}

	for i := 0; i < t.N; i++ {
		store.CheckServiceNodes("db")
	}
}
