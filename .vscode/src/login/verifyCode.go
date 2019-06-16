//
// veric.go
// veric
//
// 1.注意在redis中可能出现的错误,
// 2.网络请求是同步,数据返回来之后gin才返回.
// 3.注意redis的多线程问题(暂未处理).
// Created by Stephanie on 2019/06/16.
// Copyright © 2019 Stephanie. All rights reserved.
//

 package login

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//SmsVerificationCode  获取验证码
func SmsVerificationCode(phoneNum string) (code string, err error) {
	return smsReqeust(phoneNum)
}

//fetchCode 获取短信验证码
func testCode(phoneNum string) {
	formdata := url.Values{
		"name":    {"masum"},
		"content": []string{"li", "wang"},
	}
	resp, err := http.PostForm("http://yun.tim.qq.com/v5/tlssmssvr/sendsms?sid=sid", formdata)
	if err != nil {

	}
	fmt.Println(resp)
}

func smsReqeust(phoneNum string) (code string, err error) {

	expireTime := 3

	exists, _ := redisHelperExists(phoneNum)
	if exists {
		return "1", fmt.Errorf("请在%d分钟后重新请求验证码", expireTime)
	}

	appid := "1400220829"
	appkey := "0a3931f269be610318cdd545f325ca71"
	random := strconv.Itoa(rand.Int())
	time := strconv.FormatInt(time.Now().Unix(), 10)
	mobile := phoneNum

	url, err := url.Parse("https://yun.tim.qq.com/v5/tlssmssvr/sendsms")
	if err != nil {
		return "1", fmt.Errorf("%s", "无法获取验证码地址!")
	}
	qury := url.Query()
	qury.Set("sdkappid", appid)
	qury.Set("random", random)
	url.RawQuery = qury.Encode()

	// 6 width verify code
	verifyCode := genValidateCode(6)
	//签名
	signedStr := "appkey=" + appkey + "&random=" + random + "&time=" + time + "&mobile=" + mobile
	sined := sha256.Sum256([]byte(signedStr))
	fmt.Println("yy网络请求中......")

	//6位验证码
	message := map[string]interface{}{
		"params": []string{
			"您的",       // {1}验证码:{2} {3}分钟内有效
			verifyCode, //6位验证码
			strconv.Itoa(expireTime),
		},
		"tel": map[string]interface{}{
			"mobile":     mobile,
			"nationcode": "86",
		},
		"time":   time,
		"dd":     signedStr,
		"sig":    fmt.Sprintf("%x", sined),
		"tpl_id": "351611",
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return "1", fmt.Errorf("%s", "json解析失败")
	}
	// return "1", fmt.Errorf("%s","json解析失败")
	fmt.Println("网络请求中......")
	resp, err := http.Post(url.String(), "application/json", bytes.NewBuffer(bytesRepresentation))
	fmt.Println("网络请求完成......")
	if err != nil {
		return "1", fmt.Errorf("%s", "获取验证码失败")
	}
	fmt.Println("处理返回结果......")
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "1", err
	}
	re := result["result"].(float64)
	if re == 0 {
		//save to memory.....
		fmt.Println("处理redis......")
		redisHelperSet(phoneNum, verifyCode, expireTime)
		fmt.Println("完成处理redis......")
	}
	// 注意errmsg;这个字段不要写错
	return fmt.Sprint(result["result"]), errors.New(result["errmsg"].(string))
}

// 产生随机验证码
func genValidateCode(width int) string {
	numeric := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	len := len(numeric)

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(len)])
	}
	return sb.String()
}

func redisHelperExists(key string) (exist bool, err error) {
	// 判断当前短信是否有效
	c, err := redis.Dial("tcp", "192.168.1.105:6379")
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()
	return redis.Bool(c.Do("EXISTS", key))
}

func redisHelperSet(key, value string, ex int) (reply interface{}, err error) {
	// redis 本地地址会随着环境不同而不同.
	c, err := redis.Dial("tcp", "192.168.1.105:6379")
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()
	if ex != 0 {
		return c.Do("SET", key, value, "EX", ex*60)
	}
	return c.Do("SET", key, value)
}

func redisHelperGet(key string) (reply interface{}, err error) {
	c, err := redis.Dial("tcp", "192.168.1.105:6379")
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()
	return c.Do("GET", key)
}

/****
user_client
client_version
timestamp
time_code
global_province_id
union_id
device_token
****/
