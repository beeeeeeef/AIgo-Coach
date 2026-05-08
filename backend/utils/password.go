package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 对密码进行 bcrypt 哈希
func HashPassword(password string)(string,error){
	bytes,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
// CheckPasswordHash 比较密码和哈希值是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}