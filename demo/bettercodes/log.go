package bettercodes

import "fmt"

type ILogger interface {
	Log(msg string)
}

type Logger struct{}

func (*Logger) Log(msg string) {
	fmt.Println(msg)
}
