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

import "time"

// I2CsDevice can communicate with Version 2 (checksum) Insteon Engines
type I2CsDevice struct {
	*I2Device
}

// NewI2CsDevice will initialize a new I2CsDevice object and make
// it ready for use
func NewI2CsDevice(address Address, sendCh chan<- *MessageRequest, recvCh <-chan *Message, timeout time.Duration) *I2CsDevice {
	return &I2CsDevice{NewI2Device(address, sendCh, recvCh, timeout)}
}

// EnterLinkingMode will put the device into linking mode. This is
// equivalent to holding down the set button until the device
// beeps and the indicator light starts flashing
func (i2cs *I2CsDevice) EnterLinkingMode(group Group) (err error) {
	return extractError(i2cs.SendCommand(CmdEnterLinkingModeExt.SubCommand(int(group)), make([]byte, 14)))
}

// String returns the string "I2CS Device (<address>)" where <address> is the destination
// address of the device
func (i2cs *I2CsDevice) String() string {
	return sprintf("I2CS Device (%s)", i2cs.Address())
}

// SendCommand will send the given command bytes to the device including
// a payload (for extended messages). If payload length is zero then a standard
// length message is used to deliver the commands. The command bytes from the
// response ack are returned as well as any error
func (i2cs *I2CsDevice) SendCommand(command Command, payload []byte) (response Command, err error) {
	if command[1] == CmdSetOperatingFlags[1] && len(payload) == 0 {
		payload = make([]byte, 14)
	}
	return i2cs.I2Device.SendCommand(command, payload)
}
