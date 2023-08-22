package bson

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

var structInfoCache sync.Map // map[reflect.Type]*structInfo

type structInfo struct {
	Fields []fieldInfo
}

type fieldInfo struct {
	Key       string
	Num       int
	OmitEmpty bool
}

func (si *structInfo) asDoc(val reflect.Value) docRefl {
	doc := make(docRefl, 0, len(si.Fields))

	for _, info := range si.Fields {
		value := val.Field(info.Num)
		if info.OmitEmpty && isZero(value) {
			continue
		}

		doc = append(doc, pairRefl{
			Key: info.Key,
			Val: value,
		})
	}
	return doc
}

func getStruct(val reflect.Value) *structInfo {
	typ := val.Type()
	if info, ok := structInfoCache.Load(typ); ok {
		return info.(*structInfo)
	}

	info, _ := getStructInfo(typ)
	structInfoCache.Store(typ, info)
	return info
}

func getStructInfo(typ reflect.Type) (*structInfo, error) {
	n := typ.NumField()
	fields := make([]fieldInfo, 0, n)
	fieldsMap := make(map[string]fieldInfo, n)

	for i := 0; i < n; i++ {
		field := typ.Field(i)
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}

		info := fieldInfo{Num: i}

		tag := field.Tag.Get("bson")
		if tag == "" && strings.Index(string(field.Tag), ":") == -1 {
			tag = string(field.Tag)
		}
		if tag == "-" {
			continue
		}

		tagsParts := strings.Split(tag, ",")
		if len(tagsParts) > 1 {
			for _, flag := range tagsParts[1:] {
				switch flag {
				case "omitempty":
					info.OmitEmpty = true
				default:
					panic("Unsupported flag: " + flag)
				}
			}
			tag = tagsParts[0]
		}

		if tag != "" {
			info.Key = tag
		} else {
			info.Key = strings.ToLower(field.Name)
		}

		if _, ok := fieldsMap[info.Key]; ok {
			return nil, errors.New("Duplicated key: " + info.Key)
		}

		fields = append(fields, info)
		fieldsMap[info.Key] = info
	}

	info := &structInfo{
		Fields: fields,
	}
	return info, nil
}
