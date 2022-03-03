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

func FloatPercentageValidator(f string) (err error) {
	if f, _ := strconv.ParseFloat(f, 64); f < 0 || f > 100 {
		err = errors.New("value must be in range 0-100")
	}
	return
}
