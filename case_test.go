package jsonschema

import "reflect"

type ExampleCase struct {
	Type string `json:"type" jsonschema:"required;enum="int|string|bool"`
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

func (ex ExampleCase) Case() SchemaSwitch {
	cases := make(map[string]interface{})
	cases["int"] = IntPayload{}
	cases["string"] = StringPayload{}
	cases["bool"] = BoolPayload{}

	ByField, _ := reflect.TypeOf(ExampleCase{}).FieldByName("Type")
	return SchemaSwitch{
		ByField: ByField,
		Cases:   cases,
	}
}
