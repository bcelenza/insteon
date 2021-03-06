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

var dimmer insteon.Dimmer

func init() {
	cmd := Commands.Register("dimmer", "<command> <device id>", "Interact with a specific dimmer", dimmerCmd)
	cmd.Register("config", "", "retrieve dimmer configuration information", dimmerConfigCmd)
	cmd.Register("on", "<level>", "turn the dimmer on", dimmerOnCmd)
	cmd.Register("off", "", "turn the dimmer off", switchOffCmd)
	cmd.Register("onfast", "<level>", "turn the dimmer on fast", dimmerOnFastCmd)
	cmd.Register("brighten", "", "brighten the dimmer one step", dimmerBrightenCmd)
	cmd.Register("dim", "", "dim the dimmer one step", dimmerDimCmd)
	cmd.Register("startBrighten", "", "", dimmerStartBrightenCmd)
	cmd.Register("startDim", "", "", dimmerStartDimCmd)
	cmd.Register("stopChange", "", "", dimmerStopChangeCmd)
	cmd.Register("instantChange", "<level>", "instantly set the dimmer to the desired level (0-255)", dimmerInstantChangeCmd)
	cmd.Register("status", "", "get the switch status", switchStatusCmd)
	cmd.Register("setstatus", "<level>", "set the dimmer switch status LED to <level> (0-31)", dimmerSetStatusCmd)
	cmd.Register("onramp", "<level> <ramp>", "turn the dimmer on to the desired level (0-15) at the given ramp rate (0-15)", dimmerOnRampCmd)
	cmd.Register("offramp", "<ramp>", "turn the dimmer off at the givem ramp rate (0-31)", dimmerOffRampCmd)
	cmd.Register("setramp", "<ramp>", "set default ramp rate (0-31)", dimmerSetRampCmd)
	cmd.Register("setlevel", "<level>", "set default on level (0-255)", dimmerSetOnLevelCmd)
}

func dimmerCmd(args []string, next cli.NextFunc) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("device id and action must be specified")
	}

	var addr insteon.Address
	err = addr.UnmarshalText([]byte(args[0]))
	if err != nil {
		return fmt.Errorf("invalid device address: %v", err)
	}

	device, err := devConnect(modem.Network, addr)
	if err == nil {
		var ok bool
		if dimmer, ok = device.(insteon.Dimmer); ok {
			sw = dimmer
			err = next()
		} else {
			err = fmt.Errorf("Device %s is not a dimmer", addr)
		}
	}
	return err
}

func dimmerConfigCmd(args []string, next cli.NextFunc) error {
	config, err := dimmer.DimmerConfig()
	if err == nil {
		fmt.Printf("           X10 Address: %02x.%02x\n", config.HouseCode, config.UnitCode)
		fmt.Printf("          Default Ramp: %d\n", config.Ramp)
		fmt.Printf("         Default Level: %d\n", config.OnLevel)
		fmt.Printf("                   SNR: %d\n", config.SNT)
	}

	flags, err := dimmer.OperatingFlags()
	if err == nil {
		fmt.Printf("          Program Lock: %v\n", flags.ProgramLock())
		fmt.Printf("             LED on Tx: %v\n", flags.TxLED())
		fmt.Printf("            Resume Dim: %v\n", flags.ResumeDim())
		fmt.Printf("                LED On: %v\n", flags.LED())
		fmt.Printf("         Load Sense On: %v\n", flags.LoadSense())
		fmt.Printf("              DB Delta: %v\n", flags.DBDelta())
		fmt.Printf("                   SNR: %v\n", flags.SNR())
		fmt.Printf("          SNR Failures: %v\n", flags.SNRFailCount())
		fmt.Printf("           X10 Enabled: %v\n", flags.X10Enabled())
		fmt.Printf("        Blink on Error: %v\n", flags.ErrorBlink())
		fmt.Printf("Cleanup Report Enabled: %v\n", flags.CleanupReport())
	}
	return err
}

func dimmerOnCmd(args []string, next cli.NextFunc) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("no on level given")
	}
	level, err := strconv.Atoi(args[1])
	if err == nil {
		err = dimmer.OnLevel(level)
	}
	return err
}

func dimmerOnFastCmd(args []string, next cli.NextFunc) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("no on level given")
	}
	level, err := strconv.Atoi(args[1])
	if err == nil {
		err = dimmer.OnFast(level)
	}
	return err
}

func dimmerBrightenCmd(args []string, next cli.NextFunc) (err error) {
	return dimmer.Brighten()
}

func dimmerDimCmd(args []string, next cli.NextFunc) (err error) {
	return dimmer.Dim()
}

func dimmerStartBrightenCmd(args []string, next cli.NextFunc) (err error) {
	return dimmer.StartBrighten()
}

func dimmerStartDimCmd(args []string, next cli.NextFunc) (err error) {
	return dimmer.StartDim()
}

func dimmerStopChangeCmd(args []string, next cli.NextFunc) (err error) {
	return dimmer.StopChange()
}

func dimmerInstantChangeCmd(args []string, next cli.NextFunc) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("no on level given")
	}
	level, err := strconv.Atoi(args[1])
	if err == nil {
		err = dimmer.InstantChange(level)
	}
	return err
}

func dimmerSetStatusCmd(args []string, next cli.NextFunc) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("no on level given")
	}
	level, err := strconv.Atoi(args[1])
	if err == nil {
		err = dimmer.SetStatus(level)
	}
	return err
}

func dimmerOnRampCmd(args []string, next cli.NextFunc) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("no on level given")
	}

	if len(args) < 3 {
		return fmt.Errorf("no ramp rate given")
	}
	level, err := strconv.Atoi(args[1])
	if err == nil {
		var ramp int
		ramp, err = strconv.Atoi(args[2])
		if err == nil {
			err = dimmer.OnAtRamp(level, ramp)
		}
	}
	return err
}

func dimmerOffRampCmd(args []string, next cli.NextFunc) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("no ramp rate given")
	}
	ramp, err := strconv.Atoi(args[1])
	if err == nil {
		err = dimmer.OffAtRamp(ramp)
	}
	return err
}

func dimmerSetRampCmd(args []string, next cli.NextFunc) error {
	if len(args) < 2 {
		return fmt.Errorf("No ramp rate given")
	}

	ramp, err := strconv.Atoi(args[1])
	if err == nil {
		err = dimmer.SetDefaultRamp(ramp)
	}
	return err
}

func dimmerSetOnLevelCmd(args []string, next cli.NextFunc) error {
	if len(args) < 2 {
		return fmt.Errorf("No on level given")
	}

	level, err := strconv.Atoi(args[1])
	if err == nil {
		err = dimmer.SetDefaultOnLevel(level)
	}
	return err
}
