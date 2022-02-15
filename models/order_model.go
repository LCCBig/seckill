package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

/**
订单
*/
type Order struct {
	//ID
	Id int `json:"id"`
	//用户ID
	UserId int `json:"userId"`
	//商品ID
	GoodId int `json:"goodId"`
	//收货地址ID
	DeliveryAddrId int `json:"deliveryAddrId"`
	//冗余过来的商品名称
	GoodName string `json:"goodName"`
	//商品数量
	GoodCount int `json:"goodCount"`
	//商品单价
	GoodsPrice decimal.Decimal `json:"goodsPrice"`
	//1pc,2android,3ios
	OrderChannel int `json:"orderChannel"`
	//订单状态，0新建未支付，1已支付，2已发货，3已收货，4已退款，5已完成
	OrderStatus int `json:"orderStatus"`
	//支付时间
	PayDate time.Time `json:"payDate"`
	//订单的创建时间
	CreateDate time.Time `json:"createDate"`
}

/**
添加新订单
*/
func (order *Order) InsertOrder(tx *sqlx.Tx) *sql.Result {
	//INSERT INTO t_order(user_id,good_id,delivery_addr_id,good_name,good_count,good_price,order_channal,order_status) VALUES(1,1,1,"iphone13",1,5333.99,1,1);
	result := tx.MustExec("INSERT INTO t_order(user_id,good_id,delivery_addr_id,good_name,good_count,good_price,order_channal,order_status) VALUES(?,?,?,?,?,?,?,?)",
		order.UserId, order.GoodId, order.DeliveryAddrId, order.GoodName, order.GoodCount, order.GoodsPrice, order.OrderChannel, order.OrderStatus)
	return &result
}
