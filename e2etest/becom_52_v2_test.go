package e2e

import (
	"fmt"
	"github.com/DRK-Blutspende-BaWueHe/go-astm/astm1384"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

// (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
type Header struct {
	Delimiters              string    `astm:"2.1.1,delimiter"` // 6.2
	MessageControlID        string    `astm:"3"`               // 6.3
	AccessPassword          string    `astm:"4"`               // 6.4
	SenderNameOrID          string    `astm:"5"`               // 6.5
	SenderStreetAddress     string    `astm:"6"`               // 6.6
	Field6                  string    `astm:"7"`               // 6.7
	SenderTelephone         string    `astm:"8"`               // 6.8
	CharacteristicsOfSender string    `astm:"9"`               // 6.9
	ReceiverID              string    `astm:"10"`              // 6.10
	Comment                 string    `astm:"11"`              // 6.11
	ProcessingID            string    `astm:"12"`              // 6.12
	Version                 string    `astm:"13"`              // 6.13
	DateAndTime             time.Time `astm:"14,longdate"`     // 6.14
}

// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type Patient struct {
	SequenceNumber                     int       `astm:"2"`   // 7.2 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	PracticeAssignedPatientID          string    `astm:"3"`   // 7.3
	LabAssignedPatientID               string    `astm:"4"`   // 7.4
	ID3                                string    `astm:"5"`   // 7.5
	LastName                           string    `astm:"6,1"` // 7.6.1
	FirstName                          string    `astm:"6,2"` // 7.6.2
	MothersMaidenName                  string    `astm:"7"`   // 7.7
	DOB                                time.Time `astm:"8"`   // 7.8
	Gender                             string    `astm:"9"`   // 7.9
	Race                               string    `astm:"10"`  // 7.10
	Address                            string    `astm:"11"`  // 7.11
	F12                                string    `astm:"12"`  // 7.12
	Telephone                          string    `astm:"13"`  // 7.13
	AttendingPhysicianID               string    `astm:"14"`  // 7.14
	SpecialField1                      string    `astm:"15"`  // 7.15
	SpecialField2                      string    `astm:"16"`  // 7.16
	Height                             string    `astm:"17"`  // 7.17
	Weight                             string    `astm:"18"`  // 7.18
	SuspectedDiagnosis                 string    `astm:"19"`  // 7.19
	ActiveMedication                   string    `astm:"20"`  // 7.20
	Diet                               string    `astm:"21"`  // 7.21
	PracticeField1                     string    `astm:"22"`  // 7.22
	PracticeField2                     string    `astm:"23"`  // 7.23
	AdmissionAndDischargeDates         string    `astm:"24"`  // 7.24
	AdmissionStatus                    string    `astm:"25"`  // 7.25
	Location                           string    `astm:"26"`  // 7.26
	NatureOfAlternativeDiagnosticCodes string    `astm:"27"`  // 7.27
	AlternativeDiagnosticCodes         string    `astm:"28"`  // 7.28
	Religion                           string    `astm:"29"`  // 7.29
	MaritalStatus                      string    `astm:"30"`  // 7.30
	IsolationStatus                    string    `astm:"31"`  // 7.31
	Language                           string    `astm:"32"`  // 7.32
	HospitalService                    string    `astm:"33"`  // 7.33
	HospitalInstitution                string    `astm:"34"`  // 7.34
	DosageCategory                     string    `astm:"35"`  // 7.35
}

// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type Order struct {
	SequenceNumber                      int       `astm:"2"`     // 8.4.2 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	SpecimenID                          string    `astm:"3"`     // 8.4.3
	CodeOfSpecimen1                     string    `astm:"4,1:1"` // 8.4.4
	TypeOfSpecimen1                     string    `astm:"4,1:2"`
	CodeOfDonorSample1                  string    `astm:"4,1:3"`
	TypeOfDonorSample1                  string    `astm:"4,1:4"`
	CodeOfSpecimen2                     string    `astm:"4,2:1"`
	TypeOfSpecimen2                     string    `astm:"4,2:2"`
	CodeOfDonorSample2                  string    `astm:"4,2:3"`
	TypeOfDonorSample2                  string    `astm:"4,2:4"`
	UniversalTestID_LOINC               string    `astm:"5,1"`         // 8.4.5.1
	UniversalTestID_TestName            string    `astm:"5,2"`         // 8.4.5.2
	UniversalTestID_AlternativeTestName string    `astm:"5,3"`         // 8.4.5.2
	UniversalTestID_ManufacturerCode    string    `astm:"5,4"`         // 8.4.5.3
	UniversalTestID_Custom1             string    `astm:"5,5"`         // 8.4.5.4
	UniversalTestID_Custom2             string    `astm:"5,6"`         // 8.4.5.4
	UniversalTestID_Custom3             string    `astm:"5,7"`         // 8.4.5.4
	UniversalTestID_Custom4             string    `astm:"5,8"`         // 8.4.5.4
	UniversalTestID_Custom5             string    `astm:"5,9"`         // 8.4.5.4
	UniversalTestID_Custom6             string    `astm:"5,10"`        // 8.4.5.4
	UniversalTestID_Custom7             string    `astm:"5,11"`        // 8.4.5.4
	UniversalTestID_Custom8             string    `astm:"5,12"`        // 8.4.5.4
	UniversalTestID_Custom9             string    `astm:"5,13"`        // 8.4.5.4
	Priority                            string    `astm:"6"`           // 8.4.6
	RequestedOrderDateTime              time.Time `astm:"7,longdate"`  // 8.4.7
	SpecimenCollectionDateTime          time.Time `astm:"8,longdate"`  // 8.4.8
	CollectionEndTime                   time.Time `astm:"9,longdate"`  // 8.4.9
	CollectionVolume                    string    `astm:"10"`          // 8.4.10
	CollectionID                        string    `astm:"11"`          // 8.4.11
	ActionCode                          string    `astm:"12"`          // 8.4.12
	DangerCode                          string    `astm:"13"`          // 8.4.13
	RelevantClinicalInformation         string    `astm:"14"`          // 8.4.14
	DateTimeSpecimenReceived            string    `astm:"15"`          // 8.4.15
	SpecimenDescriptor                  string    `astm:"16"`          // 8.4.16
	OrderingPhysician                   string    `astm:"17"`          // 8.4.17
	PhysicianTelephone                  string    `astm:"18"`          // 8.4.18
	UserField1                          string    `astm:"19"`          // 8.4.19
	UserField2                          string    `astm:"20"`          // 8.4.20
	LaboratoryField1                    string    `astm:"21"`          // 8.4.21
	LaboratoryField2                    string    `astm:"22"`          // 8.4.22
	DateTimeResultsReported             time.Time `astm:"23,longdate"` // 8.4.23
	InstrumentCharge                    string    `astm:"24"`          // 8.4.24
	InstrumentSectionID                 string    `astm:"25"`          // 8.4.25
	ReportType                          string    `astm:"26"`          // 8.4.26
	Reserved                            string    `astm:"27"`          // 8.4.27
	LocationOfSpecimenCollection        string    `astm:"28"`          // 8.4.28
	NosocomialInfectionFlag             string    `astm:"29"`          // 8.4.29
	SpecimenService                     string    `astm:"30"`          // 8.4.30
	SpecimenInstitution                 string    `astm:"31"`          // 8.4.31
}

// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type Result struct {
	SequenceNumber                           int       `astm:"2"`           // 9.2 https://samson-rus.com/wp-content/files/LIS2-A2.pdf
	UniversalTestID                          string    `astm:"3,1"`         // 9.3
	UniversalTestIDName                      string    `astm:"3,2"`         // 9.3
	UniversalTestIDType                      string    `astm:"3,3"`         // 9.3
	ManufacturersTestType                    string    `astm:"3,4"`         // 9.3
	ManufacturersTestName                    string    `astm:"3,5"`         // 9.3
	ManufacturersTestCode                    string    `astm:"3,6"`         // 9.3
	TestCode                                 string    `astm:"3,7"`         // 9.3
	DataMeasurementValue                     string    `astm:"4,1"`         // 9.4
	InitialMeasurementValue                  string    `astm:"4,2"`         // 9.4
	MeasurementValueOfDevice                 string    `astm:"4,3"`         // 9.4
	Units                                    string    `astm:"5"`           // 9.5
	ReferenceRange                           string    `astm:"6"`           // 9.6
	ResultAbnormalFlag                       string    `astm:"7"`           // 9.7
	NatureOfAbnormalTesting                  string    `astm:"8"`           // 9.8
	ResultStatus                             string    `astm:"9"`           // 9.9
	DateOfChangeInInstrumentNormativeTesting time.Time `astm:"10,longdate"` // 9.10
	OperatorIDPerformed                      string    `astm:"11,1"`        // 9.11
	OperatorIDVerified                       string    `astm:"11,2"`        // 9.11
	OperatorIDLastModified                   string    `astm:"11,3"`        // 9.11
	OperatorIDWellReading                    string    `astm:"11,4"`        // 9.11
	OperatorIDFirstViewing                   string    `astm:"11,5"`        // 9.11
	DateTimeTestStarted                      time.Time `astm:"12,longdate"` // 9.12
	DateTimeCompleted                        time.Time `astm:"13,longdate"` // 9.13
	IntstrumentIdentification                string    `astm:"14"`          // 9.14
	InstrumentName                           string    `astm:"15"`          // out of standard
	InstrumentSerialNumber                   string    `astm:"16"`          // out of standard
	InstrumentOperatorID                     string    `astm:"17"`          // out of standard
}

// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type Comment struct {
	SequenceNumber int `astm:"2"` // 10.2 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	//CommentSource  [][]string `astm:"3"` // 10.3
	DescriptionOfReagent1     string    `astm:"3,1:1"` // 10.3
	BarcodeOfReagent1         string    `astm:"3,1:2"` // 10.3
	LotNumberOfReagent1       string    `astm:"3,1:3"` // 10.3
	ExpirationDateOfReagent1  time.Time `astm:"3,1:4"` // 10.3
	DescriptionOfReagent2     string    `astm:"3,2:1"` // 10.3
	BarcodeOfReagent2         string    `astm:"3,2:2"` // 10.3
	LotNumberOfReagent2       string    `astm:"3,2:3"` // 10.3
	ExpirationDateOfReagent2  time.Time `astm:"3,2:4"` // 10.3
	TypeOfTestMedia           string    `astm:"4,1"`   // 10.3
	PlateID                   string    `astm:"4,2"`   // 10.3
	LotNumber                 string    `astm:"4,3"`   // 10.3
	ExpirationOfCassette      time.Time `astm:"4,4"`   // 10.3
	IdCardPlate               string    `astm:"4,5"`   // 10.3
	Comment                   string    `astm:"5"`     // 10.5
	FileNameOfReactionPicture string    `astm:"6"`     // 10.6 (not standard)
}

// Lis2RequestInformation - RequestInformation
// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type RequestInformation struct {
	SequenceNumber                  int    `astm:"2"`  // 11.2 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	StartingRangeIDNumber           string `astm:"3"`  // 11.3
	EndingRangeIDNumber             string `astm:"4"`  // 11.4
	UniversalTestID                 string `astm:"5"`  // 11.5
	NatureOfRequestTimeLimits       string `astm:"6"`  // 11.6
	BeginningRequestResultsDateTime string `astm:"7"`  // 11.7
	EndingRequestResultsDateTime    string `astm:"8"`  // 11.8
	RequestingPhysicianName         string `astm:"9"`  // 11.9
	RequestingPhysicianTelephone    string `astm:"10"` // 11.10
	UserField1                      string `astm:"11"` // 11.11
	UserField2                      string `astm:"12"` // 11.12
	RequestInformationStatus        string `astm:"13"` // 11.13
}

// Lis2Terminator - Terminator Record (Hasta la vista....)
// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type Terminator struct {
	SequenceNumber int    `astm:"2"` // 12.2 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	TerminatorCode string `astm:"3"` // 12.3
}

// Lis2Manufacturer -Manufacturer Record
// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type Manufacturer struct {
	SequenceNumber string `astm:"2"`  // 14.2 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	F2             string `astm:"3"`  // 14.3
	F3             string `astm:"4"`  // 14.4
	F4             string `astm:"5"`  // 14.5
	F5             string `astm:"6"`  // 14.6
	F6             string `astm:"7"`  // 14.7
	F7             string `astm:"8"`  // 14.8
	F8             string `astm:"9"`  // 14.9
	F9             string `astm:"10"` // 14.10
	F10            string `astm:"11"` // 14.11
	F11            string `astm:"12"` // 14.12
	F12            string `astm:"13"` // 14.13
	F13            string `astm:"14"` // 14.14
}

type Message struct {
	Header       Header       `astm:"H"`
	Manufacturer Manufacturer `astm:"M,optional"`
	Patient      Patient      `astm:"P"`
	Order        Order        `astm:"O"`
	Result       Result       `astm:"R"`
	Terminator   Terminator   `astm:"L"`
}

func TestReadFileBeCom52_v2(t *testing.T) {
	fileData, err := ioutil.ReadFile("../protocoltest/becom/5.2/bloodtype_test.astm")
	if err != nil {
		fmt.Println("Failed : ", err)
		return
	}

	var message Message
	err = astm1384.Unmarshal2(fileData, &message,
		astm1384.EncodingUTF8, astm1384.TimezoneEuropeBerlin, astm1384.LIS2A2)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Printf("Decoded:\n%+v\n", message)

}

type MinimalMessage struct {
	Header     Header     `astm:"H"`
	Terminator Terminator `astm:"L"`
	//Manufacturer Manufacturer `astm:"M,optional"`
	//Patient      Patient      `astm:"P"`
	//Order        Order        `astm:"O"`
	//Result       Result       `astm:"R"`
}

func TestReadMinimalMessage_v2(t *testing.T) {
	fileData, err := ioutil.ReadFile("../protocoltest/becom/5.2/minimal.astm")
	if err != nil {
		fmt.Println("Failed : ", err)
		return
	}

	var message MinimalMessage
	err = astm1384.Unmarshal2(fileData, &message,
		astm1384.EncodingUTF8, astm1384.TimezoneEuropeBerlin, astm1384.LIS2A2)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	assert.Equal(t, "Bio-Rad", message.Header.SenderNameOrID)

	fmt.Printf("Decoded:\n%+v\n", message)

}

/* tstcae array */
type MessagePOR struct {
	Header       Header       `astm:"H"`
	Manufacturer Manufacturer `astm:"M,optional"`
	OrderResults []struct {
		Patient         Patient `astm:"P"`
		Order           Order   `astm:"O"`
		CommentedResult []struct {
			result  Result    `astm:"R"`
			comment []Comment `astm:"C,optional"`
		}
	}
	Terminator Terminator `astm:"L"`
}

func TestReadFileBeComPOR_v2(t *testing.T) {
	fileData, err := ioutil.ReadFile("../protocoltest/becom/5.2/bloodtype_test_por.astm")
	if err != nil {
		fmt.Println("Failed : ", err)
		return
	}
	var message MessagePOR
	err = astm1384.Unmarshal2(fileData, &message,
		astm1384.EncodingUTF8, astm1384.TimezoneEuropeBerlin, astm1384.LIS2A2)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

}
