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
	if u == "" {
		return "", ErrInvalidUnit
	}
	// Accept any non-empty unit - units are defined in the database
	return u, nil
}

// IsInteger returns true if the quantity represents a whole number.
func (q Quantity) IsInteger() bool {
	return float64(q) == math.Trunc(float64(q))
}

// ValidateQuantityForIntegerOnly validates that the quantity is an integer
// if integerOnly is true. Returns an error if validation fails.
func ValidateQuantityForIntegerOnly(qty Quantity, integerOnly bool) error {
	if integerOnly && !qty.IsInteger() {
		return ErrQuantityMustBeInteger
	}
	return nil
}
