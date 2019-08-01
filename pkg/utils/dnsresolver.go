// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

package utils

import (
	"net"
)

// DNSResolverIface is an interface for the DNSResolver
type DNSResolverIface interface {
	Resolve(entry string) ([]string, error)
}

// DNSResolver is the DNS resolver
type DNSResolver struct{}

// NewDNSResolver instanciate a new DNSResolver
func NewDNSResolver() *DNSResolver {
	return &DNSResolver{}
}

// Resolve takes a domain name and returns the
func (r *DNSResolver) Resolve(entry string) ([]string, error) {
	addrs, err := net.LookupHost(entry)
	return addrs, err
}
