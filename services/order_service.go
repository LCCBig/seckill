package services

import (
	"context"
	"encoding/json"
	"seckill/dbs"
	"seckill/models"
	"strconv"
)

/**
通过用户ID和商品ID获取秒杀记录缓存
*/
func GetSecKillOrderCache(goodIduserId int, goodId int) *models.SecKillGood {
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
	json.Unmarshal([]byte(get.Val()), secKillOrder)

	return &secKillOrder
}
