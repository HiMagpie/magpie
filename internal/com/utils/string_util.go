package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

/**
 * Md5加密字符串
 */
func Md5Str(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}


/**
 * 截取字符串
 */
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

/**
 * offset 从0开始
 */
func Md5AndSub(s string, offset, length int) string {
	str := Md5Str(s)
	if len(str) <= offset {
		return ""
	}

	end := offset + length
	if end > len(str) - 1 {
		end = len(str) - 1
	}
	return str[offset: end]
}


/**
 * 字符串转int64
 */
func AtoInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}