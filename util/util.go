package util

import (
	"crypto/md5"
	rand1 "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Strtomd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

//转为64位十进制整数
func StrToInt64(value string) (int64, error) {
	if len(value) == 0 {
		return int64(0), nil
	}
	return strconv.ParseInt(value, 10, 64)
}

func GenerateSalt() string {
	const randomLength = 16
	var salt []byte
	var asciiPad int64
	asciiPad = 32

	for i := 0; i < randomLength; i++ {
		salt = append(salt, byte(rand.Int63n(94)+asciiPad))
	}
	return string(salt)
}

func GenrateHash(salt string, password string) string {
	var hash string
	fullString := salt + password
	sha := sha256.New()
	sha.Write([]byte(fullString))
	hash = base64.URLEncoding.EncodeToString(sha.Sum(nil))

	return hash
}

//判断字符串是否为邮箱
func IsEmail(email string) bool {
	var e = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")
	return e.MatchString(email)
}

func IsMatch(v1 string, v2 string) bool {
	reg, err := regexp.Compile(v1)
	if err != nil {
		return false
	}
	if reg.MatchString(v2) {
		return true
	}
	return false
}

//判断路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//判断文件是否存在
func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// int64 数组 去重
func RemoveDuplicateInt64s(a []int64) (ret []int64) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if i > 0 && a[i-1] == a[i] {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

func GenerateSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand1.Reader, b); err != nil {
		return ""
	}
	return Strtomd5(base64.URLEncoding.EncodeToString(b))
}

func IntArrToStr(ints []int, joinStr string) string {
	str := make([]string, len(ints))
	for i, v := range ints {
		str[i] = strconv.Itoa(v)
	}
	return strings.Join(str, joinStr)
}
