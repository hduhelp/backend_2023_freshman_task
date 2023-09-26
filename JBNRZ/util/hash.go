package util

import (
	"crypto/md5"
	"fmt"
)

func Str2md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
