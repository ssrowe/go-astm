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
  - Timezone conversion on marshal and unmarshal
  - Marshalling and Unmarshalling supported
  - Custom delimiters are recognized in the Header and appplied (defaults are \^&)
  - Supported Types : string, float32, float64, time.Time, enums, int

## Installation

Install the package with the following command.

``` shell
go get github.com/DRK-Blutspende-BaWueHe/go-astm/...
```

## Messaage Structure

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
## Message Structure and Annotation

### Optional records
``` go
type SimpleMessage struct  {
	Header       standardlis2a2.Header       `astm:"H"`
	Manufacturer standardlis2a2.Manufacturer `astm:"M,optional"`
	Patient      standardlis2a2.Patient      `astm:"P"`
	Order        standardlis2a2.Order        `astm:"O"`
	Result       standardlis2a2.Result       `astm:"R"`
	Terminator   standardlis2a2.Terminator   `astm:"L"`
}
```

### Nested arrays
``` go
type MessagePORC struct {
	Header       standardlis2a2.Header       `astm:"H"`
	Manufacturer standardlis2a2.Manufacturer `astm:"M,optional"`
	OrderResults []struct {
		Patient         standardlis2a2.Patient `astm:"P"`
		Order           standardlis2a2.Order   `astm:"O"`
		CommentedResult []struct {
			Result  standardlis2a2.Result    `astm:"R"`
			Comment []standardlis2a2.Comment `astm:"C,optional"`
		}
	}
	Terminator standardlis2a2.Terminator `astm:"L"`
}
```

## Custom Record Structure

### Annotations
Often the default is not enough. You can customize any record with annotation. 

Fields are adressed by their Field # :
``` go
   ...
   Filed string `astm:"3"`  // Select 3rd field, start counting with 1
   ...
```

by Field.Component#
``` go
   ...
   Filed string `astm:"3.2"`  // Select 3rd field, 2nd component, start counting with 1
   ...
```

by Field.Repeat.Component#
``` go
   ...
   Filed string `astm:"3.2.2"`  // Select 3rd field, 2nd array index, 2nd component, start counting with 1
   ...
```

### Example Custom structure
``` go
type Result struct {
	SequenceNumber                           int       `astm:"2,sequence"`   // sequence generates numbers when value is 0 
	UniversalTestID                          string    `astm:"3.1"`         
	UniversalTestIDName                      string    `astm:"3.2"`         
	UniversalTestIDType                      string    `astm:"3.3"`         
	ManufacturersTestType                    string    `astm:"3.4"`         
	ManufacturersTestName                    string    `astm:"3.5"`         
	ManufacturersTestCode                    string    `astm:"3.6"`         
	TestCode                                 string    `astm:"3.7"`         
	DataMeasurementValue                     string    `astm:"4.1"`         
	InitialMeasurementValue                  string    `astm:"4.2"`         
	MeasurementValueOfDevice                 string    `astm:"4.3"`         
	Units                                    string    `astm:"5"`           
	ReferenceRange                           string    `astm:"6"`           
	ResultAbnormalFlag                       string    `astm:"7"`           
	NatureOfAbnormalTesting                  string    `astm:"8"`           
	ResultStatus                             string    `astm:"9"`           
	DateOfChangeInInstrumentNormativeTesting time.Time `astm:"10,longdate"` 
	OperatorIDPerformed                      string    `astm:"11.1"`        
	OperatorIDVerified                       string    `astm:"11.2"`        
	DateTimeTestStarted                      time.Time `astm:"12,longdate"` 
	DateTimeCompleted                        time.Time `astm:"13,longdate"` 
	IntstrumentIdentification                string    `astm:"14"`          
}
```
