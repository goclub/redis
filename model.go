package red

import (
	"errors"
	xconv "github.com/goclub/conv"
	"reflect"
)

func StructToFieldValue(v interface{}) (fieldValues []FieldValue, err error)  {
	rValue := reflect.ValueOf(v)
	err = coreStructToFieldValue(rValue, &fieldValues) ; if err != nil {
		return
	}
	return
}
func coreStructToFieldValue(rValue reflect.Value, fieldValues *[]FieldValue) error {
	rType := rValue.Type()
	for i:=0;i<rType.NumField();i++ {
		structField := rType.Field(i)
		tag, ok := structField.Tag.Lookup("redis")
		if structField.Type.Kind() == reflect.Struct && ok == false {
			return coreStructToFieldValue(rValue.Field(i), fieldValues)
		}
		if ok == false {
			continue
		}
		rItem := rValue.Field(i)
		var value string
		if structField.Type.Kind() == reflect.String {
			value = rItem.String()
		} else if valuer, asValuer := rItem.Interface().(Valuer) ; asValuer {
			value = valuer.RedisValue()
		} else {
			return errors.New("goclub/redis:" + structField.Type.Name() + " not string or not implements red.Valuer")
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
		_, ok := structField.Tag.Lookup("redis")
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
		if scaner, asScaner := rItemAddr.Interface().(Scanner); asScaner {
			err := scaner.RedisScan(value) ; if err != nil {
				return err
			}
		} else {
			err := xconv.StringReflect(value, rItem) ; if err != nil {
				return err
			}
		}
		*offset++
	}
	return nil
}

type Scanner interface {
	RedisScan(value string) error
}
type Valuer interface {
	RedisValue() string
}