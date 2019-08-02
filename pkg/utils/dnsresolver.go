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

// NewDNSResolver instantiate a new DNSResolver
func NewDNSResolver() *DNSResolver {
	re := regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]\.?)$`)
	return &DNSResolver{
		re: re,
	}
}

// Resolve takes a domain name and returns the
func (r *DNSResolver) Resolve(entry string) ([]string, error) {
	if !r.isValidEntry(entry) {
		return nil, fmt.Errorf("invalid host: %s", entry)
	}
	addrs, err := net.LookupHost(entry)
	sort.Strings(addrs)
	return addrs, err
}

// isValidEntry verifies if an entry is a valid dns hostname
func (r *DNSResolver) isValidEntry(host string) bool {
	host = strings.Trim(host, " ")
	return r.re.MatchString(host)
}
