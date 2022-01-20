package home

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ToLogin(c *gin.Context) {

	//utils.DoSetCookie(c,"","",0,true)
	c.HTML(http.StatusOK, "login.html", nil)
}
