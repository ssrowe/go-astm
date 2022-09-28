package lis2a2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifyOrderMessage(t *testing.T) {

	astm := "H|\\^&|||LIS|||||NEO|||LIS2A2|20220928182311\n"
	astm = astm + "P|1||||^|||||||||||||||||||||||||||||\n"
	astm = astm + "O|1|idk1||^^^Pool_Cell||||R|||N||||Blood^Product|||||||||||||||\n"
	astm = astm + "L|1|N\n"

	messageType, err := IdentifyMessage([]byte(astm), EncodingUTF8)
	assert.Nil(t, err)

	assert.Equal(t, MessageTypeOrdersOnly, messageType)
}

func TestIdentifyQuery(t *testing.T) {

	astm := `H|\^&|||RVT|||||LIS|||LIS2-A2|20200302132021
Q|1|VALI200301||ALL
Q|2|VALI200302||ALL
Q|3|VALI200303||ALL
Q|4|VALI200304||ALL
Q|5|VALI200305||ALL
L|1|N`

	messageType, err := IdentifyMessage([]byte(astm), EncodingUTF8)
	assert.Nil(t, err)

	assert.Equal(t, MessageTypeQuery, messageType)
}

func TestIdentifyOrderAndResult(t *testing.T) {

	astm := `H|\^&|||RVT|||||LIS|||LIS2-A2|20200302131145
P|1||||^^^^|||U|||||||||||||||||Main||||||||||
O|1|VAL99999903||^^^Pool_Cell|R||||||||||^||||||||||F||||||
R|1|^^^Pool_Cell 1|0^0^8.8|||||F||Immucor||20200226153444|5030100389|
R|2|^^^Pool_Cell|Negative|||||F||immucor||20200226153444|5030100389|
L|1|N`
	messageType, err := IdentifyMessage([]byte(astm), EncodingUTF8)

	assert.Nil(t, err)
	assert.Equal(t, MessageTypeOrdersAndResults, messageType)
}

func TestIdentifyWithEmptyLines(t *testing.T) {

	astm := `H|\^&|||RVT|||||LIS|||LIS2-A2|20200302132021
Q|1|VALI200301||ALL
Q|2|VALI200302||ALL

Q|4|VALI200304||ALL
Q|5|VALI200305||ALL
L|1|N

`

	messageType, err := IdentifyMessage([]byte(astm), EncodingUTF8)
	assert.Nil(t, err)

	assert.Equal(t, MessageTypeQuery, messageType)
}
