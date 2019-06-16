package hosts

import (
	"log"
	"login"
	"net/http"
	"unicode/utf8"

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

	/* 获取验证码
	   - phoneNum 手机号
	*/
	group.POST("/verifyCode", func(c *gin.Context) {
		phoneNum := c.PostForm("phone_num")
		if utf8.RuneCountInString(phoneNum) != 11 {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "手机号位数不对!",
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
