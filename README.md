# nfcUID
nfcUID - terminal program for reading NFC tag UID by NFC reader and write it as keyboard output.

## Tech
Program used Go bindings to the PC/SC API.  
https://github.com/ebfe/scard  
It is works with ACS readers through their PC/SC driver. So driver should be installed first.  
You can find driver in manufacturer site https://www.acs.com.hk/  

## Supported readers
Program tested with readers:
  - ACR122U
  - ACR1281U-C1
  - ACR1252U-M1
## Supported NFC tags
Program tested with tags:
  - Mifare Classic
  - Mifare Ultralight
  - NTAG203
  - NTAG213
  - NTAG216

## Download/Install
1. Intall PC/SC driver for your reader (https://www.acs.com.hk/)
2. Download archive from release page and extract it. No installation needed.
3. Connect reader to PC.
4. In terminal windows launch program binary from extracted archive.
5. Open text field where you want to write tag UID (e.g. Notepad window, Excel cell etc)
6. Touch NFC tag to reader. Tag UID will be writted to current cursor postion in text field.
7. Release tag from reader.

## Build
To build your own binary just clone repo and use golang build command
```
go build
```

## Demo
You can see how it works in this video 
https://www.youtube.com/watch?v=38UCCXbQPL0
