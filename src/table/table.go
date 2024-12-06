package table

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/olekukonko/tablewriter"
)

type Options struct {
	Headers  []string
	Border   bool
	Centered bool
}

type Data struct {
	Headers []string
	Rows    [][]string
}

func NewData(headers []string, rows [][]string) Data {
	return Data{
		Headers: parseHeaders(headers),
		Rows:    rows,
	}
}

func PrintTable(data any, opts Options) error {
	table := tablewriter.NewWriter(os.Stdout)
	configureTable(table, &opts)

	tableData, err := parseData(data, opts.Headers)

	if err != nil {
		return err
	}

	table.SetHeader(tableData.Headers)
	table.SetRowLine(true)
	table.AppendBulk(tableData.Rows)
	table.Render()
	return nil
}

// configureTable applies the provided TableOptions to configure the table's appearance.
// It sets up borders, column separators, text alignment, and header colors based on the options.
func configureTable(table *tablewriter.Table, options *Options) {
	table.SetBorder(options.Border)
	table.SetColumnSeparator(" ")

	if options.Centered {
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	}
}

// parseData converts the input data into a table-friendly format.
// It accepts any slice of structs or slice of slices and optional headers.
// Returns a Data struct containing the parsed headers and rows, or an error if the input is invalid.
func parseData(data any, headers []string) (Data, error) {
	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Slice {
		return Data{}, fmt.Errorf("data is not a slice")
	}

	if v.Len() == 0 {
		return Data{}, nil
	}

	first := v.Index(0)
	switch first.Kind() {
	case reflect.Struct:
		return parseStructSlice(v, headers), nil
	case reflect.Slice:
		return parseSliceOfSlices(v, headers), nil
	default:
		return Data{}, fmt.Errorf("unsupported data type: %v", first.Kind())
	}
}

// parseStructSlice converts a slice of structs into table data.
// It extracts field names as headers (if not provided) and field values as rows.
// The struct fields are converted to strings and organized into rows.
func parseStructSlice(structs reflect.Value, headers []string) Data {
	structType := structs.Index(0).Type()

	var finalHeaders []string
	var rows [][]string

	if len(headers) == 0 {
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			finalHeaders = append(finalHeaders, field.Name)
		}
	} else {
		finalHeaders = parseHeaders(headers)
	}

	for i := 0; i < structs.Len(); i++ {
		item := structs.Index(i)
		fieldCnt := item.NumField()
		row := make([]string, fieldCnt)

		for j := 0; j < fieldCnt; j++ {
			field := item.Field(j)
			row[j] = toString(field.Interface())
		}
		rows = append(rows, row)
	}

	return NewData(finalHeaders, rows)
}

// parseSliceOfSlices converts a slice of slices into table data.
// The first row is treated as headers (if not provided) and subsequent rows as data.
// All values are converted to strings.
func parseSliceOfSlices(slices reflect.Value, headers []string) Data {
	var rows [][]string
	var finalHeaders []string

	if len(headers) == 0 && slices.Len() > 0 {
		first := slices.Index(0)
		for i := 0; i < first.Len(); i++ {
			h := first.Index(i)
			finalHeaders = append(finalHeaders, toString(h.Interface()))
		}
	} else {
		finalHeaders = parseHeaders(headers)
	}

	for i := 1; i < slices.Len(); i++ {
		item := slices.Index(i)
		row := make([]string, item.Len())
		for j := 0; j < item.Len(); j++ {
			row[j] = toString(item.Index(j).Interface())
		}
		rows = append(rows, row)
	}

	return NewData(finalHeaders, rows)
}

// parseHeaders processes the header strings by converting them to title case
// and replacing spaces with underscores.
func parseHeaders(headers []string) []string {
	parsed := make([]string, len(headers))
	for i, header := range headers {
		s := strings.ToLower(header)
		s = cases.Title(language.English).String(s)
		s = strings.ReplaceAll(s, " ", "_")
		parsed[i] = s
	}
	return parsed
}

// toString converts any value to its string representation.
// Returns an empty string for nil values.
func toString(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
