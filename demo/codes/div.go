package codes

import "errors"

func Div(a, b int) (float32, error) {
	if b == 0 {
		return 0.0, errors.New("division by zero")
	}
	return float32(a) / float32(b), nil
}
