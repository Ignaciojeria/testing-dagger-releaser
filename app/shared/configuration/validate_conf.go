package configuration

import (
	"fmt"
	"reflect"
	"strings"
)

func validateConfig[T any](conf T) (T, error) {
	var validationErrors []error

	val := reflect.ValueOf(conf)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i).String()

		requiredTag := field.Tag.Get("required")
		if requiredTag == "true" && value == "" {
			validationErrors = append(validationErrors, fmt.Errorf("%s is required but not set", field.Name))
		}
	}
	if len(validationErrors) > 0 {
		// Convert errors to strings
		var errorStrings []string
		for _, err := range validationErrors {
			errorStrings = append(errorStrings, err.Error())
		}
		return conf, fmt.Errorf("configuration errors:\n%s", strings.Join(errorStrings, "\n"))
	}
	return conf, nil
}
