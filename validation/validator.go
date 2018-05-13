package validation

import (
	"fmt"
	"strings"
)

type Validator struct {
	errors  map[string][]error
	context string
}

func NewValidator() *Validator {
	return &Validator{
		errors:  map[string][]error{},
		context: "",
	}
}

func (v *Validator) WithContext(ctx string) *Validator {
	childValidator := NewValidator()
	childValidator.errors = v.errors
	childValidator.context = fmt.Sprintf("%s.%s", v.context, ctx)
	return childValidator
}

func (v *Validator) Add(err error) {
	v.errors[v.context] = append(v.errors[v.context], err)
}

func (v *Validator) HasError() bool {
	return len(v.errors) != 0
}

func (v *Validator) Error() string {
	var errorStrs []string
	for ctx, errors := range v.errors {
		for _, err := range errors {
			errorStrs = append(errorStrs, fmt.Sprintf("%s: %s", ctx, err.Error()))
		}
	}
	return strings.Join(errorStrs, "\n")
}
