package jsonschema_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/discovery-digital/jsonschema"
	"github.com/discovery-digital/jsonschema/internal/testmodels"
)

type testSet struct {
	reflector *jsonschema.Reflector
	fixture   string
	actual    interface{}
}

var schemaGenerationTests = []testSet{
	{&jsonschema.Reflector{}, "fixtures/defaults.json", testmodels.TestUser{}},
	{&jsonschema.Reflector{AllowAdditionalProperties: true}, "fixtures/allow_additional_props.json", testmodels.TestUser{}},
	{&jsonschema.Reflector{RequiredFromJSONSchemaTags: true}, "fixtures/required_from_jsontags.json", testmodels.TestUser{}},
	{&jsonschema.Reflector{ExpandedStruct: true}, "fixtures/defaults_expanded_toplevel.json", testmodels.TestUser{}},
	{&jsonschema.Reflector{}, "fixtures/test_one_of_default.json", testmodels.TestUserOneOf{}},
	{&jsonschema.Reflector{}, "fixtures/test_versioned_packages.json", testmodels.TestVersionedPackages{}},
	{&jsonschema.Reflector{}, "fixtures/if_then_else.json", testmodels.Application{}},
	{&jsonschema.Reflector{}, "fixtures/case.json", testmodels.ExampleCase{}},
	{&jsonschema.Reflector{}, "fixtures/test_min_max_items.json", testmodels.SliceTestType{}},
}

func TestSchemaGeneration(t *testing.T) {
	for _, tt := range schemaGenerationTests {
		runTests(t, tt)
	}
}
func TestOverrides(t *testing.T) {
	override := jsonschema.GetSchemaTagOverride()
	override.Set(testmodels.Hardware{}, "Brand", "enum=microsoft|apple|lenovo|dell")

	test := testSet{
		reflector: &jsonschema.Reflector{Overrides: override},
		fixture:   "fixtures/override_jsonschema_tag.json",
		actual:    testmodels.TestUserOneOf{},
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
