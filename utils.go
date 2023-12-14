package bson

import (
	"reflect"
	"sort"
	"strings"
	"unicode"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

type pairRefl struct {
	Key string
	Val any
}

type docRefl []pairRefl

func (doc docRefl) Len() int           { return len(doc) }
func (doc docRefl) Less(i, j int) bool { return doc[i].Key < doc[j].Key }
func (doc docRefl) Swap(i, j int)      { doc[i], doc[j] = doc[j], doc[i] }

func sortPairRefl(doc docRefl) {
	// TODO(cristaloleg): 20 might be too high/low.
	if len(doc) <= 20 {
		// for i := 1; i < len(ss); i++ {
		// 	for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
		// 		ss[j], ss[j-1] = ss[j-1], ss[j]
		// 	}
		// }
		// below is the code above, but (almost) without bound checks

		for i := 1; i < len(doc); i++ {
			for j := i; j > 0; j-- {
				if doc[j].Key >= doc[j-1].Key {
					break
				}
				doc[j], doc[j-1] = doc[j-1], doc[j]
			}
		}
	} else {
		sort.Sort(doc)
	}
}

func walkStruct(dict docRefl, v reflect.Value) (docRefl, error) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		strField := t.Field(i)
		field := v.FieldByIndex(strField.Index)

		if !field.CanInterface() || isNil(field) {
			continue
		}

		tag, ok := fieldTag(strField, field)
		if !ok {
			continue
		}

		if tag == "" && strField.Anonymous &&
			strField.Type.Kind() == reflect.Struct {

			var err error
			dict, err = walkStruct(dict, field)
			if err != nil {
				return nil, err
			}
		} else {
			dict = append(dict, pairRefl{Key: tag, Val: field})
		}
	}
	return dict, nil
}

func fieldTag(field reflect.StructField, v reflect.Value) (string, bool) {
	tag := field.Tag.Get("bson")

	var opts string
	switch {
	case tag == "":
		return field.Name, true
	case tag == "-":
		return "", false
	default:
		if idx := strings.Index(tag, ","); idx != -1 {
			tag, opts = tag[:idx], tag[idx:]
		}
	}

	switch {
	case strings.Contains(opts, ",omitempty") && isZero(v):
		return "", false
	case !isValidTag(tag):
		return field.Name, true
	default:
		return tag, true
	}
}

func isValidTag(key string) bool {
	if key == "" {
		return false
	}

	for _, c := range key {
		if c != ' ' && c != '$' && c != '-' && c != '_' && c != '.' &&
			!unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return false
	}
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	// TODO(cristaloleg): supporting reflect.Struct might be hard.
	default:
		return false
	}
}
