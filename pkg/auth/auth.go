package auth

import "golang.org/x/crypto/bcrypt"

// @title	Encrypt
// @description	Encrypt encrypts the plain text with bcrypt
// @auth	高宏宇
// @param	source string 用户密码
// @return	string 加密用户密码	error 错误信息
func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// @title	Compare
// @description	Compare compares the encrypted text with the plain text if it's the same
// @auth	高宏宇
// @param	hashedPassword 用户加密密码	password string	用户密码
// @return	error 错误信息
func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
