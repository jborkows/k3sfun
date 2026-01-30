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
	return u, nil
}

// IsInteger returns true if the quantity represents a whole number.
func (q Quantity) IsInteger() bool {
	return float64(q) == math.Trunc(float64(q))
}
