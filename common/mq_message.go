package common

import "seckill/models"

type SecKillMassage struct {
	User   *models.User `json:"user"`
	GoodId int          `json:"goodId"`
}
