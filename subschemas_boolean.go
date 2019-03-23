package jsonschema

import (
	"reflect"
)

// File named in respect to https://json-schema.org/latest/json-schema-validation.html#rfc.section.6.7
var andAnyOfType = reflect.TypeOf((*andAnyOf)(nil)).Elem()
var anyOfType = reflect.TypeOf((*anyOf)(nil)).Elem()
var andOneOfType = reflect.TypeOf((*andOneOf)(nil)).Elem()
var oneOfType = reflect.TypeOf((*oneOf)(nil)).Elem()
var andAllOfType = reflect.TypeOf((*andAllOf)(nil)).Elem()
var allOfType = reflect.TypeOf((*allOf)(nil)).Elem()

// AndAnyOf will generate the anyOf rule and retain the jsonschema rules for the struct that implements it
// `anyOf` is used to ensure that the data must be valid against *at least one* of the cases *or more*
// { "type": "number", "anyOf": [ { "multipleOf": 5 }, { "multipleOf": 3 } ]}
// In the example above, the input must be a number and can be either a multiple of 5 or 3 or both but never neither
type andAnyOf interface {
	AndAnyOf() []reflect.StructField
}

// AnyOf will overrule all jsonschema rules for the struct that implements it
// `anyOf` is used to ensure that the data must be valid against *at least one* of the cases *or more*
// { "anyOf": [ { "type": "number", "multipleOf": 5 }, { "type": "number", "multipleOf": 3 } ] }
// In the example above, the input must be a number and can be either a multiple of 5 or 3 or both but never neither
type anyOf interface {
	AnyOf() []reflect.StructField
}

// AndOneOf will generate the oneOf rule and retain the jsonschema rules for the struct that implements it
// `oneOf` can be used to factor out common parts of subschema and when *only one case* must be valid
// { "type": "number", "oneOf": [ { "multipleOf": 5 }, { "multipleOf": 3 } ]}
// In the example above, the input must be a number and must be either a multiple of 5 or 3 but not both and never neither
type andOneOf interface {
	AndOneOf() []reflect.StructField
}

// OneOf will overrule all jsonschema rules for the struct that implements it
// `oneOf` can be used to factor out common parts of subschema and when *only one case* must be valid
// { "oneOf": [ { "type": "number", "multipleOf": 5 }, { "type": "number", "multipleOf": 3 } ] }
// In the example above, the input must be a number and must be either a multiple of 5 or 3 but not both and never neither
type oneOf interface {
	OneOf() []reflect.StructField
}

// AllOf will generate the allOf rule and retain the jsonschema rules for the struct that implements it
// `allOf` is used to ensure that the data must be valid against *all cases*
// { "type": "number", "allOf": [ { "multipleOf": 5 }, { "multipleOf": 3 } ]}
// In the example above, the input must be a number and a multiple of 5 *and* 3
type andAllOf interface {
	AndAllOf() []reflect.StructField
}

// AllOf will overrule all jsonschema rules for the struct that implements it
// `allOf` is used to ensure that the data must be valid against *all cases*
// { "allOf": [ { "type": "number", "multipleOf": 5 }, { "type": "number", "multipleOf": 3 } ] }
// In the example above, the input must be a number and a multiple of 5 *and* 3
type allOf interface {
	AllOf() []reflect.StructField
}

// When AnyOf/OneOf/AllOf are implemented, the jsonschema for the implementing struct will be supplanted with
// exclusive anyOf/oneOf/allOf rules
func (r *Reflector) getExclusiveSubschemaForBooleanCases(definitions Definitions, t reflect.Type) *Type {

	var nonNilPointer interface{}
	t, nonNilPointer = getNonNilPointerTypeAndInterface(t)

	if t.Implements(anyOfType) {
		s := nonNilPointer.(anyOf).AnyOf()
		return &Type{AnyOf: r.getSubschemasForBooleanCases(definitions, s)}
	}
	if t.Implements(oneOfType) {
		s := nonNilPointer.(oneOf).OneOf()
		return &Type{OneOf: r.getSubschemasForBooleanCases(definitions, s)}
	}
	if t.Implements(allOfType) {
		s := nonNilPointer.(allOf).AllOf()
		return &Type{AllOf: r.getSubschemasForBooleanCases(definitions, s)}
	}

	return nil

}

// Append jsonschema rules from AndOneOf/AndAnyOf/AndAllOf interfaces
// to the jsonschema for the struct that implements them
func (r *Reflector) addSubschemasForBooleanCases(schema *Type, definitions Definitions, t reflect.Type) {
	if schema == nil {
		return
	}

	var nonNilPointer interface{}
	t, nonNilPointer = getNonNilPointerTypeAndInterface(t)

	if t.Implements(andAnyOfType) {
		s := nonNilPointer.(andAnyOf).AndAnyOf()
		schema.AnyOf = r.getSubschemasForBooleanCases(definitions, s)
	}
	if t.Implements(andOneOfType) {
		s := nonNilPointer.(andOneOf).AndOneOf()
		schema.OneOf = r.getSubschemasForBooleanCases(definitions, s)
	}
	if t.Implements(andAllOfType) {
		s := nonNilPointer.(andAllOf).AndAllOf()
		schema.AllOf = r.getSubschemasForBooleanCases(definitions, s)
	}
}

func (r *Reflector) getSubschemasForBooleanCases(definitions Definitions, s []reflect.StructField) []*Type {
	oneOfList := make([]*Type, 0)
	for _, oneType := range s {
		if oneType.Type == nil {
			oneOfList = append(oneOfList, &Type{Type: "null"})
		} else {
			oneOfList = append(oneOfList, r.reflectTypeToSchema(definitions, oneType.Type))
		}
	}
	return oneOfList
}
