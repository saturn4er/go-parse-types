package test_examples

import "github.com/saturn4er/go-parse-types/test_examples/sub_pack"

type D interface {
	Hello()
}

type A struct {
	sub_pack.B
	ALD *int
}

func (a *A) SomeMethod(b int) error {
	return nil
}

func TestFunc(a int) error {
	return nil
}
