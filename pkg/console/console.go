package console

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// header defines the name and width of a column displayed on screen.
type header struct {
	Title  string // Title is the name of the column as displayed on screen.
	Length int    // Length is the fixed width of the column (e.g. spacing)
}

// row is a row within the console table output.
type row map[string]string

// headers is a collection of header objects with special functions for Add-ing and Merge-ing.
type headers []header

// Add adds a new header hr to the array if a header of the same title does not exist.
// If a header with the same title already exists in the list then it will be
// replaced if the length of the new header (hr) is bigger than the current one that exists.
func (h headers) Add(hr header) headers {
	for i, v := range h {
		if v.Title == hr.Title {
			if v.Length < hr.Length {
				h[i] = hr
				return h
			}
			return h
		}
	}

	return append(h, hr)
}

// Merge adds multiple new headers hs to the array if headers of the same title do not exist.
// If one or more headers already exist with the same title then they will be replaced if the
// length of the new header being added is bigger than the one already existing.
func (h headers) Merge(hs headers) headers {
	nh := h
	for _, v := range hs {
		nh = nh.Add(v)
	}
	return nh
}

// rowFromMap reads a given input which is assumed to be a map of some kind and produces
// headers and row objects. The headers will be the keys from the map converted to strings
// and the row objects will be the key value pairs with both keys and values being strings.
func rowFromMap(input interface{}) (headers, row) {
	hs := make(headers, 0)
	r := make(row)

	// the input can either be of type reflect.Value or
	// an actual map type. Check which is the case here
	var val reflect.Value
	if v, ok := input.(reflect.Value); ok {
		val = v
	} else {
		val = reflect.ValueOf(input)
	}

	// iterate over the map keys for the value.
	for _, k := range val.MapKeys() {
		// convert the keys to an uppercase string and create a header
		// object from each key.
		h := strings.ToUpper(fmt.Sprintf("%v", k))
		hs = hs.Add(header{Title: h, Length: len(h) + 5})

		// convert the value to a string representation and check the
		// length of the string against the header to determine whether
		// the column needs to be widened.
		v := val.MapIndex(k)
		r[h] = fmt.Sprintf("%v", v)

		// if the column needs to be widened then update the header object.
		if len(r[h]) > len(h) {
			hs = hs.Add(header{Title: h, Length: len(r[h]) + 5})
		}
	}
	return hs, r
}

// rowFromStruct reads a given input which is assumed to be a struct of some kind and produces
// headers and row objects. The headers will be the name of the fields within the struct or
// the 'console' tag associated with the field. The row objects will be the key value pairs with
// both keys and values being converted into string representations.
func rowFromStruct(input interface{}) (headers, row) {
	hs := make(headers, 0)
	r := make(row)

	// the input can either be of type reflect.Value or
	// an actual map type. Check which is the case here
	var val reflect.Value
	if v, ok := input.(reflect.Value); ok {
		val = v
	} else {
		val = reflect.ValueOf(input)
	}

	// iterate over the fields available in the type.
	for i := 0; i < val.Type().NumField(); i++ {
		// extract the field name and convert it to uppercase
		n := val.Type().Field(i).Name
		h := strings.ToUpper(n)

		// create a header object with the uppercase title and a length
		// 5 spaces bigger than the title for spacing.
		hs = hs.Add(header{Title: h, Length: len(h) + 5})

		// select the field from the value and print it in string representation
		v := val.FieldByName(n)
		r[h] = fmt.Sprintf("%v", v)

		// if the column needs to be widened then update the header object.
		if len(r[h]) > len(h) {
			hs = hs.Add(header{Title: h, Length: len(r[h]) + 5})
		}
	}
	return hs, r

}

// rowsFromSlice reads a given input which is assumed to be a slice and produce headers and
// row objects. The headers will be the name of all fields extracted from the items within the slice
// and the row array will be an array of the key value pairs with both keys and values converted into
// string representations.
func rowsFromSlice(input interface{}) (headers, []row) {
	hs := make(headers, 0)
	rs := make([]row, 0)

	// iterate over every item in the slice
	for i := 0; i < reflect.ValueOf(input).Len(); i++ {
		// depending on the type of value then apply the appropriate logic.
		switch reflect.ValueOf(input).Index(i).Kind() {
		case reflect.Map:

			fmt.Println("struct")
			// if the slice/array item is a map then apply the rowFromMap function
			// to extract a row and the corresponding headers.
			h, r := rowFromMap(reflect.ValueOf(input).Index(i))
			rs = append(rs, r)

			// merge the headers together
			hs = hs.Merge(h)
		case reflect.Struct:

			// if the slice/array item is a struct then apply the rowFromStruct function
			// to extract a row and the corresponding headers.
			h, r := rowFromStruct(reflect.ValueOf(input).Index(i))
			rs = append(rs, r)

			// merge the headers together
			hs = hs.Merge(h)

		case reflect.Interface:
			// if the slice/array item is a interface then apply the rowFromStruct function
			// to extract a row and the corresponding headers.
			h, r := rowFromStruct(reflect.ValueOf(input).Index(i).Elem())
			rs = append(rs, r)

			// merge the headers together
			hs = hs.Merge(h)

		default:
			fmt.Println(reflect.ValueOf(input).Index(i).Kind())
		}
	}

	return hs, rs
}

// Marshal takes an object of any type and attempts to turn it into a console table
// represented as a []byte. Supported object types are maps, slices, arrays, structs and
// pointers (assuming the pointer is to a struct, slice, map or array.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	switch reflect.TypeOf(v).Kind() {
	case reflect.Map:
		headers, r := rowFromMap(v)
		err = printRows(headers, []row{r}, &buf)
	case reflect.Slice:
		headers, rs := rowsFromSlice(v)
		err = printRows(headers, rs, &buf)
	case reflect.Array:
		headers, rs := rowsFromSlice(v)
		err = printRows(headers, rs, &buf)
	case reflect.Struct:
		headers, rs := rowFromStruct(v)
		err = printRows(headers, []row{rs}, &buf)
	case reflect.Ptr:
		return Marshal(reflect.ValueOf(v).Elem().Interface())
	default:
		return nil, fmt.Errorf("unable to marshall type")
	}

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// printRows takes a headers object, an array of row objects and an io.Writer and attempts
// to print the console table to the specified io.Writer. If any part of this fails then an
// error is returned.
func printRows(hs headers, rs []row, buf io.Writer) error {
	// print the headers by iterating over the headers hs object
	// and printing the title and applying padding to ensure that the
	// space taken up by the title is the same width as the column.
	for _, h := range hs {
		_, err := buf.Write(padColumn(h.Title, h.Length))
		if err != nil {
			return err
		}
	}

	// write the line break after the header row
	_, err := buf.Write([]byte("\n"))
	if err != nil {
		return err
	}

	// iterate over each row and prepare to write it in turn.
	for _, r := range rs {
		// in order to ensure spacing correctly iterate over each header
		// and print the appropriate value from this row. If the row does
		// not contain the header then an empty string will be printed
		// and spaced appropriately to take up the correct width of the column.
		for _, h := range hs {
			_, err := buf.Write(padColumn(r[h.Title], h.Length))
			if err != nil {
				return err
			}
		}

		// write the line break after the row is done
		_, err := buf.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

// padColumn takes a string value and adds padding to the end to ensure it meets the specified length.
func padColumn(value string, length int) []byte {
	o := value
	for i := len(value); i < length; i++ {
		o = o + " "
	}

	return []byte(o)
}
