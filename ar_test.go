package jsarray

import (
	"fmt"
	//	"testing"
)

func ExampleMarshal() {
	type Person struct {
		Name   string
		ID     uint
		Rating uint
	}

	f := []Person{
		{"Andrew", 0, 0},
		{"", 1, 0},
		{":¬)", 4, 0},
		{"", 0, 10},
		{"", 0, 0},
	}

	bt, err := Marshal(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bt))

	// Output:
	// [["Andrew"],[,1],[":¬)",4],[,,10],]
}
