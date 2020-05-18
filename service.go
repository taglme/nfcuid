package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/ebfe/scard"
	"github.com/taglme/string2keyboard"
)

type Service interface {
	Start()
	Flags() Flags
}

func NewService(flags Flags) Service {
	return &service{flags}
}

type Flags struct {
	CapsLock bool
	Reverse  bool
	Decimal  bool
	EndChar  CharFlag
	InChar   CharFlag
	Device   int
}

type service struct {
	flags Flags
}

func (s *service) Start() {
	//Establish a context
	ctx, err := scard.EstablishContext()
	if err != nil {
		errorExit(err)
	}
	defer ctx.Release()

	//List available readers
	readers, err := ctx.ListReaders()
	if err != nil {
		errorExit(err)
	}

	if len(readers) < 1 {
		errorExit(errors.New("Devices not found. Try to plug-in new device and restart"))
	}

	fmt.Printf("Found %d device:\n", len(readers))
	for i, reader := range readers {
		fmt.Printf("[%d] %s\n", i+1, reader)
	}

	if s.flags.Device == 0 {
		//Device should be selected by user input
		for {
			fmt.Print("Enter device number to start: ")
			inputReader := bufio.NewReader(os.Stdin)
			deviceStr, _ := inputReader.ReadString('\n')

			if runtime.GOOS == "windows" {
				deviceStr = strings.Replace(deviceStr, "\r\n", "", -1)
			} else {
				deviceStr = strings.Replace(deviceStr, "\n", "", -1)
			}
			deviceInt, err := strconv.Atoi(deviceStr)
			if err != nil {
				fmt.Println("Please input integer value")
				continue
			}
			if deviceInt < 0 {
				fmt.Println("Please input positive integer value")
				continue
			}
			if deviceInt > len(readers) {
				fmt.Printf("Value should be less than or equal to %d\n", len(readers))
				continue
			}
			s.flags.Device = deviceInt
			break
		}
	} else if s.flags.Device < 0 {
		errorExit(errors.New("Device flag should positive integer"))
		return
	} else if s.flags.Device > len(readers) {
		errorExit(errors.New("Device flag should not exceed the number of available devices"))
		return
	}

	fmt.Println("Selected device:")
	fmt.Printf("[%d] %s\n", s.flags.Device, readers[s.flags.Device-1])
	selectedReaders := []string{readers[s.flags.Device-1]}

	for {
		fmt.Println("Waiting for a Card")
		index, err := waitUntilCardPresent(ctx, selectedReaders)
		if err != nil {
			errorExit(err)
		}

		//Connect to card
		fmt.Println("Connecting to card...")
		card, err := ctx.Connect(selectedReaders[index], scard.ShareShared, scard.ProtocolAny)
		if err != nil {
			errorExit(err)
		}
		defer card.Disconnect(scard.ResetCard)

		//GET DATA command
		var cmd = []byte{0xFF, 0xCA, 0x00, 0x00, 0x00}

		rsp, err := card.Transmit(cmd)
		if err != nil {
			errorExit(err)
		}

		if len(rsp) < 2 {
			fmt.Println("Not enough bytes in answer. Try again")
			card.Disconnect(scard.ResetCard)
			continue
		}

		//Check response code - two last bytes of response
		rspCodeBytes := rsp[len(rsp)-2 : len(rsp)]
		successResponseCode := []byte{0x90, 0x00}
		if !bytes.Equal(rspCodeBytes, successResponseCode) {
			fmt.Printf("Operation failed to complete. Error code % x\n", rspCodeBytes)
			card.Disconnect(scard.ResetCard)
			continue
		}

		uidBytes := rsp[0 : len(rsp)-2]
		fmt.Printf("UID is: % x\n", uidBytes)
		fmt.Printf("Writting as keyboard input...")
		err = string2keyboard.KeyboardWrite(s.formatOutput(uidBytes))
		if err != nil {
			fmt.Printf("Could write as keyboard output. Error: %s\n", err.Error())
		} else {
			fmt.Printf("Success!\n")
		}

		card.Disconnect(scard.ResetCard)

		//Wait while card will be released
		fmt.Print("Waiting for card release...")
		err = waitUntilCardRelease(ctx, selectedReaders, index)
		fmt.Println("Card released")

	}

}

func (s *service) Flags() Flags {
	return s.flags
}

func errorExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func (s *service) formatOutput(rx []byte) string {
	var output string
	//Reverse UID in flag set
	if s.flags.Reverse {
		for i, j := 0, len(rx)-1; i < j; i, j = i+1, j-1 {
			rx[i], rx[j] = rx[j], rx[i]
		}
	}

	for i, rxByte := range rx {
		var byteStr string
		if s.flags.Decimal {
			byteStr = fmt.Sprintf("%03d", rxByte)
		} else {
			if s.flags.CapsLock {
				byteStr = fmt.Sprintf("%02X", rxByte)
			} else {
				byteStr = fmt.Sprintf("%02x", rxByte)

			}

		}
		output = output + byteStr
		if i < len(rx)-1 {
			output = output + s.flags.InChar.Output()
		}

	}

	output = output + s.flags.EndChar.Output()
	return output
}

func waitUntilCardPresent(ctx *scard.Context, readers []string) (int, error) {
	rs := make([]scard.ReaderState, len(readers))
	for i := range rs {
		rs[i].Reader = readers[i]
		rs[i].CurrentState = scard.StateUnaware
	}

	for {
		for i := range rs {
			if rs[i].EventState&scard.StatePresent != 0 {
				return i, nil
			}
			rs[i].CurrentState = rs[i].EventState
		}
		err := ctx.GetStatusChange(rs, -1)
		if err != nil {
			return -1, err
		}
	}
}

func waitUntilCardRelease(ctx *scard.Context, readers []string, index int) error {
	rs := make([]scard.ReaderState, 1)

	rs[0].Reader = readers[index]
	rs[0].CurrentState = scard.StatePresent

	for {

		if rs[0].EventState&scard.StateEmpty != 0 {
			return nil
		}
		rs[0].CurrentState = rs[0].EventState

		err := ctx.GetStatusChange(rs, -1)
		if err != nil {
			return err
		}
	}
}
