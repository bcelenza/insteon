// Copyright 2018 Andrew Bates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package insteon

import (
	"encoding/json"
	"fmt"
)

// Address is a 3 byte insteon address
type Address [3]byte

// String will format the Address object into a form
// common to Insteon devices: 00.00.00 where each byte
// is represented in hexadecimal form (e.g. 01.b4.a5) the
// string will always be 8 characters long, bytes are zero
// padded
func (a Address) String() string { return sprintf("%02x.%02x.%02x", a[0], a[1], a[2]) }

// UnmarshalText converts a human readable string into an
// Insteon address. If the address cannot be parsed then
// UnmarshalText returns an ErrAddressFormat error
func (a *Address) UnmarshalText(text []byte) error {
	var b1, b2, b3 byte
	_, err := fmt.Sscanf(string(text), "%2x.%2x.%2x", &b1, &b2, &b3)
	if err == nil {
		a[0] = b1
		a[1] = b2
		a[2] = b3
	} else {
		err = ErrAddrFormat
	}
	return err
}

// MarshalJSON will convert the address to a JSON string
func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON will populate the address from the input JSON string
func (a *Address) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err == nil {
		var n int
		n, err = fmt.Sscanf(s, "%02x.%02x.%02x", &a[0], &a[1], &a[2])
		if n < 3 {
			err = fmt.Errorf("Expected Scanf to parse 3 digits, got %d", n)
		}
	}
	return err
}
