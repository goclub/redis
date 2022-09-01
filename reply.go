package red

import (
	"fmt"
	xerr "github.com/goclub/error"
	"strconv"
)

type Reply struct {
	Value interface{}
}
func (r Reply) String() (s string, err error) {
	defer func() { if err != nil { err = xerr.WithStack(err) } }()
	switch value := r.Value.(type) {
	case int64:
		return strconv.FormatInt(value, 10), nil
	case string:
		return value, nil
	default:
		err := fmt.Errorf("goclub/redis: unexpected type(%T) value(%+v) convert string", value, value)
		return "", err
	}
}
func (r Reply) Int64() (v int64, err error) {
	defer func() { if err != nil { err = xerr.WithStack(err) } }()
	switch value := r.Value.(type) {
	case int64:
		return int64(value), nil
	case string:
		return strconv.ParseInt(value, 10, 64)
	default:
		err := fmt.Errorf("goclub/redis: unexpected type(%T) value(%+v) convert int64", value, value)
		return 0, err
	}
}
func (r Reply) Uint64() (v uint64, err error) {
	int64Value, err :=  r.Int64() ; if err != nil {
	    return
	}
	if int64Value < 0 {
		return 0, xerr.New(fmt.Sprintf("goclub/redis: %#v can not convert to uint64", int64Value))
	}
	v = uint64(int64Value)
	return
}
// func (r Reply) InterfaceSlice() (interfaceSlice []interface{}, err error) {
// 	switch value := r.Value.(type) {
// 	case []interface{}:
// 		return value, nil
// 	default:
// 		return nil, fmt.Errorf("redis: unexpected type(%T) value(%+v) convert []interface", value, value)
// 	}
// }
func (r Reply) StringSlice() (stringSlice []OptionString, err error) {
	defer func() { if err != nil { err = xerr.WithStack(err) } }()
	switch v := r.Value.(type) {
	case []interface{}:
	default:
		_=v
		err := fmt.Errorf("goclub/redis: unexpected type(%T) value(%+v) convert []OptionString", r.Value, r.Value)
		return nil, err
	}
	values := r.Value.([]interface{})
	for _, v := range values {
		if v == nil {
			stringSlice = append(stringSlice, OptionString{})
		} else {
			var item string
			switch value := v.(type) {
			case int64:
				item = strconv.FormatInt(value, 10)
			case string:
				item = value
			default:
				err := fmt.Errorf("goclub/redis: unexpected type(%T) value(%+v) convert string", value, value)
				return nil, err
			}
			stringSlice = append(stringSlice, NewOptionString(item))
		}
	}
	return
}
func (r Reply) Int64Slice() (int64Slice []OptionInt64, err error) {
	defer func() { if err != nil { err = xerr.WithStack(err) } }()
	switch v := r.Value.(type) {
	case []interface{}:
	default:
		_=v
		err := fmt.Errorf("goclub/redis: unexpected type(%T) value(%+v) convert []OptionInt64", r.Value, r.Value)
		return nil, err
	}
	values := r.Value.([]interface{})
	for _, v := range values {
		if v == nil {
			int64Slice = append(int64Slice, OptionInt64{})
		} else {
			var item int64
			switch value := v.(type) {
			case int64:
				item = value
			case string:
				item, err = strconv.ParseInt(value, 10, 64) ; if err != nil {
				    return
				}
			default:
				err := fmt.Errorf("goclub/redis: unexpected type(%T) value(%+v) convert int64", value, value)
				return nil, err
			}
			int64Slice = append(int64Slice, NewOptionInt64(item))
		}
	}
	return
}