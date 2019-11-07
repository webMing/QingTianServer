package tools

import (
	"github.com/gomodule/redigo/redis"
)

// 地址可以更改
// const redisConnAddr = "192.168.1.100:6379"
//const redisConnAddr = "192.168.0.236:6379"
const redisConnAddr = "localhost:6379"

func redisHelpBase(f func(conn redis.Conn, er error)) {
	c, err := redis.Dial("tcp", redisConnAddr)
	if err != nil {
		// log.Fatalln(err);这里可以出错信息重定向到 file 中,同时重启
		//log.Println(err)
		f(c, err)
		return
	}
	f(c, err)
	defer c.Close() //if c is nil will crash
}

// RedisHelperExists 判断键是否存在
func RedisHelperExists(key string) (bool, error) {
	var res bool
	var err error
	f := func(conn redis.Conn, er error) {
		if err == nil {
			res, err = redis.Bool(conn.Do("EXISTS", key))
		} else {
			res, err = false, er
		}
	}
	redisHelpBase(f)
	return res, err
}

// RedisHelperSet 设置键
func RedisHelperSet(key, value string, ex int) (interface{}, error) {
	//重新设置key会覆盖之前的设置
	var res interface{}
	var err error
	f := func(conn redis.Conn, er error) {
		if er == nil {
			if ex != 0 {
				//ex 默认是秒
				res, err = conn.Do("SET", key, value, "EX", ex)
			} else {
				res, err = conn.Do("SET", key, value)
			}
		} else {
			res, err = nil, er
		}
	}
	redisHelpBase(f)
	return res, err
}

// RedisHelperGet 获取键值
func RedisHelperGet(key string) (interface{}, error) {
	var res interface{}
	var err error
	f := func(conn redis.Conn, er error) {
		if er == nil {
			res, err = conn.Do("GET", key)
		} else {
			res, err = nil, er
		}
	}
	redisHelpBase(f)
	return res, err
}

// RedisHelperTTL 判断键的过期时间
func RedisHelperTTL(key string) (interface{}, error) {
	// -2 该键不存在
	// -1 改建永久存在
	// other 倒计时时间
	var res interface{}
	var err error
	f := func(conn redis.Conn, er error) {
		if er == nil {
			res, err = conn.Do("TTL", key)
		} else {
			res, err = nil, er
		}
	}
	redisHelpBase(f)
	return res, err
}

// RedisHelperDel 删除某个键
func RedisHelperDel(key string) (interface{}, error) {
	var res interface{}
	var err error
	f := func(conn redis.Conn, er error) {
		if er == nil {
			res, err = conn.Do("DEL", key)
		} else {
			res, err = nil, er
		}
	}
	redisHelpBase(f)
	return res, err
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
