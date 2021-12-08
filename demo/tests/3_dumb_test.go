package tests

import (
	"testing"

	"github.com/vietky/golang-tmpl/demo/codes"
)

func TestDumb(t *testing.T) {
	dumb := codes.NewDumb()
	v := dumb.Calculate(1, 2)
	if v != 3 {
		t.Errorf("v should be 3. Current Value is %v", v)
	}
}
