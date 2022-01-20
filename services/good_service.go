package services

import (
	"seckill/dbs"
	"seckill/models"
	"seckill/vo"
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
	good.GetGoodById(secKillGood.GoodId, clinet)

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
