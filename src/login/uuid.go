package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
	"stephanie.io/constv"
	"stephanie.io/tools"
)

// GitHub : https://github.com/satori/go.uuid
// GoDoc  : https://godoc.org/github.com/satori/go.uuid

// UUID 可导出
func UUID(c *gin.Context) interface{} {
	type request struct {
		Client string `form:"client" json:"client" binding:"required"`
	}
	re := new(request)
	err := c.ShouldBind(re)
	outerRes := tools.OuterSucssStruct()
	if err != nil {
		outerRes.Code = constv.QTFetchFailtCode
		if _, ok := err.(*validator.InvalidValidationError); ok {
			outerRes.Msg = constv.QTInvalidValidationMsg
		} else {
			for _, err := range err.(validator.ValidationErrors){
				switch err.Field() {
					case "Client":
						outerRes.Msg = "client 不能为空"
					default:
						outerRes.Msg = fmt.Sprintf("%s 该字段存在问题",err.Field())
				}
			}
		}
		return outerRes
	}

	id := uuid.NewV4().String()
	if id == "" {
		outerRes.Code = constv.QTFetchFailtCode
		outerRes.Msg = "无法获取UUID"
		return outerRes
	}

	type inerRes struct {
		UUID string `form:"uuid" json:"uuid"`
	}
	uid := new(inerRes)
	uid.UUID = id
	outerRes.Data = uid
	return outerRes

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
