package bettercodes

import "github.com/vietky/golang-tmpl/demo/codes"

type AnotherDumb struct {
	logger ILogger
}

func NewAnotherDumb(logger ILogger) codes.IDumb {
	return &AnotherDumb{
		logger,
	}
}

func (d *AnotherDumb) Hello() {
	d.logger.Log("hello")
}

func (d *AnotherDumb) Calculate(a, b int) int {
	return a + b
}
