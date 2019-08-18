package login

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestVerifyCode(t *testing.T) {
	t.Log("本地获取验证码")
	SmsVerificationCode("17796648357")
}

func TestVerifyHttp(t *testing.T) {
	t.Log("通过网络请求获取验证码")

	message := map[string]string{
		"phone_num": "13555",
	}
	bytesP, err := json.Marshal(message)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post("http://localhost:8080/verifyCode ", "application/json", bytes.NewBuffer(bytesP))
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

func TestYY(t *testing.T){
	t.Log("通过网络请求获取验证码")
	resp, err := http.Post("http://localhost:8080/api/v1/YY", "application/json",nil)
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
