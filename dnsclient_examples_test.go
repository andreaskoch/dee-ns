// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deens

import (
	"fmt"
	"os"
)

// Create a new DNS client instance and fetch all domain names.
func ExampleNewDNSClient() {
	// assemble the API credentials
	credentials := APICredentials{"john.doe@example.com", "ApItOken"}

	// create a new DNS client
	dnsClient, clientError := NewDNSClient(credentials)
	if clientError != nil {
		fmt.Fprintf(os.Stderr, "Unable to create DNS client: %s", clientError.Error())
		os.Exit(1)
	}

	// fetch all available domains
	domains, domainsError := dnsClient.GetDomains()
	if domainsError != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch domains: %s", domainsError.Error())
		os.Exit(1)
	}

	// print domain names to stdout
	for _, domain := range domains {
		fmt.Fprintf(os.Stdout, "%s\n", domain.Name)
	}
}
