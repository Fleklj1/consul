package agent

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/consul/consul/structs"
	"github.com/hashicorp/consul/logger"
	"github.com/hashicorp/consul/testutil"
	"github.com/hashicorp/consul/types"
	"github.com/hashicorp/serf/serf"
	"github.com/mitchellh/cli"
)

func TestHTTPAgentServices(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	srv1 := &structs.NodeService{
		ID:      "mysql",
		Service: "mysql",
		Tags:    []string{"master"},
		Port:    5000,
	}
	srv.agent.state.AddService(srv1, "")

	obj, err := srv.AgentServices(nil, nil)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	val := obj.(map[string]*structs.NodeService)
	if len(val) != 2 {
		t.Fatalf("bad services: %v", obj)
	}
	if val["mysql"].Port != 5000 {
		t.Fatalf("bad service: %v", obj)
	}
}

func TestHTTPAgentChecks(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	chk1 := &structs.HealthCheck{
		Node:    srv.agent.config.NodeName,
		CheckID: "mysql",
		Name:    "mysql",
		Status:  structs.HealthPassing,
	}
	srv.agent.state.AddCheck(chk1, "")

	obj, err := srv.AgentChecks(nil, nil)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	val := obj.(map[types.CheckID]*structs.HealthCheck)
	if len(val) != 1 {
		t.Fatalf("bad checks: %v", obj)
	}
	if val["mysql"].Status != structs.HealthPassing {
		t.Fatalf("bad check: %v", obj)
	}
}

func TestHTTPAgentSelf(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	req, err := http.NewRequest("GET", "/v1/agent/self", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	obj, err := srv.AgentSelf(nil, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	val := obj.(AgentSelf)
	if int(val.Member.Port) != srv.agent.config.Ports.SerfLan {
		t.Fatalf("incorrect port: %v", obj)
	}

	if int(val.Config.Ports.SerfLan) != srv.agent.config.Ports.SerfLan {
		t.Fatalf("incorrect port: %v", obj)
	}

	c, err := srv.agent.server.GetLANCoordinate()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(c, val.Coord) {
		t.Fatalf("coordinates are not equal: %v != %v", c, val.Coord)
	}

	srv.agent.config.DisableCoordinates = true
	obj, err = srv.AgentSelf(nil, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	val = obj.(AgentSelf)
	if val.Coord != nil {
		t.Fatalf("should have been nil: %v", val.Coord)
	}
}

func TestHTTPAgentReload(t *testing.T) {
	conf := nextConfig()
	tmpDir, err := ioutil.TempDir("", "consul")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write initial config, to be reloaded later
	tmpFile, err := ioutil.TempFile(tmpDir, "config")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	_, err = tmpFile.WriteString(`{"service":{"name":"redis"}}`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	tmpFile.Close()

	doneCh := make(chan struct{})
	shutdownCh := make(chan struct{})

	defer func() {
		close(shutdownCh)
		<-doneCh
	}()

	cmd := &Command{
		ShutdownCh: shutdownCh,
		Ui:         new(cli.MockUi),
	}

	args := []string{
		"-server",
		"-data-dir", tmpDir,
		"-http-port", fmt.Sprintf("%d", conf.Ports.HTTP),
		"-config-file", tmpFile.Name(),
	}

	go func() {
		cmd.Run(args)
		close(doneCh)
	}()

	testutil.WaitForResult(func() (bool, error) {
		return len(cmd.httpServers) == 1, nil
	}, func(err error) {
		t.Fatalf("should have an http server")
	})

	if _, ok := cmd.agent.state.services["redis"]; !ok {
		t.Fatalf("missing redis service")
	}

	err = ioutil.WriteFile(tmpFile.Name(), []byte(`{"service":{"name":"redis-reloaded"}}`), 0644)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	srv := cmd.httpServers[0]
	req, err := http.NewRequest("PUT", "/v1/agent/reload", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	_, err = srv.AgentReload(nil, req)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	if _, ok := cmd.agent.state.services["redis-reloaded"]; !ok {
		t.Fatalf("missing redis-reloaded service")
	}
}

func TestHTTPAgentMembers(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	req, err := http.NewRequest("GET", "/v1/agent/members", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	obj, err := srv.AgentMembers(nil, req)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	val := obj.([]serf.Member)
	if len(val) == 0 {
		t.Fatalf("bad members: %v", obj)
	}

	if int(val[0].Port) != srv.agent.config.Ports.SerfLan {
		t.Fatalf("not lan: %v", obj)
	}
}

func TestHTTPAgentMembers_WAN(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	req, err := http.NewRequest("GET", "/v1/agent/members?wan=true", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	obj, err := srv.AgentMembers(nil, req)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	val := obj.([]serf.Member)
	if len(val) == 0 {
		t.Fatalf("bad members: %v", obj)
	}

	if int(val[0].Port) != srv.agent.config.Ports.SerfWan {
		t.Fatalf("not wan: %v", obj)
	}
}

func TestHTTPAgentJoin(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	dir2, a2 := makeAgent(t, nextConfig())
	defer os.RemoveAll(dir2)
	defer a2.Shutdown()

	addr := fmt.Sprintf("127.0.0.1:%d", a2.config.Ports.SerfLan)
	req, err := http.NewRequest("GET", fmt.Sprintf("/v1/agent/join/%s", addr), nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	obj, err := srv.AgentJoin(nil, req)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	if obj != nil {
		t.Fatalf("Err: %v", obj)
	}

	if len(srv.agent.LANMembers()) != 2 {
		t.Fatalf("should have 2 members")
	}

	testutil.WaitForResult(func() (bool, error) {
		return len(a2.LANMembers()) == 2, nil
	}, func(err error) {
		t.Fatalf("should have 2 members")
	})
}

func TestHTTPAgentJoin_WAN(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	dir2, a2 := makeAgent(t, nextConfig())
	defer os.RemoveAll(dir2)
	defer a2.Shutdown()

	addr := fmt.Sprintf("127.0.0.1:%d", a2.config.Ports.SerfWan)
	req, err := http.NewRequest("GET", fmt.Sprintf("/v1/agent/join/%s?wan=true", addr), nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	obj, err := srv.AgentJoin(nil, req)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	if obj != nil {
		t.Fatalf("Err: %v", obj)
	}

	if len(srv.agent.WANMembers()) != 2 {
		t.Fatalf("should have 2 members")
	}

	testutil.WaitForResult(func() (bool, error) {
		return len(a2.WANMembers()) == 2, nil
	}, func(err error) {
		t.Fatalf("should have 2 members")
	})
}

func TestHTTPAgentLeave(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	dir2, srv2 := makeHTTPServerWithConfig(t, func(c *Config) {
		c.Server = false
		c.Bootstrap = false
	})
	defer os.RemoveAll(dir2)
	defer srv2.Shutdown()

	// Join first
	addr := fmt.Sprintf("127.0.0.1:%d", srv2.agent.config.Ports.SerfLan)
	_, err := srv.agent.JoinLAN([]string{addr})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Graceful leave now
	req, err := http.NewRequest("PUT", "/v1/agent/leave", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	obj, err := srv2.AgentLeave(nil, req)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	if obj != nil {
		t.Fatalf("Err: %v", obj)
	}

	testutil.WaitForResult(func() (bool, error) {
		m := srv.agent.LANMembers()
		success := m[1].Status == serf.StatusLeft
		return success, errors.New(m[1].Status.String())
	}, func(err error) {
		t.Fatalf("member status is %v, should be left", err)
	})
}

func TestHTTPAgentForceLeave(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	dir2, a2 := makeAgent(t, nextConfig())
	defer os.RemoveAll(dir2)
	defer a2.Shutdown()

	// Join first
	addr := fmt.Sprintf("127.0.0.1:%d", a2.config.Ports.SerfLan)
	_, err := srv.agent.JoinLAN([]string{addr})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	a2.Shutdown()

	// Force leave now
	req, err := http.NewRequest("GET", fmt.Sprintf("/v1/agent/force-leave/%s", a2.config.NodeName), nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	obj, err := srv.AgentForceLeave(nil, req)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	if obj != nil {
		t.Fatalf("Err: %v", obj)
	}

	testutil.WaitForResult(func() (bool, error) {
		m := srv.agent.LANMembers()
		success := m[1].Status == serf.StatusLeft
		return success, errors.New(m[1].Status.String())
	}, func(err error) {
		t.Fatalf("member status is %v, should be left", err)
	})
}

func TestHTTPAgentRegisterCheck(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	// Register node
	req, err := http.NewRequest("GET", "/v1/agent/check/register?token=abc123", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	args := &CheckDefinition{
		Name: "test",
		CheckType: CheckType{
			TTL: 15 * time.Second,
		},
	}
	req.Body = encodeReq(args)

	obj, err := srv.AgentRegisterCheck(nil, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if obj != nil {
		t.Fatalf("bad: %v", obj)
	}

	// Ensure we have a check mapping
	checkID := types.CheckID("test")
	if _, ok := srv.agent.state.Checks()[checkID]; !ok {
		t.Fatalf("missing test check")
	}

	if _, ok := srv.agent.checkTTLs[checkID]; !ok {
		t.Fatalf("missing test check ttl")
	}

	// Ensure the token was configured
	if token := srv.agent.state.CheckToken(checkID); token == "" {
		t.Fatalf("missing token")
	}

	// By default, checks start in critical state.
	state := srv.agent.state.Checks()[checkID]
	if state.Status != structs.HealthCritical {
		t.Fatalf("bad: %v", state)
	}
}

func TestHTTPAgentRegisterCheckPassing(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	// Register node
	req, err := http.NewRequest("GET", "/v1/agent/check/register", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	args := &CheckDefinition{
		Name: "test",
		CheckType: CheckType{
			TTL: 15 * time.Second,
		},
		Status: structs.HealthPassing,
	}
	req.Body = encodeReq(args)

	obj, err := srv.AgentRegisterCheck(nil, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if obj != nil {
		t.Fatalf("bad: %v", obj)
	}

	// Ensure we have a check mapping
	checkID := types.CheckID("test")
	if _, ok := srv.agent.state.Checks()[checkID]; !ok {
		t.Fatalf("missing test check")
	}

	if _, ok := srv.agent.checkTTLs[checkID]; !ok {
		t.Fatalf("missing test check ttl")
	}

	state := srv.agent.state.Checks()[checkID]
	if state.Status != structs.HealthPassing {
		t.Fatalf("bad: %v", state)
	}
}

func TestHTTPAgentRegisterCheckBadStatus(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	// Register node
	req, err := http.NewRequest("GET", "/v1/agent/check/register", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	args := &CheckDefinition{
		Name: "test",
		CheckType: CheckType{
			TTL: 15 * time.Second,
		},
		Status: "fluffy",
	}
	req.Body = encodeReq(args)

	resp := httptest.NewRecorder()
	if _, err := srv.AgentRegisterCheck(resp, req); err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Code != 400 {
		t.Fatalf("accepted bad status")
	}
}

func TestHTTPAgentDeregisterCheck(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
	if err := srv.agent.AddCheck(chk, nil, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register node
	req, err := http.NewRequest("GET", "/v1/agent/check/deregister/test", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	obj, err := srv.AgentDeregisterCheck(nil, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if obj != nil {
		t.Fatalf("bad: %v", obj)
	}

	// Ensure we have a check mapping
	if _, ok := srv.agent.state.Checks()["test"]; ok {
		t.Fatalf("have test check")
	}
}

func TestHTTPAgentPassCheck(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
	chkType := &CheckType{TTL: 15 * time.Second}
	if err := srv.agent.AddCheck(chk, chkType, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	req, err := http.NewRequest("GET", "/v1/agent/check/pass/test", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	obj, err := srv.AgentCheckPass(nil, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if obj != nil {
		t.Fatalf("bad: %v", obj)
	}

	// Ensure we have a check mapping
	state := srv.agent.state.Checks()["test"]
	if state.Status != structs.HealthPassing {
		t.Fatalf("bad: %v", state)
	}
}

func TestHTTPAgentWarnCheck(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
	chkType := &CheckType{TTL: 15 * time.Second}
	if err := srv.agent.AddCheck(chk, chkType, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	req, err := http.NewRequest("GET", "/v1/agent/check/warn/test", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	obj, err := srv.AgentCheckWarn(nil, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if obj != nil {
		t.Fatalf("bad: %v", obj)
	}

	// Ensure we have a check mapping
	state := srv.agent.state.Checks()["test"]
	if state.Status != structs.HealthWarning {
		t.Fatalf("bad: %v", state)
	}
}

func TestHTTPAgentFailCheck(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
	chkType := &CheckType{TTL: 15 * time.Second}
	if err := srv.agent.AddCheck(chk, chkType, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	req, err := http.NewRequest("GET", "/v1/agent/check/fail/test", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	obj, err := srv.AgentCheckFail(nil, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if obj != nil {
		t.Fatalf("bad: %v", obj)
	}

	// Ensure we have a check mapping
	state := srv.agent.state.Checks()["test"]
	if state.Status != structs.HealthCritical {
		t.Fatalf("bad: %v", state)
	}
}

func TestHTTPAgentUpdateCheck(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
	chkType := &CheckType{TTL: 15 * time.Second}
	if err := srv.agent.AddCheck(chk, chkType, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	cases := []checkUpdate{
		checkUpdate{structs.HealthPassing, "hello-passing"},
		checkUpdate{structs.HealthCritical, "hello-critical"},
		checkUpdate{structs.HealthWarning, "hello-warning"},
	}

	for _, c := range cases {
		req, err := http.NewRequest("PUT", "/v1/agent/check/update/test", nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		req.Body = encodeReq(c)

		resp := httptest.NewRecorder()
		obj, err := srv.AgentCheckUpdate(resp, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if obj != nil {
			t.Fatalf("bad: %v", obj)
		}
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}

		state := srv.agent.state.Checks()["test"]
		if state.Status != c.Status || state.Output != c.Output {
			t.Fatalf("bad: %v", state)
		}
	}

	// Make sure abusive levels of output are capped.
	{
		req, err := http.NewRequest("PUT", "/v1/agent/check/update/test", nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		update := checkUpdate{
			Status: structs.HealthPassing,
			Output: strings.Repeat("-= bad -=", 5*CheckBufSize),
		}
		req.Body = encodeReq(update)

		resp := httptest.NewRecorder()
		obj, err := srv.AgentCheckUpdate(resp, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if obj != nil {
			t.Fatalf("bad: %v", obj)
		}
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}

		// Since we append some notes about truncating, we just do a
		// rough check that the output buffer was cut down so this test
		// isn't super brittle.
		state := srv.agent.state.Checks()["test"]
		if state.Status != structs.HealthPassing || len(state.Output) > 2*CheckBufSize {
			t.Fatalf("bad: %v", state)
		}
	}

	// Check a bogus status.
	{
		req, err := http.NewRequest("PUT", "/v1/agent/check/update/test", nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		update := checkUpdate{
			Status: "itscomplicated",
		}
		req.Body = encodeReq(update)

		resp := httptest.NewRecorder()
		obj, err := srv.AgentCheckUpdate(resp, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if obj != nil {
			t.Fatalf("bad: %v", obj)
		}
		if resp.Code != 400 {
			t.Fatalf("expected 400, got %d", resp.Code)
		}
	}

	// Check a bogus verb.
	{
		req, err := http.NewRequest("POST", "/v1/agent/check/update/test", nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		update := checkUpdate{
			Status: structs.HealthPassing,
		}
		req.Body = encodeReq(update)

		resp := httptest.NewRecorder()
		obj, err := srv.AgentCheckUpdate(resp, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if obj != nil {
			t.Fatalf("bad: %v", obj)
		}
		if resp.Code != 405 {
			t.Fatalf("expected 405, got %d", resp.Code)
		}
	}
}

func TestHTTPAgentRegisterService(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	// Register node
	req, err := http.NewRequest("GET", "/v1/agent/service/register?token=abc123", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	args := &ServiceDefinition{
		Name: "test",
		Tags: []string{"master"},
		Port: 8000,
		Check: CheckType{
			TTL: 15 * time.Second,
		},
		Checks: CheckTypes{
			&CheckType{
				TTL: 20 * time.Second,
			},
			&CheckType{
				TTL: 30 * time.Second,
			},
		},
	}
	req.Body = encodeReq(args)

	obj, err := srv.AgentRegisterService(nil, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if obj != nil {
		t.Fatalf("bad: %v", obj)
	}

	// Ensure the servie
	if _, ok := srv.agent.state.Services()["test"]; !ok {
		t.Fatalf("missing test service")
	}

	// Ensure we have a check mapping
	checks := srv.agent.state.Checks()
	if len(checks) != 3 {
		t.Fatalf("bad: %v", checks)
	}

	if len(srv.agent.checkTTLs) != 3 {
		t.Fatalf("missing test check ttls: %v", srv.agent.checkTTLs)
	}

	// Ensure the token was configured
	if token := srv.agent.state.ServiceToken("test"); token == "" {
		t.Fatalf("missing token")
	}
}

func TestHTTPAgentDeregisterService(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	service := &structs.NodeService{
		ID:      "test",
		Service: "test",
	}
	if err := srv.agent.AddService(service, nil, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register node
	req, err := http.NewRequest("GET", "/v1/agent/service/deregister/test", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	obj, err := srv.AgentDeregisterService(nil, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if obj != nil {
		t.Fatalf("bad: %v", obj)
	}

	// Ensure we have a check mapping
	if _, ok := srv.agent.state.Services()["test"]; ok {
		t.Fatalf("have test service")
	}

	if _, ok := srv.agent.state.Checks()["test"]; ok {
		t.Fatalf("have test check")
	}
}

func TestHTTPAgent_ServiceMaintenanceEndpoint_BadRequest(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	// Fails on non-PUT
	req, _ := http.NewRequest("GET", "/v1/agent/service/maintenance/test?enable=true", nil)
	resp := httptest.NewRecorder()
	if _, err := srv.AgentServiceMaintenance(resp, req); err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Code != 405 {
		t.Fatalf("expected 405, got %d", resp.Code)
	}

	// Fails when no enable flag provided
	req, _ = http.NewRequest("PUT", "/v1/agent/service/maintenance/test", nil)
	resp = httptest.NewRecorder()
	if _, err := srv.AgentServiceMaintenance(resp, req); err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Code != 400 {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	// Fails when no service ID provided
	req, _ = http.NewRequest("PUT", "/v1/agent/service/maintenance/?enable=true", nil)
	resp = httptest.NewRecorder()
	if _, err := srv.AgentServiceMaintenance(resp, req); err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Code != 400 {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	// Fails when bad service ID provided
	req, _ = http.NewRequest("PUT", "/v1/agent/service/maintenance/_nope_?enable=true", nil)
	resp = httptest.NewRecorder()
	if _, err := srv.AgentServiceMaintenance(resp, req); err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Code != 404 {
		t.Fatalf("expected 404, got %d", resp.Code)
	}
}

func TestHTTPAgent_EnableServiceMaintenance(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	// Register the service
	service := &structs.NodeService{
		ID:      "test",
		Service: "test",
	}
	if err := srv.agent.AddService(service, nil, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Force the service into maintenance mode
	req, _ := http.NewRequest("PUT", "/v1/agent/service/maintenance/test?enable=true&reason=broken&token=mytoken", nil)
	resp := httptest.NewRecorder()
	if _, err := srv.AgentServiceMaintenance(resp, req); err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Code != 200 {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	// Ensure the maintenance check was registered
	checkID := serviceMaintCheckID("test")
	check, ok := srv.agent.state.Checks()[checkID]
	if !ok {
		t.Fatalf("should have registered maintenance check")
	}

	// Ensure the token was added
	if token := srv.agent.state.CheckToken(checkID); token != "mytoken" {
		t.Fatalf("expected 'mytoken', got '%s'", token)
	}

	// Ensure the reason was set in notes
	if check.Notes != "broken" {
		t.Fatalf("bad: %#v", check)
	}
}

func TestHTTPAgent_DisableServiceMaintenance(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	// Register the service
	service := &structs.NodeService{
		ID:      "test",
		Service: "test",
	}
	if err := srv.agent.AddService(service, nil, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Force the service into maintenance mode
	if err := srv.agent.EnableServiceMaintenance("test", "", ""); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Leave maintenance mode
	req, _ := http.NewRequest("PUT", "/v1/agent/service/maintenance/test?enable=false", nil)
	resp := httptest.NewRecorder()
	if _, err := srv.AgentServiceMaintenance(resp, req); err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Code != 200 {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	// Ensure the maintenance check was removed
	checkID := serviceMaintCheckID("test")
	if _, ok := srv.agent.state.Checks()[checkID]; ok {
		t.Fatalf("should have removed maintenance check")
	}
}

func TestHTTPAgent_NodeMaintenanceEndpoint_BadRequest(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	// Fails on non-PUT
	req, _ := http.NewRequest("GET", "/v1/agent/self/maintenance?enable=true", nil)
	resp := httptest.NewRecorder()
	if _, err := srv.AgentNodeMaintenance(resp, req); err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Code != 405 {
		t.Fatalf("expected 405, got %d", resp.Code)
	}

	// Fails when no enable flag provided
	req, _ = http.NewRequest("PUT", "/v1/agent/self/maintenance", nil)
	resp = httptest.NewRecorder()
	if _, err := srv.AgentNodeMaintenance(resp, req); err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Code != 400 {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestHTTPAgent_EnableNodeMaintenance(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	// Force the node into maintenance mode
	req, _ := http.NewRequest(
		"PUT", "/v1/agent/self/maintenance?enable=true&reason=broken&token=mytoken", nil)
	resp := httptest.NewRecorder()
	if _, err := srv.AgentNodeMaintenance(resp, req); err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Code != 200 {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	// Ensure the maintenance check was registered
	check, ok := srv.agent.state.Checks()[structs.NodeMaint]
	if !ok {
		t.Fatalf("should have registered maintenance check")
	}

	// Check that the token was used
	if token := srv.agent.state.CheckToken(structs.NodeMaint); token != "mytoken" {
		t.Fatalf("expected 'mytoken', got '%s'", token)
	}

	// Ensure the reason was set in notes
	if check.Notes != "broken" {
		t.Fatalf("bad: %#v", check)
	}
}

func TestHTTPAgent_DisableNodeMaintenance(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	// Force the node into maintenance mode
	srv.agent.EnableNodeMaintenance("", "")

	// Leave maintenance mode
	req, _ := http.NewRequest("PUT", "/v1/agent/self/maintenance?enable=false", nil)
	resp := httptest.NewRecorder()
	if _, err := srv.AgentNodeMaintenance(resp, req); err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Code != 200 {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	// Ensure the maintenance check was removed
	if _, ok := srv.agent.state.Checks()[structs.NodeMaint]; ok {
		t.Fatalf("should have removed maintenance check")
	}
}

func TestHTTPAgentRegisterServiceCheck(t *testing.T) {
	dir, srv := makeHTTPServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	// First register the service
	req, err := http.NewRequest("GET", "/v1/agent/service/register", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	args := &ServiceDefinition{
		Name: "memcache",
		Port: 8000,
		Check: CheckType{
			TTL: 15 * time.Second,
		},
	}
	req.Body = encodeReq(args)

	if _, err := srv.AgentRegisterService(nil, req); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Now register an additional check
	req, err = http.NewRequest("GET", "/v1/agent/check/register", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	checkArgs := &CheckDefinition{
		Name:      "memcache_check2",
		ServiceID: "memcache",
		CheckType: CheckType{
			TTL: 15 * time.Second,
		},
	}
	req.Body = encodeReq(checkArgs)

	if _, err := srv.AgentRegisterCheck(nil, req); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ensure we have a check mapping
	result := srv.agent.state.Checks()
	if _, ok := result["service:memcache"]; !ok {
		t.Fatalf("missing memcached check")
	}
	if _, ok := result["memcache_check2"]; !ok {
		t.Fatalf("missing memcache_check2 check")
	}

	// Make sure the new check is associated with the service
	if result["memcache_check2"].ServiceID != "memcache" {
		t.Fatalf("bad: %#v", result["memcached_check2"])
	}
}

func TestHTTPAgent_Monitor(t *testing.T) {
	logWriter := logger.NewLogWriter(512)
	expectedLogs := bytes.Buffer{}
	logger := io.MultiWriter(os.Stdout, &expectedLogs, logWriter)

	dir, srv := makeHTTPServerWithConfigLog(t, nil, logger, logWriter)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer srv.agent.Shutdown()

	// Try passing an invalid log level
	req, _ := http.NewRequest("GET", "/v1/agent/monitor?loglevel=invalid", nil)
	resp := newClosableRecorder()
	if _, err := srv.AgentMonitor(resp, req); err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Code != 400 {
		t.Fatalf("bad: %v", resp.Code)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Unknown log level") {
		t.Fatalf("bad: %s", body)
	}

	// Begin streaming logs from the monitor endpoint
	req, _ = http.NewRequest("GET", "/v1/agent/monitor?loglevel=debug", nil)
	resp = newClosableRecorder()
	go func() {
		if _, err := srv.AgentMonitor(resp, req); err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	// Write the incoming logs from http to a channel for comparison
	logCh := make(chan string, 5)

	// Block until the first log entry from http
	testutil.WaitForResult(func() (bool, error) {
		line, err := resp.Body.ReadString('\n')
		if err != nil && err != io.EOF {
			return false, fmt.Errorf("err: %v", err)
		}
		if line == "" {
			return false, fmt.Errorf("blank line")
		}
		logCh <- line
		return true, nil
	}, func(err error) {
		t.Fatal(err)
	})

	go func() {
		for {
			line, err := resp.Body.ReadString('\n')
			if err != nil && err != io.EOF {
				t.Fatalf("err: %v", err)
			}
			if line != "" {
				logCh <- line
			}
		}
	}()

	// Verify that the first 5 logs we get match the expected stream
	for i := 0; i < 5; i++ {
		select {
		case log := <-logCh:
			expected, err := expectedLogs.ReadString('\n')
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if log != expected {
				t.Fatalf("bad: %q %q", expected, log)
			}
		case <-time.After(10 * time.Second):
			t.Fatalf("failed to get log within timeout")
		}
	}
}

type closableRecorder struct {
	*httptest.ResponseRecorder
	closer chan bool
}

func newClosableRecorder() *closableRecorder {
	r := httptest.NewRecorder()
	closer := make(chan bool)
	return &closableRecorder{r, closer}
}

func (r *closableRecorder) CloseNotify() <-chan bool {
	return r.closer
}
