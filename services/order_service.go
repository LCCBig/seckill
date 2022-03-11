package services

import (
	"context"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"seckill/dbs"
	"seckill/models"
	"strconv"
	"time"
)

/**
通过用户ID和商品ID获取秒杀记录缓存
*/
func GetSecKillOrderCache(userId int, goodId int) *models.SecKillGood {
	//获取redis连接
	redisClinet := dbs.GetRedisClinet()

	//获取秒杀订单缓存
	background := context.Background()
	get := redisClinet.Get(background, "order:"+strconv.Itoa(userId)+":"+strconv.Itoa(goodId))
	if get.Val() == "" {
		return nil
	}
	//将json反序列化
	var secKillOrder models.SecKillGood
	json.Unmarshal([]byte(get.Val()), &secKillOrder)

	return &secKillOrder
}

func MakeSecKillOrder(user *models.User, secKillGoodId int) *models.Order {
	//异常处理
	defer func() {
		if p := recover(); p != nil {
			fmt.Println(p)
		}
	}()

	//获取redis连接
	redisClinet := dbs.GetRedisClinet()

	//获取数据库连接
	mysqlClinet := dbs.GetMysqlClinet()
	//查询商品
	var secKillGood models.SecKillGood
	secKillGood.GetSecKillGood(secKillGoodId, mysqlClinet)

	var good models.Good
	good.Id = secKillGood.GoodId
	successful := good.GetGoodById(mysqlClinet)
	if !successful {
		return nil
	}
	//设置事务超时时间
	background := context.Background()
	timeout, cancelFunc := context.WithTimeout(background, time.Second*2)
	defer cancelFunc()
	//开启事务
	beginTx := mysqlClinet.MustBeginTx(timeout, nil)
	//秒杀商品库存减一
	secKillGood.SaleOne(beginTx)
	//生成订单
	var order models.Order
	order.UserId = user.UserId
	order.GoodId = secKillGood.GoodId
	order.DeliveryAddrId = 1
	order.GoodName = good.GoodName
	order.GoodCount = 1
	order.GoodsPrice = secKillGood.SecKillPrice
	order.OrderChannel = 1
	order.OrderStatus = 0
	insertOrder := order.InsertOrder(beginTx)
	insertId, err := (*insertOrder).LastInsertId()
	if err != nil {
		return nil
	}
	order.Id = int(insertId)

	//生成秒杀订单
	var secKillOrder models.SeckillOrder
	secKillOrder.UserId = user.UserId
	secKillOrder.GoodsId = good.Id
	secKillOrder.OrderId = order.Id
	secKillOrder.InsertSecKillOrder(beginTx)

	//序列化对象
	secKillOrderJson, err := jsoniter.Marshal(secKillOrder)
	if err != nil {
		beginTx.Rollback()
		return nil
	}
	//提交事务
	beginTx.Commit()
	//获取秒杀持续时间
	validTime := secKillGood.EndDate.Second() - secKillGood.StartDate.Second()
	redisClinet.Set(background, "order:"+strconv.Itoa(user.UserId)+"_"+strconv.Itoa(secKillGood.Id), string(secKillOrderJson), time.Duration(validTime)*time.Second)
	return &order
}
