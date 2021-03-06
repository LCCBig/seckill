package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

type SecKillGood struct {
	//秒杀商品ID
	Id int `json:"id"              db:"id"`
	//商品ID
	GoodId int `json:"goodId"          db:"good_id"`
	//秒杀价
	SecKillPrice decimal.Decimal `json:"secKillPrice"    db:"seckill_price"`
	//库存数量
	StockCount int `json:"stockCount"      db:"stock_count"`
	//秒杀开始时间
	StartDate time.Time `json:"startDate"       db:"start_date"`
	//秒杀结束时间
	EndDate time.Time `json:"endDate"         db:"end_date"`
}

/**
根据ID获取秒杀商品
*/
func (secKill *SecKillGood) GetSecKillGood(id int, db *sqlx.DB) {
	row := db.QueryRow("select id,good_id,seckill_price,stock_count,start_date,end_date from t_seckill_goods where id = ?", id)
	err := row.Scan(&secKill.Id, &secKill.GoodId, &secKill.SecKillPrice, &secKill.StockCount, &secKill.StartDate, &secKill.EndDate)

	if err != nil {
		panic(err)
	}
}

/**
获取未过期的秒杀商品列表
*/
func (secKill *SecKillGood) GetSecKillGoodList(db *sqlx.DB) *[]SecKillGood {
	var goodList []SecKillGood
	goodList = make([]SecKillGood, 0, 10)
	query, err := db.Query("SELECT id,good_id,seckill_price,stock_count,start_date,end_date FROM t_seckill_goods WHERE end_date > ?", time.Now().Format("2006-01-02 15:04:05"))
	//err := db.Select(&goodList, "SELECT id,good_id,seckill_price,stock_count,start_date,end_date FROM t_seckill_goods WHERE end_date  > ?", "2022-02-17 11:26:40")

	//defer query.Close()

	if err != nil {
		panic(err)
	}

	for query.Next() {
		var goodItem SecKillGood
		query.Scan(&goodItem.Id, &goodItem.GoodId, &goodItem.SecKillPrice, &goodItem.StockCount, &goodItem.StartDate, &goodItem.EndDate)
		goodList = append(goodList, goodItem)
	}
	return &goodList
}

/**
减库存操作
*/

func (secKill *SecKillGood) SaleOne(tx *sqlx.Tx) *sql.Result {
	result := tx.MustExec("UPDATE t_seckill_goods SET stock_count = stock_count - 1 WHERE good_id = ? AND stock_count - 1 >= 0", secKill.Id)
	return &result
}
