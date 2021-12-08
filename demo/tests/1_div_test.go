package tests

import (
	"testing"

	"github.com/vietky/golang-tmpl/demo/codes"
)

func TestDiv(t *testing.T) {
	testDiv(t, 6, 4, 1.5)
}

func testDiv(t *testing.T, a, b int, want float32) {
	v, err := codes.Div(a, b)
	if err != nil {
		t.Errorf("6/4 should not throw error %+v", err)
	}
	if v != want {
		t.Errorf("6/4 should be 1.5. Current value is %v", v)
	}

}
