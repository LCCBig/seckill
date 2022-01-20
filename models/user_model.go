package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type User struct {
	UserId   int    `json:"userId"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	NikeName string `json:"nike_name"`
}

func (user *User) GetUserByUserName(userName string, db *sqlx.DB) *User {
	queryRow := db.QueryRow("select id, username, password, salt from t_user where username = ?", userName)

	err := queryRow.Scan(&user.UserId, &user.UserName, &user.Password, &user.Salt)
	if err != nil {
		fmt.Println(err)
	}

	return user
}

/**
插入Uesr
*/
func (user *User) InsertUser(db *sqlx.DB) {
	exec, err := db.Exec("INSERT INTO t_user(username, password, salt, nike_name) VALUE (?,?,?,?)", user.UserName, user.Password, user.Salt, user.NikeName)
	if err != nil {
		fmt.Println("第" + strconv.Itoa(user.UserId) + "个用户")
		fmt.Println(err)
		return
	}
	id, err := exec.LastInsertId()
	if err != nil {
		println(err)
		return
	}
	fmt.Println("test" + strconv.FormatInt(id, 10))

	//18000004130
	//18012343453
}
