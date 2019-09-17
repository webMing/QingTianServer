package tools

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gomodule/redigo/redis"
)

func TestVerifyHttp(t *testing.T) {
	t.Log("通过网络请求获取验证码")

	message := map[string]string{
		"phone_num": "136637",
	}
	bytesP, err := json.Marshal(message)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post("http://localhost:8080/api/v1/verifyCode", "application/json", bytes.NewBuffer(bytesP))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", result)

}

func TestYY(t *testing.T) {
	t.Log("通过网络请求获取验证码")
	resp, err := http.Post("http://localhost:8080/api/v1/YY", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	t.Logf("%v", *resp)

}

func TestRedisHelpSet(t *testing.T) {
	//设置新值,如果设置成功返回字符串OK,
	reply, err := redis.String(RedisHelperSet("my", "my", 0))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
}
func TestRedisHelpGet(t *testing.T) {
	reply, err := RedisHelperGet("my")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", reply.(string))
}
func TestRedisHelpExist(t *testing.T) {
	reply, err := RedisHelperExists("my")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
}

func TestRedisHelpTTL(t *testing.T) {
	reply, err := RedisHelperTTL("my")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply.(int64))
}
