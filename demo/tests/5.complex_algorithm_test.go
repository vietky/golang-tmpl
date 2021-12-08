package tests

import (
	"context"
	"testing"
	"time"

	"github.com/vietky/golang-tmpl/demo/bettercodes"
)

func TestComplexAlgorithm(t *testing.T) {
	complex := bettercodes.NewComplexAlgorithm()
	// complex := bettercodes.NewSimplerComplexAlgorithm()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	testDone := make(chan int)
	go func() {
		v := complex.DoComplexThings()
		testDone <- v
	}()

	select {
	case <-ctx.Done():
		t.Error("Timeout. Tests take too long")
	case v := <-testDone:
		if v != 7 {
			t.Error("ok")
		}
	}
}
