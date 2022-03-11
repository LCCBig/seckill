package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"seckill/dbs"
	"seckill/routers"
	"seckill/services"
	"seckill/utils"
	"sync"
)

func main() {
	//读取配置文件
	initConfig()
	//初始化Mysql连接
	initMysql()
	//初始化redis
	initRedis()
	//启动边加载秒杀商品进入redis
	InitGoodToRedis()

	//初始化消息队列消费者（消费者初始话请放在最后，调用该方法后将阻塞在该方法）
	go InitMQConsumer()

	//go utils.UserLogin(31)

	//TestMap()

	route := gin.Default()
	//初始化路由
	routers.Init(route)
	//设置端口 别忘了 端口设置前面加 ":"
	route.Run(viper.GetString("web.post"))

	//用户密码：12345566
	//good := models.Good{Id:123}
	//id := good.GetGoodById(dbs.GetMysqlClinet())
	//fmt.Println(id)

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

func InitGoodToRedis() {
	services.GoodIntoRedis()
}

/**
消息队列消费者初始化
*/
func InitMQConsumer() {
	fmt.Println("消息队列启动")
	utils.ReceiveSecKillMessage()
}

func TestMap() {
	test := sync.Map{}

	test.Store(1, true)
	_, ok := test.Load(2)
	fmt.Println(ok)
	_, ok = test.Load(1)
	fmt.Println(ok)
}
