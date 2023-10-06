package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWTsecret = []byte("ABAB")

type Claims struct {
	Id       uint   `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// 生成token
func GenerateToken(id uint, username string, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour) //设置token有效期为24小时
	claims := Claims{
		Id:       id,
		UserName: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "todo_list",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //使用指定的签名方法创建签名对象
	token, err := tokenClaims.SignedString(JWTsecret)                //使用指定的secret签名并获得完整的编码后的字符串token
	return token, err
}

// 验证token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) { //解析token
		return JWTsecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid { //验证token
			return claims, nil
		}
	}
	return nil, err
}
