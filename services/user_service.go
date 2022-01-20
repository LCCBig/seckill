package services

import (
	"seckill/dbs"
	"seckill/models"
)

func Login(userName string) *models.User {
	//获取数据库连接
	clinet := dbs.GetMysqlClinet()
	var user models.User
	user.GetUserByUserName(userName, clinet)
	return &user
}
