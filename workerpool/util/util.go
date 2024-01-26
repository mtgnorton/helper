package util

import (
	"math/rand"
	"time"
)

// GenerateRandomName
func GenerateRandomName(length int) string {
	// 定义字符集
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 种子值，保证每次运行都能得到不同的随机数序列
	rand.Seed(time.Now().UnixNano())

	// 构建随机字符串
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
