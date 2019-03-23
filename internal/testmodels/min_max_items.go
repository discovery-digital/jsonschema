package testmodels

type SliceTestType []string

func (st SliceTestType) MinItems() int {
	return 2
}

func (st SliceTestType) MaxItems() int {
	return 2
}
