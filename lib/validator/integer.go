package validator

import (
	"errors"
	"math"
	"strconv"
)

func PositiveIntegerValidator(i string) (err error) {
	if i, _ := strconv.Atoi(i); i <= 0 {
		err = errors.New("value must be a positive int")
	}
	return
}

func PortValidator(i string) (err error) {
	if i, _ := strconv.Atoi(i); i <= 0 || i > int(math.Pow(2, 16)-1) {
		err = errors.New("value must be in the valid port range")
	}
	return
}

func IntPercentageValidator(i string) (err error) {
	if i, _ := strconv.Atoi(i); i < 0 || i > 100 {
		err = errors.New("value must be in range 0-100")
	}
	return
}
