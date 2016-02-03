// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deens

import (
	"fmt"
	"net"
	"os"
)

// Create a new DNS editor instance.
func ExampleNewDNSEditor() {

	// create a DNS client
	credentials := APICredentials{"john.doe@example.com", "ApItOken"}
	dnsClient, clientError := NewDNSClient(credentials)
	if clientError != nil {
		fmt.Fprintf(os.Stderr, "Unable to create DNS client: %s", clientError.Error())
		os.Exit(1)
	}

	// create a new DNS info provider instance
	dnsInfoProvider := NewDNSInfoProvider(dnsClient)

	// create a new DNS editor instance
	dnsEditor := NewDNSEditor(dnsClient, dnsInfoProvider)

	// create an DNS A record for www.example.com pointing to 127.0.0.1
	domain := "example.com"
	subDomainName := "www"
	timeToLive := 600
	ip := net.ParseIP("127.0.0.1")

	createSubdomainError := dnsEditor.CreateSubdomain(domain, subDomainName, timeToLive, ip)
	if createSubdomainError != nil {
		fmt.Fprintf(os.Stderr, "Failed to create subdomain %s.%s: %s", subDomainName, domain, clientError.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Created subdomain %s.%s", subDomainName, domain)
}
