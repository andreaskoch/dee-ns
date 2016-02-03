// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deens

import (
	"fmt"
	"github.com/pearkes/dnsimple"
)

type testDNSClient struct {
	updateRecordFunc  func(domain string, id string, opts *dnsimple.ChangeRecord) (string, error)
	getRecordsFunc    func(domain string) ([]dnsimple.Record, error)
	getDomainsFunc    func() ([]dnsimple.Domain, error)
	createRecordFunc  func(domain string, opts *dnsimple.ChangeRecord) (string, error)
	destroyRecordFunc func(domain string, id string) error
}

func (client *testDNSClient) UpdateRecord(domain string, id string, opts *dnsimple.ChangeRecord) (string, error) {
	return client.updateRecordFunc(domain, id, opts)
}

func (client *testDNSClient) GetRecords(domain string) ([]dnsimple.Record, error) {
	return client.getRecordsFunc(domain)
}

func (client *testDNSClient) GetDomains() ([]dnsimple.Domain, error) {
	return client.getDomainsFunc()
}

func (client *testDNSClient) CreateRecord(domain string, opts *dnsimple.ChangeRecord) (string, error) {
	return client.createRecordFunc(domain, opts)
}

func (client *testDNSClient) DestroyRecord(domain string, id string) error {
	return client.destroyRecordFunc(domain, id)
}

// testDNSClientFactory creates test DNS clients.
type testDNSClientFactory struct {
	client DNSClient
}

// CreateClient create a new DNSimple client instance.
func (clientFactory testDNSClientFactory) CreateClient() (DNSClient, error) {
	if clientFactory.client == nil {
		return nil, fmt.Errorf("No client available")
	}

	return clientFactory.client, nil
}
