package e2e

import (
	"bytes"
	"testing"
	"time"

	"github.com/DRK-Blutspende-BaWueHe/go-astm/lib/standardlis2a2"
	"github.com/DRK-Blutspende-BaWueHe/go-astm/lis2a2"
	"github.com/stretchr/testify/assert"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type MinimalMessage struct {
	Header     standardlis2a2.Header     `astm:"H"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestReadMinimalMessage(t *testing.T) {
	fileData := ""
	fileData = fileData + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\n"
	fileData = fileData + "L|1|N\n"

	var message MinimalMessage
	err := lis2a2.Unmarshal([]byte(fileData), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)

	locale, err := time.LoadLocation("Europe/Berlin")
	assert.Nil(t, err)

	assert.Equal(t, "Bio-Rad", message.Header.SenderNameOrID)
	assert.Equal(t, "IH v5.2", message.Header.SenderStreetAddress)
	assert.Equal(t, "", message.Header.Comment)

	localtime := message.Header.DateAndTime.In(locale)
	assert.Equal(t, "20220315194227", localtime.Format("20060102150405"))
}

type FullSingleASTMMessage struct {
	Header       standardlis2a2.Header       `astm:"H"`
	Manufacturer standardlis2a2.Manufacturer `astm:"M,optional"`
	Patient      standardlis2a2.Patient      `astm:"P"`
	Order        standardlis2a2.Order        `astm:"O"`
	Result       standardlis2a2.Result       `astm:"R"`
	Terminator   standardlis2a2.Terminator   `astm:"L"`
}

func TestFullSingleASTMMessage(t *testing.T) {
	fileData := ""
	fileData = fileData + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\n"
	fileData = fileData + "P|1||1010868845||Testus^Test||19400607|M||||||||||||||||||||||||^\n"
	fileData = fileData + "O|1|1122206642|specimen1^^^\\specimen2^^^|^^^MO10^^28343^|R|20220311103217|20220311103217|||||||||||11||||20220311114103|||P\n"
	fileData = fileData + "R|1|^^^AntiA^MO10^Bloodgroup: A,B,D Confirmation for Patients (DiaClon) (5005)^|40^^|C||||R||lalina^|20220311114103||11|IH-1000|0300768|lalina\n"
	fileData = fileData + "L|1|N\n"

	var message FullSingleASTMMessage
	err := lis2a2.Unmarshal([]byte(fileData), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)

	assert.Equal(t, "Testus", message.Patient.LastName)
	assert.Equal(t, "Test", message.Patient.FirstName)
	assert.Equal(t, "19400607", message.Patient.DOB.Format("20060102")) // dates are read okay
	assert.Equal(t, "specimen1", message.Order.InstrumentSpecimenID)
	assert.Equal(t, "lalina", message.Result.OperatorIDPerformed)
}

//-----------------------------------------------------------------------------------
// TEST
//-----------------------------------------------------------------------------------
// Testing a rather more complex structure with optional and arrays on the
// structure as well as on the Records
//-----------------------------------------------------------------------------------
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

func TestNestedStructure(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||1010868845||Testus^Test||19400607|M||||||||||||||||||||||||^\r"
	data = data + "O|1|1122206642|1122206642^^^\\1122206642^^^|^^^MO10^^28343^|R|20220311103217|20220311103217|||||||||||11||||20220311114103|||P\r"
	data = data + "R|1|^^^AntiA^MO10^Bloodgroup: A,B,D Confirmation for Patients (DiaClon) (5005)^|40^^|C||||R||lalina^|20220311114103||11|IH-1000|0300768|lalina\r"
	data = data + "C|1|FirstComment^^05761.03.12^20240131\\^^^|CAS^5005352062212117030^50053.52.06^20221231^4||\r"
	data = data + "C|2|SecondComment^^05761.03.12^20240131\\^^^|CAS^5005352062212117030^50053.52.06^20221231^4||\r"
	data = data + "R|2|^^^AntiB^MO10^Bloodgroup: A,B,D Confirmation for Patients (DiaClon) (5005)^|0^^|C||||R||lalina^|20220311114103||11|IH-1000|0300768|lalina\r"
	data = data + "C|1|ID-Diluent 2^^05761.03.12^20240131\\^^^|CAS^5005352062212117030^50053.52.06^20221231^5||\r"
	data = data + "R|3|^^^AntiD^MO10^Bloodgroup: A,B,D Confirmation for Patients (DiaClon) (5005)^|0^^|C||||R||lalina^|20220311114103||11|IH-1000|0300768|lalina\r"
	data = data + "P|2||1010868845||Testis^Tost||19400607|M||||||||||||||||||||||||^\r"
	data = data + "O|1|1122206642|1122206642^^^\\1122206642^^^|^^^MO10^^28343^|R|20220311103217|20220311103217|||||||||||11||||20220311114103|||P\r"
	data = data + "R|1|^^^AntiA^MO10^Bloodgroup: A,B,D Confirmation for Patients (DiaClon) (5005)^|40^^|C||||R||lalina^|20220311114103||11|IH-1000|0300768|lalina\r"
	data = data + "L|1|N\r"

	var message MessagePORC
	err := lis2a2.Unmarshal([]byte(data), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)

	// patients have been read in an array
	assert.Equal(t, "Testus", message.OrderResults[0].Patient.LastName)
	assert.Equal(t, "Testis", message.OrderResults[1].Patient.LastName)

	// the array of comments was produced in two entries into the array
	assert.Equal(t, 2, len(message.OrderResults[0].CommentedResult[0].Comment))
	assert.Equal(t, "FirstComment", message.OrderResults[0].CommentedResult[0].Comment[0].CommentSource)
	assert.Equal(t, "SecondComment", message.OrderResults[0].CommentedResult[0].Comment[1].CommentSource)
}

//-----------------------------------------------------------------------------------
// TEST
//-----------------------------------------------------------------------------------
// Custom Delimiters : In the header there is a delimiter-field that can change
// the default delimiters
//-----------------------------------------------------------------------------------
type MessageCustomDelimiterTest struct {
	Header     standardlis2a2.Header     `astm:"H"`
	Patient    standardlis2a2.Patient    `astm:"P"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestCustomDelimiters(t *testing.T) {
	data := ""
	data = data + "H|\\#&|||Bio-Rad|IH v5.2||||||||20220315194227\n"
	data = data + "P|1||1010868845||Testus#Test||19400607|M||||||||||||||||||||||||^\r"
	data = data + "L|1|N\n" // ! mixed line-endings (should not matter)

	var message MessageCustomDelimiterTest
	err := lis2a2.Unmarshal([]byte(data), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)

	// the delimiter is now "#" instead of "^"; this should have been recognizes
	assert.Equal(t, "Testus", message.Patient.LastName)
	assert.Equal(t, "Test", message.Patient.FirstName)

}

//-----------------------------------------------------------------------------------
// TEST
//-----------------------------------------------------------------------------------
// Custom Structures can be defined and mapped to records
//  aside:    also testing float values
//-----------------------------------------------------------------------------------
type CompleteOutOfStandardCustomRecord struct {
	SequenceNumber string  `astm:"2"`
	F2             string  `astm:"3"`
	F3             string  `astm:"4"`
	Float32Value   float32 `astm:"5"`
	Float64Value   float64 `astm:"6"`
}

type MessageWithOutOfStandardCustomRecord struct {
	Header       standardlis2a2.Header             `astm:"H"`
	CustomRecord CompleteOutOfStandardCustomRecord `astm:"X"`
	Terminator   standardlis2a2.Terminator         `astm:"L"`
}

func TestCustomRecord(t *testing.T) {
	data := ""
	data = data + "H|\\#&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "X|1|A|B|4.14159|2.172\r"
	data = data + "L|1|N\r" // ! mixed line-endings (should not matter)

	var message MessageWithOutOfStandardCustomRecord
	err := lis2a2.Unmarshal([]byte(data), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)
	assert.Equal(t, float32(4.14159), message.CustomRecord.Float32Value)
	assert.Equal(t, float64(2.172), message.CustomRecord.Float64Value)
}

// test timezone ok
// test delimiters ok
// test optionals ok
// test encoding
// line ending 0a or 0d or 0d0a all okay ? ok

type SubMessageRecord struct {
	Field11 string `astm:"2.1.1"`
	Field12 string `astm:"2.1.2"`
	Field13 string `astm:"2.1.3"`
	Field21 string `astm:"2.2.1"`
	Field22 string `astm:"2.2.2"`
	Field23 string `astm:"2.2.3"`
}

type MessageWithSubArrayRecord struct {
	Anonymous struct { // Testing wether this annoynmous structure is recused into
		Rec SubMessageRecord `astm:"?"`
	}

	AnonymousArray []struct { // This anynymous array of structures needs to be created and scanned
		Rec SubMessageRecord `astm:"!"`
	}
}

func TestArrayMapping(t *testing.T) {

	data := "?|a^^c\\d^e^f|\r"
	data = data + "!|x^y\\z^^|\r"
	data = data + "!|1^2^3\\4^5^6|\r"

	var message MessageWithSubArrayRecord
	err := lis2a2.Unmarshal([]byte(data), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)

	assert.Equal(t, "a", message.Anonymous.Rec.Field11)
	assert.Equal(t, "c", message.Anonymous.Rec.Field13)

	assert.Equal(t, "d", message.Anonymous.Rec.Field21)
	assert.Equal(t, "e", message.Anonymous.Rec.Field22)
	assert.Equal(t, "f", message.Anonymous.Rec.Field23)

	// now test that the subarray values have been read
	assert.Equal(t, 2, len(message.AnonymousArray))
	assert.Equal(t, "x", message.AnonymousArray[0].Rec.Field11)
	assert.Equal(t, "y", message.AnonymousArray[0].Rec.Field12)
	assert.Equal(t, "z", message.AnonymousArray[0].Rec.Field21)
	assert.Equal(t, "", message.AnonymousArray[0].Rec.Field22)

	assert.Equal(t, "1", message.AnonymousArray[1].Rec.Field11)
	assert.Equal(t, "2", message.AnonymousArray[1].Rec.Field12)
	assert.Equal(t, "4", message.AnonymousArray[1].Rec.Field21)
	assert.Equal(t, "5", message.AnonymousArray[1].Rec.Field22)
}

type SomeEnum string

const (
	EnumValue1 SomeEnum = "EnumValue1"
	EnumValue2 SomeEnum = "EnumValue2"
)

type SomeEnumRecord struct {
	Value SomeEnum `astm:"2"`
}

type TestUnmarshalEnumMessage struct {
	Record SomeEnumRecord `astm:"E"`
}

// TODO enum value
func TestEnumEncoding(t *testing.T) {
	data := "E|EnumValue1|\r"

	var message TestUnmarshalEnumMessage
	err := lis2a2.Unmarshal([]byte(data), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)

	assert.Equal(t, EnumValue1, message.Record.Value)
}

//-----------------------------------------------------------------------------------
// TEST
//-----------------------------------------------------------------------------------
// Access time.Time in a struct with multiple components
//  aside:    also testing time.time values
//-----------------------------------------------------------------------------------
type AccessTimeComment struct {
	SequenceNumber              int       `astm:"2,sequence" db:"sequence_number"`            // 3.2.8 - ih_com_host_connection_manual_astm_1394_en_h009164_v1_8.pdf
	DescriptionOfReagent        string    `astm:"3.1.1"  db:"description_of_reagent"`         //
	BarcodeOfReagent            string    `astm:"3.1.2" db:"barcode_of_reagent"`              //
	LotNumberOfReagent          string    `astm:"3.1.3" db:"lot_number_of_reagent"`           //
	ExpirationDateOfReagent     time.Time `astm:"3.1.4" db:"expiration_date_of_reagent"`      //
	DescriptionOfReagent2       string    `astm:"3.2.1" db:"description_of_reagent_2"`        //
	BarcodeOfReagent2           string    `astm:"3.2.2" db:"barcode_of_reagent_2"`            //
	LotNumberOfReagent2         string    `astm:"3.2.3" db:"lot_number_of_reagent_2"`         //
	ExpirationDateOfReagent2    time.Time `astm:"3.2.4" db:"expiration_date_of_reagent_2"`    //
	TypeOfTestMedia             string    `astm:"4.1" db:"type_of_test_media"`                //
	PlateOrIDCardBarcode        string    `astm:"4.2" db:"plate_or_id_card_barcode"`          //
	LotNumberOfCassetteOrPlate  string    `astm:"4.3" db:"lot_number_of_cassette_or_plate"`   //
	ExpDateForIDCardOrPlate     time.Time `astm:"4.4" db:"exp_date_for_id_card_or_plate"`     //
	IDCardOrPlateRealWellNumber int       `astm:"4.5" db:"id_card_or_plate_real_well_number"` //
	Comment                     string    `astm:"5" db:"comment"`                             //
	FileName                    string    `astm:"6" db:"file_name"`                           //
}
type MessageTimeAccess struct {
	Header     standardlis2a2.Header     `astm:"H"`
	Comment    AccessTimeComment         `astm:"C"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestComponentAccessOnTime(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "C|1|FirstComment^^05761.03.12^20240131\\^^^|CAS^5005352062212117030^50053.52.06^20221231^4||\r"
	data = data + "L|1|N\r"

	var message MessageTimeAccess
	err := lis2a2.Unmarshal([]byte(data), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)

	location, err := time.LoadLocation(string(lis2a2.TimezoneEuropeBerlin))
	assert.Nil(t, err, "Can not parse timezone")

	expDate, err := time.ParseInLocation("20060102", "20240131", location)
	assert.Nil(t, err, "Can not parse date")
	assert.Equal(t, expDate, message.Comment.ExpirationDateOfReagent)
}

type TestCommentNoneBugComment struct {
	SequenceNumber int       `astm:"2,sequence"`
	Field1         time.Time `astm:"3.1.4"` // out of bounds with component index
	Field2         time.Time `astm:"3.2.4"` // out of bounds with repeat index
	Field3         time.Time `astm:"4.4"`
}
type TestCommentNoneBugMessage struct {
	Field TestCommentNoneBugComment `astm:"C"`
}

type TestCommentNoneBugCommentCrash struct {
	SequenceNumber int       `astm:"2,sequence"`
	Field1         time.Time `astm:"3.1.4,required"` // out of bounds with component index
	Field2         time.Time `astm:"3.2.4"`          // out of bounds with repeat index
	Field3         time.Time `astm:"4.4,required"`
}

type TestCommentNoneBugMessageCrash struct {
	Field TestCommentNoneBugComment `astm:"C"`
}

func TestCommentNoneBug(t *testing.T) {
	data := ""
	data = data + "C|1|^^^||\r"

	var message TestCommentNoneBugMessage
	err := lis2a2.Unmarshal([]byte(data), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)

	assert.Equal(t, time.Time{}, message.Field.Field1)
	assert.Equal(t, time.Time{}, message.Field.Field2)
	assert.Equal(t, time.Time{}, message.Field.Field3)

	/* var crash TestCommentNoneBugMessageCrash
	err := lis2a2.Unmarshal([]byte(data), &crash,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)
	assert.NotNil(t, err) */
}

//-----------------------------------------------------------------------------------
// TEST a german message Win1252 Encoded
//-----------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------
type MessageGermanLanguageTest struct {
	Header     standardlis2a2.Header     `astm:"H"`
	Patient    standardlis2a2.Patient    `astm:"P"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestGermanLanguage(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\n"
	data = data + "P|1||1010868845||König^#$§?/+öäüß||19400607|M||||||||||||||||||||||||^\r"
	data = data + "L|1|N\n" // ! mixed line-endings (should not matter)

	var message MessageGermanLanguageTest

	encdata := helperEncode(charmap.Windows1252, []byte(data))
	err := lis2a2.Unmarshal([]byte(encdata), &message, lis2a2.EncodingWindows1252, lis2a2.TimezoneEuropeBerlin)
	assert.Nil(t, err)
	assert.Equal(t, "König", message.Patient.LastName)
	assert.Equal(t, "#$§?/+öäüß", message.Patient.FirstName)

	encdata = helperEncode(charmap.ISO8859_1, []byte(data))
	err = lis2a2.Unmarshal([]byte(encdata), &message, lis2a2.EncodingISO8859_1, lis2a2.TimezoneEuropeBerlin)
	assert.Nil(t, err)
	assert.Equal(t, "König", message.Patient.LastName)
	assert.Equal(t, "#$§?/+öäüß", message.Patient.FirstName)
}

func TestTransmissionWithoutLTerminator(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||\r"
	data = data + "P|1||DIA-27-079-5-1\r"

	var message standardlis2a2.DefaultMessage
	err := lis2a2.Unmarshal([]byte(data), &message, lis2a2.EncodingWindows1252, lis2a2.TimezoneEuropeBerlin)
	assert.NotNil(t, err)
}

func TestFullMultipleASTMMessage(t *testing.T) {
	var data string

	// Message 1
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSQVIGG3||20220715071219\r"
	data = data + "R|1|^^^SARSQVIGG3|2598,88|BAU/ml|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSQVIGG3||20220715071219\r"
	data = data + "R|1|^^^SARSQVIGG3|3636,64|BAU/ml|\r"
	data = data + "L|1|N\r"

	// Message 2
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSNCPIGG||20220715071219\r"
	data = data + "R|1|^^^SARSNCPIGG|0,08|Ratio|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSNCPIGG||20220715071219\r"
	data = data + "R|1|^^^SARSNCPIGG|0,20|Ratio|\r"
	data = data + "L|1|N\r"

	// Message 3
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSNEUTRA||20220715071219\r"
	data = data + "R|1|^^^SARSNEUTRA|99,39|% IH|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSNEUTRA||20220715071219\r"
	data = data + "R|1|^^^SARSNEUTRA|99,23|% IH|\r"
	data = data + "L|1|N\r"

	// Message 4
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|>10|Ratio|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|>10|Ratio|\r"
	data = data + "P|3||DIA-01-061-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|4,87|Ratio|\r"
	data = data + "P|4||DIA-01-055-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|4,14|Ratio|\r"
	data = data + "L|1|N"

	var message standardlis2a2.DefaultMultiMessage
	err := lis2a2.Unmarshal(
		[]byte(data), &message,
		lis2a2.EncodingUTF8,
		lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, 4, len(message.Messages))
}

func TestFullMultipleASTMMessageWithWrongInput(t *testing.T) {
	var data string

	// Message 1
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSQVIGG3||20220715071219\r"
	data = data + "R|1|^^^SARSQVIGG3|2598,88|BAU/ml|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSQVIGG3||20220715071219\r"
	data = data + "R|1|^^^SARSQVIGG3|3636,64|BAU/ml|\r"
	data = data + "L|1|N\r"

	// Message 2
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSNCPIGG||20220715071219\r"
	data = data + "R|1|^^^SARSNCPIGG|0,08|Ratio|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSNCPIGG||20220715071219\r"
	data = data + "R|1|^^^SARSNCPIGG|0,20|Ratio|\r"
	data = data + "L|1|N\r"

	// Message 3
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSNEUTRA||20220715071219\r"
	data = data + "R|1|^^^SARSNEUTRA|99,39|% IH|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSNEUTRA||20220715071219\r"
	data = data + "R|1|^^^SARSNEUTRA|99,23|% IH|\r"
	data = data + "L|1|N\r"

	// Message 4
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|>10|Ratio|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|>10|Ratio|\r"
	data = data + "P|3||DIA-01-061-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|4,87|Ratio|\r"
	data = data + "P|4||DIA-01-055-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|4,14|Ratio|\r"
	data = data + "L|1|N"

	var message standardlis2a2.DefaultMessage
	err := lis2a2.Unmarshal(
		[]byte(data),
		&message,
		lis2a2.EncodingUTF8,
		lis2a2.TimezoneEuropeBerlin)

	assert.NotNil(t, err)
}

func helperEncode(charmap *charmap.Charmap, data []byte) []byte {
	e := charmap.NewEncoder()
	var b bytes.Buffer
	writer := transform.NewWriter(&b, e)
	writer.Write([]byte(data))
	resultdata := b.Bytes()
	writer.Close()
	return resultdata
}

func TestFailOnUndisciplinedMultipleCRCRatEndOfLine(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||\u000d\u000d"
	data = data + "P|1||DIA-04-066-7-1\u000d\u000d"
	data = data + "O|1|||^^^SARS-CoV-2 NeutraLISA||20220715071342\u000d\u000d"
	data = data + "R|1|^^^SARS-CoV-2 NeutraLISA|99,66|% IH|\u000d\u000d"
	data = data + "L|1|N\u000d\u000d"

	var message standardlis2a2.DefaultMessage
	err := lis2a2.Unmarshal([]byte(data), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)
}

func TestMultipleMessagesInOne(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||\u000d\u000d"
	data = data + "P|1||DIA-04-066-7-1\u000d\u000d"
	data = data + "O|1|||^^^SARS-CoV-2 NeutraLISA||20220715071342\u000d\u000d"
	data = data + "R|1|^^^SARS-CoV-2 NeutraLISA|12,5|% IH|\u000d\u000d"
	data = data + "L|1|N\u000d\u000d"
	data = data + "H|\\^&|||\u000d\u000d"
	data = data + "P|1||DIA-04-066-7-2\u000d\u000d"
	data = data + "O|1|||^^^SARS-CoV-2 NeutraLISA||20220715071343\u000d\u000d"
	data = data + "R|1|^^^SARS-CoV-2 NeutraLISA|99,66|% IH|\u000d\u000d"
	data = data + "L|1|N\u000d\u000d"

	var message standardlis2a2.DefaultMultiMessage
	err := lis2a2.Unmarshal([]byte(data), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(message.Messages))
	assert.Equal(t, "DIA-04-066-7-2", message.Messages[1].OrderResults[0].Patient.LabAssignedPatientID)

	assert.Equal(t, "12,5", message.Messages[0].OrderResults[0].CommentedResult[0].Result.DataMeasurementValue)
	assert.Equal(t, "99,66", message.Messages[1].OrderResults[0].CommentedResult[0].Result.DataMeasurementValue)
}
