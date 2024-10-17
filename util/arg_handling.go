package util

import (
	"errors"
	"fmt"
)

func ExpectOneOptional[T any](optionalArgs []T, argName string, funcName string) (*T, error) {

	if len(optionalArgs) == 0 {
		return nil, nil
	} else if len(optionalArgs) > 1 {
		return nil, errors.New(fmt.Sprintf("Cannot provide more than one %T to %s", optionalArgs[0], funcName))
	}
	return &optionalArgs[0], nil
}
