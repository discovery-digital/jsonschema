package testmodels

type MostInner struct {
	Foo string `json:"foo,omitempty"`
	Bar string `json:"bar,omitempty"`
}
type Inner struct {
	MostInner
	Foo string `json:"foo,omitempty"`
	Bar string `json:"bar,omitempty"`
}

type Root struct {
	Inner
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}
