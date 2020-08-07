package utils

import (
	"errors"
	"fmt"
	"github.com/fatih/structs"
	st "github.com/golang/protobuf/ptypes/struct"
	"reflect"
)

// ToStruct converts a map[string]interface{} to a ptypes.Struct
func ToStruct(v map[string]interface{}) *st.Struct {
	size := len(v)
	if size == 0 {
		return nil
	}
	fields := make(map[string]*st.Value, size)
	for k, v := range v {
		fields[k] = ToValue(v)
	}
	return &st.Struct{
		Fields: fields,
	}
}

// Struct converts a map[string]interface{} to a ptypes.Struct
func Struct(v interface{}) (*st.Struct, error) {
	return Map2Struct(structs.Map(v))
}

func StringToStruct(v map[string]string) *st.Struct {
	size := len(v)
	if size == 0 {
		return nil
	}
	fields := make(map[string]*st.Value, size)
	for k, v := range v {
		fields[k] = ToStringValue(v)
	}
	return &st.Struct{
		Fields: fields,
	}
}

func ToStringValue(s string) *st.Value {
	return &st.Value{
		Kind: &st.Value_StringValue{
			StringValue: s,
		},
	}
}

// ToValue converts an interface{} to a ptypes.Value
func ToValue(v interface{}) *st.Value {
	switch v := v.(type) {
	case nil:
		return nil
	case bool:
		return &st.Value{
			Kind: &st.Value_BoolValue{
				BoolValue: v,
			},
		}
	case int:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int8:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int32:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint8:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint32:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float32:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: v,
			},
		}
	case string:
		return &st.Value{
			Kind: &st.Value_StringValue{
				StringValue: v,
			},
		}
	case error:
		return &st.Value{
			Kind: &st.Value_StringValue{
				StringValue: v.Error(),
			},
		}
	default:
		// Fallback to reflection for other types
		return toValue(reflect.ValueOf(v))
	}
}
func toValue(v reflect.Value) *st.Value {
	switch v.Kind() {
	case reflect.Bool:
		return &st.Value{
			Kind: &st.Value_BoolValue{
				BoolValue: v.Bool(),
			},
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v.Int()),
			},
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v.Uint()),
			},
		}
	case reflect.Float32, reflect.Float64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: v.Float(),
			},
		}
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return toValue(reflect.Indirect(v))
	case reflect.Array, reflect.Slice:
		size := v.Len()
		if size == 0 {
			return nil
		}
		values := make([]*st.Value, size)
		for i := 0; i < size; i++ {
			values[i] = toValue(v.Index(i))
		}
		return &st.Value{
			Kind: &st.Value_ListValue{
				ListValue: &st.ListValue{
					Values: values,
				},
			},
		}
	case reflect.Struct:
		t := v.Type()
		size := v.NumField()
		if size == 0 {
			return nil
		}
		fields := make(map[string]*st.Value, size)
		for i := 0; i < size; i++ {
			name := t.Field(i).Name
			// Better way?
			if len(name) > 0 && 'A' <= name[0] && name[0] <= 'Z' {
				fields[name] = toValue(v.Field(i))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &st.Value{
			Kind: &st.Value_StructValue{
				StructValue: &st.Struct{
					Fields: fields,
				},
			},
		}
	case reflect.Map:
		keys := v.MapKeys()
		if len(keys) == 0 {
			return nil
		}
		fields := make(map[string]*st.Value, len(keys))
		for _, k := range keys {
			if k.Kind() == reflect.String {
				fields[k.String()] = toValue(v.MapIndex(k))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &st.Value{
			Kind: &st.Value_StructValue{
				StructValue: &st.Struct{
					Fields: fields,
				},
			},
		}
	default:
		// Last resort
		return &st.Value{
			Kind: &st.Value_StringValue{
				StringValue: fmt.Sprint(v),
			},
		}
	}
}
func elabValue(value *st.Value) (interface{}, error) {
	var err error
	if value == nil {
		return nil, nil
	}
	if structValue, ok := value.GetKind().(*st.Value_StructValue); ok {
		result := make(map[string]interface{})
		for k, v := range structValue.StructValue.Fields {
			result[k], err = elabValue(v)
			if err != nil {
				return nil, err
			}
		}
		return result, err
	}
	if listValue, ok := value.GetKind().(*st.Value_ListValue); ok {
		result := make([]interface{}, len(listValue.ListValue.Values))
		for i, el := range listValue.ListValue.Values {
			result[i], err = elabValue(el)
			if err != nil {
				return nil, err
			}
		}
		return result, err
	}
	if _, ok := value.GetKind().(*st.Value_NullValue); ok {
		return nil, nil
	}
	if numValue, ok := value.GetKind().(*st.Value_NumberValue); ok {
		return numValue.NumberValue, nil
	}
	if strValue, ok := value.GetKind().(*st.Value_StringValue); ok {
		return strValue.StringValue, nil
	}
	if boolValue, ok := value.GetKind().(*st.Value_BoolValue); ok {
		return boolValue.BoolValue, nil
	}
	return nil, fmt.Errorf("Cannot convert the value %+v", value)
}
func elabEntry(entry interface{}) (*st.Value, error) {
	var err error
	if entry == nil {
		return &st.Value{Kind: &st.Value_NullValue{}}, nil
	}
	rt := reflect.TypeOf(entry)
	switch rt.Kind() {
	case reflect.String:
		if realValue, ok := entry.(string); ok {
			return &st.Value{Kind: &st.Value_StringValue{StringValue: realValue}}, nil
		}
		return nil, errors.New("cannot convert string value")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &st.Value{Kind: &st.Value_NumberValue{NumberValue: float64(reflect.ValueOf(entry).Int())}}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &st.Value{Kind: &st.Value_NumberValue{NumberValue: float64(reflect.ValueOf(entry).Uint())}}, nil
	case reflect.Float32, reflect.Float64:
		return &st.Value{Kind: &st.Value_NumberValue{NumberValue: reflect.ValueOf(entry).Float()}}, nil
	case reflect.Bool:
		if realValue, ok := entry.(bool); ok {
			return &st.Value{Kind: &st.Value_BoolValue{BoolValue: realValue}}, nil
		}
		return nil, errors.New("cannot convert boolean value")
	case reflect.Array, reflect.Slice:
		lstEntry := reflect.ValueOf(entry)
		lstValue := &st.ListValue{Values: make([]*st.Value, lstEntry.Len(), lstEntry.Len())}
		for i := 0; i < lstEntry.Len(); i++ {
			lstValue.Values[i], err = elabEntry(lstEntry.Index(i).Interface())
			if err != nil {
				return nil, err
			}
		}
		return &st.Value{Kind: &st.Value_ListValue{ListValue: lstValue}}, nil
	case reflect.Struct:
		return elabEntry(structs.Map(entry))
	case reflect.Map:
		mapEntry := make(map[string]interface{})
		entryValue := reflect.ValueOf(entry)
		for _, k := range entryValue.MapKeys() {
			mapEntry[k.String()] = entryValue.MapIndex(k).Interface()
		}
		structVal, err := Map2Struct(mapEntry)
		return &st.Value{Kind: &st.Value_StructValue{StructValue: structVal}}, err
	}
	return nil, fmt.Errorf("Cannot convert [%+v] kind:%s", entry, rt.Kind())
}

// Map2Struct ...
func Map2Struct(input map[string]interface{}) (*st.Struct, error) {
	var err error
	result := &st.Struct{Fields: make(map[string]*st.Value)}
	for k, v := range input {
		result.Fields[k], err = elabEntry(v)
		if err != nil {
			return nil, err
		}
	}
	return result, err
}

// Struct2Map converts the protobuf struct to go map interface
func Struct2Map(str *st.Struct) (map[string]interface{}, error) {
	var err error
	result := make(map[string]interface{})
	for k, v := range str.Fields {
		result[k], err = elabValue(v)
		if err != nil {
			return nil, err
		}
	}
	return result, err
}
