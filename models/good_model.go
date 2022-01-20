package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type Good struct {
	//商品ID
	Id int `json:"id"          db:"id"`
	//商品名称
	GoodName string `json:"goodName"    db:"good_name"`
	//商品标题
	GoodTitle string `json:"goodTitle"   db:"good_title"`
	//商品图片
	GoodImg string `json:"goodImg"     db:"good_img"`
	//商品详情
	GoodDetail string `json:"goodDetail"  db:"good_detail"`
	//商品价格
	GoodPrice decimal.Decimal `json:"goodPrice"   db:"good_price"`
	//商品库存，-1表示没有限制
	GoodStock int `json:"goodStock"   db:"good_stock"`
}

func (good *Good) GetGoodById(id int, db *sqlx.DB) {
	row := db.QueryRow("select id,good_name,good_title,good_img,good_detail,good_price,good_stock from t_good where id = ?", id)
	err := row.Scan(&good.Id, &good.GoodName, &good.GoodTitle, &good.GoodImg, &good.GoodDetail, &good.GoodPrice, &good.GoodStock)

	if err != nil {
		print(err)
	}
}

/**
通过多个ID获取商品
*/
func (good *Good) GetGoodByIds(ids []int, db *sqlx.DB) *[]Good {
	var goodList []Good
	goodList = make([]Good, 0, 10)
	//拼接sql语句
	query, args, err := sqlx.In("select id,good_name,good_title,good_img,good_detail,good_price,good_stock from t_good where id in (?)", ids)
	db.Rebind(query)
	if err != nil {
		print(err)
	}

	err = db.Select(&goodList, query, args...)
	if err != nil {
		print(err)
	}
	//for query.Next() {
	//	var goodItem Good
	//	query.Scan(&goodItem.Id,&goodItem.GoodName,&goodItem.GoodTitle,&goodItem.GoodImg,&goodItem.GoodDetail,&goodItem.GoodPrice,&goodItem.GoodStock)
	//	goodList = append(goodList,goodItem)
	//}
	return &goodList
}

func (good *Good) GetGoodList(db *sqlx.DB) *[]Good {
	var goodList []Good
	goodList = make([]Good, 0, 10)
	query, err := db.Query("select id,good_name,good_title,good_img,good_detail,good_price,good_stock from t_good")

	defer query.Close()

	if err != nil {
		print(err)
	}

	for query.Next() {
		var goodItem Good
		query.Scan(&goodItem.Id, &goodItem.GoodName, &goodItem.GoodTitle, &goodItem.GoodImg, &goodItem.GoodDetail, &goodItem.GoodPrice, &goodItem.GoodStock)
		goodList = append(goodList, goodItem)
	}
	return &goodList
}
