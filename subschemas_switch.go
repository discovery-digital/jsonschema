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
	// ByField = the name of the field you wish to evaluate (ex: "species")
	ByField string
	// Each key = the value for the field being evaluated (ex: "turtle")
	// Each value = the struct that holds the jsonschema tags to validate against when it is that value (ex: Turtle{})
	Cases map[string]interface{}
	// Order - the fields from `Cases` can be provided here to guarantee output in specified order, otherwise this will be seeded internally
	Order []string
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
	//Build order when not provided my the user of this library
	if len(sc.Order) == 0 {
		for key := range sc.Cases {
			sc.Order = append(sc.Order, key)
		}
	}
	casesList := make([]*Type, 0)
	for _, value := range sc.Order {
		t := &Type{}
		t.If = &Type{
			Properties: map[string]*Type{
				sc.ByField: &Type{
					Enum: []interface{}{value},
				},
			},
		}
		t.Then = r.reflectTypeToSchema(definitions, reflect.TypeOf(sc.Cases[value]))
		t.Else = t.If
		casesList = append(casesList, t)
	}
	return casesList
}
