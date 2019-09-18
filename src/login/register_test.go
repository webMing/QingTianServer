package login

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Log("通过网络测试注册接口")

	message := map[string]string{
		"phone_num": "136637",
		"passwd": "12345",
		"user_client":"iOS",
	}
	bytesP, err := json.Marshal(message)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post("http://localhost:8080/api/v1/register", "application/json", bytes.NewBuffer(bytesP))
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
