package imp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/command/flags"
	"github.com/hashicorp/consul/command/kv/impexp"
	"github.com/mitchellh/cli"
)

func New(ui cli.Ui) *cmd {
	c := &cmd{UI: ui}
	c.init()
	return c
}

type cmd struct {
	UI     cli.Ui
	flags  *flag.FlagSet
	http   *flags.HTTPFlags
	help   string
	prefix string

	// testStdin is the input for testing.
	testStdin io.Reader
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.StringVar(&c.prefix, "prefix", "", "Key prefix for imported data")
	c.http = &flags.HTTPFlags{}
	flags.Merge(c.flags, c.http.ClientFlags())
	flags.Merge(c.flags, c.http.ServerFlags())
	flags.Merge(c.flags, c.http.NamespaceFlags())
	c.help = flags.Usage(help, c.flags)
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}

	// Check for arg validation
	args = c.flags.Args()
	data, err := c.dataFromArgs(args)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error! %s", err))
		return 1
	}

	// Create and test the HTTP client
	client, err := c.http.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error connecting to Consul agent: %s", err))
		return 1
	}

	var entries []*impexp.Entry
	if err := json.Unmarshal([]byte(data), &entries); err != nil {
		c.UI.Error(fmt.Sprintf("Cannot unmarshal data: %s", err))
		return 1
	}

	for _, entry := range entries {
		value, err := base64.StdEncoding.DecodeString(entry.Value)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error base 64 decoding value for key %s: %s", entry.Key, err))
			return 1
		}

		pair := &api.KVPair{
			Key:   filepath.Join(c.prefix, entry.Key),
			Flags: entry.Flags,
			Value: value,
		}

		w := api.WriteOptions{Namespace: entry.Namespace}
		if _, err := client.KV().Put(pair, &w); err != nil {
			c.UI.Error(fmt.Sprintf("Error! Failed writing data for key %s: %s", pair.Key, err))
			return 1
		}

		c.UI.Info(fmt.Sprintf("Imported: %s", pair.Key))
	}

	return 0
}

func (c *cmd) dataFromArgs(args []string) (string, error) {
	var stdin io.Reader = os.Stdin
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	switch len(args) {
	case 0:
		return "", errors.New("Missing DATA argument")
	case 1:
	default:
		return "", fmt.Errorf("Too many arguments (expected 1, got %d)", len(args))
	}

	data := args[0]

	if len(data) == 0 {
		return "", errors.New("Empty DATA argument")
	}

	switch data[0] {
	case '@':
		data, err := ioutil.ReadFile(data[1:])
		if err != nil {
			return "", fmt.Errorf("Failed to read file: %s", err)
		}
		return string(data), nil
	case '-':
		if len(data) > 1 {
			return data, nil
		}
		var b bytes.Buffer
		if _, err := io.Copy(&b, stdin); err != nil {
			return "", fmt.Errorf("Failed to read stdin: %s", err)
		}
		return b.String(), nil
	default:
		return data, nil
	}
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return c.help
}

const synopsis = "Imports a tree stored as JSON to the KV store"
const help = `
Usage: consul kv import [DATA]

  Imports key-value pairs to the key-value store from the JSON representation
  generated by the "consul kv export" command.

  The data can be read from a file by prefixing the filename with the "@"
  symbol. For example:

      $ consul kv import @filename.json

  Or it can be read from stdin using the "-" symbol:

      $ cat filename.json | consul kv import -

  Alternatively the data may be provided as the final parameter to the command,
  though care must be taken with regards to shell escaping.

  For a full list of options and examples, please see the Consul documentation.
`
