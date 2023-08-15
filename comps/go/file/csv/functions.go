package csv

import (
	"encoding/csv"
	"reflect"
	"strconv"
	"strings"
)

type fieldMismatch struct {
	expected, found int
}

func (e *fieldMismatch) Error() string {
	return "CSV line fields mismatch. Expected " + strconv.Itoa(e.expected) + " found " + strconv.Itoa(e.found)
}

type unsupportedType struct {
	Type string
}

func (e *unsupportedType) Error() string {
	return "Unsupported type: " + e.Type
}

// CsvStringToObject convert csv in string to object passed by parameter
func CsvStringToObject(csvstr string, object interface{}) error {

	var reader = csv.NewReader(strings.NewReader(csvstr))
	reader.Comma = ','
	record, err := reader.Read()
	if err != nil {
		return err
	}
	s := reflect.ValueOf(object).Elem()
	if s.NumField() != len(record) {
		return &fieldMismatch{s.NumField(), len(record)}
	}
	for i := 0; i < s.NumField(); i++ {
		file := s.Field(i)
		switch file.Type().String() {
		case "string":
			file.SetString(record[i])
		case "int":
			ival, err := strconv.ParseInt(record[i], 10, 0)
			if err != nil {
				return err
			}
			file.SetInt(ival)
		default:
			return &unsupportedType{file.Type().String()}
		}
	}
	return nil
}
