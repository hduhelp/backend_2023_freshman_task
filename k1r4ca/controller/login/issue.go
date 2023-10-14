package login

import (
	ulits "HDUhelper_Todo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IssueJWT(ctx *gin.Context) {

	username := ctx.Param("username")
	jwt, err := ulits.IssueToken(username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to issue JWT. Please check and try again",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Done!",
		"Your JWT": jwt,
	})
}
