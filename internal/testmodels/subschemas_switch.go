package testmodels

import (
	"github.com/discovery-digital/jsonschema"
)

// These are models used for the subschema_switch test, but the actual test cases are in reflect_test.go
type ExampleCase struct {
	Type string `json:"type" jsonschema:"optional"`
}

type IntPayload struct {
	Payload int `json:"payload"`
}

type StringPayload struct {
	Payload string `json:"payload"`
}

type BoolPayload struct {
	Payload bool `json:"payload"`
}

func (ex ExampleCase) Case() jsonschema.SchemaSwitch {
	cases := make(map[string]interface{})
	order := []string{}
	cases["bool"] = BoolPayload{}
	order = append(order, "bool")
	cases["int"] = IntPayload{}
	order = append(order, "int")
	cases["string"] = StringPayload{}
	order = append(order, "string")

	return jsonschema.SchemaSwitch{
		ByField: "type",
		Cases:   cases,
		Order:   order,
	}
}
