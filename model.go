package red

import (
	"errors"
	xconv "github.com/goclub/conv"
	"reflect"
)

func StructFields(v interface{}) (fields []string, err error) {
	rValue := reflect.ValueOf(v)
	err = coreStructFields(rValue, &fields) ; if err != nil {
		return
	}
	return
}
func coreStructFields(rValue reflect.Value, fields *[]string) error {
	rType := rValue.Type()
	for i:=0;i<rType.NumField();i++ {
		structField := rType.Field(i)
		tag, ok := structField.Tag.Lookup("red")
		if structField.Type.Kind() == reflect.Struct && ok == false {
			return coreStructFields(rValue.Field(i), fields)
		}
		if ok == false {
			continue
		}
		*fields = append(*fields, tag)
	}
	return nil
}

func StructFieldValues(v interface{}) (fieldValues []FieldValue, err error)  {
	rValue := reflect.ValueOf(v)
	err = coreStructFieldValues(rValue, &fieldValues) ; if err != nil {
		return
	}
	return
}
func coreStructFieldValues(rValue reflect.Value, fieldValues *[]FieldValue) error {
	rType := rValue.Type()
	for i:=0;i<rType.NumField();i++ {
		structField := rType.Field(i)
		tag, ok := structField.Tag.Lookup("red")
		if structField.Type.Kind() == reflect.Struct && ok == false {
			return coreStructFieldValues(rValue.Field(i), fieldValues)
		}
		if ok == false {
			continue
		}
		rItem := rValue.Field(i)
		value, convErr := func () (string, error) {
			if valuer, asValuer := rItem.Interface().(Marshaler) ; asValuer {
				data, err := valuer.MarshalText()
				return string(data), err
			}
			return xconv.ReflectToString(rItem)
		}() ; if convErr != nil {
			return errors.New("goclub/redis: name:" + structField.Name + " kind:" +structField.Type.Kind().String() + " not string or not implements red.Marshaler")
		}
		*fieldValues = append(*fieldValues, FieldValue{
			Field: tag,
			Value: value,
		})
	}
	return nil
}

func StructScan(ptr interface{},  values []string) error {
	rValue := reflect.ValueOf(ptr)
	rType := rValue.Type()
	if rType.Kind() != reflect.Ptr {
		return errors.New("goclub/redis: StructScan(ptr interface{}) ptr must be pointer")
	}
	rValue = rValue.Elem()
	rType = rType.Elem()
	length := len(values)
	offset := 0
	return coreStructScan(rValue, rType, values, length, &offset)
}

func coreStructScan(rValue reflect.Value, rType reflect.Type, values []string, length int, offset *int) error {
	if *offset >= length {
		return nil
	}
	for i:=0;i<rType.NumField();i++ {
		structField := rType.Field(i)
		_, ok := structField.Tag.Lookup("red")
		if structField.Type.Kind() == reflect.Struct && ok == false {
			rItem := rValue.Field(i)
			return coreStructScan(rItem, rItem.Type(), values, length, offset)
		}
		if ok == false {
			continue
		}
		value := values[*offset]
		rItem := rValue.Field(i)
		rItemAddr := rItem
		if rItem.CanAddr() {
			rItemAddr = rItem.Addr()
		}
		if scaner, asScaner := rItemAddr.Interface().(Unmarshaler); asScaner {
			scaner.UnmarshalText([]byte(value))
		} else {
			err := xconv.StringToReflect(value, rItem) ; if err != nil {
				return errors.New("goclub/redis: name:" + structField.Name + " kind:" +structField.Type.Kind().String() + " not string or not implements red.Unmarshaler")
			}
		}
		*offset++
	}
	return nil
}

type Unmarshaler interface {
	UnmarshalText(data []byte) error
}
type Marshaler interface {
	MarshalText() ([]byte, error)
}

