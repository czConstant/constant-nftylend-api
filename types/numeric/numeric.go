package numeric

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

type Decimal struct {
	decimal.Decimal
}

func (n *Decimal) ToDecimal() decimal.Decimal {
	return n.Decimal
}

func (n *Decimal) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		n = nil
		return nil
	}
	s = strings.Trim(s, `"`)
	d, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	*n = Decimal{
		d,
	}
	return nil
}

func (n *Decimal) MarshalJSON() ([]byte, error) {
	if n == nil {
		return []byte("null"), nil
	}
	s := n.Decimal.String()
	return []byte(fmt.Sprintf(`"%s"`, s)), nil
}

func (n *Decimal) Scan(src interface{}) error {
	if src == nil {
		n = nil
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return errors.New("invalid data type")
	}
	s := string(b)
	if s == "" {
		n = nil
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	*n = Decimal{
		d,
	}
	return nil
}

func (n Decimal) Value() (driver.Value, error) {
	return n.Decimal.String(), nil
}

// BigFloat

type BigFloat struct {
	big.Float
}

func (n *BigFloat) BigFloat() *big.Float {
	return &n.Float
}

func (n *BigFloat) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		n = nil
		return nil
	}
	s = strings.Trim(s, `"`)
	m, ok := big.NewFloat(0).SetString(s)
	if !ok {
		return errors.New("invalid data type")
	}
	*n = BigFloat{
		*m,
	}
	return nil
}

func (n *BigFloat) MarshalJSON() ([]byte, error) {
	if n == nil {
		return []byte("null"), nil
	}
	return []byte(n.BigFloat().Text('f', 10)), nil
}

func (n *BigFloat) Scan(src interface{}) error {
	if src == nil {
		n = nil
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return errors.New("invalid data type")
	}
	s := string(b)
	if s == "" {
		n = nil
		return nil
	}
	m, ok := big.NewFloat(0).SetString(s)
	if !ok {
		return errors.New("invalid data type")
	}
	*n = BigFloat{
		*m,
	}
	return nil
}

func (n BigFloat) Value() (driver.Value, error) {
	return n.BigFloat().Text('f', 10), nil
}

// BigFloat

type BigInt struct {
	big.Int
}

func (n *BigInt) BigInt() *big.Int {
	return &n.Int
}

func (n *BigInt) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		n = nil
		return nil
	}
	s = strings.Trim(s, `"`)
	m, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return errors.New("invalid data type")
	}
	*n = BigInt{
		*m,
	}
	return nil
}

func (n *BigInt) MarshalJSON() ([]byte, error) {
	if n == nil {
		return []byte("null"), nil
	}
	s := n.String()
	return []byte(fmt.Sprintf(`"%s"`, s)), nil
}

func (n *BigInt) Scan(src interface{}) error {
	if src == nil {
		n = nil
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return errors.New("invalid data type")
	}
	s := string(b)
	if s == "" {
		n = nil
		return nil
	}
	m, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return errors.New("invalid data type")
	}
	*n = BigInt{
		*m,
	}
	return nil
}

func (n BigInt) Value() (driver.Value, error) {
	return n.String(), nil
}
