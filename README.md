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

## Define your Structure

Every instrument is a little different. Use the lis2a2-default implmenetation provided with this library as a starting point.

``` go
type CommentedResult struct {
	Result  Result    `astm:"R"`
	Comment []Comment `astm:"C,optional"`
}

type PORC struct {
	Patient         Patient   `astm:"P"`
	Comment         []Comment `astm:"C,optional"`
	Order           Order     `astm:"O"`
	CommentedResult []CommentedResult
}

// https://samson-rus.com/wp-content/files/LIS2-A2.pdf Page 30 : Logial Structure of Message
type DefaultMessage struct {
	Header       Header       `astm:"H"`
	Manufacturer Manufacturer `astm:"M,optional"`
	OrderResults []PORC
	Terminator   Terminator `astm:"L"`
}
```

## Reading ASTM

The following Go code decodes a ASTM provided as a string and stores all its information in the &message.

``` go
var message lis2a2.DefaultMessage

err := lis2a2.Unmarshal([]byte(textdata), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)
if err != nil {
  log.Fatal(err)		
}
```

## Writing ASM

Converting an annotated Structure (see above) to an enocded bytestream. 

The bytestream is encoded by-row, lacking the CR code at the end. 

``` go
lines, err := lis2a2.Marshal(msg, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)

// output on screen
for _, line := range lines {
		linestr := string(line)
		fmt.Println(linestr)
}
```
