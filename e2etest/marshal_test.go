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

	for _, line := range lines {
		linestr := string(line)
		fmt.Println(linestr)
	}
	assert.Nil(t, err)
}
