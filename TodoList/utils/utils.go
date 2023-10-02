package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

// MyPayload 定义负载继承jwt的标准负载
type MyPayload struct {
	Username string
	jwt.StandardClaims
}

// 定义secret签名
var signatureKey []byte = []byte(viper.GetString("secret.signatureKey"))

// MakeUserToken 生成加密token
func MakeUserToken(username string) string {
	//传入用户信息生成负载实例
	payload := MyPayload{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 3600,
		},
	}

	//生成加密Signature
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString(signatureKey)
	if err != nil {
		panic(err)
	}
	return token
}

func ParserUserToken(token string) (*MyPayload, error) {
	//解密后jwt.Token对象，从该对象可以获取Header，Payload，Signature（claims）等信息
	unsafeToken, err1 := jwt.ParseWithClaims(token, &MyPayload{}, func(token *jwt.Token) (interface{}, error) {
		return signatureKey, nil
	})

	//将负载转化为结构体
	claims, ok := unsafeToken.Claims.(*MyPayload)

	if ok && unsafeToken.Valid {
		return claims, nil
	} else {
		return claims, err1
	}
}
