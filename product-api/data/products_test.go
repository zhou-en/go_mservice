package data

import (
	"testing"
)

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "test",
		Price: 1.0,
		SKU:   "abs-abc-wer",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
