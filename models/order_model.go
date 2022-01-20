package models

import (
	"github.com/shopspring/decimal"
	"time"
)

/**
订单
*/
type Order struct {
	//ID
	Id int
	//用户ID
	UserId int
	//收货地址ID
	DeliveryAddrId int
	//冗余过来的商品名称
	GoodName string
	//商品数量
	GoodCount int
	//商品单价
	GoodsPrice decimal.Decimal
	//1pc,2android,3ios
	OrderChannel int
	//订单状态，0新建未支付，1已支付，2已发货，3已收货，4已退款，5已完成
	Status int
	//支付时间
	PayDate time.Time
	//订单的创建时间
	CreateDate time.Time
}
