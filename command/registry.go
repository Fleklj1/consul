package command

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mitchellh/cli"
)

// Factory is a function that returns a new instance of a CLI-sub command.
type Factory func(cli.Ui) (cli.Command, error)

// Register adds a new CLI sub-command to the registry.
func Register(name string, fn Factory) {
	if registry == nil {
		registry = make(map[string]Factory)
	}

	if registry[name] != nil {
		panic(fmt.Errorf("Command %q is already registered", name))
	}
	registry[name] = fn
}

// RegisterHidden adds a new CLI sub-command to the registry that won't show up
// in help or autocomplete.
func RegisterHidden(name string, fn Factory) {
	if hiddenRegistry == nil {
		hiddenRegistry = make(map[string]Factory)
	}

	if hiddenRegistry[name] != nil {
		panic(fmt.Errorf("Command %q is already registered", name))
	}
	hiddenRegistry[name] = fn
}

// Map returns a realized mapping of available CLI commands in a format that
// the CLI class can consume. This should be called after all registration is
// complete.
func Map(ui cli.Ui) map[string]cli.CommandFactory {
	return makeCommands(ui, registry)
}

// Map returns a realized mapping of available but hidden CLI commands in a
// format that the CLI class can consume. This should be called after all
// registration is complete.
func MapHidden(ui cli.Ui) map[string]cli.CommandFactory {
	return makeCommands(ui, hiddenRegistry)
}

func makeCommands(ui cli.Ui, reg map[string]Factory) map[string]cli.CommandFactory {
	m := make(map[string]cli.CommandFactory)
	for name, fn := range reg {
		thisFn := fn
		m[name] = func() (cli.Command, error) {
			return thisFn(ui)
		}
	}
	return m
}

// registry has an entry for each available CLI sub-command, indexed by sub
// command name. This should be populated at package init() time via Register().
var registry map[string]Factory

// hiddenRegistry behaves identically to registry but is for commands that are
// hidden - i.e. not publically documented in the help or autocomplete.
var hiddenRegistry map[string]Factory

// MakeShutdownCh returns a channel that can be used for shutdown notifications
// for commands. This channel will send a message for every interrupt or SIGTERM
// received.
func MakeShutdownCh() <-chan struct{} {
	resultCh := make(chan struct{})
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()

	return resultCh
}
