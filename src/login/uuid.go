package login

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"unicode/utf8"
)

// UUID 自定义key
func UUID() (user map[string]interface{}, err error) {
	random := strconv.Itoa(rand.Int())
	// return strconv.Itoa(rand.Intn(10))
	uuid := fmt.Sprintf("%x", md5.Sum([]byte(random)))
	if utf8.RuneCountInString(uuid) == 0 {
		return nil, errors.New("无法获取UUID")
	}
	u := map[string]interface{}{
		"code": 0,
		"msg":  "获取成功",
		"uuid": uuid,
	}
	return u, nil
}
