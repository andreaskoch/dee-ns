// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deens

import (
	"net"
	"testing"
)

func Test_isEmpty_EmptyString_ResultIsTrue(t *testing.T) {
	// arrange
	inputs := []string{
		"",
		" ",
		"    ",
		" ",
		" ",
		" ",
	}

	// act
	for _, input := range inputs {
		result := isEmpty(input)

		// assert
		if result == false {
			t.Fail()
			t.Logf("isEmpty(%q) should return true", input)
		}
	}
}

func Test_isEmpty_NotEmptyString_ResultIsFalse(t *testing.T) {
	// arrange
	inputs := []string{
		"-",
		".",
		" a ",
		" _ ",
	}

	// act
	for _, input := range inputs {
		result := isEmpty(input)

		// assert
		if result == true {
			t.Fail()
			t.Logf("isEmpty(%q) should return false", input)
		}
	}
}

func Test_isValidSubdomain_GivenTextIsValid_ResultIsTrue(t *testing.T) {

	// arrange
	inputs := []string{
		"",
		"www",
		"w-w-w",
		"w.w.w",
		"a",
		"1",
		"123",
		"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk",
		"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk",
	}

	for _, input := range inputs {
		// act
		result := isValidSubdomain(input)

		// assert
		if result == false {
			t.Fail()
			t.Logf("isValidSubdomain(%q) should have returned true", input)
		}
	}
}

func Test_isValidSubdomain_GivenTextIsInvalid_ResultIsFalse(t *testing.T) {

	// arrange
	inputs := []string{
		" www",
		"www ",
		"w ww",
		"-a",
		"-hi-",
		"_hi_",
		"*hi*",
		"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk",
	}

	for _, input := range inputs {
		// act
		result := isValidSubdomain(input)

		// assert
		if result == true {
			t.Fail()
			t.Logf("isValidSubdomain(%q) should have returned false", input)
		}
	}
}

// If the given IP is an IPv4 address, "A" should be returned as the record type.
func Test_getDNSRecordTypeByIP_IPisIPv4_AIsReturned(t *testing.T) {
	// arrange
	ip := net.ParseIP("127.0.0.1")

	// act
	result := getDNSRecordTypeByIP(ip)

	// assert
	if result != "A" {
		t.Fail()
		t.Logf("getDNSRecordTypeByIP(%s) should return %q", ip, "A")
	}

}

func Test_getFormattedDomainName(t *testing.T) {
	// arrange
	inputs := []struct {
		subdomain      string
		domain         string
		expectedResult string
	}{
		{"", "", ""},
		{"www", "", ""},
		{"", "example.com", "example.com"},
		{"www", "example.com", "www.example.com"},
	}

	for _, input := range inputs {

		// act
		result := getFormattedDomainName(input.subdomain, input.domain)

		// assert
		if result != input.expectedResult {
			t.Fail()
			t.Logf("getFormattedDomainName(%q, %q) returned %q but should have returned %q.", input.subdomain, input.domain, result, input.expectedResult)
		}
	}
}
