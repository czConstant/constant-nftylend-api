package models

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const (
	AddressCheckTxDuraion = -120 * time.Second
	// AxieClaimDuration     = 14 * 24 * time.Hour // for live
	AxieClaimDuration = 1 * 24 * time.Hour // for tesing

	MpOrderSellFeeRate = 0.01
)

type NftTermSetting struct {
	Term float64 `json:"term"`
	Rate float64 `json:"rate"`
}

func ConvertFloatToReserveAmount(amt float64) uint64 {
	if amt < 0 {
		panic("invalid amount")
	}
	return decimal.NewFromFloat(amt).Shift(2).Truncate(0).BigInt().Uint64()
}

func ConvertPriceAmount(amt float64) float64 {
	if amt < 0 {
		panic("invalid amount")
	}
	amt, _ = decimal.NewFromFloat(amt).Round(8).Float64()
	return amt
}

func ValidateFiatAmount(amt float64) error {
	if amt < 0 {
		return errors.New("invalid amount")
	}
	if amt != ConvertFiatAmount(amt) {
		return errors.New("fiat amount is invalid")
	}
	return nil
}

func ConvertFiatAmount(amt float64) float64 {
	if amt < 0 {
		panic("invalid amount")
	}
	amt, _ = decimal.NewFromFloat(amt).Round(2).Float64()
	return amt
}

func ValidateNftCurrencyAmount(amt float64, decimals uint) error {
	if amt < 0 {
		return errors.New("amount is small than 0")
	}
	newAmt, _ := decimal.NewFromFloat(amt).Round(int32(decimals)).Float64()
	if newAmt != amt {
		return errors.New("amount is invalid decimals")
	}
	return nil
}

func ConvertNumberFloat(amt float64, decimals uint) float64 {
	if amt < 0 {
		panic(errors.New("amount is small than 0"))
	}
	newAmt, _ := decimal.NewFromFloat(amt).Round(int32(decimals)).Float64()
	return newAmt
}

func ConvertWeiToCollateralFloatAmount(amt *big.Int, decimals uint) float64 {
	if amt.Cmp(big.NewInt(0)) < 0 {
		panic(errors.New("amount is small than 0"))
	}
	newAmt, err := decimal.NewFromString(amt.String())
	if err != nil {
		panic(err)
	}
	newAmt = newAmt.Shift(-int32(decimals)).Truncate(int32(8))
	val, _ := newAmt.Float64()
	return val
}

func ConvertWeiToBigFloat(amt *big.Int, decimals uint) *big.Float {
	if amt.Cmp(big.NewInt(0)) < 0 {
		panic(errors.New("amount is small than 0"))
	}
	newAmt, err := decimal.NewFromString(amt.String())
	if err != nil {
		panic(err)
	}
	newAmt = newAmt.Shift(-int32(decimals)).Truncate(int32(decimals))
	return newAmt.BigFloat()
}

func ConvertBigFloatToWei(amt *big.Float, decimals uint) *big.Int {
	if amt.Cmp(big.NewFloat(0)) < 0 {
		panic(errors.New("amount is small than 0"))
	}
	newAmt, err := decimal.NewFromString(amt.String())
	if err != nil {
		panic(err)
	}
	newAmt = newAmt.Shift(int32(decimals)).Truncate(0)
	return newAmt.BigInt()
}

func ConvertCryptoCurrencyAmount(amt float64) float64 {
	if amt < 0 {
		panic("invalid amount")
	}
	amt, _ = decimal.NewFromFloat(amt).Round(8).Float64()
	return amt
}

func ConvertReserveAmountToFloat(reserveAmt uint64) float64 {
	amt, _ := decimal.New(int64(reserveAmt), -2).Float64()
	return amt
}

func ConvertFloatToCollateralAmount(amt float64) uint64 {
	if amt < 0 {
		panic("amount invalid")
	}
	return decimal.NewFromFloat(amt).Shift(8).Truncate(0).BigInt().Uint64()
}

func ParseString2FloatAmountArr(s, sep string) []float64 {
	rets := []float64{}
	s = strings.TrimSpace(s)
	if s != "" {
		ss := strings.Split(s, sep)
		if len(ss) > 0 {
			for _, n := range ss {
				dm, err := decimal.NewFromString(n)
				if err != nil {
					panic(err)
				}
				val, _ := dm.Truncate(2).Float64()
				rets = append(rets, val)
			}
		}
	}
	return rets
}

func MulFloats(val1 float64, vals ...float64) float64 {
	val := decimal.NewFromFloat(val1)
	for _, v := range vals {
		val = val.Mul(decimal.NewFromFloat(v))
	}
	num, _ := val.Float64()
	return num
}

func DivFloats(val1 float64, vals ...float64) float64 {
	val := decimal.NewFromFloat(val1)
	for _, v := range vals {
		val = val.Div(decimal.NewFromFloat(v))
	}
	num, _ := val.Float64()
	return num
}

func AddFloats(val1 float64, vals ...float64) float64 {
	val := decimal.NewFromFloat(val1)
	for _, v := range vals {
		val = val.Add(decimal.NewFromFloat(v))
	}
	num, _ := val.Float64()
	return num
}

func SubFloats(val1 float64, vals ...float64) float64 {
	val := decimal.NewFromFloat(val1)
	for _, v := range vals {
		val = val.Sub(decimal.NewFromFloat(v))
	}
	num, _ := val.Float64()
	return num
}

func ConvertStringToFloat(s string) (float64, error) {
	num, err := decimal.NewFromString(s)
	if err != nil {
		return 0, err
	}
	amount, _ := num.Float64()
	return amount, nil
}

func ConvertString2BigInt(s string) (*big.Int, error) {
	n, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("%s is not number", s)
	}
	return n, nil
}

func ToEtherPriceFiatAmount(num big.Float, pr float64) float64 {
	rs, _ := big.NewFloat(0).Mul(&num, big.NewFloat(pr)).Float64()
	return ConvertFiatAmount(rs)
}

func ToEtherAmount(e *big.Int) big.Float {
	if e == nil {
		return big.Float{}
	}
	return *decimal.NewFromBigInt(e, -18).BigFloat()
}

func ToEtherWeiAmount(num big.Float) big.Int {
	dn, err := decimal.NewFromString(big.NewFloat(0).Mul(&num, big.NewFloat(1e18)).String())
	if err != nil {
		panic(err)
	}
	return *dn.BigInt()
}

func ToBigInt(s string) big.Int {
	if s == "" {
		return big.Int{}
	}
	n, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		panic(errors.New("numer is invalid"))
	}
	return *n
}

func MulBigFloats(val1 *big.Float, vals ...*big.Float) *big.Float {
	val := val1
	for _, v := range vals {
		val = big.NewFloat(0).Mul(val, v)
	}
	return val
}

func AddBigFloats(val1 *big.Float, vals ...*big.Float) *big.Float {
	val := val1
	for _, v := range vals {
		val = big.NewFloat(0).Add(val, v)
	}
	return val
}

func SubBigFloats(val1 *big.Float, vals ...*big.Float) *big.Float {
	val := val1
	for _, v := range vals {
		val = big.NewFloat(0).Sub(val, v)
	}
	return val
}

func QuoBigFloats(val1 *big.Float, vals ...*big.Float) *big.Float {
	val := val1
	for _, v := range vals {
		if v.Cmp(big.NewFloat(0)) == 0 {
			panic(errors.New("divide zero"))
		}
		val = big.NewFloat(0).Quo(val, v)
	}
	return val
}

func FormatFloatNumber(f string, amt float64) string {
	return FormatStringNumber(fmt.Sprintf(f, amt))
}

func FormatStringNumber(amt string) string {
	if strings.Contains(amt, ".") {
		amt = strings.TrimRight(amt, "0")
		amt = strings.TrimRight(amt, ".")
	}
	return amt
}
