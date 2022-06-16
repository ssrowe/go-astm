package e2e

import (
	"fmt"
	"github.com/DRK-Blutspende-BaWueHe/go-astm/astm1384"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEuroImmunResults(t *testing.T) {
	fileData, err := ioutil.ReadFile("../protocoltest/euroimmun/sampleigg.astm")
	if err != nil {
		fmt.Println("Failed : ", err)
		return
	}

	message, err := astm1384.Unmarshal(fileData,
		astm1384.EncodingWindows1252,
		astm1384.TimezoneUTC,
		astm1384.LIS2A2)
	if err != nil {
		fmt.Println("Error in unmarshaling the message ", err)
		return
	}

	assert.Equal(t, "TEST-27-079-5-1", message.Records[0].Patient.LabAssignedPatientID)
	assert.Equal(t, "SARSCOV2IGA", message.Records[0].OrdersAndResults[0].Order.UniversalTestID_ManufacturerCode)
	assert.Equal(t, "20220218080737", message.Records[0].OrdersAndResults[0].Order.RequestedOrderDateTime.Format("20060102150405"))
	assert.Equal(t, "SARSCOV2IGA", message.Records[0].OrdersAndResults[0].Results[0].Result.UniversalTestID[0][3])
	assert.Equal(t, ">8", message.Records[0].OrdersAndResults[0].Results[0].Result.Data[0])
	assert.Equal(t, "Ratio", message.Records[0].OrdersAndResults[0].Results[0].Result.Units[0])

	assert.Equal(t, 20 /*Records in Resultset*/, len(message.Records))
}
