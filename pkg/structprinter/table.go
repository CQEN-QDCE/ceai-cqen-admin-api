package structprinter

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
)

const MORE_DETAILS_NOTICE = "Le champs %v contient plus de détails en utilisant les sorties (-o) json ou yaml."
const FIELD_TOO_LONG_NOTICE = "Le champs %v a été tronqué, pour obtenir la valeur complète utiliser les sorties (-o) json ou yaml. "

const COL_MAX_LENGTH = 100

type Table struct {
	Writer  *tabwriter.Writer
	Notices []string
}

func NewTable() *Table {
	return &Table{
		Writer: tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0),
	}
}

func (t *Table) AddNotice(notice string) {
	for _, n := range t.Notices {
		if n == notice {
			return
		}
	}

	t.Notices = append(t.Notices, notice)
}

func (t *Table) PrintNotice() {
	if len(t.Notices) > 0 {
		fmt.Fprint(t.Writer, "\nNote(s): \n")

		for _, notice := range t.Notices {
			fmt.Fprint(t.Writer, notice)
			fmt.Fprint(t.Writer, "\n")
		}
	}
}

func (t *Table) PrintTableHeaderLabels(data interface{}) (headerString string, separatorString string) {
	v := reflect.ValueOf(data)

	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Type.Kind() == reflect.Struct {
			//Handle struct inheritance (recursive)
			subHeader, subSeparator := t.PrintTableHeaderLabels(v.Field(i).Interface())
			headerString = headerString + subHeader
			separatorString = separatorString + subSeparator
		} else {
			headerString = headerString + v.Type().Field(i).Name
			separatorString = separatorString + strings.Repeat("-", len(v.Type().Field(i).Name))
		}

		if i+1 < v.NumField() {
			headerString = headerString + "\t"
			separatorString = separatorString + "\t"
		}
	}

	return headerString, separatorString
}

func (t *Table) PrintTableHeader(data interface{}) {
	//Transform header field to tab separated string
	var header bytes.Buffer
	var separator bytes.Buffer

	headerString, separatorString := t.PrintTableHeaderLabels(data)

	header.WriteString(headerString)
	separator.WriteString(separatorString)

	header.WriteString("\n")
	separator.WriteString("\n")

	fmt.Fprint(t.Writer, header.String())
	fmt.Fprint(t.Writer, separator.String())
}

func (t *Table) PrintTableLineValues(data interface{}) (valueString string) {
	v := reflect.ValueOf(data)

	for i := 0; i < v.NumField(); i++ {
		//Extract value if field is pointer
		var field reflect.Value

		if v.Field(i).Kind() == reflect.Ptr {
			field = v.Field(i).Elem()
		} else {
			field = v.Field(i)
		}

		//Print field value
		if field.Kind() == reflect.Struct {
			//Print time struct as is
			if field.Type().String() == "time.Time" {
				valueString = valueString + fmt.Sprintf("%.19s", field.Interface())
			} else {
				//Handle struct inheritance (recursive)
				valueString = valueString + t.PrintTableLineValues(field.Interface())
			}
		} else if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
			valueString = valueString + fmt.Sprint(field.Len())
			t.AddNotice(fmt.Sprintf(MORE_DETAILS_NOTICE, v.Type().Field(i).Name))
		} else if field.IsValid() {
			fieldValue := fmt.Sprint(field.Interface())

			if len(fieldValue) > 50 {
				valueString = valueString + fmt.Sprintf("%.47s...", fieldValue)
				t.AddNotice(fmt.Sprintf(FIELD_TOO_LONG_NOTICE, v.Type().Field(i).Name))
			} else {
				valueString = valueString + fieldValue
			}
		} else {
			valueString = valueString + "-"
		}

		if i+1 < v.NumField() {
			valueString = valueString + "\t"
		}
	}

	return valueString
}

func (t *Table) PrintTableLine(data interface{}) {
	var values bytes.Buffer

	values.WriteString(t.PrintTableLineValues(data))

	values.WriteString("\n")

	fmt.Fprint(t.Writer, values.String())
}

func (t *Table) PrintTable(data interface{}) error {
	//Test array to be an array
	reflectValue := reflect.ValueOf(data)

	if reflectValue.Kind() == reflect.Array || reflectValue.Kind() == reflect.Slice {
		t.PrintTableHeader(reflectValue.Index(0).Interface())

		for i := 0; i < reflectValue.Len(); i++ {
			t.PrintTableLine(reflectValue.Index(i).Interface())
		}
	} else if reflectValue.Kind() == reflect.Struct {
		t.PrintTableHeader(data)
		t.PrintTableLine(data)
	} else {
		return fmt.Errorf("PrintTable error: Unsupported data type")
	}

	t.PrintNotice()

	return nil
}

func (t *Table) OutputTable() error {
	err := t.Writer.Flush()

	if err != nil {
		return err
	}

	fmt.Println()

	return nil
}

func PrintTable(data interface{}) error {
	table := NewTable()

	err := table.PrintTable(data)

	if err != nil {
		return err
	}

	return table.OutputTable()
}
