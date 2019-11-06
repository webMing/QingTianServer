package constv

import "errors"

const (
	//QTInvalidValidationMsg 无效内容
	QTInvalidValidationMsg = "无效内容"
	//QTInvalidPhoneNumberMsg 手机号不是11位
	QTInvalidPhoneNumberMsg = "手机号不是11位"
	//QTInvalidCheckNumberMsg 验证码位数不对
	QTInvalidCheckNumberMsg = "验证码位数不对"
	
	//QTFetchSucssCode 请求成功
	QTFetchSucssCode = "0"
	//QTFetchSucssMsg 请求成功
	QTFetchSucssMsg = "请求成功"

	//QTFetchFailtCode 请求失败
	QTFetchFailtCode = "1"
	//QTFetchFailtMsg 请求失败
	QTFetchFailtMsg = "请求失败"
)

// QTERR 自定义消息
func QTERR(msg string) error {
	return errors.New(msg)
}
