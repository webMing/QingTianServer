package constv

import "errors"

const (
	//QTInvalidValidationError 无效内容
	QTInvalidValidationError = "无效内容"
	//QTInvalidPhoneNumberError 手机号不是11位
	QTInvalidPhoneNumberError = "手机号不是11位"
	//QTInvalidCheckNumberError 验证码位数不对
	QTInvalidCheckNumberError = "验证码位数不对"
)

// QTERR 自定义消息
func QTERR(msg string) error {
	return errors.New(msg)
}
