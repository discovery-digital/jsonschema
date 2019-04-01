package testmodels

type MostInner struct {
	Foo string `json:"foo,omitempty"`
	Bar string `json:"bar"`
}
type Inner struct {
	Foo string `json:"foo"`
	Bar string `json:"bar,omitempty"`
	MostInner
}

type Root struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
	Inner
}