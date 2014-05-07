package command

import (
	"fmt"
	"github.com/hashicorp/consul/testutil"
	"github.com/hashicorp/serf/serf"
	"github.com/mitchellh/cli"
	"strings"
	"testing"
	"errors"
)

func TestForceLeaveCommand_implements(t *testing.T) {
	var _ cli.Command = &ForceLeaveCommand{}
}

func TestForceLeaveCommandRun(t *testing.T) {
	a1 := testAgent(t)
	a2 := testAgent(t)
	defer a1.Shutdown()
	defer a2.Shutdown()

	addr := fmt.Sprintf("127.0.0.1:%d", a2.config.Ports.SerfLan)
	_, err := a1.agent.JoinLAN([]string{addr})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Forcibly shutdown a2 so that it appears "failed" in a1
	a2.Shutdown()

	ui := new(cli.MockUi)
	c := &ForceLeaveCommand{Ui: ui}
	args := []string{
		"-rpc-addr=" + a1.addr,
		a2.config.NodeName,
	}

	code := c.Run(args)
	if code != 0 {
		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
	}

	m := a1.agent.LANMembers()
	if len(m) != 2 {
		t.Fatalf("should have 2 members: %#v", m)
	}

	testutil.WaitForResult(func() (bool, error) {
		m = a1.agent.LANMembers()
		success := m[1].Status == serf.StatusLeft
		err := errors.New(m[1].Status.String())
		return success, err
	}, func(err error) {
		t.Fatalf("member status is %v, not left", err)
	})
}

func TestForceLeaveCommandRun_noAddrs(t *testing.T) {
	ui := new(cli.MockUi)
	c := &ForceLeaveCommand{Ui: ui}
	args := []string{"-rpc-addr=foo"}

	code := c.Run(args)
	if code != 1 {
		t.Fatalf("bad: %d", code)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "node name") {
		t.Fatalf("bad: %#v", ui.ErrorWriter.String())
	}
}
