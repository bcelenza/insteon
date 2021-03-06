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
	"bytes"
	"testing"
)

func TestRecordControlFlags(t *testing.T) {
	tests := []struct {
		input              byte
		expectedInUse      bool
		expectedController bool
		expectedString     string
	}{
		{0x40, false, true, "AC"},
		{0x00, false, false, "AR"},
		{0xc0, true, true, "UC"},
		{0x80, true, false, "UR"},
	}

	for i, test := range tests {
		flags := RecordControlFlags(test.input)
		if flags.InUse() != test.expectedInUse {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedInUse, flags.InUse())
		}

		if flags.Available() == test.expectedInUse {
			t.Errorf("tests[%d] expected %v got %v", i, !test.expectedInUse, flags.Available())
		}

		if flags.Controller() != test.expectedController {
			t.Errorf("tests[%d] expected %v got %v", i, !test.expectedController, flags.Controller())
		}

		if flags.Responder() == test.expectedController {
			t.Errorf("tests[%d] expected %v got %v", i, !test.expectedController, flags.Responder())
		}

		if flags.String() != test.expectedString {
			t.Errorf("tests[%d] expected %q got %q", i, test.expectedString, flags.String())
		}
	}
}

func TestRecordControlFlagsUnmarshalText(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
		expected    RecordControlFlags
	}{
		{"A", "Expected 2 characters got 1", RecordControlFlags(0x00)},
		{"AR", "", RecordControlFlags(0x00)},
		{"UR", "", RecordControlFlags(0x80)},
		{"AC", "", RecordControlFlags(0x40)},
		{"UC", "", RecordControlFlags(0xc0)},
		{"FR", "Invalid value for Available flag", RecordControlFlags(0x00)},
		{"AZ", "Invalid value for Controller flag", RecordControlFlags(0x00)},
	}

	for i, test := range tests {
		var flags RecordControlFlags
		err := flags.UnmarshalText([]byte(test.input))
		if err == nil {
			if test.expectedErr != "" {
				t.Errorf("tests[%d] expected %q got nil", i, test.expectedErr)
			} else if flags != test.expected {
				t.Errorf("tests[%d] expected 0x%02x got 0x%02x", i, test.expected, flags)
			}
		} else if err.Error() != test.expectedErr {
			t.Errorf("tests[%d] expected %q got %q", i, test.expectedErr, err)
		}
	}
}

func TestSettingRecordControlFlags(t *testing.T) {
	flags := RecordControlFlags(0xff)
	tests := []struct {
		set      func()
		expected byte
	}{
		{flags.setAvailable, 0x7f},
		{flags.setInUse, 0xff},
		{flags.setResponder, 0xbf},
		{flags.setController, 0xff},
	}

	for i, test := range tests {
		test.set()
		if byte(flags) != test.expected {
			t.Errorf("tests[%d] expected 0x%02x got 0x%02x", i, test.expected, byte(flags))
		}
	}
}

func TestLinkEqual(t *testing.T) {
	availableController := RecordControlFlags(0x40)
	availableResponder := RecordControlFlags(0x00)
	usedController := RecordControlFlags(0xc0)
	usedResponder := RecordControlFlags(0x80)

	newLink := func(flags RecordControlFlags, group Group, address Address) *LinkRecord {
		buffer := []byte{byte(flags), byte(group), address[0], address[1], address[2], 0x00, 0x00, 0x00}
		link := &LinkRecord{}
		link.UnmarshalBinary(buffer)
		return link
	}

	l1 := newLink(availableController, Group(0x01), Address{0x01, 0x02, 0x03})
	l2 := l1

	tests := []struct {
		link1    *LinkRecord
		link2    *LinkRecord
		expected bool
	}{
		{newLink(usedController, Group(0x01), Address{0x01, 0x02, 0x03}), newLink(availableController, Group(0x01), Address{0x01, 0x02, 0x03}), false},
		{newLink(usedResponder, Group(0x01), Address{0x01, 0x02, 0x03}), newLink(availableController, Group(0x01), Address{0x01, 0x02, 0x03}), false},
		{newLink(availableController, Group(0x01), Address{0x01, 0x02, 0x03}), newLink(availableController, Group(0x01), Address{0x01, 0x02, 0x03}), true},
		{newLink(availableResponder, Group(0x01), Address{0x01, 0x02, 0x03}), newLink(availableController, Group(0x01), Address{0x01, 0x02, 0x03}), false},
		{newLink(availableController, Group(0x01), Address{0x01, 0x02, 0x03}), newLink(availableController, Group(0x01), Address{0x01, 0x02, 0x04}), false},
		{newLink(availableController, Group(0x01), Address{0x01, 0x02, 0x03}), nil, false},
		{l1, l2, true},
	}

	for i, test := range tests {
		if test.link1.Equal(test.link2) != test.expected {
			t.Errorf("tests[%d] expected %v got %v", i, test.expected, !test.expected)
		}
	}
}

func TestLinkMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		input           []byte
		expectedFlags   RecordControlFlags
		expectedGroup   Group
		expectedAddress Address
		expectedData    [3]byte
		expectedString  string
		expectedError   error
	}{
		{
			input:           []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			expectedFlags:   RecordControlFlags(0x01),
			expectedGroup:   Group(0x02),
			expectedAddress: Address{0x03, 0x04, 0x05},
			expectedData:    [3]byte{0x06, 0x07, 0x08},
			expectedString:  "AR 2 03.04.05 0x06 0x07 0x08",
		},
		{
			input:         []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
			expectedError: ErrBufferTooShort,
		},
	}

	for i, test := range tests {
		link := &LinkRecord{}
		err := link.UnmarshalBinary(test.input)
		if !isError(err, test.expectedError) {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedError, err)
			continue
		} else if err != nil {
			continue
		}

		if link.Flags != test.expectedFlags {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedFlags, link.Flags)
		}

		if link.Group != test.expectedGroup {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedGroup, link.Group)
		}

		if link.Address != test.expectedAddress {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedAddress, link.Address)
		}

		if link.Data != test.expectedData {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedData, link.Data)
		}

		if link.String() != test.expectedString {
			t.Errorf("tests[%d] expected %q got %q", i, test.expectedString, link.String())
		}

		buf, _ := link.MarshalBinary()

		if !bytes.Equal(buf, test.input) {
			t.Errorf("tests[%d] expected %v got %v", i, test.input, buf)
		}
	}
}

func TestLinkRecordMarshalText(t *testing.T) {
	tests := []struct {
		expectedString string
		expected       LinkRecord
		expectedErr    string
	}{
		{"UC        1 01.01.01   00 00 00", LinkRecord{0, RecordControlFlags(0xc0), Group(1), Address{1, 1, 1}, [3]byte{0, 0, 0}}, ""},
		{"UC        1 01.01.01   00 00", LinkRecord{}, "Expected 6 fields got 5"},
	}

	for i, test := range tests {
		if test.expectedErr == "" {
			buf, _ := test.expected.MarshalText()
			if !bytes.Equal([]byte(test.expectedString), buf) {
				t.Errorf("tests[%d] expected %q got %q", i, test.expectedString, string(buf))
			}
		}

		var linkRecord LinkRecord
		err := linkRecord.UnmarshalText([]byte(test.expectedString))
		if err == nil {
			if test.expectedErr != "" {
				t.Errorf("tests[%d] expected %q got nil", i, test.expectedErr)
			} else if test.expected != linkRecord {
				t.Errorf("tests[%d] expected %v got %v", i, test.expected, linkRecord)
			}
		} else if test.expectedErr != err.Error() {
			t.Errorf("tests[%d] expected %q got %q", i, test.expectedErr, err.Error())
		}
	}
}

func TestGroupUnmarshalText(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
		expected    Group
	}{
		{"1", "", Group(1)},
		{"wxyz", "invalid number format", Group(0)},
		{"256", "valid groups are between 1 and 255 (inclusive)", Group(0)},
		{"-1", "valid groups are between 1 and 255 (inclusive)", Group(0)},
	}

	for i, test := range tests {
		var group Group
		err := group.UnmarshalText([]byte(test.input))
		if err == nil {
			if test.expectedErr != "" {
				t.Errorf("tests[%d] expected %q got %q", i, test.expectedErr, err)
			} else if group != test.expected {
				t.Errorf("tests[%d] expected %d got %d", i, test.expected, group)
			}
		} else if test.expectedErr != err.Error() {
			t.Errorf("tests[%d] expected %q got %q", i, test.expectedErr, err.Error())
		}
	}
}
