package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RootPage(ctx *gin.Context) {
	ctx.String(http.StatusOK, "index.html")
}

func LoginPage(ctx *gin.Context) {
	ctx.String(http.StatusOK, "login.html")
}

func RegisterPage(ctx *gin.Context) {
	ctx.String(http.StatusOK, "register.html")
}
