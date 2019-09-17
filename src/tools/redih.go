package tools

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

// 地址可以更改
// const redisConnAddr = "192.168.1.100:6379"
// const redisConnAddr = "192.168.0.248:6379"
const redisConnAddr = "localhost:6379"

func redisHelpBase(f func(conn redis.Conn)) {
	c, err := redis.Dial("tcp", redisConnAddr)
	if err != nil {
		log.Fatalln(err)
	}
	f(c)
	defer c.Close()
}

// RedisHelperExists 判断键是否存在
func RedisHelperExists(key string) (exist bool, err error) {
	f := func(conn redis.Conn) {
		exist, err = redis.Bool(conn.Do("EXISTS", key))
	}
	redisHelpBase(f)
	return
}

// RedisHelperSet 设置键
func RedisHelperSet(key, value string, ex int) (reply interface{}, err error) {
	//重新设置key会覆盖之前的设置
	f := func(conn redis.Conn) {
		if ex != 0 {
			//ex 默认是秒
			reply, err = conn.Do("SET", key, value, "EX", ex)
		} else {
			reply, err = conn.Do("SET", key, value)
		}
	}
	redisHelpBase(f)
	return
}

// RedisHelperGet 获取键值
func RedisHelperGet(key string) (reply interface{}, err error) {
	f := func(conn redis.Conn) {
		reply, err = conn.Do("GET", key)
	}
	redisHelpBase(f)
	return
}

// RedisHelperTTL 判断键的过期时间
func RedisHelperTTL(key string) (reply interface{}, err error) {
	// -2 该键不存在
	// -1 改建永久存在
	// other 倒计时时间
	f := func(conn redis.Conn) {
		reply, err = conn.Do("TTL", key)
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