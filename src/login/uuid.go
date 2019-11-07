package login

import (
	"stephanie.io/tools"
	"github.com/satori/go.uuid"
)

// GitHub : https://github.com/satori/go.uuid
// GoDoc  : https://godoc.org/github.com/satori/go.uuid

// UUID 可导出
func UUID() interface{} {
	type  res struct {
		UUID string `form:"uuid" json:"uuid"`
	}
	id := uuid.NewV4().String()
	if id == "" {
		return tools.OuterFailtStruct()
	}
	uid := new(res)
	uid.UUID = id
	base := tools.OuterSucssStruct()
	base.Data = uid
	return base

}

/***
func uuid() (user map[string]interface{}, err error) {
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
***/
