package jsonschema_test

import (
	"bytes"
	"encoding/json"
	"github.com/discovery-digital/jsonschema"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

type testSet struct {
	reflector *jsonschema.Reflector
	fixture   string
	actual    interface{}
}

var schemaGenerationTests = []testSet{
	{&jsonschema.Reflector{}, "fixtures/defaults.json", TestUser{}},
	{&jsonschema.Reflector{AllowAdditionalProperties: true}, "fixtures/allow_additional_props.json", TestUser{}},
	{&jsonschema.Reflector{RequiredFromJSONSchemaTags: true}, "fixtures/required_from_jsontags.json", TestUser{}},
	{&jsonschema.Reflector{ExpandedStruct: true}, "fixtures/defaults_expanded_toplevel.json", TestUser{}},
	{&jsonschema.Reflector{}, "fixtures/test_one_of_default.json", TestUserOneOf{}},
	{&jsonschema.Reflector{}, "fixtures/test_versioned_packages.json", TestVersionedPackages{}},
	{&jsonschema.Reflector{}, "fixtures/if_then_else.json", Application{}},
	{&jsonschema.Reflector{}, "fixtures/case.json", ExampleCase{}},
	{&jsonschema.Reflector{}, "fixtures/test_min_max_items.json", SliceTestType{}},
}

func TestSchemaGeneration(t *testing.T) {
	for _, tt := range schemaGenerationTests {
		runTests(t, tt)
	}
}
func TestOverrides(t *testing.T) {
	override := jsonschema.GetSchemaTagOverride()
	override.Set(Hardware{}, "Brand", "enum=microsoft|apple|lenovo|dell")

	test := testSet{
		reflector: &jsonschema.Reflector{Overrides: override},
		fixture:   "fixtures/override_jsonschema_tag.json",
		actual:    TestUserOneOf{},
	}

	runTests(t, test)
}

func runTests(t *testing.T, tt testSet) {
	name := strings.TrimSuffix(filepath.Base(tt.fixture), ".json")
	t.Run(name, func(t *testing.T) {
		f, err := ioutil.ReadFile(tt.fixture)
		if err != nil {
			t.Errorf("ioutil.ReadAll(%s): %s", tt.fixture, err)
			return
		}

		actualSchema := tt.reflector.Reflect(tt.actual)

		actualJSON, err := json.Marshal(actualSchema)
		if err != nil {
			t.Errorf("json.MarshalIndent(%v, \"\", \"  \"): %v", actualJSON, err)
			return
		}
		actualJSON = sanitizeExpectedJson(actualJSON)
		cleanExpectedJSON := sanitizeExpectedJson(f)

		if !bytes.Equal(cleanExpectedJSON, actualJSON) {

			t.Errorf("reflector %+v wanted schema %s, got %s", tt.reflector, cleanExpectedJSON, actualJSON)
		}
	})
}

func sanitizeExpectedJson(expectedJSON []byte) []byte {
	var js interface{}
	json.Unmarshal(expectedJSON, &js)
	clean, _ := json.MarshalIndent(js, "", "  ")
	return clean
}
