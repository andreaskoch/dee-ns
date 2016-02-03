// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deens

import (
	"fmt"
	"os"
)

// Create a new DNS info provider.
func ExampleNewDNSInfoProvider() {

	// create a DNS client
	credentials := APICredentials{"john.doe@example.com", "ApItOken"}
	dnsClient, clientError := NewDNSClient(credentials)
	if clientError != nil {
		fmt.Fprintf(os.Stderr, "Unable to create DNS client: %s", clientError.Error())
		os.Exit(1)
	}

	// create a new DNS info provider instance
	dnsInfoProvider := NewDNSInfoProvider(dnsClient)

	// get all domain names
	domainNames, domainNamesError := dnsInfoProvider.GetDomainNames()
	if domainNamesError != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch domain names: %s", domainNamesError.Error())
		os.Exit(1)
	}

	// print a list all domain names
	for _, domainName := range domainNames {
		fmt.Fprintf(os.Stdout, "%s\n", domainName)
	}
}
