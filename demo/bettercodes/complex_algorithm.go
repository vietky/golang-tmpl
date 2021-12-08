package bettercodes

import (
	"time"

)

const Timeout = 1

type IComplexAlgorithm interface {
	DoComplexThings() int
}

type ComplexAlgorithm struct {
}

func NewComplexAlgorithm() IComplexAlgorithm {
	return &ComplexAlgorithm{}
}

func (c *ComplexAlgorithm) DoComplexThings() int {
	c.a()
	c.b()
	c.c()
	return c.complex()
}

func (*ComplexAlgorithm) a() {
	time.Sleep(time.Duration(Timeout)*time.Second)
	return
}

func (*ComplexAlgorithm) b() {
	time.Sleep(time.Duration(Timeout)*time.Second)
	return
}

func (*ComplexAlgorithm) c() {
	time.Sleep(time.Duration(Timeout)*time.Second)
	return
}

func (*ComplexAlgorithm) complex() int {
	return 7
}