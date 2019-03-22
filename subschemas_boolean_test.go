package jsonschema_test

import (
	"reflect"
)

// These are models used for the subschema_boolean test, but the actual test cases are in reflect_test.go
type TestUserOneOf struct {
	Tester    Tester    `json:"tester"`
	Developer Developer `json:"developer"`
}

func (user TestUserOneOf) OneOf() []reflect.StructField {
	tester, _ := reflect.TypeOf(user).FieldByName("Tester")
	developer, _ := reflect.TypeOf(user).FieldByName("Developer")
	return []reflect.StructField{
		tester,
		developer,
	}
}

// Tester  struct
type Tester struct {
	Experience StringOrNull `json:"experience"`
}

// Developer  struct
type Developer struct {
	Experience     StringOrNull `json:"experience" jsonschema:"minLength=1"`
	Language       StringOrNull `json:"language" jsonschema:"pattern=\\S+"`
	HardwareChoice Hardware     `json:"hardware"`
}

type StringOrNull struct {
	String string
	IsNull bool
}

type Hardware struct {
	Brand  string `json:"brand" jsonschema:"notEmpty"`
	Memory int    `json:"memory"`
}

type Laptop struct {
	Brand           string `json:"brand" jsonschema:"pattern=^(apple|lenovo|dell)$"`
	NeedTouchScreen bool   `json:"need_touchscreen"`
}

type Desktop struct {
	FormFactor   string `json:"form_factor" jsonschema:"pattern=^(standard|micro|mini|nano)"`
	NeedKeyboard bool   `json:"need_keyboard"`
}

func (p StringOrNull) OneOf() []reflect.StructField {
	strings, _ := reflect.TypeOf(p).FieldByName("String")
	return []reflect.StructField{
		strings,
		reflect.StructField{Type: nil},
	}
}

func (h Hardware) AndOneOf() []reflect.StructField {
	return []reflect.StructField{
		reflect.StructField{Type: reflect.TypeOf(Laptop{})},
		reflect.StructField{Type: reflect.TypeOf(Desktop{})},
	}
}



