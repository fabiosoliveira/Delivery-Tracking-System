package valueObject

import (
	"errors"
	"regexp"
)

type Email string

func (e *Email) ValidateEmail() error {
	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	if !re.MatchString(string(*e)) {
		return errors.New("invalid email")
	}

	return nil
}
