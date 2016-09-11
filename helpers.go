package types

import (
	"strings"
	"os"
)

func GetGOPATH() string {
	gopaths := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
	if len(gopaths) == 0 {
		return ""
	}
	return gopaths[0]
}