package testmodels

import (
	"github.com/discovery-digital/jsonschema"
)

// These are models used for the subschema_switch test, but the actual test cases are in reflect_test.go
type ExampleCaseWithoutOrder struct {
	Type string `json:"type" jsonschema:"optional"`
}

type IntPayloadWithoutOrder struct {
	Payload int `json:"payload"`
}

type StringPayloadWithoutOrder struct {
	Payload string `json:"payload"`
}

type BoolPayloadWithoutOrder struct {
	Payload bool `json:"payload"`
}

func (ex ExampleCaseWithoutOrder) Case() jsonschema.SchemaSwitch {
	cases := make(map[string]interface{})
	cases["bool"] = BoolPayload{}
	cases["int"] = IntPayload{}
	cases["string"] = StringPayload{}
	return jsonschema.SchemaSwitch{
		ByField: "type",
		Cases:   cases,
	}
}
