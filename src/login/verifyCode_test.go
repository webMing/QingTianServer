package login

import (
	"testing"
)

func TestVerifyCode(t *testing.T) {
	t.Log("本地获取验证码")
	SmsVerificationCode("17796648357")
}
