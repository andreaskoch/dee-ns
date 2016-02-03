// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deens

import (
	"testing"
)

type testCredentialsStore struct {
	saveFunc   func(credentials APICredentials) error
	getFunc    func() (APICredentials, error)
	deleteFunc func() error
}

func (credStore testCredentialsStore) SaveCredentials(credentials APICredentials) error {
	return credStore.saveFunc(credentials)
}

func (credStore testCredentialsStore) GetCredentials() (APICredentials, error) {
	return credStore.getFunc()
}

func (credStore testCredentialsStore) DeleteCredentials() error {
	return credStore.deleteFunc()
}

func Test_newAPICredentials_ValidEmailAndToken_NoErrorIsReturned(t *testing.T) {
	// arrange
	var inputs = []struct {
		email string
		token string
	}{
		{"example@example.com", "1234"},
		{"example@example", "a"},
		{"test+test@example.co.uk", "ölö23p4k23lö4köl23k4öä"},
	}

	// act
	for _, input := range inputs {
		_, err := NewAPICredentials(input.email, input.token)

		// assert
		if err != nil {
			t.Fail()
			t.Logf("NewAPICredentials(%q, %q) should not return an error because the input is valid. But it returned: %s", input.email, input.token, err.Error())
		}
	}
}

func Test_newAPICredentials_InvalidValidEmailOrToken_ErrorIsReturned(t *testing.T) {
	// arrange
	var inputs = []struct {
		email string
		token string
	}{
		{"example@example.com", ""},
		{"", "12456"},
		{"", ""},
		{" ", " "},
	}

	// act
	for _, input := range inputs {
		_, err := NewAPICredentials(input.email, input.token)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("NewAPICredentials(%q, %q) should return an error because the given input is invalid.", input.email, input.token)
		}
	}
}
