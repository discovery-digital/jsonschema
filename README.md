# Go JSON Schema Reflection

[![Build Status](https://travis-ci.org/alecthomas/jsonschema.png)](https://travis-ci.org/alecthomas/jsonschema)
[![Gitter chat](https://badges.gitter.im/alecthomas.png)](https://gitter.im/alecthomas/Lobby)
[![Go Report Card](https://goreportcard.com/badge/github.com/alecthomas/jsonschema)](https://goreportcard.com/report/github.com/alecthomas/jsonschema)
[![GoDoc](https://godoc.org/github.com/alecthomas/jsonschema?status.svg)](https://godoc.org/github.com/alecthomas/jsonschema)

This package can be used to generate [JSON Schemas](http://json-schema.org/latest/json-schema-validation.html) from Go types through reflection.

It supports arbitrarily complex types, including `interface{}`, maps, slices, etc.
And it also supports json-schema features such as minLength, maxLength, pattern, format and etc.

  * [Basic Example](#basic-example)
  * [Configurable behaviour](#configurable-behaviour)
    + [ExpandedStruct](#expandedstruct)
    + [Overrides](#overrides)
      - [type SchemaTagOverride](#type-schematagoverride)
      - [func GetSchemaTagOverride](#func--getschematagoverride)
      - [Example](#example)
  * [Subschema Support](#subschema-support)
    + [Boolean cases: `oneOf` / `anyOf` / `allOf`](#boolean-cases-oneof--anyof--allof)
      - [Inclusive usage (most common)](#inclusive-usage-most-common)
        * [Example](#example-1)
      - [Exclusive usage](#exclusive-usage)
        * [Example](#example-2)
    + [Conditional cases: `if/then/else`](#conditional-cases-ifthenelse)
      - [Example](#example-3)
  * [Other features](#other-features)
    + [Slice min/maxItems support](#slice-minmaxitems-support)
    + [`optional` tag value](#optional-tag-value)
    + [`switch` construct](#switch-construct)
      - [Example](#example-4)

## Basic Example

The following Go type:

```go
type TestUser struct {
  ID        int                    `json:"id"`
  Name      string                 `json:"name"`
  Nickname  *string                `json:"nickname",jsonschema="allowNull"`
  Friends   []int                  `json:"friends,omitempty"`
  Tags      map[string]interface{} `json:"tags,omitempty"`
  BirthDate time.Time              `json:"birth_date,omitempty"`
}
```

Results in following JSON Schema:

```go
jsonschema.Reflect(&TestUser{})
```

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$ref": "#/definitions/TestUser",
  "definitions": {
    "TestUser": {
      "type": "object",
      "properties": {
        "birth_date": {
          "type": "string",
          "format": "date-time"
        },
        "friends": {
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "id": {
          "type": "integer"
        },
        "name": {
          "type": "string"
        },
        "nickname": {
          "oneOf": [
            {
              "type": "string"
            },
            {
              "type": "null"
            }
          ]
        },
        "tags": {
          "type": "object",
          "patternProperties": {
            ".*": {
              "type": "object",
              "additionalProperties": true
            }
          }
        }
      },
      "additionalProperties": false,
      "required": ["id", "name", "nickname"]
    }
  }
}
```
## Configurable behaviour

The behaviour of the schema generator can be altered with parameters when a `jsonschema.Reflector`
instance is created.

```go
type Reflector struct {
	// AllowAdditionalProperties will cause the Reflector to generate a schema
	// with additionalProperties to 'true' for all struct types. This means
	// the presence of additional keys in JSON objects will not cause validation
	// to fail. Note said additional keys will simply be dropped when the
	// validated JSON is unmarshaled.
	AllowAdditionalProperties bool

	// RequiredFromJSONSchemaTags will cause the Reflector to generate a schema
	// that requires any key tagged with `jsonschema:required`, overriding the
	// default of requiring any key *not* tagged with `json:,omitempty`.
	RequiredFromJSONSchemaTags bool

	// ExpandedStruct will cause the toplevel definitions of the schema not
	// be referenced itself to a definition.
	ExpandedStruct bool

	// Overrides is of interface SchemaTagOverride and will be used to override any jsonschema tags on existing fields
	// The expected use case is for shared nested structs where validation is stricter on certain fields
	// For example a shared nested struct with field `Species` and tag `enum=Human|Dog|Alien` may be used by
	// applications that want to declare a stricter tag `required,enum=Dog`
	Overrides SchemaTagOverride
}
```

### ExpandedStruct

If set to ```true```, makes the top level struct not to reference itself in the definitions. But type passed should be a struct type.

eg.

```go
type GrandfatherType struct {
	FamilyName string `json:"family_name"`
}

type SomeBaseType struct {
	SomeBaseProperty int `json:"some_base_property"`
	somePrivateBaseProperty            string `json:"i_am_private"`
	SomeIgnoredBaseProperty            string `json:"-"`
	SomeSchemaIgnoredProperty          string `jsonschema:"-"`
	SomeUntaggedBaseProperty           bool   
	someUnexportedUntaggedBaseProperty bool
	Grandfather                        GrandfatherType `json:"grand"`
}
```

will output:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "required": [
    "some_base_property",
    "grand",
    "SomeUntaggedBaseProperty"
  ],
  "properties": {
    "SomeUntaggedBaseProperty": {
      "type": "boolean"
    },
    "grand": {
      "$schema": "http://json-schema.org/draft-07/schema#",
      "$ref": "#/definitions/GrandfatherType"
    },
    "some_base_property": {
      "type": "integer"
    }
  },
  "type": "object",
  "definitions": {
    "GrandfatherType": {
      "required": [
        "family_name"
      ],
      "properties": {
        "family_name": {
          "type": "string"
        }
      },
      "additionalProperties": false,
      "type": "object"
    }
  }
}
```

### Overrides
There are cases where you have a model that works with existing jsonschema tags that are overly strict.

An implementation of `SchemaTagOverride` can be passed to `jsonschema.Reflector` to do this.

#### type SchemaTagOverride

```go
type SchemaTagOverride interface {
	// Set should be given:
	// targetStruct - struct that contains the field to be overridden
	// targetField - name of the field that is to be overridden
	// tag - the provided jsonschema tag
	Set(targetStruct interface{}, targetField string, tag string) error
	// Get is used by this library to retrieve overrides
	Get(targetStructType reflect.Type, targetField string) string
}
```
SchemaTagOverride is a mechanism to allow jsonschema tag overrides


#### func  GetSchemaTagOverride

```go
func GetSchemaTagOverride() SchemaTagOverride
```
GetSchemaTagOverride returns initialized SchemaTagOverride


#### Example
Given model
```go
type Human struct {
	Name   string `json:"name" jsonschema:"notEmpty"`
	Sex    string `json:"sex" jsonschema:"enum=a|b|c|d|e|f|g"`
}
```

To override the rules for `Sex`, we would do the following:
```go
	sto := jsonschema.GetSchemaTagOverride()
	_ = sto.Set(Human{}, "Sex", "enum=foo|bar|baz")

	a := jsonschema.Reflector{
		AllowAdditionalProperties: true,
		Overrides: sto,
	}

	schema := a.Reflect(&Human{})
	out, _ := json.MarshalIndent(&schema, "", "\t")
	fmt.Println(string(out))
```

The output:
```
{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"$ref": "#/definitions/main.Human",
	"definitions": {
		"main.Human": {
			"required": [
				"name",
				"sex"
			],
			"properties": {
				"name": {
					"pattern": "^\\S",
					"type": "string"
				},
				"sex": {
					"enum": [
						"foo",
						"bar",
						"baz"
					],
					"type": "string"
				}
			},
			"additionalProperties": true,
			"type": "object"
		}
	}
}
```

## Subschema Support
### Boolean cases: `oneOf` / `anyOf` / `allOf`
* `oneOf` can be used to factor out common parts of subschema and when *only one case* must be valid
* `anyOf` is used to ensure that the data must be valid against *at least one* of the cases *or more*
* `allOf` is used to ensure that the data must be valid against *all cases*

There are two interfaces for each subschema that can be implemented due to behavior differences.
#### Inclusive usage (most common)
When `AndOneOf` is used, jsonschema is generated from the struct and the output of the method
```go
	AndOneOf() []reflect.StructField
	AndAnyOf() []reflect.StructField
	AndAllOf() []reflect.StructField
```


##### Example
In this example, we have a common schema and then require **mutually exclusive** schema based on the value of `type`

Note that `omitempty` or `jsonschema:"optional"` should be specified on `Human` due to `Renter` not requiring it.
```go
type Human struct {
	Name string `json:"name" jsonschema:"notEmpty"`
	Age int `json:"age"`
	Type string `json:"type"`
	Networth int `json:"networth"`
	SocialSecurity string `json:"socialSecurity,omitempty"`
}

type Owner struct {
    Type string `json:"type" jsonschema:"enum=owner"`
	Networth int `json:"networth" jsonschema:"minimum=200000"`
	SocialSecurity string `json:"socialSecurity" jsonschema:"notEmpty"`
}

type Renter struct {
    Type string `json:"type" jsonschema:"enum=renter"`
	Networth int `json:"networth" jsonschema:"minimum=100"`
}


func (h Human) AndOneOf() []reflect.StructField {
	return []reflect.StructField{
		reflect.StructField{ Type: reflect.TypeOf(Owner{}) },
		reflect.StructField{ Type: reflect.TypeOf(Renter{}) },
	}
}
```

The generated jsonschema expects when the `type` value is:
* `owner` - have a networth higher than 200,000 and their social security number
* `renter` - have a networth higher than 100

But it will require all to provide a value for `name`, `age`, `type`, and `networth`.

```
{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"$ref": "#/definitions/main.Human",
	"definitions": {
		"main.Human": {
			"required": [
				"name",
				"age",
				"type",
				"networth"
			],
			"properties": {
				"age": {
					"type": "integer"
				},
				"name": {
					"pattern": "^\\S",
					"type": "string"
				},
				"networth": {
					"type": "integer"
				},
				"socialSecurity": {
					"type": "string"
				},
				"type": {
					"type": "string"
				}
			},
			"additionalProperties": true,
			"type": "object",
			"oneOf": [
				{
					"$schema": "http://json-schema.org/draft-07/schema#",
					"$ref": "#/definitions/main.Owner"
				},
				{
					"$schema": "http://json-schema.org/draft-07/schema#",
					"$ref": "#/definitions/main.Renter"
				}
			]
		},
		"main.Owner": {
			"required": [
				"type",
				"networth",
				"socialSecurity"
			],
			"properties": {
				"networth": {
					"minimum": 200000,
					"type": "integer"
				},
				"socialSecurity": {
					"pattern": "^\\S",
					"type": "string"
				},
				"type": {
					"enum": [
						"owner"
					],
					"type": "string"
				}
			},
			"additionalProperties": true,
			"type": "object"
		},
		"main.Renter": {
			"required": [
				"type",
				"networth"
			],
			"properties": {
				"networth": {
					"minimum": 100,
					"type": "integer"
				},
				"type": {
					"enum": [
						"renter"
					],
					"type": "string"
				}
			},
			"additionalProperties": true,
			"type": "object"
		}
	}
}
```

#### Exclusive usage
When `OneOf` is used, jsonschema is *exclusively* generated from the output of the method call and ignores any jsonschema rules on the struct that implemented the method.
```
	OneOf() []reflect.StructField
	AnyOf() []reflect.StructField
	AllOf() []reflect.StructField
```

##### Example
Let's say we have a payload, but we only want **mutually exclusive** schemas with no common factor
```go
type Payload struct {
	Contents string `json:"contents" jsonschema:"enum=hello"`
}

type Registration struct {
	Email string `json:"contents" jsonschema:"format=email"`
}

type Comment struct {
	Text string `json:"contents" jsonschema:"notEmpty"`
}

func (p Payload) OneOf() []reflect.StructField {
	return []reflect.StructField{
		reflect.StructField{ Type: reflect.TypeOf(Registration{}) },
		reflect.StructField{ Type: reflect.TypeOf(Comment{}) },
	}
}
```

As you can see in the `jsonschema` output below, the `Payload` jsonschema tags are ignored.
```
{
	"oneOf": [
		{
			"$schema": "http://json-schema.org/draft-07/schema#",
			"$ref": "#/definitions/main.Registration"
		},
		{
			"$schema": "http://json-schema.org/draft-07/schema#",
			"$ref": "#/definitions/main.Comment"
		}
	],
	"definitions": {
		"main.Comment": {
			"required": [
				"contents"
			],
			"properties": {
				"contents": {
					"pattern": "^\\S",
					"type": "string"
				}
			},
			"additionalProperties": true,
			"type": "object"
		},
		"main.Registration": {
			"required": [
				"contents"
			],
			"properties": {
				"contents": {
					"type": "string",
					"format": "email"
				}
			},
			"additionalProperties": true,
			"type": "object"
		}
	}
}
```

### Conditional cases: `if/then/else`
The struct must implement the method below:

```go
    IfThenElse() SchemaCondition
```

```go
type SchemaCondition struct {
	If   reflect.StructField
	Then interface{}
	Else interface{}
}
```

* **If**: A `reflect.StructField` that defines the condition to be met.
* **Then**: A type that will be converted to a jsonschema subschema and evaluated if the condition is met
* **Else**: A type that will be converted to a jsonschema subschema and evaluated if the condition is not met

#### Example 
If a payload has type set to web, we evaluate the payload against the jsonschema generated from the WebApp struct otherwise we evaluate it against the jsonschema generated from the MobileApp struct.

Note `browser` and `device` on `Application` are optional due to either value being provided dependent on the type.
```go
type Application struct {
	Type string `json:"type" jsonschema:"required"`
	Browser string `json:"browser,omitempty"`
	Device string `json:"device,omitempty"`
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

func (app Application) IfThenElse() SchemaCondition {
	conditionField, _ := reflect.TypeOf(ApplicationValidation{}).FieldByName("Type")
	return SchemaCondition{
		If: conditionField,
		Then: WebApp{},browser
		Else: MobileApp{},
	}
}
```

Output:
```
{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"$ref": "#/definitions/main.Application",
	"definitions": {
		"main.Application": {
			"required": [
				"type"
			],
			"properties": {
				"browser": {
					"type": "string"
				},
				"device": {
					"type": "string"
				},
				"type": {
					"type": "string"
				}
			},
			"additionalProperties": true,
			"type": "object",
			"if": {
				"properties": {
					"type": {
						"enum": [
							"web"
						]
					}
				}
			},
			"then": {
				"$schema": "http://json-schema.org/draft-07/schema#",
				"$ref": "#/definitions/main.WebApp"
			},
			"else": {
				"$schema": "http://json-schema.org/draft-07/schema#",
				"$ref": "#/definitions/main.MobileApp"
			}
		},
		"main.MobileApp": {
			"required": [
				"device"
			],
			"properties": {
				"device": {
					"type": "string"
				}
			},
			"additionalProperties": true,
			"type": "object"
		},
		"main.WebApp": {
			"required": [
				"browser"
			],
			"properties": {
				"browser": {
					"type": "string"
				}
			},
			"additionalProperties": true,
			"type": "object"
		}
	}
}
```


## Other features
### Slice min/maxItems support
If jsonschema.Reflector is provided a typed slice collection that implements MinItems and/or MaxItems as depicted below, jsonschema will be generated to validate the slice and its contents.

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/discovery-digital/jsonschema"
)

type Goo struct {
	Bar string `json:"bar" jsonschema:"required"`
}

type GooCollection []Goo
func (fc GooCollection) MinItems() int {
	return 2
}

func (fc GooCollection) MaxItems() int {
	return 7
}

func main() {
	goo := GooCollection{}

	s := jsonschema.Reflect(goo)
	schema, _ := json.MarshalIndent(s, "", "\t")
	fmt.Println(string(schema))
}
```

Generated jsonschema
```go
{
  "definitions": {
    "main.Goo": {
      "additionalProperties": false,
      "properties": {
        "bar": {
          "type": "string"
        }
      },
      "required": [
        "bar"
      ],
      "type": "object"
    }
  },
  "items": {
    "$ref": "#/definitions/main.Goo",
    "$schema": "http://json-schema.org/draft-07/schema#"
  },
  "maxItems": 7,
  "minItems": 2,
  "type": "array"
}
```

### `optional` tag value
The `optional` jsonschema tag value can be used when you are taking json input where validation on a field should be optional
but you do not want to declare `omitempty` because you serialize the struct to json to a third party
and the fields must exist (such as a field that's an int)

### `switch` construct
The struct must implement the method below:
```
	Case() SchemaSwitch
```

```
// SchemaSwitch holds data for emulating switch case over some field value
type SchemaSwitch struct {
	// ByField = the name of the field you wish to evaluate (ex: "species")
	ByField string
	// Each key = the value for the field being evaluated (ex: "turtle")
	// Each value = the struct that holds the jsonschema tags to validate against when it is that value (ex: Turtle{})
	Cases   map[string]interface{}
}
```

#### Example
Imagine that we want to evaluate a payload against different schema depending on the value of a specific field.

We have a payload where the value of type can be `int`, `string`, or `bool`. Depending on what value this field is, we want to apply different validation logic to the value of `payload`.

```go
type ExampleCase struct {
	Type string `json:"type" jsonschema:"required;enum="int|string|bool"`
	Payload interface{} `json:"payload" jsonschema:"-"`
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

	return SchemaSwitch{ ByField: "type", Cases:   cases}
}
```

Generated jsonschema:
```go
{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"$ref": "#/definitions/main.ExampleCase",
	"definitions": {
		"main.BoolPayload": {
			"required": [
				"payload"
			],
			"properties": {
				"payload": {
					"type": "boolean"
				}
			},
			"additionalProperties": true,
			"type": "object"
		},
		"main.ExampleCase": {
			"required": [
				"type"
			],
			"properties": {
				"type": {
					"type": "string"
				}
			},
			"additionalProperties": true,
			"type": "object",
			"oneOf": [
				{
					"if": {
						"properties": {
							"type": {
								"enum": [
									"string"
								]
							}
						}
					},
					"then": {
						"$schema": "http://json-schema.org/draft-07/schema#",
						"$ref": "#/definitions/main.StringPayload"
					},
					"else": {
						"properties": {
							"type": {
								"enum": [
									"string"
								]
							}
						}
					}
				},
				{
					"if": {
						"properties": {
							"type": {
								"enum": [
									"bool"
								]
							}
						}
					},
					"then": {
						"$schema": "http://json-schema.org/draft-07/schema#",
						"$ref": "#/definitions/main.BoolPayload"
					},
					"else": {
						"properties": {
							"type": {
								"enum": [
									"bool"
								]
							}
						}
					}
				},
				{
					"if": {
						"properties": {
							"type": {
								"enum": [
									"int"
								]
							}
						}
					},
					"then": {
						"$schema": "http://json-schema.org/draft-07/schema#",
						"$ref": "#/definitions/main.IntPayload"
					},
					"else": {
						"properties": {
							"type": {
								"enum": [
									"int"
								]
							}
						}
					}
				}
			]
		},
		"main.IntPayload": {
			"required": [
				"payload"
			],
			"properties": {
				"payload": {
					"type": "integer"
				}
			},
			"additionalProperties": true,
			"type": "object"
		},
		"main.StringPayload": {
			"required": [
				"payload"
			],
			"properties": {
				"payload": {
					"type": "string"
				}
			},
			"additionalProperties": true,
			"type": "object"
		}
	}
}
```

The way the `if/then/else` construct works in jsonschema is that it will evaluate to `true` when:
- the `if` or `else` conditional is met
- the `else` statement is absent (regardless if `if` evaluates to true or false)

By restating the `if` specification for `else`, we force a failure in jsonschema when the `if` evaluation fails as the `else` evaluation will also fails. Combining this functionality with `oneOf` gives us effectively a switch statement.

```
{
   "type": "bool",
   "payload": true
}
```
This example payload will be evaluated in the `oneOf` block for each `if/then/else`. It will satisfy the condition where the value of `type` equals the string `bool` which will subject the payload to be evaluated against the generated jsonschema for `BoolPayload`.

```
{
  "if": {
    "properties": {
      "type": { "enum": [ "bool" ] }
    }
  },
  "then": {
    "$ref": "#/definitions/jsonschema.BoolPayload",
    "$schema": "http://json-schema.org/draft-07/schema#"
  },
  "else": {
    "properties": {
      "type": { "enum": [ "bool" ] }
    }
  }
}
```

`oneOf` is generally adequate for most conditional evaluations, but validators will validate the payload against all cases and provide validation errors for all cases which may be confusing. For example, we can get by with simple cases if we make the following adjustments: 

`oneOf` would evaluate each schema against our payload:
```
  "oneOf": [
      {"$ref": "#/definitions/jsonschema.IntPayload","$schema": "http://json-schema.org/draft-07/schema#"},
      {"$ref": "#/definitions/jsonschema.StringPayload","$schema": "http://json-schema.org/draft-07/schema#"},
      {"$ref": "#/definitions/jsonschema.BoolPayload","$schema": "http://json-schema.org/draft-07/schema#"}
   ]
```

We also expand our individual schema to require the value of `type` to be a certain string:
```
        "jsonschema.StringPayload": {
            "additionalProperties": false,
            "properties": {
                "payload": {
                    "type": "string"
                },
                "type": {
                     "enum": ["int"]
                }
            },
            "required": [
                "payload",
                "type"
            ],
            "type": "object"
        }
```

This construct would satisfy simple cases where we want to make sure a different schema is evaluated depending on the value of `type`. However, since the validator will evaluate the given payload against *each* case, as there is no mechanism to rule out its evaluation completely, we will receive validation errors for `StringPayload`, `IntPayload`, and `BoolPayload` even when we satisfy `BoolPayload` partially. When we add `if/then/else` to `oneOf`, we provide the mechanism to rule out the evaluation of a schema completely and return better validation errors to clients as a result.


