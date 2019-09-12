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

func smsReqeust(phoneNum string) (code string, err error) {

	expireTime := 3  * 60 //过期时间 3 分钟

	num,err := redis.Int64(redisHelperTTL(phoneNum))
	if err != nil {
		return "1",err
	}
	// -1 没有设置过期时间.-2 该键目前不存在, other ttl time
	if num == -1 && num != -2{
		return "1", fmt.Errorf("请在%d秒后重新请求验证码", num)
	}

	appid := "1400220829"
	appkey := "0a3931f269be610318cdd545f325ca71"
	random := strconv.Itoa(rand.Int())
	time := strconv.FormatInt(time.Now().Unix(), 10)
	mobile := phoneNum

	url, err := url.Parse("https://yun.tim.qq.com/v5/tlssmssvr/sendsms")
	if err != nil {
		return "1", fmt.Errorf("%s", "验证码地址错误")
	}
	qury := url.Query()
	qury.Set("sdkappid", appid)
	qury.Set("random", random)
	url.RawQuery = qury.Encode()

	// 获取6位长度的验证码 
	verifyCode := genValidateCode(6)
	//签名
	signedStr := "appkey=" + appkey + "&random=" + random + "&time=" + time + "&mobile=" + mobile
	sined := sha256.Sum256([]byte(signedStr))

	//6位验证码
	message := map[string]interface{}{
		"params": []string{
			"您的",       // {1}验证码:{2} {3}分钟内有效
			verifyCode, //  6位验证码
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
		return "1", fmt.Errorf("%s", "获取验证码请求参数转换json出错")
	}
	fmt.Println("网络请求中......")
	resp, err := http.Post(url.String(), "application/json", bytes.NewBuffer(bytesRepresentation))
	fmt.Println("网络请求完成......")
	if err != nil {
		return "1", fmt.Errorf("%s", "获取验证码网络请求失败")
	}
	defer resp.Body.Close()
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
		redisHelperSet(phoneNum,verifyCode, expireTime)
		fmt.Println("完成redis......")
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

// 地址可以更改
// const redisConnAddr = "192.168.1.100:6379"
const redisConnAddr = "192.168.0.248:6379"
func redisHelpBase(f func(conn redis.Conn)){
	c, err := redis.Dial("tcp",redisConnAddr)
	if err != nil {
		log.Fatalln(err)
	}
	f(c)
	defer c.Close()
}
// redisHelperExists 判断键是否存在
func redisHelperExists(key string) (exist bool, err error) {
	f := func(conn redis.Conn){
		exist,err = redis.Bool(conn.Do("EXISTS", key))
	}
	redisHelpBase(f)
	return
}
// redisHelperSet 设置键
func redisHelperSet(key, value string, ex int) (reply interface{}, err error) {
	//重新设置key会覆盖之前的设置
	f := func(conn redis.Conn){
		if ex != 0{
			//ex 默认是秒
			reply,err = conn.Do("SET",key,value,"EX",ex)
		}else{
			reply,err = conn.Do("SET",key,value)
		}
	}
	redisHelpBase(f)
	return
}
//  redisHelperSet 获取键值
func redisHelperGet(key string) (reply interface{}, err error) {
	f := func(conn redis.Conn){
		reply,err = conn.Do("GET", key)
	}
	redisHelpBase(f)
	return 
}
//  redisHelperTTL 判断键的过期时间
func redisHelperTTL(key string) (reply interface{}, err error){
	// -2 该键不存在
	// -1 改建永久存在
	// other 倒计时时间
	f := func(conn redis.Conn){
		reply,err = conn.Do("TTL",key)
	}
	redisHelpBase(f)
	return
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
