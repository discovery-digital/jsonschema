package jsonschema

import (
	"reflect"
	"strings"
)

// Ensure we get a non-pointer type
func getNonPointerType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		return getNonPointerType(t)
	}

	return t
}

// We need to ensure we get the pointer type to a *value struct* and the interface from a *non-nil* pointer value
func getNonNilPointerTypeAndInterface(t reflect.Type) (reflect.Type, interface{}) {
	t = getNonPointerType(t)

	// This gives us a pointer value from the *value struct*
	pv := reflect.New(t)
	t = pv.Type()
	pvi := pv.Interface()

	return t, pvi
}

// `json:"-"` and `json:"omitempty"` make the field optional
func requiredFromJSONTags(tags []string) bool {
	if ignoredByJSONTags(tags) {
		return false
	}

	for _, tag := range tags[1:] {
		if tag == "omitempty" {
			return false
		}
	}
	return true
}

// Setting RequiredFromJSONSchemaTags to true allows usage of the `required` tag to make a field required
func requiredFromJSONSchemaTags(tags []string) bool {
	if ignoredByJSONSchemaTags(tags) {
		return false
	}
	for _, tag := range tags {
		if tag == "required" {
			return true
		}
	}
	return false
}

// `jsonschema:"optional"` will make the field optional
//
// The use case for this is when you are taking json input where validation on a field should be optional
// but you do not want to declare `omitempty` because you serialize the struct to json to a third party
// and the fields must exist (such as a field that's an int)
func remainsRequiredFromJSONSchemaTags(tags []string, currentlyRequired bool) bool {
	for _, tag := range tags {
		if tag == "optional" {
			return false
		}
	}
	return currentlyRequired
}

func ignoredByJSONTags(tags []string) bool {
	return tags[0] == "-"
}

func ignoredByJSONSchemaTags(tags []string) bool {
	return tags[0] == "-"
}

// getPackageNameFromPath splits path to struct and return last element which is package name
func getPackageNameFromPath(path string) string {
	pathSlices := strings.Split(path, "/")
	return pathSlices[len(pathSlices)-1]
}

// bool2bytes serializes bool to JSON
func bool2bytes(val bool) []byte {
	if val {
		return []byte("true")
	}
	return []byte("false")
}
