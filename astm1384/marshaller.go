package astm1384

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	charmap "github.com/aglyzov/charmap"
)

type ProtocolVersion int

const LIS2A2 ProtocolVersion = 2

type Encoding int

const EncodingUTF8 Encoding = 1
const EncodingASCII Encoding = 2
const EncodingWindows1250 Encoding = 3
const EncodingWindows1251 Encoding = 4
const EncodingWindows1252 Encoding = 5
const EncodingDOS852 Encoding = 6
const EncodingDOS855 Encoding = 7
const EncodingDOS866 Encoding = 8

type Timezone string

const TimezoneUTC Timezone = "UTC"
const TimezoneEuropeBerlin Timezone = "Europe/Berlin"
const TimezoneEuropeBudapest Timezone = "Europe/Budapest"
const TimezoneEuropeLondon Timezone = "Europe/London"

func Unmarshal(messageData []byte, enc Encoding, tz Timezone, pv ProtocolVersion) (*ASTMMessage, error) {

	switch pv {
	case ProtocolVersion(LIS2A2):
	default:
		return nil, errors.New("protocol Not implemented")
	}

	location, err := time.LoadLocation(string(tz))
	if err != nil {
		return nil, err
	}

	var messageStr string
	switch enc {
	case EncodingUTF8:
		// do nothing, this is correct
		messageStr = string(messageData)
	case EncodingASCII:
		messageStr = string(messageData)
	case EncodingDOS866:
		messageBytes, err := charmap.ANY_to_UTF8(messageData, "DOS866")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid input : %s", err))
		}
		messageStr = string(messageBytes)
	case EncodingDOS855:
		messageBytes, err := charmap.ANY_to_UTF8(messageData, "DOS855")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid input : %s", err))
		}
		messageStr = string(messageBytes)
	case EncodingDOS852:
		messageBytes, err := charmap.ANY_to_UTF8(messageData, "DOS852")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid input : %s", err))
		}
		messageStr = string(messageBytes)
	case EncodingWindows1250:
		messageBytes, err := charmap.ANY_to_UTF8(messageData, "CP1250")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid input : %s", err))
		}
		messageStr = string(messageBytes)
	case EncodingWindows1251:
		messageBytes, err := charmap.ANY_to_UTF8(messageData, "CP1251")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid input : %s", err))
		}
		messageStr = string(messageBytes)
	case EncodingWindows1252:
		messageBytes, err := charmap.ANY_to_UTF8(messageData, "CP1252")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid input : %s", err))
		}
		messageStr = string(messageBytes)
	default:
		return nil, errors.New(fmt.Sprintf("Invalid Codepage %d", enc))
	}

	tokenInput, err2 := astm1384Scanner(messageStr, location, "\n")
	if err2 != nil {
		return nil, err2
	}

	message, err := parseAST(tokenInput)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func Marshal(message *ASTMMessage, enc Encoding, tz Timezone, pv ProtocolVersion) ([]byte, error) {
	var buffer bytes.Buffer

	location, err := time.LoadLocation(string(tz))
	if err != nil {
		return []byte{}, err
	}

	err = convertToASTMFileRecord("H", message.Header, []string{"\\", "^", "&"}, location, &buffer)
	if err != nil {
		log.Println(err)
		return []byte{}, errors.New(fmt.Sprintf("Failed to marshal header: %s", err))
	}
	buffer.Write([]byte{10, 13})

	if message.Manufacturer != nil {
		err := convertToASTMFileRecord("M", message.Manufacturer, []string{"\\", "^", "&"}, location, &buffer)
		if err != nil {
			log.Println(err)
			return []byte{}, errors.New(fmt.Sprintf("Failed to marshal manufacturer-record: %s", err))
		}
		buffer.Write([]byte{10, 13})
	}

	for i, record := range message.Records {
		if record.Patient != nil {
			record.Patient.SequenceNumber = i + 1
			err := convertToASTMFileRecord("P", record.Patient, []string{"|", "^", "&"}, location, &buffer)
			if err != nil {
				log.Println(err)
				return []byte{}, errors.New(fmt.Sprintf("Failed to marshal header: %s", err))
			}
			buffer.Write([]byte{10, 13})

			for orderResultsIdx, orderResults := range record.OrdersAndResults {
				orderResults.Order.SequenceNumber = orderResultsIdx + 1
				err := convertToASTMFileRecord("O", orderResults.Order, []string{"|", "^", "&"}, location, &buffer)
				if err != nil {
					log.Println(err)
					return []byte{}, errors.New(fmt.Sprintf("Failed to marshal order-records: %s", err))
				}
				buffer.Write([]byte{10, 13})

				for resultsIdx, result := range orderResults.Results {
					result.Result.SequenceNumber = resultsIdx + 1
					err := convertToASTMFileRecord("R", result.Result, []string{"|", "^", "&"}, location, &buffer)
					if err != nil {
						log.Println(err)
						return []byte{}, errors.New(fmt.Sprintf("Failed to marshal result-record %s", err))
					}
					buffer.Write([]byte{10, 13})
					for commentIdx, comment := range result.Comments {
						comment.SequenceNumber = commentIdx + 1
						err := convertToASTMFileRecord("C", comment, []string{"|", "^", "&"}, location, &buffer)
						if err != nil {
							log.Println(err)
							return []byte{}, errors.New(fmt.Sprintf("Failed to marshal result-comment %s", err))
						}
						buffer.Write([]byte{10, 13})
					}
				}
			}
		}
	}
	buffer.Write([]byte("L|1|N"))
	buffer.Write([]byte{10, 13})

	return buffer.Bytes(), nil
}

func convertToASTMFileRecord(recordType string, target interface{}, delimiter []string, tz *time.Location, buffer *bytes.Buffer) error {

	t := reflect.TypeOf(target).Elem()

	entries := make(map[int]string, 0)

	maxIdx := 0

	for i := 0; i < t.NumField(); i++ {
		astmTag := t.Field(i).Tag.Get("astm")
		astmTagsList := strings.Split(astmTag, ",")
		if len(astmTagsList) == 0 || astmTag == "" {
			continue // nothing to process when someone requires astm:
		}
		idx, err := strconv.Atoi(astmTagsList[0])
		idx = idx - 1
		if idx < 0 {
			return errors.New(fmt.Sprintf("Illegal annotation <%s> in for field %s", astmTag, t.Name()))
		}
		if err != nil {
			return err
		}
		if idx > maxIdx {
			maxIdx = idx
		}

		isLongDate := false
		for i := 0; i < len(astmTagsList); i++ {
			if astmTagsList[i] == "longdate" {
				isLongDate = true
			}
		}

		field := reflect.ValueOf(target).Elem().Field(i)
		fieldValue := field.Interface()

		switch fieldValue.(type) {
		case int:
			entries[idx] = strconv.Itoa(int(field.Int()))
		case string:
			entries[idx] = string(field.String())
		case []string:
			arry := fieldValue.([]string)
			outString := ""
			for i := 0; i < len(arry); i++ {
				outString = outString + arry[i]
				if i < len(arry)-1 {
					outString = outString + "^"
				}
			}
			entries[idx] = outString
		case [][]string:
			arrym := fieldValue.([][]string)
			outString := ""
			for i := 0; i < len(arrym); i++ {
				subarray := arrym[i]
				for j := 0; j < len(subarray); j++ {
					outString = outString + subarray[j]
					if j < len(subarray)-1 {
						outString = outString + "^"
					}
				}
				if i < len(arrym)-1 {
					outString = outString + "\\"
				}
			}
			entries[idx] = outString
		case time.Time:
			if fieldValue.(time.Time).IsZero() {
				// dates can be zero = no output
				break
			}
			if isLongDate {
				entries[idx] = fieldValue.(time.Time).In(tz).Format("20060102150405")
			} else {
				entries[idx] = fieldValue.(time.Time).Format("20060102")
			}
		default:
			return errors.New(fmt.Sprintf("Unsupported field type %s", field.Type()))
		}
	}

	output := recordType + "|"
	for i := 0; i <= maxIdx; i++ {
		value := entries[i]
		output = output + value
		if i < maxIdx {
			output = output + "|"
		}
	}
	buffer.Write([]byte(output))
	return nil
}
