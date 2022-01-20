package vo

import (
	"seckill/models"
)

type GoodDetailVo struct {
	//用户
	User models.User `json:"user"`
	//商品信息
	Good GoodVo `json:"good"`
	//秒杀状态
	SecKillStatus int `json:"secKillStatus"`
	//开始秒杀剩余时间
	RemainSeconds int `json:"remainSeconds"`
}
