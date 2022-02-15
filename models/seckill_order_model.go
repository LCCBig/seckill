package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type SeckillOrder struct {
	//秒杀订单ID
	Id int `json:"id"`
	//用户ID
	UserId int `json:"userId"`
	//订单ID
	OrderId int `json:"orderId"`
	//商品ID
	GoodsId int `json:"goodsId"`
}

func (secKillOrder *SeckillOrder) InsertSecKillOrder(tx *sqlx.Tx) *sql.Result {
	//INSERT INTO t_order(user_id,good_id,delivery_addr_id,good_name,good_count,good_price,order_channal,order_status) VALUES(1,1,1,"iphone13",1,5333.99,1,1);
	result := tx.MustExec("INSERT INTO t_seckill_order(user_id,order_id,good_id) VALUES(?,?,?);",
		secKillOrder.UserId, secKillOrder.OrderId, secKillOrder.GoodsId)
	return &result
}
