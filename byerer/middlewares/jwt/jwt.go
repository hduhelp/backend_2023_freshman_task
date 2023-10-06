package jwt

import (
	"TODOlist/dao/mysql"
	"TODOlist/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var jwtSecret = []byte("hduHelp")

type Claims struct {
	UserID   int64  `json:"userID"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func MiddleWareJWT(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "please log in",
		})
		c.Abort()
		return
	} //redirect?

	//parse token
	claims, err := ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "login expired",
		})
		c.Abort()
		return
	}
	var user models.User
	result := mysql.DB.Where("userID = ? AND username = ?", claims.UserID, claims.Username).Find(&user)
	if result.Error != nil {
		//query failed
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "user does not exist",
		})
	}
	c.Set("userID", claims.UserID)
	c.Next()
}

func GenerateToken(userid int64, username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		userid,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "hduHelp",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
