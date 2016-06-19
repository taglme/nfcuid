package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fuzxxl/nfc"
	"github.com/taglme/string2keyboard"
)

func main() {
	var devices []string
	var err error
	var reader nfc.Device
	var information string
	var nfcTarget *nfc.ISO14443aTarget
	var target []nfc.Target
	var uid string

	res := nfc.Version()
	fmt.Printf("Using libnfc version: %s\n", res)
	devices, err = nfc.ListDevices()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found devices: \n %s", devices[0])
	reader, err = nfc.Open(devices[0])
	if err != nil {
		log.Fatal(err)
	}
	information, err = reader.Information()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nInformation about reader: %s\n", information)
	modul := nfc.Modulation{nfc.ISO14443a, nfc.Nbr106}

	for {

		fmt.Printf("Wait for tag ...\n")
		for {

			target, err = reader.InitiatorListPassiveTargets(modul)
			if err != nil {
				log.Fatal(err)
			}
			if len(target) == 0 {
				time.Sleep(100 * time.Millisecond)
			} else {
				break
			}

		}
		nfcTarget = target[0].(*nfc.ISO14443aTarget)
		uid = string(nfcTarget.UID[:nfcTarget.UIDLen])
		uidS := fmt.Sprintf("%x", uid)
		fmt.Printf(" Tag UID is: %s\n", uidS)
		fmt.Printf("Writting as keyboard input...")
		string2keyboard.KeyboardWrite(uidS)
		fmt.Printf("Done.\n")

	}

	err = reader.Close()
	if err != nil {
		log.Fatal(err)
	}

}
