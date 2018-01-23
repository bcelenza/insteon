package insteon

import (
	"fmt"
)

// I2CsDevice will correctly communicate with Insteon version 2 CS
// (checksum) devices.  The primary difference between I2Device and
// I2CsDevice is that I2CsDevice sets the message version to `2` which
// forces the message marshaler to compute the message checksum. Version
// 2 devices also requre that the EnterLinkingMode command is sent
// as an Extended Direct message (as opposed to standard direct) also
// forcing a checksum computation
type I2CsDevice struct {
	*I2Device
}

// NewI2CsDevice will initialize a new I2CsDevice object and make
// it ready for use
func NewI2CsDevice(address Address, connection Connection) *I2CsDevice {
	return &I2CsDevice{NewI2Device(address, NewI2CsConnection(connection))}
}

// EnterLinkingMode will put the device into linking mode. This is
// equivalent to holding down the set button until the device
// beeps and the indicator light starts flashing
func (i2cs *I2CsDevice) EnterLinkingMode(group Group) (err error) {
	_, err = SendExtendedCommand(i2cs.Connection(), CmdEnterLinkingModeExt.SubCommand(int(group)), NewBufPayload(14))
	return err
}

func (i2cs *I2CsDevice) String() string {
	return fmt.Sprintf("I2CS Device (%s)", i2cs.Address())
}
