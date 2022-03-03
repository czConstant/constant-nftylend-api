package helpers

import "fmt"

func GetSignMsg(msg string) string {
	return fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg), msg)
}
