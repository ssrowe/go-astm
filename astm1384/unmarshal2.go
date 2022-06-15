package astm1384

import (
	"errors"
	"fmt"
	"github.com/aglyzov/charmap"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Unmarshal2(messageData []byte, target interface{}, enc Encoding, tz Timezone, pv ProtocolVersion) error {

	var (
		messageBytes []byte
		err          error
	)
	switch enc {
	case EncodingUTF8:
		// do nothing, this is correct
		messageBytes = messageData
	case EncodingASCII:
		messageBytes = messageData
	case EncodingDOS866:
		messageBytes, err = charmap.ANY_to_UTF8(messageData, "DOS866")
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid input : %s", err))
		}
	case EncodingDOS855:
		messageBytes, err = charmap.ANY_to_UTF8(messageData, "DOS855")
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid input : %s", err))
		}
	case EncodingDOS852:
		messageBytes, err = charmap.ANY_to_UTF8(messageData, "DOS852")
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid input : %s", err))
		}
	case EncodingWindows1250:
		messageBytes, err = charmap.ANY_to_UTF8(messageData, "CP1250")
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid input : %s", err))
		}
	case EncodingWindows1251:
		messageBytes, err = charmap.ANY_to_UTF8(messageData, "CP1251")
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid input : %s", err))
		}
	case EncodingWindows1252:
		messageBytes, err = charmap.ANY_to_UTF8(messageData, "CP1252")
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid input : %s", err))
		}
	default:
		return errors.New(fmt.Sprintf("Invalid Codepage %d", enc))
	}

	// currentLine := 0
	buffer := strings.Split(string(messageBytes), string([]byte{0x0A})) // copy

	// pre cautiously strip the 0A Linefeed
	for i := 0; i < len(buffer); i++ {
		buffer[i] = strings.Trim(buffer[i], string([]byte{0x0D}))
	}

	_, _, err = seqScan(buffer, 1 /*recursion-depth*/, 0, target, enc, tz, pv)
	if err != nil {
		return err
	}

	return nil
}

type RETV int

const (
	OK         RETV = 1
	UNEXPECTED RETV = 2
	ERROR      RETV = 3

	delimiter = "delimiter"
)

var (
	FieldDelimiter     = "|"
	RepeatDelimiter    = "\\"
	ComponentDelimiter = "^"
	EscapeDelimiter    = "&"
)

// TODO: Not working fully atm
/* Scan Structure recursive. Note there are only 10b type of people: those that understand recursions, and those who dont */
func seqScan(buffer []string, depth int, currentLine int, target interface{}, enc Encoding, tz Timezone, pv ProtocolVersion) (int, RETV, error) {

	outerStructureType := reflect.TypeOf(target).Elem()

	fmt.Printf("seqScan (%s, %d)\n", reflect.TypeOf(target).Name(), depth)

	for i := 0; i < outerStructureType.NumField(); i++ {

		astmTag := outerStructureType.Field(i).Tag.Get("astm")
		astmTagsList := strings.Split(astmTag, ",")

		if len(astmTagsList) < 1 {
			continue // not annotated
		}

		// no tags provided means a nested array with more records or ignore
		if len(astmTagsList[0]) < 1 {
			// Not annotated array. If it's a struct have to recurse, otherwise skip
			if outerStructureType.Field(i).Type.Kind() == reflect.Slice {

				// What is the type of the Slice? (Struct or string ?)
				sliceFieldType := reflect.TypeOf(outerStructureType.Field(i))

				// Array of Structs
				if sliceFieldType.Kind() == reflect.Struct {
					innerStructureType := outerStructureType.Field(i).Type.Elem()

					sliceForNestedStructure := reflect.MakeSlice(outerStructureType.Field(i).Type, 0, 0)

					for {
						allocatedElement := reflect.New(innerStructureType)
						var err error
						var retv RETV
						currentLine, retv, err = seqScan(buffer, depth+1, currentLine, allocatedElement.Interface(), enc, tz, pv)
						if err != nil {
							if retv == UNEXPECTED {
								if depth > 0 {
									// if nested structures abort due to unexpected records that does not create an error
									// as the parse will be continued one level higher
									break
								} else {
									return currentLine, ERROR, err
								}
							}
						}

						sliceForNestedStructure = reflect.Append(sliceForNestedStructure, allocatedElement.Elem())
						reflect.ValueOf(target).Elem().Field(i).Set(sliceForNestedStructure)
					}
					continue
				}
			}
		}

		expectRecordType := astmTagsList[0][0] // Expected Record type

		optional := false
		if contains(astmTagsList, "optional") {
			optional = true
		}

		if expectRecordType == buffer[currentLine][0] {
			location, err := time.LoadLocation(string(tz))
			if err != nil {
				return currentLine, ERROR, err
			}

			err = MapRecordFromString(expectRecordType, buffer[currentLine], outerStructureType.Field(i), UseDelimiter, location)
			if err != nil {
				return currentLine, ERROR, err
			}
			currentLine = currentLine + 1
		} else {
			if optional {
				fmt.Printf("Skipped optional %c\n", expectRecordType)
				continue
			} else {
				return currentLine, UNEXPECTED, errors.New(fmt.Sprintf("Expected Record-Type '%c' input was '%c' in depth (%d) (Abort)", expectRecordType, buffer[currentLine][0], depth))
			}
		}

		if currentLine >= len(buffer) {
			break
		}
	}

	return currentLine, OK, nil
}

// possible inputs:
// "4"
// "4.1"
// "4.1.1"
// "whereas field indexes should be 1-99 (check plz)

func readFieldAddressAnnotation(annotation string) (field, repeat, component int, error) {
	fieldSplitted := strings.Split(annotation, ".")

	if len(fieldSplitted) > 3 {
		return 0, 0, 0, errors.New("invalid annotation")
	}

	return 1, 1, 1, nil
}

func MapRecordFromString(recordType byte, inputStr string, target interface{}, useDelimiter string, timezone *time.Location) error {
	if reflect.ValueOf(target).Type().Kind() != reflect.Struct {
		return errors.New("invalid type of target")
	}

	fields := strings.Split(inputStr, "|")
	if len(fields) < 1 {
		return errors.New("invalid length of input string") //TODO: Maybe other message
	}

	t := reflect.TypeOf(target).Elem()

	for i := 0; i < t.NumField(); i++ {
		overrideDelimter := false
		astmTag := t.Field(i).Tag.Get("astm")
		astmTagsList := strings.Split(astmTag, ",")

		if len(astmTagsList) == 0 || astmTag == "" {
			continue // nothing to process when someone requires astm:
		}

		if contains(astmTagsList, delimiter) {
			overrideDelimter = true
		}

		//mapFieldNo, err := strconv.Atoi(astmTagsList[0]) // just a number
		mapFieldNo, index, component, err := readFieldAddressAnnotation(astmTagsList[0])
		if err != nil {
			return err
		}

		if mapFieldNo >= len(fields) {
			continue // mapped field is beyond the data
		}

		// gleich hier weiter

		field := reflect.ValueOf(target).Elem().Field(i)
		fieldValue := field.Interface()

		switch fieldValue.(type) {
		case string:
			// string wasmichwirklich interessiert = a^b^\1^2^3
			// TODO: Implement Delimiter override!!
			if len(astmTagsList) > 1 {
				// further subdivide like this part "|^^^MO10^^28343^|"
				subFields := strings.Split(fields[mapFieldNo], ComponentDelimiter)
				subFieldNo, err := strconv.Atoi(astmTagsList[1])
				if err != nil {
					return errors.New(fmt.Sprintf("Invalid annotation astm:%s. %s", astmTag, err))
				}
				if len(subFields) >= subFieldNo && subFieldNo >= 1 {
					field.SetString(subFields[subFieldNo-1])
				} else {
					// when fields are not present they just dont get mapped = skipping here
				}
			} else {
				field.SetString(fields[mapFieldNo])
			}
		case int:
			num, err := strconv.Atoi(fields[mapFieldNo])
			if err != nil {
				return err
			}
			field.SetInt(int64(num))
		case []string:
			instr := fields[mapFieldNo]
			list := splitAny(instr, useDelimiter)
			field.Set(reflect.ValueOf(list))
		case [][]string:
			fieldFromFile := fields[mapFieldNo]
			// the amount of repeat-separators is the first dimension, then each repeats the patters
			arry := make([][]string, 0)
			sequences := strings.Split(fieldFromFile, "\\")
			for _, sequence := range sequences {
				data := strings.Split(sequence, "^")
				arry = append(arry, data)
			}
			field.Set(reflect.ValueOf(arry))
		case time.Time:
			instr := fields[mapFieldNo]
			if instr == "" {
				field.Set(reflect.ValueOf(time.Time{}))
			} else if len(instr) == 8 { // YYYYMMDD See Section 5.6.2 https://samson-rus.com/wp-content/files/LIS2-A2.pdf
				timeLocated, err := time.ParseInLocation("20060102", instr, timezone)
				if err != nil {
					return errors.New(fmt.Sprintf("Invalid time format <%s>", instr))
				}
				field.Set(reflect.ValueOf(timeLocated))
			} else if len(instr) == 14 { // YYYYMMDDHHMMSS
				timeLocated, err := time.ParseInLocation("20060102150405", instr, timezone)
				if err != nil {
					return errors.New(fmt.Sprintf("Invalid time format <%s>", instr))
				}
				field.Set(reflect.ValueOf(timeLocated.UTC()))
			} else {
				return errors.New(fmt.Sprintf("Unrecognized time format <%s>", instr))
			}
		default:
			return errors.New(fmt.Sprintf("Invalid field-Type '%s' for mapping (not implemented)", t.Field(i).Type))
		}
	}

	fmt.Printf("DECODE %s : %s\n, %s \n", string(recordType), inputStr, reflect.ValueOf(target).Type().Kind())
	return nil
}

func Scan(input string, target interface{}, expectRecord byte, optional bool) error {
	if len(input) < 1 {
		return errors.New(fmt.Sprintf("Empty input. Excpected Record:%c", expectRecord))
	}
	if input[0] != expectRecord {
		return errors.New(fmt.Sprintf("Input Stream presentet Record type '%c' but expected '%c'", input[0], expectRecord))
	}
	fmt.Println("Scanning: ", input)
	return nil
}

func contains(list []string, search string) bool {
	for _, x := range list {
		if x == search {
			return true
		}
	}
	return false
}
