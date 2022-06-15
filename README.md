# go-astm ![build status](https://travis-ci.org/78bit/uuid.svg?branch=master)

Golang library for handling ASTM lis2a2 Procotol

###### Install
`go get github.com/DRK-Blutspende-BaWueHe/go-astm`

## Features
  - Encoding 
    - UTF8 
    - ASCII
    - Windows1250 
    - Windows1251 
    - Windows1252 
    - DOS852 
    - DOS855 
    - DOS866 
  - Timezone Support
  - Marshal/Unmarshal function

## Installation

Install the package with the following command.

``` shell
go get github.com/DRK-Blutspende-BaWueHe/go-astm/...
```
## Quick Start

The following Go code decodes a ASTM read from a File.

``` go
fileData, err := ioutil.ReadFile("protocoltest/becom/5.2/bloodtype.astm")
if err != nil {
  log.Fatal(err)		
}

message, err := astm1384.Unmarshal(fileData,
 astm1384.EncodingWindows1252, 
 astm1384.TimezoneEuropeBerlin, 
 astm1384.LIS2A2)

if err != nil {
   log.Fatal(err)		
}
```
