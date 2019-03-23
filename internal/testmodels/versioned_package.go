package testmodels

import (
	"github.com/discovery-digital/jsonschema/internal/testmodels/v1"
	v2 "github.com/discovery-digital/jsonschema/internal/testmodels/v2"
	"reflect"
)

// These are models used for the versioned package test, but the actual test cases are in reflect_test.go
type TestVersionedPackages struct {
	Tester    TesterPackage    `json:"tester"`
	Developer DeveloperPackage `json:"developer"`
}

func (user TestVersionedPackages) OneOf() []reflect.StructField {
	tester, _ := reflect.TypeOf(user).FieldByName("Tester")
	developer, _ := reflect.TypeOf(user).FieldByName("Developer")
	return []reflect.StructField{
		tester,
		developer,
	}
}

// Tester  struct
type TesterPackage struct {
	Experience StringOrNull1 `json:"experience"`
}

// Developer  struct
type DeveloperPackage struct {
	Experience     StringOrNull1 `json:"experience" jsonschema:"minLength=1"`
	Language       StringOrNull1 `json:"language" jsonschema:"pattern=\\S+"`
	HardwareChoice v1.Hardware   `json:"hardware"`
	HardwareChoic  v2.Hardware   `json:"hardware1"`
}

type StringOrNull1 struct {
	String string
	IsNull bool
}

type DesktopPackage struct {
	FormFactor   string `json:"form_factor" jsonschema:"pattern=^(standard|micro|mini|nano)"`
	NeedKeyboard bool   `json:"need_keyboard"`
}

func (p StringOrNull1) OneOf() []reflect.StructField {
	strings, _ := reflect.TypeOf(p).FieldByName("String")
	return []reflect.StructField{
		strings,
		reflect.StructField{Type: nil},
	}
}
