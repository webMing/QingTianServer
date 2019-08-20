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
	"log"
    _ "github.com/go-sql-driver/mysql" //a blank import should be only in a main or test package, or have a comment justifying it
	"github.com/gin-gonic/gin"
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

	db, err := sql.Open("mysql", "root@tcp(localhost)/Ste?timeout=90s&charset=utf8&collation=utf8mb4_unicode_ci")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO user(passwd,phone_num,user_client,user_client_type) VALUES(?,?,?,?)")
	if err != nil {
		log.Fatalln(err)
	}

	passwd := c.PostForm("passwd")
	phoneNum := c.PostForm("phone_num")
	userClient := c.PostForm("user_client")
	clientVersion := c.PostForm("client_version")

	res, err := stmt.Exec(passwd, phoneNum, userClient, clientVersion)
	if err != nil {
		log.Fatalln(err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatalln(err)
	}

	var u map[string]interface{}
	if lastID == 0 {
		u["code"] = 1
		u["msg"] = "插入数据出现错误~"
	} else {
		u["code"] = 0
		u["msg"] = "OK"
	}
	u["user_id"] = lastID
	u["phone_num"] = phoneNum
	u["passwd"] = passwd
	u["user_client"] = userClient
	u["client_version"] = clientVersion

	return u, nil
}
