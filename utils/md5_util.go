package utils

import (
	"encoding/hex"
	"fmt"
	"hash"
)

type MD5Util struct {
	Seed hash.Hash
}

var salt string = "1a2b3c4d"

/**
统一salt加密
*/
func (md5Util *MD5Util) InputPassToFromPass(inputPass string) string {
	defer md5Util.Seed.Reset()
	str := string(salt[0]) + string(salt[2]) + inputPass + string(salt[5]) + string(salt[4])
	md5Util.Seed.Write([]byte(str))
	password := hex.EncodeToString(md5Util.Seed.Sum(nil))
	return password
}

//func (md5Util *MD5Util)Test(inputPass string) string {
//	str := string(salt[0]) + string(salt[2]) + inputPass + string(salt[5]) + string(salt[4])
//	fmt.Println([]byte(str))
//	md5Util.Seed.Write([]byte(str))
//	cipherStr := md5Util.Seed.Sum(nil)
//	fmt.Println(cipherStr)
//	md5Util.Seed.Reset()
//	return hex.EncodeToString(md5Util.Seed.Sum(nil))
//}
/**
用户salt加密
*/
func (md5Util *MD5Util) FormPassToDBPass(formPass string, saltByUser string) string {
	defer md5Util.Seed.Reset()
	str := string(saltByUser[0]) + string(saltByUser[2]) + formPass + string(saltByUser[5]) + string(saltByUser[4])
	md5Util.Seed.Write([]byte(str))
	password := fmt.Sprintf("%x", md5Util.Seed.Sum(nil))
	return password
}

//双加密
func (md5Util *MD5Util) InputPassToDBPass(inputPass string, saltByUser string) string {
	formPass := md5Util.InputPassToFromPass(inputPass)
	dbPass := md5Util.FormPassToDBPass(formPass, saltByUser)
	return dbPass
}
