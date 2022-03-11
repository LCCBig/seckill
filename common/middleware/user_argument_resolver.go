package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"seckill/dbs"
	"seckill/models"
)

func AuthMd5() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.PostForm("password")
		context.Next()
	}
}

func GetUserByCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userTicket, err := ctx.Cookie("userTicket")
		if err != nil {
			panic(err)
		}
		//获取redis连接
		redisClinet := dbs.GetRedisClinet()
		//根据ticket获取user
		background := context.Background()
		userJson := redisClinet.Get(background, "user:"+userTicket)
		if userJson.Val() == "" {
			return
		}
		var user models.User
		jsoniter.Unmarshal([]byte(userJson.Val()), &user)
		ctx.Set("user", user)
	}
}
