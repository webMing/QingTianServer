package tools

import (
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
