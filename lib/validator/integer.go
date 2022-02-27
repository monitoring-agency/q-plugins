package validator

import (
	"errors"
	"strconv"
)

func PositiveIntegerValidator(i string) (err error) {
	if i, _ := strconv.Atoi(i); i <= 0 {
		err = errors.New("value must be a positive int")
	}
	return
}
