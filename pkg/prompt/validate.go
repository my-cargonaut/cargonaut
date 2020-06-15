package prompt

import (
	"errors"
	"strings"
)

func valdatorChain(skipOnEmpty bool, validators ...ValidatorFunc) ValidatorFunc {
	return func(input string) error {
		// SKip validation if no input is present.
		if skipOnEmpty && input == "" {
			return nil
		}
		for _, v := range validators {
			if v == nil {
				continue
			} else if err := v(input); err != nil {
				return err
			}
		}
		return nil
	}
}

func validateRequired(input string) error {
	if input = strings.TrimSpace(input); input == "" {
		return errors.New("cannot be empty")
	}
	return nil
}
