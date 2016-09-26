package command

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

// KVPutCommand is a Command implementation that is used to write data to the
// key-value store.
type KVPutCommand struct {
	Ui cli.Ui

	// testStdin is the input for testing.
	testStdin io.Reader
}

func (c *KVPutCommand) Help() string {
	helpText := `
Usage: consul kv put [options] KEY [DATA]

  Writes the data to the given path in the key-value store. The data can be of
  any type, but it will be transported as a base64-encoded string for safe
  transport:

      $ consul kv put config/redis/maxconns 5

  The data can also be consumed from a file on disk by prefixing with the "@"
  symbol. For example:

      $ consul kv put config/program/license @license.lic

  The DATA argument itself is optional. If omitted, this will create an empty
  key-value pair at the specified path:

      $ consul kv put webapp/beta/active

  To perform a Check-And-Set operation, specify the -cas flag with the
  appropriate -modify-index flag corresponding to the key you want to perform
  the CAS operation on:

      $ consul kv put -cas -modify-index=844 config/redis/maxconns 5

  Additional flags and more advanced use cases are detailed below.

` + apiOptsText + `

KV Put Options:

  -acquire                Obtain a lock on the key. If the key does not exist,
                          this operation will create the key and obtain the
                          lock. The session must already exist and be specified
                          via the -session flag. The default value is false.

  -cas                    Perform a Check-And-Set operation. If this value is
                          specified without -modify-index, the key will first be
                          fetched and the resulting ModifyIndex will be used on
                          the next query. The default value is false.

  -flags=<int>            Unsigned integer value to assign to this key-value
                          pair. This value is not read by Consul, so clients can
                          use this value however makes sense for their use case.
                          The default value is 0 (no flags).

  -modify-index=<int>     Unsigned integer representing the ModifyIndex of the
                          key. This is often combined with the -cas flag, but it
                          can be specified for any key. The default value is 0.

  -release                Forfeit the lock on the key at the givne path. This
                          requires the -session flag to be set. The key must be
                          held by the session in order to be unlocked. The
                          default value is false.

  -session=<string>       User-defined identifer for this session as a string.
                          This is commonly used with the -acquire and -release
                          operations to build robust locking, but it can be set
                          on any key. The default value is empty (no session).
`
	return strings.TrimSpace(helpText)
}

func (c *KVPutCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("get", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	httpAddr := HTTPAddrFlag(cmdFlags)
	datacenter := cmdFlags.String("datacenter", "", "")
	token := cmdFlags.String("token", "", "")
	stale := cmdFlags.Bool("stale", false, "")
	cas := cmdFlags.Bool("cas", false, "")
	flags := cmdFlags.Uint64("flags", 0, "")
	modifyIndex := cmdFlags.Uint64("modify-index", 0, "")
	session := cmdFlags.String("session", "", "")
	acquire := cmdFlags.Bool("acquire", false, "")
	release := cmdFlags.Bool("release", false, "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	// Check for arg validation
	args = cmdFlags.Args()
	key, data, err := c.dataFromArgs(args)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error! %s", err))
		return 1
	}

	if (*release || *acquire) && *session == "" {
		c.Ui.Error("Error! Missing -session (required with -acquire and -release)")
		return 1
	}

	// Create and test the HTTP client
	conf := api.DefaultConfig()
	conf.Address = *httpAddr
	conf.Token = *token
	client, err := api.NewClient(conf)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error connecting to Consul agent: %s", err))
		return 1
	}

	pair := &api.KVPair{
		Key:         key,
		ModifyIndex: *modifyIndex,
		Flags:       *flags,
		Value:       []byte(data),
		Session:     *session,
	}

	wo := &api.WriteOptions{
		Datacenter: *datacenter,
		Token:      *token,
	}

	switch {
	case *cas:
		// If the user did not supply a -modify-index, but wants a check-and-set,
		// grab the current modify index and store that on the key.
		if pair.ModifyIndex == 0 {
			currentPair, _, err := client.KV().Get(key, &api.QueryOptions{
				Datacenter: *datacenter,
				Token:      *token,
				AllowStale: *stale,
			})
			if err != nil {
				c.Ui.Error(fmt.Sprintf("Error! Could not get current key: %s", err))
				return 1
			}
			if currentPair != nil {
				pair.ModifyIndex = currentPair.ModifyIndex
			}
		}

		ok, _, err := client.KV().CAS(pair, wo)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error! Did not write to %s: %s", key, err))
			return 1
		}
		if !ok {
			c.Ui.Error(fmt.Sprintf("Error! Did not write to %s: CAS failed", key))
			return 1
		}

		c.Ui.Info(fmt.Sprintf("Success! Data written to: %s", key))
		return 0
	case *acquire:
		ok, _, err := client.KV().Acquire(pair, wo)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error! Failed writing data: %s", err))
			return 1
		}
		if !ok {
			c.Ui.Error("Error! Did not acquire lock")
			return 1
		}

		c.Ui.Info(fmt.Sprintf("Success! Lock acquired on: %s", key))
		return 0
	case *release:
		ok, _, err := client.KV().Release(pair, wo)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error! Failed writing data: %s", key))
			return 1
		}
		if !ok {
			c.Ui.Error("Error! Did not release lock")
			return 1
		}

		c.Ui.Info(fmt.Sprintf("Success! Lock released on: %s", key))
		return 0
	default:
		if _, err := client.KV().Put(pair, wo); err != nil {
			c.Ui.Error(fmt.Sprintf("Error! Failed writing data: %s", err))
			return 1
		}

		c.Ui.Info(fmt.Sprintf("Success! Data written to: %s", key))
		return 0
	}
}

func (c *KVPutCommand) Synopsis() string {
	return "Sets or updates data in the KV store"
}

func (c *KVPutCommand) dataFromArgs(args []string) (string, string, error) {
	var stdin io.Reader = os.Stdin
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	switch len(args) {
	case 0:
		return "", "", fmt.Errorf("Missing KEY argument")
	case 1:
		return args[0], "", nil
	case 2:
	default:
		return "", "", fmt.Errorf("Too many arguments (expected 1 or 2, got %d)", len(args))
	}

	key := args[0]
	data := args[1]

	switch data[0] {
	case '@':
		data, err := ioutil.ReadFile(data[1:])
		if err != nil {
			return "", "", fmt.Errorf("Failed to read file: %s", err)
		}
		return key, string(data), nil
	case '-':
		var b bytes.Buffer
		if _, err := io.Copy(&b, stdin); err != nil {
			return "", "", fmt.Errorf("Failed to read stdin: %s", err)
		}
		return key, b.String(), nil
	default:
		return key, data, nil
	}
}
