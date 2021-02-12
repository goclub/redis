package red

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

type InvalidType struct {
	Name string
	Some struct{
		A string
	} `redis:"some"`
}
type AB struct {
	A string
	B string
}
func (data AB) RedisValue() string {
	return data.A + data.B
}
func (data *AB) RedisScan(value string) error {
	log.Print(value)
	runes := []rune(value)
	data.A = string(runes[0])
	data.B = string(runes[1])
	return nil
}
func TestStructToFieldValue(t *testing.T) {
	{
		type Data struct {
			Name string `redis:"name"`
			Age string
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
			Like string `redis:"like"`
		}
		type Data struct {
			Name string `redis:"name"`
			Age string `redis:"age"`
			Sub Sub
		}
		data := Data{
			Name:"nimo",
			Age: "18",
			Sub: Sub{Like: "read"},
		}
		fieldValues, err := StructToFieldValue(data)
		assert.NoError(t, err)
		assert.Equal(t, fieldValues, []FieldValue{{"name","nimo"},{"age","18"},{"like","read"}})
	}
	{

		type Data struct {
			Name string `redis:"name"`
			Age string `redis:"age"`
			AB AB `redis:"ab"`
		}
		data := Data{
			Name:"nimo",
			Age: "18",
			AB: AB{"1","2"},
		}
		fieldValues, err := StructToFieldValue(data)
		assert.NoError(t, err)
		assert.Equal(t, fieldValues, []FieldValue{{"name","nimo"},{"age","18"}, {"ab", "12"}})
	}
	{
		_, err := StructToFieldValue(InvalidType{})
		assert.EqualError(t, err, "goclub/redis: not string or not implements red.Valuer")
	}
}

func TestStructScan(t *testing.T) {
	{
		type Data struct {
			Name string `redis:"name"`
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
			Like string `redis:"like"`
		}
		type Data struct {
			Name string `redis:"name"`
			Age int `redis:"age"`
			Sub Sub
		}
		var data Data
		assert.NoError(t, StructScan(&data, []string{"nimoc", "18"}))
		assert.Equal(t, data, Data{Name: "nimoc", Age: 18})
	}
	{
		type Sub struct {
			Like string `redis:"like"`
			Count uint `redis:"count"`
		}
		type Data struct {
			Name string `redis:"name"`
			Age float64 `redis:"age"`
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
			AB AB `redis:"ab"`
		}
		var data Data
		assert.NoError(t, StructScan(&data, []string{"12"}))
		assert.Equal(t, data, Data{AB:AB{"1","2",}})
	}

}