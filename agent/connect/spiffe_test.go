package connect

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testSpiffeIDCases contains the test cases for parsing and encoding
// the SPIFFE IDs. This is a global since it is used in multiple test functions.
var testSpiffeIDCases = []struct {
	Name       string
	URI        string
	Struct     interface{}
	ParseError string
}{
	{
		"invalid scheme",
		"http://google.com/",
		nil,
		"scheme",
	},

	{
		"basic service ID",
		"spiffe://1234.consul/ns/default/dc/dc01/svc/web",
		&SpiffeIDService{
			Host:       "1234.consul",
			Namespace:  "default",
			Datacenter: "dc01",
			Service:    "web",
		},
		"",
	},
}

func TestParseSpiffeID(t *testing.T) {
	for _, tc := range testSpiffeIDCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert := assert.New(t)

			// Parse the URI, should always be valid
			uri, err := url.Parse(tc.URI)
			assert.Nil(err)

			// Parse the ID and check the error/return value
			actual, err := ParseSpiffeID(uri)
			assert.Equal(tc.ParseError != "", err != nil, "error value")
			if err != nil {
				assert.Contains(err.Error(), tc.ParseError)
				return
			}
			assert.Equal(tc.Struct, actual)
		})
	}
}
