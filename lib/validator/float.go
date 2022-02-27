package validator

import (
	"errors"
	"strconv"
)

func PositiveFloatValidator(f string) (err error) {
	if f, _ := strconv.ParseFloat(f, 64); f <= 0 {
		err = errors.New("value must be a positive float")
	}
	return
}
