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

package main

import (
	"fmt"
	"strconv"

	"github.com/abates/cli"
	"github.com/abates/insteon"
)

var sw insteon.Switch

func init() {
	cmd := Commands.Register("switch", "<command> <device id>", "Interact with a specific switch", swCmd)
	cmd.Register("config", "", "retrieve switch configuration information", switchConfigCmd)
	cmd.Register("on", "", "turn the switch/light on", switchOnCmd)
	cmd.Register("off", "", "turn the switch/light off", switchOffCmd)
	cmd.Register("status", "", "get the switch status", switchStatusCmd)
	cmd.Register("setled", "", "set operating flags", switchSetLedCmd)
}

func swCmd(args []string, next cli.NextFunc) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("device id and action must be specified")
	}

	var addr insteon.Address
	err = addr.UnmarshalText([]byte(args[0]))
	if err != nil {
		return fmt.Errorf("invalid device address: %v", err)
	}

	device, err = devConnect(modem.Network, addr)
	if err == nil {
		var ok bool
		if sw, ok = device.(insteon.Switch); ok {
			err = next()
		} else {
			err = fmt.Errorf("Device at %s is a %T not a switch", addr, device)
		}
	}
	return err
}

func switchConfigCmd([]string, cli.NextFunc) error {
	config, err := sw.SwitchConfig()
	if err == nil {
		err = printDevInfo(device, fmt.Sprintf("  X10 Address: %02x.%02x", config.HouseCode, config.UnitCode))
	}
	return err
}

func switchOnCmd([]string, cli.NextFunc) error {
	return sw.On()
}

func switchOffCmd([]string, cli.NextFunc) error {
	return sw.Off()
}

func switchStatusCmd([]string, cli.NextFunc) error {
	level, err := sw.Status()
	if err == nil {
		if level == 0 {
			fmt.Printf("Switch is off\n")
		} else if level == 255 {
			fmt.Printf("Switch is on\n")
		} else {
			fmt.Printf("Switch is on at level %d\n", level)
		}
	}
	return err
}

func switchSetCmd(args []string, next cli.NextFunc) error {
	return next()
}

func switchSetLedCmd(args []string, next cli.NextFunc) error {
	if len(args) < 2 {
		return fmt.Errorf("Expected device address and flag value")
	}
	b, err := strconv.ParseBool(args[1])
	if err == nil {
		err = sw.SetLED(b)
	}
	return err
}
