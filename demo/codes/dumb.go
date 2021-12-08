package codes

import "fmt"

type IDumb interface {
	Hello()
	Calculate(a, b int) int
}

type Dumb struct {
}

func NewDumb() IDumb {
	return &Dumb{}
}

func (d *Dumb) Hello() {
	fmt.Println("hello")
}

func (d *Dumb) Calculate(a, b int) int {
	result := a
	for i := 0; i < b; i++ {
		result += a
	}
	return result
}
