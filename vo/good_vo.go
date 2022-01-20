package vo

import (
	"github.com/shopspring/decimal"
	"time"
)

/**
秒杀界面商品展示vo
*/
type GoodVo struct {
	//商品ID
	Id int `json:"id"`
	//秒杀商品ID
	SecKillId int `json:"secKillId"`
	//商品名称
	GoodName string `json:"goodName"`
	//商品标题
	GoodTitle string `json:"goodTitle"`
	//商品图片
	GoodImg string `json:"goodImg"`
	//商品详情
	GoodDetail string `json:"goodDetail"`
	//商品价格
	GoodPrice decimal.Decimal `json:"goodPrice"`
	//秒杀价
	SecKillPrice decimal.Decimal `json:"secKillPrice"`
	//库存数量
	StockCount int `json:"stockCount"`
	//秒杀开始时间
	StartDate time.Time `json:"startDate"`
	//秒杀结束时间
	EndDate time.Time `json:"endDate"`
}
