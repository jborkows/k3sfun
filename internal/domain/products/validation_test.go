package products

import "testing"

func TestNormalizeGroupName(t *testing.T) {
	if _, err := NormalizeGroupName("   "); err == nil {
		t.Fatalf("expected error")
	}
	got, err := NormalizeGroupName(" vegetables ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "vegetables" {
		t.Fatalf("expected trimmed name, got %q", got)
	}
}

func TestNormalizeProductName(t *testing.T) {
	if _, err := NormalizeProductName("   "); err == nil {
		t.Fatalf("expected error")
	}
	got, err := NormalizeProductName(" carrots ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "carrots" {
		t.Fatalf("expected trimmed name, got %q", got)
	}
}

func TestNormalizeUnit(t *testing.T) {
	ok := []Unit{UnitKG, UnitLiter, UnitPiece, UnitGram}
	for _, u := range ok {
		if _, err := NormalizeUnit(u); err != nil {
			t.Fatalf("unexpected error for %q: %v", u, err)
		}
	}
	if _, err := NormalizeUnit("bad"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestIsInteger(t *testing.T) {
	tests := []struct {
		value    float64
		expected bool
	}{
		{0, true},
		{1, true},
		{-1, true},
		{10, true},
		{100.0, true},
		{0.5, false},
		{1.1, false},
		{1.9, false},
		{-0.5, false},
		{0.1, false},
	}
	for _, tc := range tests {
		got := IsInteger(tc.value)
		if got != tc.expected {
			t.Errorf("IsInteger(%v) = %v, want %v", tc.value, got, tc.expected)
		}
	}
}

func TestValidateQuantityForIntegerOnly(t *testing.T) {
	// When integerOnly is false, any value should be valid
	if err := ValidateQuantityForIntegerOnly(1.5, false); err != nil {
		t.Errorf("expected no error for non-integer-only product, got %v", err)
	}
	if err := ValidateQuantityForIntegerOnly(1, false); err != nil {
		t.Errorf("expected no error for integer value on non-integer-only product, got %v", err)
	}

	// When integerOnly is true, only integers should be valid
	if err := ValidateQuantityForIntegerOnly(1, true); err != nil {
		t.Errorf("expected no error for integer value on integer-only product, got %v", err)
	}
	if err := ValidateQuantityForIntegerOnly(10.0, true); err != nil {
		t.Errorf("expected no error for 10.0 on integer-only product, got %v", err)
	}
	if err := ValidateQuantityForIntegerOnly(0, true); err != nil {
		t.Errorf("expected no error for 0 on integer-only product, got %v", err)
	}

	// Non-integer values should fail for integer-only products
	if err := ValidateQuantityForIntegerOnly(1.5, true); err == nil {
		t.Errorf("expected error for 1.5 on integer-only product")
	}
	if err := ValidateQuantityForIntegerOnly(0.1, true); err == nil {
		t.Errorf("expected error for 0.1 on integer-only product")
	}
}
