package hosts

import (
	
	"encoding/json"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"stephanie.io/login"
)

const (
	baseURL string = "api/v1"
)

// Server  provide service
func Server() {

	//chanage v8 to v9
	binding.Validator = new(DefaultValidator)

	router := gin.Default()
	group := router.Group(baseURL)

	/* 获取UUID
	该UUID也是可以由客户端生成(争议UUID是否需要从这里获取)
	*/
	group.POST("/uuid", func(c *gin.Context) {
		c.JSON(http.StatusOK, login.UUID())
	})

	/* 获取图片验证码
	uuid string
	*/
	group.POST("/imgCheckCode", func(c *gin.Context) {
		uuid, err := login.Capthca(c)
		if err != nil {
			log.Fatalln(err)
		}
		c.JSON(http.StatusOK, uuid)
	})

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
	   phoneNum 手机号
	*/
	group.POST("/verifyCode", func(c *gin.Context) {
		/***********************************
		cc := c.PostForm("phone_num")
		phoneNum := c.Query("phone_num")
		***************************************/

		/***************************************
		// method 1
		type replyJSON struct {
			PhoneNum string `json:"phone_num"`
		}
		var re replyJSON
		err := c.BindJSON(&re)
		if err != nil {
			panic(err)
		}
		******************************************/

		len := c.Request.ContentLength
		body := make([]byte, len)
		c.Request.Body.Read(body)
		m := map[string]string{}
		err := json.Unmarshal(body, &m)
		if err != nil {
			panic(err)
		}
		phoneNum := m["phone_num"]
		uuid := m["uuid"]
		checkNum := m["check_num"]

		//检测手机号位数
		if utf8.RuneCountInString(phoneNum) != 11 {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "手机号位数不是11位",
			})
			return
		}
		//检测uuid checkNum 是否为空

		// 发送验证码网络请求
		code, err := login.SmsVerificationCode(phoneNum, uuid, checkNum)
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  err.Error(),
		})
	})
	router.Run(":8080")
}
