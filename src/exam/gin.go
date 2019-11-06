package exam 

import (
	"github.com/gin-gonic/gin/binding"
	"fmt"
    "stephanie.io/hosts"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

// GinA 测试内容
func GinA() {

	//chanage v8 to v9
	binding.Validator = new(hosts.defaultValidator)

	router := gin.Default()
	router.GET("/hello", func(context *gin.Context) {
		//返回的内容
		type User struct {
			Name string `form:"username" json:"username" binding:"required,min=6"`
		}
		user := new(User)
		err := context.ShouldBind(user)
		if err != nil {
			// 传入内容无效设置
			if _, ok := err.(*validator.InvalidValidationError); ok {
				fmt.Println(err)
				return
			}
			tp := err.(validator.ValidationErrors)[0]
			switch tp.Field() {
			case "Name":
				fmt.Println("姓名不能为空")
			}
		}

		//结果内容
		type Resl struct {
			Code string `form:"code,omitempty" json:"code,omitempty"`
			Msg  string `form:"msg" json:"msg"`
			Name string `form:"name" json:"name"`
		}
		res := new(Resl)
		res.Name = user.Name
		context.JSON(200, res)
		// context.XML(200, res)

	})
	// 指定地址和端口号
	router.Run("localhost:9090")
}
