# nfcUID
Cross-platform terminal application for reading NFC tag UID and write it as keyboard output to text field in any application.

## Overview
Application read NFC tag UID using PC/SC API.
PC/SC is a standard to interface computers with smartcards, available on most operating systems, including Windows, MacOS, Linux.
UID writed in active text input field by generating keystrokes on the keyboard.
So cursor should be in some text input field.
UID output format options are available.

## Supported readers
Application should work with any PC/SC compatible reader. It is tested with ACS readers:
  - ACR122U
  - ACR1281U-C1
  - ACR1252U-M1

## Supported NFC tags
Application should work with any NFC tag with UID. It is tested with NXP tags:
  - Mifare Classic
  - Mifare Ultralight
  - NTAG203
  - NTAG213
  - NTAG216

## Download
Binaries for Windows, MacOS and Linux platforms could be downloaded from release page.

## Build
Install the source code and build
```golang
go get github.com/taglme/nfcuid
go build
```

## Output format

There are options for application that should be specified as flags.

 - '-device' (integer) - device number to use. Set to 0 if your want to select it manually or set to the desired device number to auto-select.
 - '-caps-lock' (boolean) -  UID output with caps lock
 - '-decimal' (boolean) - UID output in decimal format
 - '-reverse' (boolean) - UID output in reverse order
 - '-end-char' (string) - character at the end of UID. Options: 'hyphen', 'enter', 'semicolon', 'colon', 'comma', 'none', 'space', 'tab',  (default "none")
 - '-in-char' (string) - character between bytes of UID. Options: 'hyphen', 'enter', 'semicolon', 'colon', 'comma', 'none', 'space', 'tab',  (default "none")

Run with '-h' flag to check usage
```
nfcuid -h
```

## Examples

```golang
//This will auto-select first available PC/SC device in system
//Output will be in direct order 
//Bytes of UID will be in hex format
//Between bytes of UID will be hyphen ("-") printed
//At end of UID newline will be printed
nfcuid -end-char=enter -in-char=hyphen -caps-lock=false -reverse=false -decimal=false -device=1

//Output 
04-ae-65-ca-82-49-80

```

## Demo
Short application [overview video](https://youtu.be/BSyxhg2iWfA) (in Russian) on Youtube

## More NFC software
More NFC software and services are available on [Tagl.me](https://tagl.me)
