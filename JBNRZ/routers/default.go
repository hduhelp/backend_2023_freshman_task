package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RootPage(ctx *gin.Context) {
	//ctx.HTML(http.StatusOK, "index.html", gin.H{
	//	"Title": "User",
	//})
	ctx.Redirect(http.StatusFound, "login")
}

func LoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", nil)
}

func RegisterPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", nil)
}
