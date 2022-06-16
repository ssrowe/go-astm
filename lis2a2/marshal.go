package lis2a2

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Marshal(message interface{}, enc Encoding, tz Timezone, linebreak LineBreak) ([]byte, error) {

	location, err := time.LoadLocation(string(tz))
	if err != nil {
		return []byte{}, err
	}

	buffer, err := digIn(message, 1, enc, location, linebreak)

	return buffer, err
}

func digIn(message interface{}, depth int, enc Encoding, location *time.Location, linebreak LineBreak) ([]byte, error) {

	//repeatDelimiter := "\\"
	//componentDelimiter := "^"
	//escapeDelimiter := "&"

	buffer := make([]byte, 0)
	// targetStructType := reflect.TypeOf(message).Elem()
	//targetStructValue := reflect.ValueOf(message).Elem()

	messageType := reflect.TypeOf(message)
	//fmt.Printf("\n\ndigging into structure %+v\n", messageType)

	for i := 0; i < messageType.NumField(); i++ {

		currentRecordType := messageType.Field(i)
		currentRecordValue := reflect.ValueOf(message).Field(i)
		fmt.Println("Field...", currentRecordType)
		fmt.Println("     ...", currentRecordValue)
		//ftype := targetStructType.Field(i)
		astmTag := currentRecordType.Tag.Get("astm")
		astmTagsList := strings.Split(astmTag, ",")

		if len(astmTagsList) <= 0 { // nothing anotated, skipping that
			continue
		}

		fmt.Println("Dealing with ", astmTagsList[0])

		if currentRecordType.Type.Kind() == reflect.Struct {
			fmt.Println("Its a struct!")
			//--------------------------------------------------------------------------------here it sucks
			for i := 0; i < currentRecordType.Type.NumField(); i++ {
				field := currentRecordType.Type.Field(i)
				value := reflect.ValueOf(message).Field(i)
				fmt.Printf("%+v = %+v\n", field.Name, value)
			}

		}
		//recordType := string(astmTagsList[0][0])
		/*		if oneRec, err := encodeStruct(recordType,
					messageType.Field(i),
					repeatDelimiter, componentDelimiter, escapeDelimiter); err == nil {
					for _, c := range oneRec {
						buffer = append(buffer, c)
					}
					fmt.Println("Havving produced : ", string(oneRec))
				} else {
					fmt.Println("Failed here with ", err)
				}*/
	}

	return buffer, nil
}

func encodeStruct(recordType string, record interface{}, repeatDelimiter, componentDelimiter, escapeDelimiter string) ([]byte, error) {

	buffer := make([]byte, 0)
	buffer = append(buffer, byte(recordType[0]))
	buffer = append(buffer, byte('|'))

	messageType := reflect.TypeOf(record).Elem()
	// messageValue := reflect.ValueOf(record).Elem()
	fmt.Printf("Fuckingshit %+v\n", messageType)
	os.Exit(-1)
	fmt.Printf("da struct: %+v\n", messageType)
	for i := 0; i < messageType.NumField(); i++ {
		fmt.Printf("Some field.... %s \n", "some")
	}

	return buffer, nil
}

func MarshalOlde(message interface{}, enc Encoding, tz Timezone) ([]byte, error) {
	return nil, nil
	/*var buffer bytes.Buffer

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

	return buffer.Bytes(), nil */
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
