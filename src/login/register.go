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
	"unicode/utf8"

	"github.com/gomodule/redigo/redis"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" //a blank import should be only in a main or test package, or have a comment justifying it
	"stephanie.io/tools"
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

	// 手机号格式是否正确
	if utf8.RuneCountInString(re.PhoneNum) != 11 {
		u := map[string]interface{}{
			"code": 1,
			"msg":  "手机号位数不是11位",
		}
		// 不写err,上层如果捕捉到err会停止服务
		// err = errors.New("手机号位数不对")
		// return u,nil
		return u, nil
	}

	// 验证码是否有效
	code, err := redis.String(tools.RedisHelperGet(re.PhoneNum))
	if utf8.RuneCountInString(code) == 0 {
		u := map[string]interface{}{
			"code": 1,
			"msg":  "验证码不存在,请重新获取验证码",
		}
		// 不写err,上层如果捕捉到err会停止服务
		//err = errors.New("验证码不存在")
		return u, nil
	}

	// 不要使用本地设置
	//db, err := sql.Open("mysql", "root:centos@tcp(192.168.0.248)/QingTian?timeout=90s&charset=utf8&collation=utf8mb4_unicode_ci")
	db, err := sql.Open("mysql", "root@tcp(localhost)/QingTian?timeout=90s&charset=utf8&collation=utf8mb4_unicode_ci")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 是否已经注册过

	/*
		err := db.QueryRow("SELECT user_id from user where phone_num=?",re.PhoneNum).Scan(nil)

		stmt, err := db.Prepare("INSERT INTO user(user_name,passwd,phone_num,user_client) VALUES(?,?,?,?)")
		if err != nil {
			log.Fatalln(err)
		}

		res, err := stmt.Exec(re.PhoneNum, re.Passwd, re.PhoneNum, re.UserClient)
		if err != nil {
			log.Fatalln(err)
		}
	*/

	// 判断如果已经插入过,就不要设置
	res, err := db.Exec("INSERT INTO user(user_name,passwd,phone_num,user_client) SELECT ?,?,?,? from dual WHERE NOT EXISTS(SELECT phone_num from user WHERE phone_num = ?)", re.PhoneNum, re.Passwd, re.PhoneNum, re.UserClient, re.PhoneNum)
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
		u["msg"] = "该账号已经注册过,请直接登录~"
		return u, nil
	}

	u["code"] = 0
	u["msg"] = "OK"
	u["user_id"] = lastID
	u["phone_num"] = re.PhoneNum
	u["passwd"] = re.Passwd
	u["user_client"] = re.UserClient

	return u, nil
}
