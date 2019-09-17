//
// register.go
// register
//
//  register
//
// Created by Stephanie on 2019/06/16.
// Copyright © 2019 Stephanie. All rights reserved.
//

package login

import (
	"database/sql"
	"errors"
	"log"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" //a blank import should be only in a main or test package, or have a comment justifying it
)

/* 用户注册
phone_num       	手机号
passwd sha265  		密码
user_client    		客服端类型
client_version 		客户端版本
time_code      		时间戳
global_province_id  省区id
union_id
device_toke
*/

// Register 用户注册
func Register(c *gin.Context) (user map[string]interface{}, err error) {

	type replyJSON struct {
		PhoneNum   string `json:"phone_num"`
		Passwd     string `json:"passwd"`
		UserClient string `json:"user_client"`
	}
	var re replyJSON
	err = c.BindJSON(&re)
	if err != nil {
		panic(err)
	}

	if utf8.RuneCountInString(re.PhoneNum) != 11 {
		user = map[string]interface{}{
			"code": 1,
			"msg":  "手机号位数不是11位",
		}
		err = errors.New("手机号位数不对")
		return
	}

	// 不要使用本地设置
	//db, err := sql.Open("mysql", "root:centos@tcp(192.168.0.248)/QingTian?timeout=90s&charset=utf8&collation=utf8mb4_unicode_ci")
	db, err := sql.Open("mysql", "root@tcp(localhost)/QingTian?timeout=90s&charset=utf8&collation=utf8mb4_unicode_ci")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO user(user_name,passwd,phone_num,user_client) VALUES(?,?,?,?)")
	if err != nil {
		log.Fatalln(err)
	}

	res, err := stmt.Exec(re.PhoneNum, re.Passwd, re.PhoneNum, re.UserClient)
	if err != nil {
		log.Fatalln(err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatalln(err)
	}

	// var u map[string]interface{}
	// 这里的字典不能为空
	u := make(map[string]interface{})
	if lastID == 0 {
		u["code"] = 1
		u["msg"] = "插入数据出现错误~"
	} else {
		u["code"] = 0
		u["msg"] = "OK"
	}
	u["user_id"] = lastID
	u["phone_num"] = re.PhoneNum
	u["passwd"] = re.Passwd
	u["user_client"] = re.UserClient

	return u, nil
}
