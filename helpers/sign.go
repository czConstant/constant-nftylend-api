package helpers

import (
	"fmt"
	"math/big"
)

func GetSignMsg(msg string) string {
	return fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg), msg)
}

func AppendHexStrings(values []string) string {
	var ret string
	for _, v := range values {
		ret = fmt.Sprintf("%s%s", ret, v)
	}
	return ret
}

func ParseHex2Hex(v string) string {
	if has0xPrefix(v) {
		v = v[2:]
	}
	return v
}

func ParseNumber2Hex(v string) string {
	n, _ := big.NewInt(0).SetString(v, 10)
	return ParseBigInt2Hex(n)
}

func ParseBigInt2Hex(v *big.Int) string {
	if v == nil {
		return "0000000000000000000000000000000000000000000000000000000000000000"
	}
	return fmt.Sprintf("%64s", fmt.Sprintf("0000000000000000000000000000000000000000000000000000000000000000%s", v.Text(16)))
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func ParseAddress2Hex(v string) string {
	if has0xPrefix(v) {
		v = v[2:]
	}
	return fmt.Sprintf("%40s", fmt.Sprintf("0000000000000000000000000000000000000000000000000000000000000000%s", v))
}
