package e2e

import (
	"io/ioutil"
	"testing"

	"github.com/DRK-Blutspende-BaWueHe/go-astm/lib/standardlis2a2"
	"github.com/DRK-Blutspende-BaWueHe/go-astm/lis2a2"
	"github.com/stretchr/testify/assert"
)

type Message struct {
	Header       standardlis2a2.Header       `astm:"H"`
	Manufacturer standardlis2a2.Manufacturer `astm:"M,optional"`
	Patient      standardlis2a2.Patient      `astm:"P"`
	OrderResults []struct {
		Order           standardlis2a2.Order `astm:"O"`
		CommentedResult []struct {
			Result  standardlis2a2.Result  `astm:"R"`
			Comment standardlis2a2.Comment `astm:"C,optional"`
		}
	}
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func Testfatalconversion_bug(t *testing.T) {
	fileData, err := ioutil.ReadFile("bug_fatalconversion.dat")
	assert.Nil(t, err)

	var message Message
	err = lis2a2.Unmarshal([]byte(fileData), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)
}
