package testmodels

type TestFamilyMember struct {
	Children []TestFamilyMember `json:"children"`
}
