package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "somename",
		Price: 2,
		SKU:   "asfd",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
