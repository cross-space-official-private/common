package math

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/shopspring/decimal"
	"math"
)

var (
	Infinitive = decimal.NewFromInt(math.MaxInt64)
)

type BigDecimal struct {
	decimal.Decimal
	Scale uint
}

func (de BigDecimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(de.StringFixed(int32(de.Scale)))
}

func (de BigDecimal) Value() (driver.Value, error) {
	return de.String(), nil
}
