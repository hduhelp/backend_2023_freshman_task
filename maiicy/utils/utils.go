package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"login-system/models"
	"time"
)

var (
	SecretKey = []byte("") // 用于签名和验证 JWT 的密钥
)

func CalculateMD5(input string) (string, error) {
	data := []byte(input)
	hasher := md5.New()
	_, err := hasher.Write(data)
	if err != nil {
		return "", err
	}
	md5Hash := hasher.Sum(nil)
	md5Str := hex.EncodeToString(md5Hash)
	return md5Str, nil
}

func GenerateJWT(user models.User) (string, error) {
	// 定义 JWT 的有效负载
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 设置过期时间为一天
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥进行签名
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (models.User, error) {
	var user models.User

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return user, fmt.Errorf("JWT 令牌解析失败: %w", err)
	}

	if !token.Valid {
		return user, fmt.Errorf("无效的JWT令牌")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user, fmt.Errorf("无效的JWT令牌有效负载")
	}

	user.ID = uint(claims["id"].(float64))
	user.Username = claims["username"].(string)

	return user, nil
}
