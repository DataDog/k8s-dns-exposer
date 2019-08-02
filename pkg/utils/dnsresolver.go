// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

package utils

import (
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"
)

// DNSResolverIface is an interface for the DNSResolver
type DNSResolverIface interface {
	Resolve(entry string) ([]string, error)
}

// DNSResolver is the DNS resolver
type DNSResolver struct {
	re *regexp.Regexp
}

// NewDNSResolver instanciate a new DNSResolver
func NewDNSResolver() *DNSResolver {
	re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]\.?)$`)
	return &DNSResolver{
		re: re,
	}
}

// Resolve takes a domain name and returns the
func (r *DNSResolver) Resolve(entry string) ([]string, error) {
	if !r.isValidEntry(entry) {
		return nil, fmt.Errorf("Invalid host: %s", entry)
	}
	addrs, err := net.LookupHost(entry)
	sort.Strings(addrs)
	return addrs, err
}

func (r *DNSResolver) isValidEntry(host string) bool {
	host = strings.Trim(host, " ")
	return r.re.MatchString(host)
}
