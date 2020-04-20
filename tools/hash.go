package tools

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

//小写的
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

func ValidatePasswd(plainpwd, salt, passwd string) bool {
	return Md5Encode(plainpwd+salt) == passwd
}
func MakePasswd(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

func Sha1(data string) string {
	_sha1 := sha1.New()
	_sha1.Write([]byte(data))
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func MakeSha1Password(password, salt string) string {

	for i := 0; i <= 10000; i++ {
		password = Sha1(password + salt)
	}
	//go println(s)
	return password
}

func VaildateSha1Passwd(plainpwd, salt, sha1pwd string) bool {
	return MakeSha1Password(plainpwd, salt) == sha1pwd
}
