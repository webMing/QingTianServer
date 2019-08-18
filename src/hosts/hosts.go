package hosts

import (
	"log"
	"net/http"
	"unicode/utf8"

	"stephanie.io/login"

	"github.com/gin-gonic/gin"
)

const (
	baseURL string = "api/v1"
)

// Server  provide service
func Server() {
	router := gin.Default()
	group := router.Group(baseURL)

	/* 用户注册
	phone_num       	手机号
	passwd sha265  		密码
	user_client    		客服端类型
	client_version 		客户端版本
	time_code      		时间戳
	global_province_id  省区id
	union_id
	device_toke
	*/
	group.POST("/register", func(c *gin.Context) {
		user, err := login.Register(c)
		if err != nil {
			log.Fatalln(err)
		}
		c.JSON(http.StatusOK, user)
	})

	group.POST("/YY", func(c *gin.Context) {
		c.JSON(http.StatusOK, "iij")
	})

	/* 获取验证码
	   phoneNum 手机号
	*/
	group.POST("/verifyCode", func(c *gin.Context) {
		phoneNum := c.Query("phone_num")
		// phoneNum := c.PostForm("phone_num")
		// cc, _ := c.GetQueryMap("phone_num")
		// ff := c.Param("phone_num")
		// fmt.Println(cc, ff)
		//检测手机号位数
		if utf8.RuneCountInString(phoneNum) != 11 {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "手机号位数不是11位",
			})
			return
		}
		// 发送验证码网络请求
		code, err := login.SmsVerificationCode(phoneNum)
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  err.Error(),
		})
	})
	router.Run(":8080")
}
