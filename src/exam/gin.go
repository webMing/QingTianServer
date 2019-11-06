package exam

import (
	"fmt"

	"stephanie.io/constv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"stephanie.io/hosts"
)

// GinA 测试内容
func GinA() {

	//chanage v8 to v9
	binding.Validator = new(hosts.DefaultValidator)

	router := gin.Default()
	router.GET("/hello", func(context *gin.Context) {

		type Resl struct {
			Code string `form:"code" json:"code"`
			Msg  string `form:"msg" json:"msg"`
			Name string `form:"name,omitempty" json:"name,omitempty"`
		}
		res := new(Resl)
		res.Code = constv.QTFetchSucssCode
		res.Msg = constv.QTFetchSucssMsg

		type User struct {
			Name string `form:"name" json:"name" binding:"required,min=6"`
		}
		user := new(User)
		err := context.ShouldBind(user)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				//fmt.Println(err)
				res.Code = constv.QTFetchFailtCode
				res.Msg = constv.QTInvalidValidationMsg
			} else {
				tp := err.(validator.ValidationErrors)[0]
				switch tp.Field() {
				case "Name":
					res.Code = constv.QTFetchFailtCode
					res.Msg = "姓名长度不够,至少需要6个字符"
				}
			}

		}
		
		context.JSON(200, res)
		// context.XML(200, res)

	})
	// 指定地址和端口号
	router.Run("localhost:9090")
}
