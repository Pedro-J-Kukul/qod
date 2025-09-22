// Filename: Internal/validator/validator.go
// Description: This file contains validation functions for sending to the database.

package validator

import "slices"

type Validator struct {
	Errors map[string]string
}

// NewValidator creates a new Validator instance.
func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Function to see if the Validator's map contains any entries
func (v *Validator) IsEmpty() bool {
	return len(v.Errors) == 0
}

// Adds an error enhtry to the Validator's error map
func (v *Validator) AddError(key string, message string) {
	_, exists := v.Errors[key]
	if !exists {
		v.Errors[key] = message
	}
}

// if any validation check returns fall, make an entry to the validator's error map
func (v *Validator) Check(acceptable bool, key string, message string) {
	if !acceptable {
		v.AddError(key, message)
	}
}

// Function to check if a value is in a list of permitted values
func (v *Validator) PermittedValue(value string, permittedValues ...string) bool {
	return slices.Contains(permittedValues, value)
}
