package products

import (
	"strconv"
	"time"
)

type ProductID int64
type GroupID int64
type Unit string

// Quantity represents a product quantity value.
// It wraps float64 to provide type safety and domain-specific methods.
type Quantity float64

// Float64 returns the underlying float64 value.
func (q Quantity) Float64() float64 {
	return float64(q)
}

// String returns the quantity formatted as a string.
func (q Quantity) String() string {
	return strconv.FormatFloat(float64(q), 'f', -1, 64)
}

const (
	UnitKG      Unit = "kg"
	UnitLiter   Unit = "litr"
	UnitPiece   Unit = "sztuk"
	UnitGram    Unit = "gramy"
	UnitPackage Unit = "opakowanie"
	UnitBunch   Unit = "pęczek"
	UnitBulb    Unit = "główki"
)

const MaxProductsPageSize int64 = 30

type Group struct {
	ID   GroupID
	Name string
}

type Product struct {
	ID          ProductID
	Name        string
	IconKey     string
	GroupID     *GroupID
	GroupName   string
	Quantity    Quantity
	Unit        Unit
	MinQuantity Quantity
	Missing     bool
	IntegerOnly bool
	UpdatedAt   time.Time
}

type ProductFilter struct {
	OnlyMissingOrLow bool
	NameQuery        string
	GroupIDs         []GroupID
	Limit            int64
	Offset           int64
}

type NewProduct struct {
	Name        string
	IconKey     string
	GroupID     *GroupID
	Quantity    Quantity
	Unit        Unit
	MinQuantity Quantity
}
