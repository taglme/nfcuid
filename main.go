package main

import (
	"errors"
	"flag"
)

func main() {
	var appFlags Flags
	var endChar, inChar string
	var ok bool
	//Read application flags
	flag.StringVar(&endChar, "end-char", "none", "Character at the end of UID. Options: "+CharFlagOptions())
	flag.StringVar(&inChar, "in-char", "none", "Сharacter between bytes of UID. Options: "+CharFlagOptions())
	flag.BoolVar(&appFlags.CapsLock, "caps-lock", false, "UID with Caps Lock")
	flag.BoolVar(&appFlags.Reverse, "reverse", false, "UID reverse order")
	flag.BoolVar(&appFlags.Decimal, "decimal", false, "UID in decimal format")
	flag.IntVar(&appFlags.Device, "device", 0, "Device number to use")
	flag.Parse()

	//Check flags
	appFlags.EndChar, ok = StringToCharFlag(endChar)
	if !ok {
		errorExit(errors.New("Unknown end character flag. Run with '-h' flag to check options"))
		return
	}
	appFlags.InChar, ok = StringToCharFlag(inChar)
	if !ok {
		errorExit(errors.New("Unknown in character flag. Run with '-h' flag to check options"))
		return
	}

	service := NewService(appFlags)
	service.Start()

}
