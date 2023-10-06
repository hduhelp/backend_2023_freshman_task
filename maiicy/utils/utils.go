package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"regexp"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

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

func GenerateJWT(user User) (string, error) {
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

func ParseJWT(tokenString string) (User, error) {
	var user User

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return user, err
	}

	if !token.Valid {
		return user, fmt.Errorf("无效的JWT令牌")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user, fmt.Errorf("无效的JWT令牌有效负载")
	}

	user.ID = int(claims["id"].(float64))
	user.Username = claims["username"].(string)

	return user, nil
}

func IsValidUsername(username string) bool {
	if len(username) < 3 {
		return false
	}

	usernamePattern := "^[a-zA-Z0-9]+$"
	match, _ := regexp.MatchString(usernamePattern, username)
	return match
}
