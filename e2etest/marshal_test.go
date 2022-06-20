package e2e

import (
	"fmt"
	"testing"

	"github.com/DRK-Blutspende-BaWueHe/go-astm/lib/standardlis2a2"
	"github.com/DRK-Blutspende-BaWueHe/go-astm/lis2a2"
	"github.com/stretchr/testify/assert"
)

type IllFormatedButLegal struct {
	GeneratedSequence int    `astm:"1,sequence"`
	ThirdfieldArray1  string `astm:"2.1.3"`
	FirstFieldArray1  string `astm:"2.1.1"`
	FirstFieldArray2  string `astm:"2.2.1"`
	SecondfieldArray2 string `astm:"2.2.2"`
	SomeEmptyField    string `astm:"3"`
}

type MinimalMessageMarshal struct {
	Header     standardlis2a2.Header     `astm:"H"`
	Ill        IllFormatedButLegal       `astm:"?"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestSimpleMarshal(t *testing.T) {
	var msg MinimalMessageMarshal

	msg.Header.AccessPassword = "password"
	msg.Header.Version = "0.1.0"
	msg.Header.SenderNameOrID = "test"

	msg.Ill.ThirdfieldArray1 = "third-arr1"
	msg.Ill.FirstFieldArray1 = "first-arr1"
	msg.Ill.FirstFieldArray2 = "first-arr2"
	msg.Ill.SecondfieldArray2 = "second-arr2"

	lines, err := lis2a2.Marshal(msg, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)

	/*	// output on screen
		for _, line := range lines {
			linestr := string(line)
			fmt.Println(linestr)
			} */
	assert.Nil(t, err)

	assert.Equal(t, string(lines[0]), "H|\\^&||password|test||||||||0.1.0|")
	assert.Equal(t, string(lines[1]), "?|1|first-arr1^^third-arr1\\first-arr2^second-arr2||")
	assert.Equal(t, string(lines[2]), "L|1||")
}

type ArrayMessageMarshal struct {
	Header     standardlis2a2.Header     `astm:"H"`
	Patient    []standardlis2a2.Patient  `astm:"P"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestGenverateSequence(t *testing.T) {

	var msg ArrayMessageMarshal

	msg.Patient = make([]standardlis2a2.Patient, 2)
	msg.Patient[0].LastName = "Firstus'"
	msg.Patient[0].FirstName = "Firstie"
	msg.Patient[1].LastName = "Secundus'"
	msg.Patient[1].FirstName = "Secundie"

	lines, err := lis2a2.Marshal(msg, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)

	assert.Nil(t, err)
	// output on screen
	for _, line := range lines {
		linestr := string(line)
		fmt.Println(linestr)
	}

	assert.Equal(t, string(lines[0]), "H|\\^&||||||||||||")
	assert.Equal(t, string(lines[1]), "P|1||||Firstus'^Firstie|||||||||||||||||||||||||||||")
	assert.Equal(t, string(lines[2]), "P|2||||Secundus'^Secundie|||||||||||||||||||||||||||||")
	assert.Equal(t, string(lines[3]), "L|1||")
}

type PatientResult struct {
	Patient standardlis2a2.Patient  `astm:"P"`
	Result  []standardlis2a2.Result `astm:"R"`
}

type ArrayNestedStructMessageMarshal struct {
	Header        standardlis2a2.Header `astm:"H"`
	PatientResult []PatientResult
	Terminator    standardlis2a2.Terminator `astm:"L"`
}

func TestNestedStruct(t *testing.T) {
	var msg ArrayNestedStructMessageMarshal

	msg.PatientResult = make([]PatientResult, 2)
	msg.PatientResult[0].Patient.FirstName = "Chuck"
	msg.PatientResult[0].Patient.LastName = "Norris"
	msg.PatientResult[0].Patient.Religion = "Binaries"
	msg.PatientResult[0].Result = make([]standardlis2a2.Result, 2)
	msg.PatientResult[0].Result[0].ManufacturersTestName = "SulfurBloodCount"
	msg.PatientResult[0].Result[0].MeasurementValueOfDevice = "100"
	msg.PatientResult[0].Result[0].Units = "%"
	msg.PatientResult[0].Result[1].ManufacturersTestName = "Catblood"
	msg.PatientResult[0].Result[1].MeasurementValueOfDevice = ">100000"
	msg.PatientResult[0].Result[1].Units = "U/l"
	msg.PatientResult[1].Patient.FirstName = "Eric"
	msg.PatientResult[1].Patient.LastName = "Cartman"
	msg.PatientResult[1].Patient.Religion = "None"
	msg.PatientResult[1].Result = make([]standardlis2a2.Result, 1)
	msg.PatientResult[1].Result[0].ManufacturersTestName = "Fungenes"
	msg.PatientResult[1].Result[0].MeasurementValueOfDevice = "present"
	msg.PatientResult[1].Result[0].Units = "none"

	lines, err := lis2a2.Marshal(msg, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)

	assert.Nil(t, err)
	// output on screen
	for _, line := range lines {
		linestr := string(line)
		fmt.Println(linestr)
	}

	assert.Equal(t, string(lines[0]), "H|\\^&||||||||||||")
	assert.Equal(t, string(lines[1]), "P|1||||Norris^Chuck||||||||||||||||||||||Binaries|||||||")
	assert.Equal(t, string(lines[2]), "R|1|^^^^SulfurBloodCount^^|^^100|%|||||^||")
	assert.Equal(t, string(lines[3]), "R|2|^^^^Catblood^^|^^>100000|U/l|||||^||")
	assert.Equal(t, string(lines[4]), "P|1||||Cartman^Eric||||||||||||||||||||||None|||||||")
	assert.Equal(t, string(lines[5]), "R|1|^^^^Fungenes^^|^^present|none|||||^||")
	assert.Equal(t, string(lines[6]), "L|1||")
}

type TestMarshalEnum string

const (
	SomeTestMarshalEnum1 TestMarshalEnum = "Something"
	SomeTestMarshalEnum2 TestMarshalEnum = "SomethingElse"
)

type TestMarshalEnumRecord struct {
	Field TestMarshalEnum
}

type TestMarshalEnumMessage struct {
	Record TestMarshalEnumRecord `astm:"X"`
}

func TestEnumMarshal(t *testing.T) {
	var msg TestMarshalEnumMessage

	msg.Record.Field = SomeTestMarshalEnum2

	lines, err := lis2a2.Marshal(msg, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)

	assert.Nil(t, err)
	// output on screen
	for _, line := range lines {
		linestr := string(line)
		fmt.Println(linestr)
	}

}
