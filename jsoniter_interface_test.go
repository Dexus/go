package jsoniter

import (
	"encoding/json"
	"github.com/json-iterator/go/require"
	"testing"
	"unsafe"
)

func Test_write_array_of_interface(t *testing.T) {
	should := require.New(t)
	array := []interface{}{"hello"}
	str, err := MarshalToString(array)
	should.Nil(err)
	should.Equal(`["hello"]`, str)
}

func Test_write_map_of_interface(t *testing.T) {
	should := require.New(t)
	val := map[string]interface{}{"hello": "world"}
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`{"hello":"world"}`, str)
}

func Test_write_map_of_interface_in_struct(t *testing.T) {
	type TestObject struct {
		Field map[string]interface{}
	}
	should := require.New(t)
	val := TestObject{map[string]interface{}{"hello": "world"}}
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`{"Field":{"hello":"world"}}`, str)
}

func Test_write_map_of_interface_in_struct_with_two_fields(t *testing.T) {
	type TestObject struct {
		Field  map[string]interface{}
		Field2 string
	}
	should := require.New(t)
	val := TestObject{map[string]interface{}{"hello": "world"}, ""}
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Contains(str, `"Field":{"hello":"world"}`)
}

type MyInterface interface {
	Hello() string
}

type MyString string

func (ms MyString) Hello() string {
	return string(ms)
}

func Test_write_map_of_custom_interface(t *testing.T) {
	should := require.New(t)
	myStr := MyString("world")
	should.Equal("world", myStr.Hello())
	val := map[string]MyInterface{"hello": myStr}
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`{"hello":"world"}`, str)
}

func Test_write_interface(t *testing.T) {
	should := require.New(t)
	var val interface{}
	val = "hello"
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`"hello"`, str)
}

func Test_read_interface(t *testing.T) {
	should := require.New(t)
	var val interface{}
	err := UnmarshalFromString(`"hello"`, &val)
	should.Nil(err)
	should.Equal("hello", val)
}

func Test_read_custom_interface(t *testing.T) {
	should := require.New(t)
	var val MyInterface
	RegisterTypeDecoderFunc("jsoniter.MyInterface", func(ptr unsafe.Pointer, iter *Iterator) {
		*((*MyInterface)(ptr)) = MyString(iter.ReadString())
	})
	err := UnmarshalFromString(`"hello"`, &val)
	should.Nil(err)
	should.Equal("hello", val.Hello())
}

func Test_decode_object_contain_empty_interface(t *testing.T) {
	type TestObject struct {
		Field interface{}
	}
	should := require.New(t)
	obj := TestObject{}
	obj.Field = 1024
	should.Nil(UnmarshalFromString(`{"Field": "hello"}`, &obj))
	should.Equal("hello", obj.Field)
}

func Test_decode_object_contain_non_empty_interface(t *testing.T) {
	type TestObject struct {
		Field MyInterface
	}
	should := require.New(t)
	obj := TestObject{}
	obj.Field = MyString("abc")
	should.Nil(UnmarshalFromString(`{"Field": "hello"}`, &obj))
	should.Equal(MyString("hello"), obj.Field)
}

func Test_encode_object_contain_empty_interface(t *testing.T) {
	type TestObject struct {
		Field interface{}
	}
	should := require.New(t)
	obj := TestObject{}
	obj.Field = 1024
	str, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{"Field":1024}`, str)
}

func Test_encode_object_contain_non_empty_interface(t *testing.T) {
	type TestObject struct {
		Field MyInterface
	}
	should := require.New(t)
	obj := TestObject{}
	obj.Field = MyString("hello")
	str, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{"Field":"hello"}`, str)
}

func Test_nil_non_empty_interface(t *testing.T) {
	ConfigDefault.cleanEncoders()
	ConfigDefault.cleanDecoders()
	type TestObject struct {
		Field []MyInterface
	}
	should := require.New(t)
	obj := TestObject{}
	b := []byte(`{"Field":["AAA"]}`)
	should.NotNil(json.Unmarshal(b, &obj))
	should.NotNil(Unmarshal(b, &obj))
}

func Test_read_large_number_as_interface(t *testing.T) {
	should := require.New(t)
	var val interface{}
	err := Config{UseNumber: true}.Froze().UnmarshalFromString(`123456789123456789123456789`, &val)
	should.Nil(err)
	output, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`123456789123456789123456789`, output)
}

func Test_nested_one_field_struct(t *testing.T) {
	should := require.New(t)
	type YetYetAnotherObject struct {
		Field string
	}
	type YetAnotherObject struct {
		Field *YetYetAnotherObject
	}
	type AnotherObject struct {
		Field *YetAnotherObject
	}
	type TestObject struct {
		Me *AnotherObject
	}
	obj := TestObject{&AnotherObject{&YetAnotherObject{&YetYetAnotherObject{"abc"}}}}
	str, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{"Me":{"Field":{"Field":{"Field":"abc"}}}}`, str)
	str, err = MarshalToString(&obj)
	should.Nil(err)
	should.Equal(`{"Me":{"Field":{"Field":{"Field":"abc"}}}}`, str)
}

func Test_struct_with_one_nil(t *testing.T) {
	type TestObject struct {
		F *float64
	}
	var obj TestObject
	should := require.New(t)
	output, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{"F":null}`, output)
}

func Test_struct_with_one_nil_embedded(t *testing.T) {
	type Parent struct {
		Field1 string
		Field2 string
	}
	type TestObject struct {
		*Parent
	}
	obj := TestObject{}
	should := require.New(t)
	bytes, err := json.Marshal(obj)
	should.Nil(err)
	should.Equal("{}", string(bytes))
	output, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{}`, output)
}

func Test_struct_with_not_nil_embedded(t *testing.T) {
	type Parent struct {
		Field0 string
		Field1 []string
		Field2 map[string]interface{}
	}
	type TestObject struct {
		*Parent
	}
	should := require.New(t)
	var obj TestObject
	err := UnmarshalFromString(`{"Field0":"1","Field1":null,"Field2":{"K":"V"}}`, &obj)
	should.Nil(err)
	should.Nil(obj.Field1)
	should.Equal(map[string]interface{}{"K": "V"}, obj.Field2)
	should.Equal("1", obj.Field0)
}

func Test_array_with_one_nil_ptr(t *testing.T) {
	obj := [1]*float64{nil}
	should := require.New(t)
	output, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`[null]`, output)
}

func Test_array_with_one_not_nil_ptr(t *testing.T) {
	two := float64(2)
	obj := [1]*float64{&two}
	should := require.New(t)
	output, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`[2]`, output)
}

func Test_embedded_array_with_one_nil(t *testing.T) {
	type TestObject struct {
		Field1 int
		Field2 [1]*float64
	}
	var obj TestObject
	should := require.New(t)
	output, err := MarshalToString(obj)
	should.Nil(err)
	should.Contains(output, `"Field2":[null]`)
}

func Test_array_with_nothing(t *testing.T) {
	var obj [2]*float64
	should := require.New(t)
	output, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`[null,null]`, output)
}
