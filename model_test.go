package red

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type InvalidType struct {
	Name string
	Some struct{
		A string
	} `red:"some"`
}
type AB struct {
	A string
	B string
}

func (data AB) MarshalText() ([]byte, error) {
	return []byte(data.A + data.B), nil
}

func (data *AB) UnmarshalText(b []byte) error {
	value := string(b)
	runes := []rune(value)
	data.A = string(runes[0])
	data.B = string(runes[1])
	return nil
}
func TestStructToFieldValue(t *testing.T) {
	{
		type Data struct {
			Name string `red:"name"`
			Age int
			UserName string
			Data *Data
		}
		data := Data {
			Name: "nimo",
		}
		fieldValues, err := StructToFieldValue(data)
		assert.NoError(t, err)
		assert.Equal(t, fieldValues, []FieldValue{{"name","nimo"}})
	}
	{
		type Sub struct {
			Like string `red:"like"`
		}
		type Data struct {
			Name string `red:"name"`
			Age int `red:"age"`
			Sub Sub
		}
		data := Data{
			Name:"nimo",
			Age: 18,
			Sub: Sub{Like: "read"},
		}
		fieldValues, err := StructToFieldValue(data)
		assert.NoError(t, err)
		assert.Equal(t, fieldValues, []FieldValue{{"name","nimo"},{"age","18"},{"like","read"}})
	}
	{

		type Data struct {
			Name string `red:"name"`
			Age float64 `red:"age"`
			AB AB `red:"ab"`
		}
		data := Data{
			Name:"nimo",
			Age: 18,
			AB: AB{"1","2"},
		}
		fieldValues, err := StructToFieldValue(data)
		assert.NoError(t, err)
		assert.Equal(t, fieldValues, []FieldValue{{"name","nimo"},{"age","18"}, {"ab", "12"}})
	}
	{
		_, err := StructToFieldValue(InvalidType{})
		assert.EqualError(t, err, "goclub/redis: name:Some kind:struct not string or not implements red.Marshaler")
	}
	{
		{

			type Data struct {
				Time time.Time `red:"time"`
			}
			data := Data{

			}
			err := data.Time.UnmarshalText([]byte("2021-02-13T00:00:00+08:00"))
			assert.NoError(t, err)
			fieldValues, err := StructToFieldValue(data)
			assert.NoError(t, err)
			assert.Equal(t, fieldValues, []FieldValue{{"time","2021-02-13T00:00:00+08:00",}})
		}
	}
}

func TestStructScan(t *testing.T) {
	{
		type Data struct {
			Name string `red:"name"`
			Age string
			UserName string
			Data *Data
		}
		var data Data
		assert.NoError(t, StructScan(&data, []string{"value"}))
		assert.Equal(t, data, Data{Name: "value"})
	}
	{
		type Sub struct {
			Like string `red:"like"`
		}
		type Data struct {
			Name string `red:"name"`
			Age int `red:"age"`
			Sub Sub
		}
		var data Data
		assert.NoError(t, StructScan(&data, []string{"nimoc", "18"}))
		assert.Equal(t, data, Data{Name: "nimoc", Age: 18})
	}
	{
		type Sub struct {
			Like string `red:"like"`
			Count uint `red:"count"`
		}
		type Data struct {
			Name string `red:"name"`
			Age float64 `red:"age"`
			Sub
		}
		var data Data
		assert.NoError(t, StructScan(&data, []string{"nimoc", "18.2", "1", "2"}))
		assert.Equal(t, data, Data{Name: "nimoc", Age: 18.2, Sub: Sub{Like: "1", Count: 2}})
	}
	{
		type Data struct {
			AB AB
		}
		var data Data
		assert.NoError(t, StructScan(&data, nil))
		assert.Equal(t, data, Data{})
	}
	{
		type Data struct {
			AB AB `red:"ab"`
		}
		var data Data
		assert.NoError(t, StructScan(&data, []string{"12"}))
		assert.Equal(t, data, Data{AB:AB{"1","2",}})
	}
	{
		type Data struct {
			Time time.Time `red:"time"`
		}
		var data Data
		assert.NoError(t, StructScan(&data, []string{"2021-02-13T00:00:00.000000+08:00"}))
		testb, err := data.Time.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, string(testb), "2021-02-13T00:00:00+08:00")
	}
}