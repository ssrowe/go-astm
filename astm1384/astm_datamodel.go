package astm1384

import "time"

// This file is not ment to be edited careless. It controls the generation of output file

// Important Note:
// When using subfields e.g. astm:x,y
// then these records must be followed with   []string astm x
// on failure : export wont work properly.

// (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
type Header struct {
	Delimiters              string    `astm:"1"`           // 6.2
	MessageControlID        string    `astm:"2"`           // 6.3
	AccessPassword          string    `astm:"3"`           // 6.4
	SenderNameOrID          string    `astm:"4"`           // 6.5
	SenderStreetAddress     string    `astm:"5"`           // 6.6
	Field6                  string    `astm:"6"`           // 6.7
	SenderTelephone         string    `astm:"7"`           // 6.8
	CharacteristicsOfSender string    `astm:"8"`           // 6.9
	ReceiverID              string    `astm:"9"`           // 6.10
	Comment                 string    `astm:"10"`          // 6.11
	ProcessingID            string    `astm:"11"`          // 6.12
	Version                 string    `astm:"12"`          // 6.13
	DateAndTime             time.Time `astm:"13,longdate"` // 6.14
}

// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type Patient struct {
	SequenceNumber                     int       `astm:"1"`   // 7.2 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	PracticeAssignedPatientID          string    `astm:"2"`   // 7.3
	LabAssignedPatientID               string    `astm:"3"`   // 7.4
	ID3                                string    `astm:"4"`   // 7.5
	LastName                           string    `astm:"5,1"` // 7.6.1
	FirstName                          string    `astm:"5,2"` // 7.6.2
	Name                               []string  `astm:"5"`   // 7.6
	MothersMaidenName                  string    `astm:"6"`   // 7.7
	DOB                                time.Time `astm:"7"`   // 7.8
	Gender                             string    `astm:"8"`   // 7.9
	Race                               string    `astm:"9"`   // 7.10
	Address                            string    `astm:"10"`  // 7.11
	F12                                string    `astm:"11"`  // 7.12
	Telephone                          string    `astm:"12"`  // 7.13
	AttendingPhysicianID               string    `astm:"13"`  // 7.14
	SpecialField1                      string    `astm:"14"`  // 7.15
	SpecialField2                      string    `astm:"15"`  // 7.16
	Height                             string    `astm:"16"`  // 7.17
	Weight                             string    `astm:"17"`  // 7.18
	SuspectedDiagnosis                 string    `astm:"18"`  // 7.19
	ActiveMedication                   string    `astm:"19"`  // 7.20
	Diet                               string    `astm:"20"`  // 7.21
	PracticeField1                     string    `astm:"21"`  // 7.22
	PracticeField2                     string    `astm:"22"`  // 7.23
	AdmissionAndDischargeDates         string    `astm:"23"`  // 7.24
	AdmissionStatus                    string    `astm:"24"`  // 7.25
	Location                           string    `astm:"25"`  // 7.26
	NatureOfAlternativeDiagnosticCodes string    `astm:"26"`  // 7.27
	AlternativeDiagnosticCodes         string    `astm:"27"`  // 7.28
	Religion                           string    `astm:"28"`  // 7.29
	MaritalStatus                      string    `astm:"29"`  // 7.30
	IsolationStatus                    string    `astm:"30"`  // 7.31
	Language                           string    `astm:"31"`  // 7.32
	HospitalService                    string    `astm:"32"`  // 7.33
	HospitalInstitution                string    `astm:"33"`  // 7.34
	DosageCategory                     string    `astm:"34"`  // 7.35
}

// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type Order struct {
	SequenceNumber                      int        `astm:"1"`           // 8.4.2 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	SpecimenID                          []string   `astm:"2"`           // 8.4.3
	InstrumentSpecimenID                [][]string `astm:"3"`           // 8.4.4
	UniversalTestID_LOINC               string     `astm:"4,1"`         // 8.4.5.1
	UniversalTestID_TestName            string     `astm:"4,2"`         // 8.4.5.2
	UniversalTestID_AlternativeTestName string     `astm:"4,3"`         // 8.4.5.2
	UniversalTestID_ManufacturerCode    string     `astm:"4,4"`         // 8.4.5.3
	UniversalTestID_Custom1             string     `astm:"4,5"`         // 8.4.5.4
	UniversalTestID_Custom2             string     `astm:"4,6"`         // 8.4.5.4
	UniversalTestID_Custom3             string     `astm:"4,7"`         // 8.4.5.4
	UniversalTestID_Custom4             string     `astm:"4,8"`         // 8.4.5.4
	UniversalTestID_Custom5             string     `astm:"4,9"`         // 8.4.5.4
	UniversalTestID_Custom6             string     `astm:"4,10"`        // 8.4.5.4
	UniversalTestID_Custom7             string     `astm:"4,11"`        // 8.4.5.4
	UniversalTestID_Custom8             string     `astm:"4,12"`        // 8.4.5.4
	UniversalTestID_Custom9             string     `astm:"4,13"`        // 8.4.5.4
	UniversalTestID                     []string   `astm:"4"`           // 8.4.5 LOINC^BatterName^AlternateName^Manufacturersname^... unlimited test-fields
	Priority                            string     `astm:"5"`           // 8.4.6
	RequestedOrderDateTime              time.Time  `astm:"6,longdate"`  // 8.4.7
	SpecimenCollectionDateTime          time.Time  `astm:"7,longdate"`  // 8.4.8
	CollectionEndTime                   time.Time  `astm:"8,longdate"`  // 8.4.9
	CollectionVolume                    string     `astm:"9"`           // 8.4.10
	CollectionID                        string     `astm:"10"`          // 8.4.11
	ActionCode                          string     `astm:"11"`          // 8.4.12
	DangerCode                          string     `astm:"12"`          // 8.4.13
	RelevantClinicalInformation         string     `astm:"13"`          // 8.4.14
	DateTimeSpecimenReceived            string     `astm:"14"`          // 8.4.15
	SpecimenDescriptor                  string     `astm:"15"`          // 8.4.16
	OrderingPhysician                   string     `astm:"16"`          // 8.4.17
	PhysicianTelephone                  string     `astm:"17"`          // 8.4.18
	UserField1                          string     `astm:"18"`          // 8.4.19
	UserField2                          string     `astm:"19"`          // 8.4.20
	LaboratoryField1                    string     `astm:"20"`          // 8.4.21
	LaboratoryField2                    string     `astm:"21"`          // 8.4.22
	DateTimeResultsReported             time.Time  `astm:"22,longdate"` // 8.4.23
	InstrumentCharge                    string     `astm:"23"`          // 8.4.24
	InstrumentSectionID                 string     `astm:"24"`          // 8.4.25
	ReportType                          string     `astm:"25"`          // 8.4.26
	Reserved                            string     `astm:"26"`          // 8.4.27
	LocationOfSpecimenCollection        string     `astm:"27"`          // 8.4.28
	NosocomialInfectionFlag             string     `astm:"28"`          // 8.4.29
	SpecimenService                     string     `astm:"29"`          // 8.4.30
	SpecimenInstitution                 string     `astm:"30"`          // 8.4.31
}

// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type Result struct {
	SequenceNumber                           int        `astm:"1"`           // 9.2 https://samson-rus.com/wp-content/files/LIS2-A2.pdf
	UniversalTestID                          [][]string `astm:"2"`           // 9.3
	Data                                     []string   `astm:"3"`           // 9.4
	Units                                    []string   `astm:"4"`           // 9.5
	ReferenceRange                           []string   `astm:"5"`           // 9.6
	ResultAbnormalFlag                       string     `astm:"6"`           // 9.7
	NatureOfAbnormalTesting                  string     `astm:"7"`           // 9.8
	ResultStatus                             string     `astm:"8"`           // 9.9
	DateOfChangeInInstrumentNormativeTesting time.Time  `astm:"9,longdate"`  // 9.10
	OperatorIdentification                   []string   `astm:"10"`          // 9.11
	DateTimeTestStarted                      time.Time  `astm:"11,longdate"` // 9.12
	DateTimeCompleted                        time.Time  `astm:"12,longdate"` // 9.13
	IntstrumentIdentification                string     `astm:"13"`          // 9.14
	InstrumentName                           string     `astm:"14"`          // out of standard
	InstrumentSerialNumber                   string     `astm:"15"`          // out of standard
	InstrumentOperatiorID                    string     `astm:"16"`          // out of standard
}

// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type Comment struct {
	SequenceNumber int        `astm:"1"` // 10.2 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	CommentSource  [][]string `astm:"2"` // 10.3
	Comment        []string   `astm:"3"` // 10.4
	CommentType    string     `astm:"4"` // 10.5
	Custom1        string     `astm:"5"` // 10.6 (not standard)
}

// Lis2RequestInformation - RequestInformation
// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type RequestInformation struct {
	SequenceNumber                  int    `astm:"1"`  // 11.2 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	StartingRangeIDNumber           string `astm:"2"`  // 11.3
	EndingRangeIDNumber             string `astm:"3"`  // 11.4
	UniversalTestID                 string `astm:"4"`  // 11.5
	NatureOfRequestTimeLimits       string `astm:"5"`  // 11.6
	BeginningRequestResultsDateTime string `astm:"6"`  // 11.7
	EndingRequestResultsDateTime    string `astm:"7"`  // 11.8
	RequestingPhysicianName         string `astm:"8"`  // 11.9
	RequestingPhysicianTelephone    string `astm:"9"`  // 11.10
	UserField1                      string `astm:"10"` // 11.11
	UserField2                      string `astm:"11"` // 11.12
	RequestInformationStatus        string `astm:"12"` // 11.13
}

// Lis2Terminator - Terminator Record (Hasta la vista....)
// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type Terminator struct {
	SequenceNumber int    `astm:"1"` // 12.2 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	TerminatorCode string `astm:"2"` // 12.3
}

// Lis2Manufacturer -Manufacturer Record
// https://samson-rus.com/wp-content/files/LIS2-A2.pdf
type Manufacturer struct {
	SequenceNumber string `astm:"1"`  // 14 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	F2             string `astm:"2"`  // 2
	F3             string `astm:"3"`  // 3
	F4             string `astm:"4"`  // 4
	F5             string `astm:"5"`  // 5
	F6             string `astm:"6"`  // 6
	F7             string `astm:"7"`  // 7
	F8             string `astm:"8"`  // 8
	F9             string `astm:"9"`  // 9
	F10            string `astm:"10"` // 10
	F11            string `astm:"11"` // 11
	F12            string `astm:"12"` // 12
	F13            string `astm:"13"` // 13
}
