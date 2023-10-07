package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"time"
)

// MyPayload 定义负载继承jwt的标准负载
type MyPayload struct {
	Username string
	jwt.RegisteredClaims
}

// 定义secret签名
var signatureKey []byte = []byte(viper.GetString("secret.signatureKey"))

// MakeUserToken 生成加密token
func MakeUserToken(username string) string {
	//传入用户信息生成负载实例
	payload := MyPayload{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))), // 过期时间3小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                       // 生效时间
		}}
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

// JWTHandler jwt拦截器
func JWTHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		//引入jwt实现登录后的会话记录,登录会话发生登录完成之后
		//header获取token
		session := sessions.Default(context)
		tokenValue := session.Get("username")
		if tokenValue == nil {
			context.String(302, "请求未携带token无法访问!")
			context.Abort()
		} else {
			token := tokenValue.(string)
			//解析token
			if token == "" {
				context.String(302, "？你的token呢？")
				context.Abort()
			}
			claims, err := ParserUserToken(token)
			if claims == nil || err != nil {
				context.String(401, "未携带有效token或已过期")
				context.Abort()
			} else {
				context.Set("user", claims.Username)
				context.Next()
			}
		}
	}
}
