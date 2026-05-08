package utils

import (
	"crypto/rand"
	"math/big"
)

//生成6位数字验证码
func GenerateCode()string{
	code :=make([]byte,6)
	for i:=0;i<6;i++{
		n,_ := rand.Int(rand.Reader,big.NewInt(10))
		code[i] = byte(n.Int64()) + byte('0')
	}
	return string(code)
}
