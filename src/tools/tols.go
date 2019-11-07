package tools

import (
	"stephanie.io/constv"
	"errors"
	"path"
	"runtime"
)

// GetCurrentFileDir 获取当前文件路径
func GetCurrentFileDir() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		return path.Dir(filename), nil
	}
	return "", errors.New("无法获取调用文件路径")
}

// BaseRelt 返回数据外层
type BaseRelt struct {
	Code string `form:"code" json:"code"`
	Msg  string `form:"msg" json:"msg"`
	Data interface{} `form:"data,omitempty" json:"data,omitempty"`
}

// OuterSucssStruct 返回成功数据的外层
func OuterSucssStruct()(res *BaseRelt) {
	res = new(BaseRelt)
	res.Code = constv.QTFetchSucssCode
	res.Msg = constv.QTFetchSucssMsg
  	return
}

// OuterFailtStruct 返回失败据的外层
func OuterFailtStruct()(res *BaseRelt) {
	res = new(BaseRelt)
	res.Code = constv.QTFetchFailtCode
	res.Msg = constv.QTFetchFailtMsg
  	return
}