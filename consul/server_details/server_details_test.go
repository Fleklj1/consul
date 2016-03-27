package server_details_test

import (
	"net"
	"testing"

	"github.com/hashicorp/consul/consul/server_details"
	"github.com/hashicorp/serf/serf"
)

func TestServerDetails_Key_Equal(t *testing.T) {
	tests := []struct {
		name  string
		k1    *server_details.Key
		k2    *server_details.Key
		equal bool
	}{
		{
			name: "IPv4 equality",
			k1: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "127.0.0.1",
			},
			k2: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "127.0.0.1",
			},
			equal: true,
		},
		{
			name: "IPv6 equality",
			k1: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "fc00::1",
			},
			k2: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "fc00::1",
			},
			equal: true,
		},
		{
			name: "IPv4 Inequality",
			k1: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "127.0.0.1",
			},
			k2: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "1.2.3.4",
			},
			equal: false,
		},
		{
			name: "IPv4 Inequality AddrString",
			k1: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "127.0.0.1ip+net",
			},
			k2: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "127.0.0.1ip",
			},
			equal: false,
		},
		{
			name: "IPv6 Inequality",
			k1: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "fc00::1",
			},
			k2: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "2001:0db8:85a3::8a2e:0370:7334",
			},
			equal: false,
		},
		{
			name: "Port Inequality",
			k1: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "127.0.0.1",
			},
			k2: &server_details.Key{
				Datacenter: "dc1",
				Port:       8500,
				AddrString: "1.2.3.4",
			},
			equal: false,
		},
		{
			name: "DC Inequality",
			k1: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "127.0.0.1",
			},
			k2: &server_details.Key{
				Datacenter: "dc2",
				Port:       8300,
				AddrString: "127.0.0.1",
			},
			equal: false,
		},
	}

	for _, test := range tests {
		if test.k1.Equal(test.k2) != test.equal {
			t.Errorf("Expected a %v result from test %s", test.equal, test.name)
		}

		// Test Key to make sure it actually works as a key
		m := make(map[server_details.Key]bool)
		m[*test.k1] = true
		if _, found := m[*test.k2]; found != test.equal {
			t.Errorf("Expected a %v result from map test %s", test.equal, test.name)
		}
	}
}

func TestServerDetails_Key(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")

	tests := []struct {
		name  string
		sd    *server_details.ServerDetails
		k     *server_details.Key
		equal bool
	}{
		{
			name: "Key equality",
			sd: &server_details.ServerDetails{
				Datacenter: "dc1",
				Port:       8300,
				Addr:       &net.IPAddr{IP: ip},
			},
			k: &server_details.Key{
				Datacenter: "dc1",
				Port:       8300,
				AddrString: "127.0.0.1ip",
			},
			equal: true,
		},
		{
			name: "Key inequality",
			sd: &server_details.ServerDetails{
				Datacenter: "dc1",
				Port:       8300,
				Addr:       &net.IPAddr{IP: ip},
			},
			k: &server_details.Key{
				Datacenter: "dc2",
				Port:       8300,
				AddrString: "127.0.0.1ip",
			},
			equal: false,
		},
	}

	for _, test := range tests {
		if test.k.Equal(test.sd.Key()) != test.equal {
			t.Errorf("Expected a %v result from test %s", test.equal, test.name)
		}
	}
}

func TestServerDetails_Key_params(t *testing.T) {
	ipv4a := net.ParseIP("127.0.0.1")
	ipv4b := net.ParseIP("1.2.3.4")

	tests := []struct {
		name  string
		sd1   *server_details.ServerDetails
		sd2   *server_details.ServerDetails
		equal bool
	}{
		{
			name: "Key equality",
			sd1: &server_details.ServerDetails{
				Name:       "s1",
				Datacenter: "dc1",
				Port:       8300,
				Addr:       &net.IPAddr{IP: ipv4a},
			},
			sd2: &server_details.ServerDetails{
				Name:       "s1",
				Datacenter: "dc1",
				Port:       8300,
				Addr:       &net.IPAddr{IP: ipv4a},
			},
			equal: true,
		},
		{
			name: "Key equality",
			sd1: &server_details.ServerDetails{
				Name:       "s1",
				Datacenter: "dc1",
				Port:       8300,
				Addr:       &net.IPAddr{IP: ipv4a},
			},
			sd2: &server_details.ServerDetails{
				Name:       "s2",
				Datacenter: "dc1",
				Port:       8300,
				Addr:       &net.IPAddr{IP: ipv4a},
			},
			equal: true,
		},
		{
			name: "Addr inequality",
			sd1: &server_details.ServerDetails{
				Name:       "s1",
				Datacenter: "dc1",
				Port:       8300,
				Addr:       &net.IPAddr{IP: ipv4a},
			},
			sd2: &server_details.ServerDetails{
				Name:       "s1",
				Datacenter: "dc1",
				Port:       8300,
				Addr:       &net.IPAddr{IP: ipv4b},
			},
			equal: false,
		},
	}

	for _, test := range tests {
		if test.sd1.Key().Equal(test.sd2.Key()) != test.equal {
			t.Errorf("Expected a %v result from test %s", test.equal, test.name)
		}

		// Test Key to make sure it actually works as a key
		m := make(map[server_details.Key]bool)
		m[*test.sd1.Key()] = true
		if _, found := m[*test.sd2.Key()]; found != test.equal {
			t.Errorf("Expected a %v result from map test %s", test.equal, test.name)
		}
	}
}

func TestIsConsulServer(t *testing.T) {
	m := serf.Member{
		Name: "foo",
		Addr: net.IP([]byte{127, 0, 0, 1}),
		Tags: map[string]string{
			"role": "consul",
			"dc":   "east-aws",
			"port": "10000",
			"vsn":  "1",
		},
	}
	ok, parts := server_details.IsConsulServer(m)
	if !ok || parts.Datacenter != "east-aws" || parts.Port != 10000 {
		t.Fatalf("bad: %v %v", ok, parts)
	}
	if parts.Name != "foo" {
		t.Fatalf("bad: %v", parts)
	}
	if parts.Bootstrap {
		t.Fatalf("unexpected bootstrap")
	}
	if parts.Expect != 0 {
		t.Fatalf("bad: %v", parts.Expect)
	}
	m.Tags["bootstrap"] = "1"
	m.Tags["disabled"] = "1"
	ok, parts = server_details.IsConsulServer(m)
	if !ok {
		t.Fatalf("expected a valid consul server")
	}
	if !parts.Bootstrap {
		t.Fatalf("expected bootstrap")
	}
	if parts.Addr.String() != "127.0.0.1:10000" {
		t.Fatalf("bad addr: %v", parts.Addr)
	}
	if parts.Version != 1 {
		t.Fatalf("bad: %v", parts)
	}
	m.Tags["expect"] = "3"
	delete(m.Tags, "bootstrap")
	delete(m.Tags, "disabled")
	ok, parts = server_details.IsConsulServer(m)
	if !ok || parts.Expect != 3 {
		t.Fatalf("bad: %v", parts.Expect)
	}
	if parts.Bootstrap {
		t.Fatalf("unexpected bootstrap")
	}

	delete(m.Tags, "role")
	ok, parts = server_details.IsConsulServer(m)
	if ok {
		t.Fatalf("unexpected ok server")
	}
}
