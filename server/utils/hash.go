package utils

import (
	"crypto/md5"
	"encoding/hex"
	
	"golang.org/x/crypto/bcrypt"
)

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// MD5V 计算给定字节切片的MD5哈希值，并返回其十六进制字符串表示。
//
// 参数:
//   - str: 要计算MD5哈希值的字节切片。
//   - b: 可选的字节切片，用于在计算哈希值时附加到输入数据之后。如果未提供，则默认为空。
//
// 返回值:
//   - string: 计算得到的MD5哈希值的十六进制字符串表示。
func MD5V(str []byte, b ...byte) string {
	// 创建一个新的MD5哈希计算器
	h := md5.New()
	
	// 将输入字节切片写入哈希计算器
	h.Write(str)
	
	// 计算最终的哈希值，并将其转换为十六进制字符串
	return hex.EncodeToString(h.Sum(b))
}
