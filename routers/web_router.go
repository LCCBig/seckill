package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seckill/common/middleware"
	"seckill/controllers/good"
	"seckill/controllers/home"
	"seckill/controllers/user"
)

func Init(router *gin.Engine) {
	//加载界面文件
	router.LoadHTMLGlob("views/templates/*")
	//配置静态文件 第一个参数是api，第二个是文件夹路径
	router.StaticFS("/static", http.Dir("views/static"))
	gin.Recovery()
	userGroup := router.Group("/user")
	{
		userGroup.POST("/login", user.Login)
	}

	homeGroup := router.Group("/home")
	{
		homeGroup.GET("/tologin", home.ToLogin)
	}

	goodGroup := router.Group("/goods").Use(middleware.GetUserByCookie())
	{
		goodGroup.GET("/toList", good.ToList)
		goodGroup.GET("/toDetail", good.ToDetail)
		goodGroup.POST("/detailInfo", good.GoodDetailInfo)
		goodGroup.POST("/seckill", good.DoSecKill)
	}
}
