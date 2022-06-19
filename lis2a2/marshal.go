package lis2a2

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

/** Marshal - wrap datastructure to code
**/
func Marshal(message interface{}, enc Encoding, tz Timezone, notation Notation) ([][]byte, error) {

	location, err := time.LoadLocation(string(tz))
	if err != nil {
		return [][]byte{}, err
	}

	// default delmiters. These will be overwritten by the first occurence of "delimter"-annotation
	repeatDelimiter := "\\"
	componentDelimiter := "^"
	escapeDelimiter := "&"

	buffer, err := iterateStructFieldsAndBuildOutput(message, 1, enc, location, &repeatDelimiter, &componentDelimiter, &escapeDelimiter)

	return buffer, err
}

type OutputRecord struct {
	Field, Repeat, Component int
	Value                    string
}

type OutputRecords []OutputRecord

func iterateStructFieldsAndBuildOutput(message interface{}, depth int, enc Encoding, location *time.Location,
	repeatDelimiter, componentDelimiter, escapeDelimiter *string) ([][]byte, error) {

	generatedSequenceNumber := 1 //TODO: move outside to apply on arrays

	buffer := make([][]byte, 0)

	messageValue := reflect.ValueOf(message)
	messageType := reflect.TypeOf(message)

	for i := 0; i < messageValue.NumField(); i++ {

		currentRecord := messageValue.Field(i)
		recordAstmTag := messageType.Field(i).Tag.Get("astm")
		recordAstmTagsList := strings.Split(recordAstmTag, ",")

		generatedSequenceNumber = 1 // for each record start fresh

		if len(recordAstmTagsList) <= 0 { // nothing anotated, skipping that
			continue
		}

		if len(recordAstmTagsList) == 0 { // no annotation = no record
			// TODO: Descend if its an array or a struct of such
			fmt.Println("Not a struct")
		} else {

			fieldList := make(OutputRecords, 0)

			for i := 0; i < currentRecord.NumField(); i++ {

				field := currentRecord.Field(i)
				fieldAstmTag := currentRecord.Type().Field(i).Tag.Get("astm")
				fieldAstmTagsList := strings.Split(fieldAstmTag, ",")

				fieldIdx, repeatIdx, componentIdx, err := readFieldAddressAnnotation(fieldAstmTagsList[0])
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Invalid field-address. (%s)", err))
				}

				//fmt.Printf("Decode %+v to %d.%d.%d for %s\n", fieldAstmTagsList, fieldIdx, repeatIdx, componentIdx, field.String())

				switch field.Type().Name() {
				case "string":
					value := ""

					if sliceContainsString(fieldAstmTagsList, ANNOTATION_SEQUENCE) {
						return nil, errors.New(fmt.Sprintf("Invalid annotation %s for string-field", ANNOTATION_SEQUENCE))
					}

					// if no delimiters are given, default is \^&
					if sliceContainsString(fieldAstmTagsList, ANNOTATION_DELIMITER) && field.String() == "" {
						value = *repeatDelimiter + *componentDelimiter + *escapeDelimiter
					} else {
						value = field.String()
					}

					fieldList = addASTMFieldToList(fieldList, fieldIdx, repeatIdx, componentIdx, value)
				case "int":
					value := fmt.Sprintf("%d", field.Int())
					if sliceContainsString(fieldAstmTagsList, ANNOTATION_SEQUENCE) {
						value = fmt.Sprintf("%d", generatedSequenceNumber)
						generatedSequenceNumber = generatedSequenceNumber + 1
					}

					fieldList = addASTMFieldToList(fieldList, fieldIdx, repeatIdx, componentIdx, value)
				case "float32":
				case "float64":
				case "Time":
					//t := time.Time(field.Interface())
					//fmt.Println("Time = ", t)
				default:
					return nil, errors.New(fmt.Sprintf("Invalid field type '%s' in struct '%s', input not processed", field.Type().Name(), currentRecord.Type().Name()))
				}

			}

			outp := generateOutputRecord(recordAstmTagsList[0], fieldList, *repeatDelimiter, *componentDelimiter, *escapeDelimiter)
			// fmt.Println(outp)
			buffer = append(buffer, []byte(outp))
		}

	}

	return buffer, nil
}

func addASTMFieldToList(data []OutputRecord, field, repeat, component int, value string) []OutputRecord {

	or := OutputRecord{
		Field:     field,
		Repeat:    repeat,
		Component: component,
		Value:     value,
	}

	data = append(data, or)
	return data
}

// used for sorting
func (or OutputRecords) Len() int { return len(or) }
func (or OutputRecords) Less(i, j int) bool {
	if or[i].Field == or[j].Field {
		if or[i].Repeat == or[j].Repeat {
			return or[i].Component < or[j].Component
		} else {
			return or[i].Repeat < or[j].Repeat
		}
	} else {
		return or[i].Field < or[j].Field
	}
}
func (or OutputRecords) Swap(i, j int) { or[i], or[j] = or[j], or[i] }

/* Converting a list of values (all string already) to the astm format. this funciton works only for one record
   example:
    (0, 0, 2) = first-arr1
    (0, 0, 0) = third-arr1
    (0, 1, 0) = first-arr2
    (0, 1, 1) = second-arr2

	-> .... "|first-arr1^^third-arr1\fist-arr2^second-arr2|"

	returns the full record for output to astm file
*/

func generateOutputRecord(recordtype string, fieldList OutputRecords, REPEAT_DELIMITER, COMPONENT_DELIMITER, ESCAPE_DELMITER string) string {

	output := ""

	sort.Sort(fieldList)

	componentbuffer := make([]string, 100)
	maxComponent := 0

	repeatbuffer := make([]string, 100)
	maxRepeat := 0

	// add a terminator to reduce abortion-spaghetti-code
	fieldList = append(fieldList, OutputRecord{Field: -1})

	fieldGroup := -1 // groupchange on every field-change
	repeatGroup := 0 // groupchange on every repeat-group (see astm-format field,repeat,component,(escape))

	output = output + recordtype + "|" // Record-ID, typical "H", "R", "O", .....

	for _, field := range fieldList {

		fieldGroupBreak := field.Field != fieldGroup && fieldGroup != -1
		repeatGroupBreak := field.Repeat != repeatGroup
		if fieldGroupBreak || repeatGroupBreak {

			buffer := ""
			for c := 0; c <= maxComponent; c++ {
				buffer = buffer + componentbuffer[c]
				if c < maxComponent {
					buffer = buffer + COMPONENT_DELIMITER
				}
			}

			repeatbuffer[repeatGroup] = buffer // sort components to repeatGroup, until no more items, then break

			if fieldGroupBreak { // new field starts = write buffer and empty
				for i := 0; i <= maxRepeat; i++ {
					output = output + repeatbuffer[i]
					if i < maxRepeat {
						output = output + REPEAT_DELIMITER
					}
				}
				output = output + "|"
				maxRepeat = 0
				repeatGroup = 0
			}

			if repeatGroupBreak {
				repeatGroup = field.Repeat
			}

			for c := 0; c < len(componentbuffer); c++ {
				componentbuffer[c] = ""
			}
			maxComponent = 0
			fieldGroup = field.Field
		}

		if fieldGroup == -1 { // starting the very first group in iteration
			fieldGroup = field.Field
		}

		componentbuffer[field.Component] = field.Value
		if field.Component > maxComponent {
			maxComponent = field.Component
		}
		if field.Repeat > maxRepeat {
			maxRepeat = field.Repeat
		}
	}

	return output
}
