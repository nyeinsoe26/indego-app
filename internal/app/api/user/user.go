package user

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Handler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	if profile == nil {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}

	ctx.HTML(http.StatusOK, "user.html", profile)
}
