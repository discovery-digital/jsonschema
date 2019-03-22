package jsonschema_test

import (
	"github.com/discovery-digital/jsonschema"
	"reflect"
)

// These are models used for the subschema_conditional test, but the actual test cases are in reflect_test.go
type Application struct {
	Type string `json:"type"`
}

type ApplicationValidation struct {
	Type string `json:"type" jsonschema:"enum=web"`
}

type WebApp struct {
	Browser string `json:"browser"`
}

type MobileApp struct {
	Device string `json:"device"`
}

func (app Application) IfThenElse() jsonschema.SchemaCondition {
	conditionField, _ := reflect.TypeOf(ApplicationValidation{}).FieldByName("Type")
	return jsonschema.SchemaCondition{
		If:   conditionField,
		Then: WebApp{},
		Else: MobileApp{},
	}
}

