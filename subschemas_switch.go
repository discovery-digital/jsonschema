package jsonschema

import "reflect"

// Case() isn't an official jsonschema rule but a shorthand to simulate switch logic (without default)
// This allows us to then evaluate a field and validate different schema based on the value of the field
// The example below evaluates the field "type" and will validate against the Apple / Bacon / Watermelon schemas
// {
//  "oneOf": [
//    { "if": { "properties": { "type": { "enum": [ "apple" ] } } },
//      "then": {
//        "$schema": "http://json-schema.org/draft-07/schema#",
//        "$ref": "#/definitions/Apple"
//      },
//      "else": { "properties": { "type": { "enum": [ "apple" ] } } }
//    },
//    { "if": { "properties": { "type": { "enum": [ "bacon" ] } } },
//      "then": {
//        "$schema": "http://json-schema.org/draft-07/schema#",
//        "$ref": "#/definitions/Bacon"
//      },
//      "else": { "properties": { "type": { "enum": [ "bacon" ] } } }
//    },
//    { "if": { "properties": { "type": { "enum": [ "watermelon" ] } } },
//      "then": {
//        "$schema": "http://json-schema.org/draft-07/schema#",
//        "$ref": "#/definitions/Watermelon"
//      },
//      "else": { "properties": { "type": { "enum": [ "watermelon" ] } } } },
//  ]
// }
type schemaCase interface {
	Case() SchemaSwitch
}

// SchemaSwitch holds data for emulating switch case over some field value
type SchemaSwitch struct {
	ByField string
	Cases   map[string]interface{}
}

var schemaCaseType = reflect.TypeOf((*schemaCase)(nil)).Elem()

// Appends jsonschema rules from Case interface to the jsonschema for the struct that implements them
func (r *Reflector) addSubschemasForSwitch(st *Type, definitions Definitions, t reflect.Type) {
	if st == nil {
		return
	}

	var nonNilPointer interface{}
	t, nonNilPointer = getNonNilPointerTypeAndInterface(t)

	if t.Implements(schemaCaseType) {
		schemaSwitch := nonNilPointer.(schemaCase).Case()
		st.OneOf = r.reflectCases(definitions, schemaSwitch)
	}
}

func (r *Reflector) reflectCases(definitions Definitions, sc SchemaSwitch) []*Type {
	casesList := make([]*Type, 0)
	for key, value := range sc.Cases {
		t := &Type{}
		t.If = &Type{
			Properties: map[string]*Type{
				sc.ByField: &Type{
					Enum: []interface{}{key},
				},
			},
		}
		t.Then = r.reflectTypeToSchema(definitions, reflect.TypeOf(value))
		t.Else = t.If
		casesList = append(casesList, t)
	}
	return casesList
}
