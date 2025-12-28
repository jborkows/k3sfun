package products

import (
	"math"
	"strings"
)

func NormalizeGroupName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrNameRequired
	}
	return name, nil
}

func NormalizeProductName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrNameRequired
	}
	return name, nil
}

func NormalizeUnit(u Unit) (Unit, error) {
	u = Unit(strings.TrimSpace(string(u)))
	switch u {
	case UnitKG, UnitLiter, UnitPiece, UnitGram:
		return u, nil
	default:
		return "", ErrInvalidUnit
	}
}

// IsInteger returns true if the float value represents a whole number.
func IsInteger(v float64) bool {
	return v == math.Trunc(v)
}

// ValidateQuantityForIntegerOnly validates that the quantity is an integer
// if integerOnly is true. Returns an error if validation fails.
func ValidateQuantityForIntegerOnly(qty float64, integerOnly bool) error {
	if integerOnly && !IsInteger(qty) {
		return ErrQuantityMustBeInteger
	}
	return nil
}
