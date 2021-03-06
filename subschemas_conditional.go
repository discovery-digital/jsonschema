package jsonschema

import "reflect"

// File and methods named in respect to https://json-schema.org/latest/json-schema-validation.html#rfc.section.6.6
var ifThenElseType = reflect.TypeOf((*ifThenElse)(nil)).Elem()

// Implement IfThenElse() when condition needs to be used
// {
//    "if": { "properties": { "power": { "minimum": 9000 } } },
//    "then": { "required": [ "disbelief" ] },
//    "else": { "required": [ "confidence" ] }
// }
type ifThenElse interface {
	IfThenElse() SchemaCondition
}

// SchemaCondition holds data for if/then/else jsonschema statements
//
// If: A reflect.StructField that defines the condition to be met.
// Then: A type that will be converted to a jsonschema subschema and evaluated if the condition is met
// Else: A type that will be converted to a jsonschema subschema and evaluated if the condition is not met
type SchemaCondition struct {
	If   reflect.StructField
	Then interface{}
	Else interface{}
}

// Append jsonschema rules from IfThenElse interface to the jsonschema for the struct that implements them
func (r *Reflector) addSubschemasForConditionalCases(schema *Type, definitions Definitions, t reflect.Type) {
	if schema == nil {
		return
	}

	t, nonNilPointer := getNonNilPointerTypeAndInterface(t)

	if t.Implements(ifThenElseType) {
		condition := nonNilPointer.(ifThenElse).IfThenElse()
		r.reflectCondition(definitions, condition, schema)
	}
}

func (r *Reflector) reflectCondition(definitions Definitions, sc SchemaCondition, t *Type) {
	conditionSchema := Type{}
	conditionSchema.structKeywordsFromTags(r.getJSONSchemaTags(sc.If, nil))

	t.If = &Type{
		Properties: map[string]*Type{
			sc.If.Tag.Get("json"): &conditionSchema,
		},
	}

	if reflect.TypeOf(sc.Then) != nil {
		t.Then = r.reflectTypeToSchema(definitions, reflect.TypeOf(sc.Then))
	}
	if reflect.TypeOf(sc.Else) != nil {
		t.Else = r.reflectTypeToSchema(definitions, reflect.TypeOf(sc.Else))
	}
}
