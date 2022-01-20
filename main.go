package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"seckill/dbs"
	"seckill/routers"
)

func main() {
	//读取配置文件
	initConfig()
	//初始化Mysql连接
	initMysql()
	//初始化redis
	initRedis()

	//go utils.UserLogin(31)

	route := gin.Default()
	//初始化路由
	routers.Init(route)
	//设置端口 别忘了 端口设置前面加 ":"
	route.Run(viper.GetString("web.post"))

	//用户密码：12345566
	//good := models.Good{}
	//ids := good.GetGoodList(dbs.GetMysqlClinet())
	//fmt.Println(*ids)

	//utils.CreateUser()

}

//初始化配置文件
func initConfig() {
	viper.SetConfigName("settings.yaml") //文件的名称
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs/") // 文件路径

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	//fmt.Println("config app:", viper.Get("database1.driver"))
	//fmt.Println("config redis:", viper.Get("database1.source"))
}

//初始化mysql连接
func initMysql() {
	dbs.IntiMysqlClient()
}

func initRedis() {
	ctx := context.Background()
	dbs.IntiRedisClient(ctx)
}
