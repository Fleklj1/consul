package agent

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/consul/consul/structs"
	"log"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// CheckType is used to create either the CheckMonitor
// or the CheckTTL. Only one of TTL or Script/Interval
// needs to be provided
type CheckType struct {
	Script   string
	Interval time.Duration

	TTL time.Duration
}

// Valid checks if the CheckType is valid
func (c *CheckType) Valid() bool {
	return c.IsTTL() || c.IsMonitor()
}

// IsTTL checks if this is a TTL type
func (c *CheckType) IsTTL() bool {
	return c.TTL != 0
}

// IsMonitor checks if this is a Monitor type
func (c *CheckType) IsMonitor() bool {
	return c.Script != "" && c.Interval != 0
}

// CheckNotifier interface is used by the CheckMonitor
// to notify when a check has a status update. The update
// should take care to be idempotent.
type CheckNotifier interface {
	UpdateCheck(checkID, status, note string)
}

// CheckMonitor is used to periodically invoke a script to
// determine the health of a given check. It is compatible with
// nagios plugins and expects the output in the same format.
type CheckMonitor struct {
	Notify   CheckNotifier
	CheckID  string
	Script   string
	Interval time.Duration
	Logger   *log.Logger

	stop     bool
	stopCh   chan struct{}
	stopLock sync.Mutex
}

// Start is used to start a check monitor.
// Monitor runs until stop is called
func (c *CheckMonitor) Start() {
	c.stopLock.Lock()
	defer c.stopLock.Unlock()
	c.stop = false
	c.stopCh = make(chan struct{})
	go c.run()
}

// Stop is used to stop a check monitor.
func (c *CheckMonitor) Stop() {
	c.stopLock.Lock()
	defer c.stopLock.Unlock()
	if !c.stop {
		c.stop = true
		close(c.stopCh)
	}
}

// run is invoked by a goroutine to run until Stop() is called
func (c *CheckMonitor) run() {
	next := time.After(0)
	for {
		select {
		case <-next:
			c.check()
			next = time.After(c.Interval)
		case <-c.stopCh:
			return
		}
	}
}

// check is invoked periodically to perform the script check
func (c *CheckMonitor) check() {
	// Determine the shell invocation based on OS
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}

	// Create the command
	cmd := exec.Command(shell, flag, c.Script)

	// Collect the output
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	// Start the check
	if err := cmd.Start(); err != nil {
		c.Logger.Printf("[ERR] agent: failed to invoke '%s': %s", c.Script, err)
		c.Notify.UpdateCheck(c.CheckID, structs.HealthUnknown, err.Error())
		return
	}

	// Wait for the check to complete
	errCh := make(chan error, 2)
	go func() {
		errCh <- cmd.Wait()
	}()
	go func() {
		time.Sleep(30 * time.Second)
		errCh <- fmt.Errorf("Timed out running check '%s'", c.Script)
	}()
	err := <-errCh

	notes := string(output.Bytes())
	c.Logger.Printf("[DEBUG] agent: check '%s' script '%s' output: %s",
		c.CheckID, c.Script, notes)

	// Check if the check passed
	if err == nil {
		c.Logger.Printf("[DEBUG] Check '%v' is passing", c.CheckID)
		c.Notify.UpdateCheck(c.CheckID, structs.HealthPassing, notes)
		return
	}

	// If the exit code is 1, set check as warning
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			code := status.ExitStatus()
			if code == 1 {
				c.Logger.Printf("[WARN] Check '%v' is now warning", c.CheckID)
				c.Notify.UpdateCheck(c.CheckID, structs.HealthWarning, notes)
				return
			}
		}
	}

	// Set the health as critical
	c.Logger.Printf("[WARN] Check '%v' is now critical", c.CheckID)
	c.Notify.UpdateCheck(c.CheckID, structs.HealthCritical, notes)
}

// CheckTTL is used to apply a TTL to check status,
// and enables clients to set the status of a check
// but upon the TTL expiring, the check status is
// automatically set to critical.
type CheckTTL struct {
	Notify  CheckNotifier
	CheckID string
	TTL     time.Duration
	Logger  *log.Logger

	timer *time.Timer

	stop     bool
	stopCh   chan struct{}
	stopLock sync.Mutex
}

// Start is used to start a check ttl, runs until Stop()
func (c *CheckTTL) Start() {
	c.stopLock.Lock()
	defer c.stopLock.Unlock()
	c.stop = false
	c.stopCh = make(chan struct{})
	c.timer = time.NewTimer(c.TTL)
	go c.run()
}

// Stop is used to stop a check ttl.
func (c *CheckTTL) Stop() {
	c.stopLock.Lock()
	defer c.stopLock.Unlock()
	if !c.stop {
		c.timer.Stop()
		c.stop = true
		close(c.stopCh)
	}
}

// run is used to handle TTL expiration and to update the check status
func (c *CheckTTL) run() {
	for {
		select {
		case <-c.timer.C:
			c.Logger.Printf("[WARN] Check '%v' missed TTL, is now critical",
				c.CheckID)
			c.Notify.UpdateCheck(c.CheckID, structs.HealthCritical, "TTL expired")

		case <-c.stopCh:
			return
		}
	}
}

// SetStatus is used to update the status of the check,
// and to renew the TTL. If expired, TTL is restarted.
func (c *CheckTTL) SetStatus(status, note string) {
	c.Logger.Printf("[DEBUG] Check '%v' status is now %v",
		c.CheckID, status)
	c.Notify.UpdateCheck(c.CheckID, status, note)
	c.timer.Reset(c.TTL)
}
