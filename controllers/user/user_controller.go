package user

import (
	"context"
	"crypto/md5"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/json-iterator/go"
	"seckill/controllers"
	"seckill/dbs"
	"seckill/services"
	"seckill/utils"
	"time"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		controllers.Response(c, 400, "数据错误", nil)
		return
	}

	user := services.Login(username)
	var md5Utile utils.MD5Util
	md5Utile.Seed = md5.New()
	password = md5Utile.FormPassToDBPass(password, user.Salt)

	if user.Password != password {
		controllers.Response(c, 400, "密码错误:"+password, nil)
		return
	}

	u := uuid.New()
	utils.DoSetCookie(c, "userTicket", u.String(), 1800, false)
	//将结构体转json
	userJson, err := jsoniter.Marshal(user)
	if err != nil {
		print(err)
	}
	//获取rerdis连接
	redisClinet := dbs.GetRedisClinet()
	//将user 放入redis中，当作分布式session
	background := context.Background()
	//将user数据放入redis中，并且设置时间为半个小时
	redisClinet.Set(background, "user:"+u.String(), string(userJson), time.Second*1800)

	controllers.Response(c, 200, "成功", nil)
}
