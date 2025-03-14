package domain

import (
	"errors"
	"strings"
)

func ErrorsJoin(errs ...error) error {

	var nonNilErrors []string
	for _, err := range errs {
		if err != nil {
			nonNilErrors = append(nonNilErrors, err.Error())
		}
	}
	if len(nonNilErrors) == 0 {
		return nil
	}
	return errors.New(strings.Join(nonNilErrors, ", "))

}
