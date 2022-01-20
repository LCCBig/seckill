package dbs

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"strings"
)

var (
	mysqlClient *sqlx.DB
)

func IntiMysqlClient() {
	driver := viper.GetString("database1.driver")
	dataSourceName := viper.GetString("database1.source")
	client, err := sqlx.Open(driver, dataSourceName)
	mysqlClient = client

	if err != nil {
		panic(err)
	}

	//最大连接数
	mysqlClient.SetMaxOpenConns(100)
	//最大空闲连接数
	mysqlClient.SetMaxIdleConns(5)

	err = mysqlClient.Ping()
	if err != nil {
		panic("无法联通！！！")
	}

}

func GetMysqlClinet() *sqlx.DB {
	return mysqlClient
}

/**
拼接sql语句的 "?"
*/
func Placeholders(n int) string {
	var b strings.Builder
	for i := 0; i < n-1; i++ {
		b.WriteString("?,")
	}
	if n > 0 {
		b.WriteString("?")
	}
	return b.String()
}
