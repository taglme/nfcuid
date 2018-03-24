package main

import (
	"fmt"
	"os"

	"github.com/ebfe/scard"
	"github.com/taglme/string2keyboard"
)

func errorExit(err error) {
	fmt.Println(err)
	os.Exit(1)
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

func main() {

	// Establish a context
	ctx, err := scard.EstablishContext()
	if err != nil {
		errorExit(err)
	}
	defer ctx.Release()

	// List available readers
	readers, err := ctx.ListReaders()
	if err != nil {
		errorExit(err)
	}

	fmt.Printf("Found %d readers:\n", len(readers))
	for i, reader := range readers {
		fmt.Printf("[%d] %s\n", i, reader)
	}

	if len(readers) > 0 {
		for {
			fmt.Println("Waiting for a Card")
			index, err := waitUntilCardPresent(ctx, readers)
			if err != nil {
				errorExit(err)
			}

			// Connect to card
			fmt.Println("Connecting to card in ", readers[index])
			card, err := ctx.Connect(readers[index], scard.ShareExclusive, scard.ProtocolAny)
			if err != nil {
				errorExit(err)
			}
			defer card.Disconnect(scard.ResetCard)

			var cmd = []byte{0xFF, 0xCA, 0x00, 0x00, 0x00}

			rsp, err := card.Transmit(cmd)
			if err != nil {
				errorExit(err)
			}
			uid := string(rsp[0:7])
			uidS := fmt.Sprintf("%x", uid)
			fmt.Printf("Tag UID is: %s\n", uidS)
			fmt.Printf("Writting as keyboard input...")
			string2keyboard.KeyboardWrite(uidS)
			fmt.Printf("Done.\n")

			card.Disconnect(scard.ResetCard)

			//Wait while card will be released
			fmt.Print("Waiting for card release...")
			err = waitUntilCardRelease(ctx, readers, index)
			fmt.Println("Card released.")

		}

	}

}
