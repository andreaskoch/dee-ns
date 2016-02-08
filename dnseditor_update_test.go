// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deens

import (
	"fmt"
	"github.com/pearkes/dnsimple"
	"net"
	"testing"
)

// testDNSUpdater updates DNSimple domain records.
type testDNSUpdater struct {
	updateSubdomainFunc func(domain, subdomain string, ip net.IP) error
}

func (editor *testDNSUpdater) UpdateSubdomain(domain, subdomain string, ip net.IP) error {
	return editor.updateSubdomainFunc(domain, subdomain, ip)
}

// If any of the given parameters is invalid UpdateSubdomain should respond with an error.
func Test_UpdateSubdomain_ParametersInvalid_ErrorIsReturned(t *testing.T) {
	// arrange
	inputs := []struct {
		domain    string
		subdomain string
		ip        net.IP
	}{
		{"example.com", " - ", net.ParseIP("::1")},
		{"", "", net.ParseIP("::1")},
		{" ", " ", net.ParseIP("::1")},
		{"example.com", "www", nil},
	}
	editor := DNSEditor{}

	for _, input := range inputs {

		// act
		err := editor.UpdateSubdomain(input.domain, input.subdomain, input.ip)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("UpdateSubdomain(%q, %q, %q) should return an error.", input.domain, input.subdomain, input.ip)
		}
	}
}

// UpdateSubdomain should return an error if the given subdomain does not exist.
func Test_UpdateSubdomain_ValidParameters_SubdomainNotFound_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ip := net.ParseIP("::1")

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (record dnsimple.Record, err error) {
			return dnsimple.Record{}, fmt.Errorf("")
		},
	}

	editor := DNSEditor{
		infoProvider: infoProvider,
	}

	// act
	err := editor.UpdateSubdomain(domain, subdomain, ip)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("UpdateSubdomain(%q, %q, %q) should return an error if the subdomain does not exist.", domain, subdomain, ip)
	}
}

func Test_UpdateSubdomain_ValidParameters_SubdomainExists_DNSRecordUpdateFails_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ip := net.ParseIP("::1")

	dnsClient := &testDNSClient{
		updateRecordFunc: func(domain string, id string, opts *dnsimple.ChangeRecord) (string, error) {
			return "", fmt.Errorf("Record update failed")
		},
	}

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (record dnsimple.Record, err error) {
			return dnsimple.Record{}, nil
		},
	}

	editor := DNSEditor{
		client:       dnsClient,
		infoProvider: infoProvider,
	}

	// act
	err := editor.UpdateSubdomain(domain, subdomain, ip)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("UpdateSubdomain(%q, %q, %q) should return an error of the record update failed at the DNS client.", domain, subdomain, ip)
	}
}

func Test_UpdateSubdomain_ValidParameters_SubdomainExists_DNSRecordUpdateSucceeds_NoErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ip := net.ParseIP("::1")

	dnsClient := &testDNSClient{
		updateRecordFunc: func(domain string, id string, opts *dnsimple.ChangeRecord) (string, error) {
			return "", nil
		},
	}

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (record dnsimple.Record, err error) {
			return dnsimple.Record{}, nil
		},
	}

	editor := DNSEditor{
		client:       dnsClient,
		infoProvider: infoProvider,
	}

	// act
	err := editor.UpdateSubdomain(domain, subdomain, ip)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("UpdateSubdomain(%q, %q, %q) should not return an error if the DNS record update succeeds.", domain, subdomain, ip)
	}
}

// If the update will not change the IP the update is aborted and an error is returned.
func Test_UpdateSubdomain_ValidParameters_SubdomainExists_ExistingIPIsTheSame_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ip := net.ParseIP("::1")

	dnsClient := &testDNSClient{
		updateRecordFunc: func(domain string, id string, opts *dnsimple.ChangeRecord) (string, error) {
			return "", nil
		},
	}

	existingRecord := dnsimple.Record{
		Name:       "example.com",
		Content:    "::1",
		RecordType: "AAAA",
		Ttl:        600,
	}

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (record dnsimple.Record, err error) {
			return existingRecord, nil
		},
	}

	editor := DNSEditor{
		client:       dnsClient,
		infoProvider: infoProvider,
	}

	// act
	err := editor.UpdateSubdomain(domain, subdomain, ip)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("UpdateSubdomain(%q, %q, %q) should return an error because the IP of the existing record is the same as in the update.", domain, subdomain, ip)
	}
}

func Test_UpdateSubdomain_ValidParameters_SubdomainExists_OnlyTheIPIsChangedOnTheDNSRecord(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ip := net.ParseIP("::2")

	existingRecord := dnsimple.Record{
		Name:       "example.com",
		Content:    "::1",
		RecordType: "AAAA",
		Ttl:        600,
	}

	dnsClient := &testDNSClient{
		updateRecordFunc: func(domain string, id string, opts *dnsimple.ChangeRecord) (string, error) {

			// assert
			if opts.Name != existingRecord.Name {
				t.Fail()
				t.Logf("The DNS name should not change during an update (Old: %q, New: %q)", existingRecord.Name, opts.Name)
			}

			if opts.Type != existingRecord.RecordType {
				t.Fail()
				t.Logf("The DNS record type should not change during an update (Old: %q, New: %q)", existingRecord.RecordType, opts.Type)
			}

			if opts.Ttl != fmt.Sprintf("%d", existingRecord.Ttl) {
				t.Fail()
				t.Logf("The DNS record TTL should not change during an update (Old: %q, New: %q)", existingRecord.Ttl, opts.Ttl)
			}

			if opts.Value != ip.String() {
				t.Fail()
				t.Logf("The DNS record value should have changed to %q", ip.String())
			}

			return "", nil
		},
	}

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (record dnsimple.Record, err error) {
			return existingRecord, nil
		},
	}

	editor := DNSEditor{
		client:       dnsClient,
		infoProvider: infoProvider,
	}

	// act
	editor.UpdateSubdomain(domain, subdomain, ip)
}
