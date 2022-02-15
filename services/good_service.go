package services

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"seckill/dbs"
	"seckill/models"
	"seckill/vo"
	"strconv"
	"time"
)

func GetGoodList() *[]vo.GoodVo {
	//获取数据库连接
	clinet := dbs.GetMysqlClinet()

	//获取秒杀商品
	var secKillGood models.SecKillGood
	secKillGoodList := secKillGood.GetSecKillGoodList(clinet)
	//设置id对应对象的map的方便后续使用
	var goodMap map[int]models.SecKillGood
	goodMap = make(map[int]models.SecKillGood)
	//获取goodID
	var gIds []int
	gIds = make([]int, len(*secKillGoodList), len(*secKillGoodList))
	for _, secKill := range *secKillGoodList {
		goodMap[secKill.GoodId] = secKill
		gIds = append(gIds, secKill.GoodId)
	}
	//获取对应的商品
	var good models.Good
	goodList := good.GetGoodByIds(gIds, clinet)
	if len(*goodList) == 0 {
		return nil
	}

	//装填vo
	var goodVo []vo.GoodVo
	goodVo = make([]vo.GoodVo, len(*goodList), len(*goodList))
	for i := 0; i < len(*goodList); i++ {
		gId := (*goodList)[i].Id
		goodVo[i].Id = (*goodList)[i].Id
		goodVo[i].GoodName = (*goodList)[i].GoodName
		goodVo[i].GoodTitle = (*goodList)[i].GoodTitle
		goodVo[i].GoodImg = (*goodList)[i].GoodImg
		goodVo[i].GoodDetail = (*goodList)[i].GoodDetail
		goodVo[i].GoodPrice = (*goodList)[i].GoodPrice

		goodVo[i].SecKillId = goodMap[gId].Id
		goodVo[i].SecKillPrice = goodMap[gId].SecKillPrice
		goodVo[i].StockCount = goodMap[gId].StockCount
		goodVo[i].StartDate = goodMap[gId].StartDate
		goodVo[i].EndDate = goodMap[gId].EndDate
	}

	return &goodVo
}

/**
根据秒杀ID获取商品信息
*/
func GetGood(secKillId int) *vo.GoodVo {
	//获取数据库连接
	clinet := dbs.GetMysqlClinet()

	//获取秒杀商品
	var secKillGood models.SecKillGood
	secKillGood.GetSecKillGood(secKillId, clinet)

	//获取对应的商品
	var good models.Good
	good.Id = secKillGood.GoodId
	successful := good.GetGoodById(clinet)
	if !successful {
		return nil
	}

	//装填vo
	var goodVo vo.GoodVo
	goodVo.Id = good.Id
	goodVo.GoodName = good.GoodName
	goodVo.GoodTitle = good.GoodTitle
	goodVo.GoodImg = good.GoodImg
	goodVo.GoodDetail = good.GoodDetail
	goodVo.GoodPrice = good.GoodPrice

	goodVo.SecKillId = secKillGood.Id
	goodVo.SecKillPrice = secKillGood.SecKillPrice
	goodVo.StockCount = secKillGood.StockCount
	goodVo.StartDate = secKillGood.StartDate
	goodVo.EndDate = secKillGood.EndDate

	return &goodVo
}

/**
通过秒杀商品ID获取秒杀商品
*/
func GetSecKillGoodByid(secKillId int) *models.SecKillGood {
	//获取数据库连接
	clinet := dbs.GetMysqlClinet()
	//获取秒杀商品
	var secKillGood models.SecKillGood
	secKillGood.GetSecKillGood(secKillId, clinet)
	return &secKillGood
}

/**
商品秒杀
*/
func SecKill(user *models.User, secKillGood *models.SecKillGood) *models.Order {
	//异常处理
	defer func() {
		if p := recover(); p != nil {
			fmt.Println(p)
		}
	}()

	//获取数据库连接
	mysqlClinet := dbs.GetMysqlClinet()
	//查询商品
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
	//获取redis连接
	//获取秒杀持续时间
	validTime := secKillGood.EndDate.Second() - secKillGood.StartDate.Second()
	redisClinet := dbs.GetRedisClinet()
	redisClinet.Set(background, "order:"+strconv.Itoa(user.UserId)+":"+strconv.Itoa(secKillGood.Id), string(secKillOrderJson), time.Duration(validTime)*time.Second)
	return &order
}
