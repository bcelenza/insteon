package insteon

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrBufferTooShort       = errors.New("Buffer is too short")
	ErrReadTimeout          = errors.New("Read Timeout")
	ErrWriteTimeout         = errors.New("Write Timeout")
	ErrAckTimeout           = errors.New("Timeout waiting for Device ACK")
	ErrNotImplemented       = errors.New("Command is not yet implemented")
	ErrUnexpectedResponse   = errors.New("Unexpected response from device")
	ErrNotLinked            = errors.New("Not in All-Link group")
	ErrNoLoadDetected       = errors.New("No load detected")
	ErrUnknownCommand       = errors.New("Unknown command")
	ErrNak                  = errors.New("NAK received")
	ErrUnknownEngineVersion = errors.New("Unknown Insteon Version number")
	ErrUnknown              = errors.New("Device returned unknown error")
	ErrIllegalValue         = errors.New("Illegal value in command")
	ErrIncorrectChecksum    = errors.New("I2CS invalid checksum")
	ErrPreNak               = errors.New("Database search took too long")
	ErrNotSupported         = errors.New("Action/command is not supported on this device")
	ErrAddrFormat           = errors.New("address format is xx.xx.xx (digits in hex)")
)

var sprintf = fmt.Sprintf

type Version int

type ProductKey [3]byte

func (p ProductKey) String() string {
	return sprintf("0x%02x%02x%02x", p[0], p[1], p[2])
}

type DevCat [2]byte

func (dc DevCat) Category() Category {
	return Category(dc[0])
}

func (dc DevCat) SubCategory() SubCategory {
	return SubCategory(dc[1])
}

type Category byte

type SubCategory byte

func (dc DevCat) String() string {
	return sprintf("%02x.%02x", dc[0], dc[1])
}

func (dc DevCat) MarshalJSON() ([]byte, error) {
	return json.Marshal(sprintf("%02x.%02x", dc[0], dc[1]))
}

func (dc *DevCat) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err == nil {
		var n int
		n, err = fmt.Sscanf(s, "%02x.%02x", &dc[0], &dc[1])
		if n < 2 {
			err = fmt.Errorf("Expected Scanf to parse 2 digits, got %d", n)
		}
	}
	return err
}

type ProductData struct {
	Key    ProductKey
	DevCat DevCat
}

func (pd *ProductData) String() string {
	return sprintf("DevCat:%s Product Key:%s", pd.DevCat, pd.Key)
}

func (pd *ProductData) UnmarshalBinary(buf []byte) error {
	if len(buf) < 14 {
		return newBufError(ErrBufferTooShort, 14, len(buf))
	}

	copy(pd.Key[:], buf[1:4])
	copy(pd.DevCat[:], buf[4:6])
	return nil
}

func (pd *ProductData) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 7)
	copy(buf[1:4], pd.Key[:])
	copy(buf[4:6], pd.DevCat[:])
	buf[6] = 0xff
	return buf, nil
}
