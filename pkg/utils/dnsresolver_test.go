// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

package utils

import (
	"net"
	"testing"
)

var resolveHostTests = []struct {
	name string
}{
	{"datadoghq.com"},
	{"datadoghq.com."},
}

func TestResolveDDHosts(t *testing.T) {

	resolver := NewDNSResolver()
	for _, tt := range resolveHostTests {
		addrs, err := resolver.Resolve(tt.name)
		if err != nil {
			t.Fatal(err)
		}
		if len(addrs) == 0 {
			t.Error("got no record")
		}
		for _, addr := range addrs {
			if net.ParseIP(addr) == nil {
				t.Errorf("got %q; want a literal IP address", addr)
			}
		}
	}
}

func TestResolveBogusHost(t *testing.T) {
	resolver := NewDNSResolver()
	addrs, err := resolver.Resolve("!!!.###.bogus..domain.")
	if err == nil {
		t.Fatalf("lookup didn't error out: %v", addrs)
	}
}

func TestResolveEmptyHost(t *testing.T) {
	resolver := NewDNSResolver()
	addrs, err := resolver.Resolve("")
	if err == nil {
		t.Fatalf("lookup didn't error out: %v", addrs)
	}
}
