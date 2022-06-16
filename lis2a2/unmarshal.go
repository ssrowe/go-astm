package lis2a2

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aglyzov/charmap"
)

func Unmarshal(messageData []byte, targetStruct interface{}, enc Encoding, tz Timezone) error {

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
		return errors.New(fmt.Sprintf("Invalid Codepage Code:%d", enc))
	}

	// first try to break by 0x0a (non-standard, but used sometimes)
	bufferedInputLines := strings.Split(string(messageBytes), string([]byte{0x0A})) // copy
	if len(bufferedInputLines) <= 1 {                                               // if it was not possible to break with non-standard 0x0a line-break try 0d (standard)
		bufferedInputLines = strings.Split(string(messageBytes), string([]byte{0x0D}))
	}
	// strip the remaining 0A and 0D Linefeed at the end
	for i := 0; i < len(bufferedInputLines); i++ {
		// 0d,0a then again as there have been files observed which had 0a0d (0d0a would be normal)
		bufferedInputLines[i] = strings.Trim(bufferedInputLines[i], string([]byte{0x0A}))
		bufferedInputLines[i] = strings.Trim(bufferedInputLines[i], string([]byte{0x0D}))
		bufferedInputLines[i] = strings.Trim(bufferedInputLines[i], string([]byte{0x0A}))
		bufferedInputLines[i] = strings.Trim(bufferedInputLines[i], string([]byte{0x0D}))

	}

	_, _, err = reflectInputToStruct(bufferedInputLines, 1 /*recursion-depth*/, 0 /*current line*/, targetStruct, enc, tz)
	if err != nil {
		return err
	}

	return nil
}

type RETV int

const (
	OK         RETV = 1
	UNEXPECTED RETV = 2 // an exit that wont abort processing. used for skipping optional records
	ERROR      RETV = 3 // a definite error that stops the process

	ANNOTATION_DELIMITER = "delimiter" // annotation that triggers the delimiters in the scanner to be reset
	ANNOTATION_REQUIRED  = "require"   // field-annotation: by default all fields are optinal
	ANNOTATION_OPTIONAL  = "optional"  // record-annotation: by default all records are mandatory
)

var (
	FieldDelimiter     = "|"
	RepeatDelimiter    = "\\"
	ComponentDelimiter = "^"
	EscapeDelimiter    = "&"
)

/* This function takes a string and a struct and matches the annotated fields to the string-input */
func reflectInputToStruct(bufferedInputLines []string, depth int, currentInputLine int, targetStruct interface{}, enc Encoding, tz Timezone) (int, RETV, error) {

	timeLocation, err := time.LoadLocation(string(tz))
	if err != nil {
		return currentInputLine, ERROR, err
	}

	targetStructType := reflect.TypeOf(targetStruct).Elem()
	targetStructValue := reflect.ValueOf(targetStruct).Elem()

	for i := 0; i < targetStructType.NumField(); i++ {

		currentRecord := targetStructValue.Field(i)
		ftype := targetStructType.Field(i)
		astmTag := ftype.Tag.Get("astm")
		astmTagsList := strings.Split(astmTag, ",")

		if len(astmTagsList) < 1 {
			continue // not annotated = no processing
		}

		// no annotation after astm:.. provided means a nested array with more records or ignore
		if len(astmTagsList[0]) < 1 {
			// Not annotated array. If it's a struct have to recurse, otherwise skip
			if targetStructType.Field(i).Type.Kind() == reflect.Slice {

				// Array of Structs
				if reflect.TypeOf(targetStructValue.Interface()).Kind() == reflect.Struct {
					innerStructureType := targetStructType.Field(i).Type.Elem()

					sliceForNestedStructure := reflect.MakeSlice(targetStructType.Field(i).Type, 0, 0)

					for {
						allocatedElement := reflect.New(innerStructureType)
						var err error
						var retv RETV
						currentInputLine, retv, err = reflectInputToStruct(bufferedInputLines, depth+1,
							currentInputLine, allocatedElement.Interface(), enc, tz)
						if err != nil {
							if retv == UNEXPECTED {
								if depth > 0 {
									// if nested structures abort due to unexpected records that does not create an error
									// as the parse will be continued one level higher
									break
								} else {
									return currentInputLine, ERROR, err
								}
							}
							if retv == ERROR { // a serious error ends the processing
								return currentInputLine, ERROR, err
							}
						}

						sliceForNestedStructure = reflect.Append(sliceForNestedStructure, allocatedElement.Elem())
						reflect.ValueOf(targetStruct).Elem().Field(i).Set(sliceForNestedStructure)
					}
					continue
				}
			}
		}

		expectInputRecordType := astmTagsList[0][0] // Expected Record type
		expectedInputRecordTypeOptional := false
		if sliceContainsString(astmTagsList, ANNOTATION_OPTIONAL) {
			expectedInputRecordTypeOptional = true
		}

		if len(bufferedInputLines[currentInputLine]) == 0 {
			continue // empty lines can only be skipped
		}

		// headers require delimiters to be disregarded
		isHeader := false
		if bufferedInputLines[currentInputLine][0] == 'H' {
			isHeader = true
		}

		if expectInputRecordType == bufferedInputLines[currentInputLine][0] {

			//Special case: its not an anotated record, it is an array of annotated records here :
			if currentRecord.Kind() == reflect.Slice {
				innerStructureType := targetStructType.Field(i).Type.Elem()
				sliceForNestedStructure := reflect.MakeSlice(targetStructType.Field(i).Type, 0, 0)
				for { // iterate for as long as the same type repeats
					allocatedElement := reflect.New(innerStructureType)

					if err = reflectAnnotatedFields(bufferedInputLines[currentInputLine], allocatedElement.Elem(), timeLocation, isHeader); err != nil {
						return currentInputLine, ERROR, errors.New(fmt.Sprintf("Failed to process input line '%s' err:%s", bufferedInputLines[currentInputLine], err))
					}

					sliceForNestedStructure = reflect.Append(sliceForNestedStructure, allocatedElement.Elem())
					reflect.ValueOf(targetStruct).Elem().Field(i).Set(sliceForNestedStructure)

					// keep reading while same elements are up
					currentInputLine = currentInputLine + 1
					if expectInputRecordType != bufferedInputLines[currentInputLine][0] {
						break
					}
					if currentInputLine >= len(bufferedInputLines) {
						break
					}
				}

			} else { // The "normal" case: scanning a string into a structure :
				if err = reflectAnnotatedFields(bufferedInputLines[currentInputLine], currentRecord, timeLocation, isHeader); err != nil {
					return currentInputLine, ERROR, errors.New(fmt.Sprintf("Failed to process input line '%s' err:%s", bufferedInputLines[currentInputLine], err))
				}
				currentInputLine = currentInputLine + 1
			}

		} else { // The expected input-record did not occur
			if expectedInputRecordTypeOptional {
				continue // skipping optional record instead of an error
			} else {
				return currentInputLine, UNEXPECTED, errors.New(fmt.Sprintf("Expected Record-Type '%c' input was '%c' in depth (%d) (Abort)", expectInputRecordType, bufferedInputLines[currentInputLine][0], depth))
			}
		}

		if currentInputLine >= len(bufferedInputLines) {
			break
		}
	}

	return currentInputLine, OK, nil
}

func reflectAnnotatedFields(inputStr string, record reflect.Value, timezone *time.Location, isHeader bool) error {

	if reflect.ValueOf(record).Type().Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("invalid type of target: '%s', expecting 'struct'", reflect.ValueOf(record).Type().Kind()))
	}

	inputFields := strings.Split(inputStr, FieldDelimiter)
	if len(inputFields) < 1 {
		return errors.New("Input contains no data")
	}

	for j := 0; j < record.NumField(); j++ {
		recordfield := record.Field(j)
		if !recordfield.CanInterface() {
			return errors.New(fmt.Sprintf("Field %s is not exported - aborting import", recordfield.Type().Name()))
		}
		recordFieldInterface := recordfield.Addr().Interface()

		hasOverrideDelimiterAnnotation := false
		inputIsRequired := false
		astmTag := record.Type().Field(j).Tag.Get("astm")
		if astmTag == "" {
			continue // nothing to process when someone requires astm:
		}
		astmTagsList := strings.Split(astmTag, ",")
		for i := 0; i < len(astmTagsList); i++ {
			astmTagsList[i] = strings.Trim(astmTagsList[i], " ")
		}
		if sliceContainsString(astmTagsList, ANNOTATION_DELIMITER) {
			// the delimiter is instantly replaced with the delimiters from the file for further parsing. By default that is "\^&"
			hasOverrideDelimiterAnnotation = true
		}
		if sliceContainsString(astmTagsList, ANNOTATION_REQUIRED) {
			inputIsRequired = true
		}
		currentInputFieldNo, repeat, component, err := readFieldAddressAnnotation(astmTagsList[0])
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid annotation for field %s. (%s)", record.Type().Field(j).Name, err))
		}
		if currentInputFieldNo >= len(inputFields) || currentInputFieldNo < 0 {
			//TODO: user should be able to toggle wether he wants an exact match = error or bestfit = skip silent
			continue // mapped field is beyond the data
		}

		switch reflect.TypeOf(recordfield.Interface()).Name() {
		case "string":
			if value, err := extractAstmFieldByRepeatAndComponent(inputFields[currentInputFieldNo], repeat, component); err == nil {

				// in headers there can be special characters, that is why the value needs to disregard the delimiters:
				if isHeader {
					value = inputFields[currentInputFieldNo]
				}

				reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(value))

				if hasOverrideDelimiterAnnotation { // the first three characters become the new delimiters
					if len(value) >= 1 {
						RepeatDelimiter = value[0:1]
					}
					if len(value) >= 2 {
						ComponentDelimiter = value[1:2]
					}
					if len(value) >= 3 {
						EscapeDelimiter = value[2:3]
					}
				}
			} else {
				if inputIsRequired { // by default we ignore missing input
					return errors.New(fmt.Sprintf("Failed to extract index (%d.%d.%d) from input line '%s' : (%s)",
						currentInputFieldNo+1, repeat+1, component+1, inputStr, err))
				}
			}
		case "int":
			if hasOverrideDelimiterAnnotation {
				return errors.New("delimiter-annotation is only allowed for string-type, not int.")
			}

			if value, err := extractAstmFieldByRepeatAndComponent(inputFields[currentInputFieldNo], repeat, component); err == nil {

				if num, err := strconv.Atoi(value); err == nil {
					reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(num))
				} else {
					if inputIsRequired { // by default we ignore missing input
						return errors.New(fmt.Sprintf("Failed to extract index (%d,%d) from field %s(%s)", repeat, component, inputFields[currentInputFieldNo], err))
					}
				}

			} else {
				return err
			}
		case "float32":
			if hasOverrideDelimiterAnnotation {
				return errors.New("delimiter-annotation is only allowed for string-type, not int.")
			}

			if value, err := extractAstmFieldByRepeatAndComponent(inputFields[currentInputFieldNo], repeat, component); err == nil {

				if num, err := strconv.ParseFloat(value, 32); err == nil {
					reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(float32(num)))
				} else {
					if inputIsRequired { // by default we ignore missing input
						return errors.New(fmt.Sprintf("Failed to extract index (%d,%d) from field %s(%s)", repeat, component, inputFields[currentInputFieldNo], err))
					}
				}

			} else {
				return err
			}
		case "float64":
			if hasOverrideDelimiterAnnotation {
				return errors.New("delimiter-annotation is only allowed for string-type, not int.")
			}

			if value, err := extractAstmFieldByRepeatAndComponent(inputFields[currentInputFieldNo], repeat, component); err == nil {

				if num, err := strconv.ParseFloat(value, 64); err == nil {
					reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(float64(num)))
				} else {
					if inputIsRequired { // by default we ignore missing input
						return errors.New(fmt.Sprintf("Failed to extract index (%d,%d) from field %s(%s)", repeat, component, inputFields[currentInputFieldNo], err))
					}
				}

			} else {
				return err
			}

			/*
				TODO: this annotation got removed because it doesnt help to have open arrays
				case "slice":
				 instr := fields[mapFieldNo]
				list := splitAny(instr, RepeatDelimiter) //CHANGEHERE
				field.Set(reflect.ValueOf(list))
				/*	case [][]string:
					fieldFromFile := fields[mapFieldNo]
					// the amount of repeat-separators is the first dimension, then each repeats the patters
					arry := make([][]string, 0)
					sequences := strings.Split(fieldFromFile, "\\")
					for _, sequence := range sequences {
						data := strings.Split(sequence, "^")
						arry = append(arry, data)
					}
					field.Set(reflect.ValueOf(arry)) */
		case "Time":
			if hasOverrideDelimiterAnnotation {
				return errors.New("delimiter-annotation is only allowed for string-type, not Time")
			}
			inputFieldValue := inputFields[currentInputFieldNo]
			if inputFieldValue == "" {
				reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(time.Time{}))
			} else if len(inputFieldValue) == 8 { // YYYYMMDD See Section 5.6.2 https://samson-rus.com/wp-content/files/LIS2-A2.pdf
				timeInLocation, err := time.ParseInLocation("20060102", inputFieldValue, timezone)
				if err != nil {
					return errors.New(fmt.Sprintf("Invalid time format <%s>", inputFieldValue))
				}
				reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(timeInLocation))

			} else if len(inputFieldValue) == 14 { // YYYYMMDDHHMMSS
				timeInLocation, err := time.ParseInLocation("20060102150405", inputFieldValue, timezone)
				if err != nil {
					return errors.New(fmt.Sprintf("Invalid time format <%s>", inputFieldValue))
				}
				reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(timeInLocation.UTC()))
			} else {
				return errors.New(fmt.Sprintf("Unrecognized time format <%s>", inputFieldValue))
			}
		default:
			return errors.New(fmt.Sprintf("Invalid field-Type '%s' for mapping (not implemented): '%s'",
				reflect.TypeOf(recordfield.Interface()).Kind(), inputStr))
		}
	}

	return nil
}

// Translating the annotation of a field to field, index/repeat, component
// Input of one value : e.g."4" -> field -> 4
// Input of two values :"4.2" -> field, compoennt -> 4,1,2
// Input of three values "4.1.1" -> field, repeat, component -> 4,1,1
// "whereas field indexes should be 1-99 (check plz)
func readFieldAddressAnnotation(annotation string) (field int, repeat int, component int, err error) {

	if annotation == "" { // no annotation will always return the first of everything
		return 0, 0, 0, nil
	}
	field = 1
	repeat = 1
	component = 1
	fieldSplitted := strings.Split(annotation, ".")

	if len(fieldSplitted) >= 1 {
		if field, err = strconv.Atoi(fieldSplitted[0]); err != nil {
			return 0, 0, 0, err
		}
	}
	if len(fieldSplitted) >= 2 {
		if component, err = strconv.Atoi(fieldSplitted[1]); err != nil {
			return 0, 0, 0, err
		}
	}
	if len(fieldSplitted) >= 3 {
		if repeat, err = strconv.Atoi(fieldSplitted[1]); err != nil {
			return 0, 0, 0, err
		}
		if component, err = strconv.Atoi(fieldSplitted[2]); err != nil {
			return 0, 0, 0, err
		}
	}

	return field - 1, repeat - 1, component - 1, nil
}

// input is an unpacked field from an astm-file free of the field delimiter ("|")
// this function ettracts the field by repeat and component-delimiter
func extractAstmFieldByRepeatAndComponent(text string, repeat int, component int) (string, error) {

	subfield := strings.Split(text, RepeatDelimiter)
	if repeat >= len(subfield) {
		return "", errors.New(fmt.Sprintf("Index (%d, %d) out of bounds '%s', delimiter '%s'", repeat, component, text, RepeatDelimiter))
	}

	subsubfield := strings.Split(subfield[repeat], ComponentDelimiter)
	if component > len(subsubfield) || component < 0 {
		return "", errors.New(fmt.Sprintf("Index (%d, %d) out of bounds '%s' delimiter '%s'", repeat, component, text, ComponentDelimiter))
	}

	if component >= len(subsubfield) {
		return "", errors.New(fmt.Sprintf("Index (%d, %d) out of bounds '%s', delimiter '%s'", repeat, component, text, RepeatDelimiter))
	}

	return subsubfield[component], nil
}

func sliceContainsString(list []string, search string) bool {
	for _, x := range list {
		if x == search {
			return true
		}
	}
	return false
}
