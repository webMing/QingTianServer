package hosts

import (
	"testing"
)

func TestConcurrentExc(t *testing.T) {
	t.Logf("测试接口是否是并发执行")
}

func TestRegister(t *testing.T) {
	t.Log("测试注册")
}
func TestLogin(t *testing.T) {
	t.Log("测试登录")
}

func TestSetUpServer(t *testing.T) {
	Server()
}
