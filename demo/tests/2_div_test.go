package tests

import (
	"testing"

	"github.com/vietky/golang-tmpl/demo/codes"
)

func TestDivWithMultipleCases(t *testing.T) {
	var tests = []struct {
		a, b int
		want float32
	}{
		{6, 4, 1.5},
		{6, 3, 2.0},
	}
	for _, r := range tests {
		v, err := codes.Div(r.a, r.b)
		if err != nil {
			t.Errorf("%v/%v should not throw error %+v", r.a, r.b, err)
		}
		if v != r.want {
			t.Errorf("6/4 should be 1.5. Current value is %v", v)
		}
	}
}
