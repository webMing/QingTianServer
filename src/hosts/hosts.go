package hosts

import (
	"stephanie.io/tools"
	"encoding/json"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"stephanie.io/login"
)

const (
	baseV1URL   = "api/v1"
	baseTestURL = "api/test"
)

/*建议每个包中如果有只有一个init函数;func init(){}*/

// Server  开始服务
func Server() {
	
	//chanage version v8 to v9
	binding.Validator = new(DefaultValidator)
	router := gin.Default()
	router.Use(gin.BasicAuth(nil))
	router.Use(userLoginChecker())
	gpV1   := router.Group(baseV1URL)
	// gpTest := router.Group(baseTestURL)

	serV1(gpV1)
	// serTest(gpTest)
	
	router.Run(":8080")	
	
}

func serV1(gp *gin.RouterGroup) {

	/* 获取UUID 该UUID也是可以由客户端生成(争议UUID是否需要从这里获取)*/
	gp.POST("/uuid", func(c *gin.Context) {
		c.JSON(http.StatusOK, login.UUID(c))
	})

	/* 获取图片验证码 */
	gp.POST("/imgCheckCode", func(c *gin.Context) {
		c.JSON(http.StatusOK, login.Capthca(c))
	})

	/* 用户注册 */
	gp.POST("/register", func(c *gin.Context) {
		user, err := login.Register(c)
		if err != nil {
			log.Fatalln(err)
		}
		c.JSON(http.StatusOK, user)
	})

	/* 获取验证码 */
	gp.POST("/verifyCode", func(c *gin.Context) {
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
		}
		//检测uuid checkNum 是否为空

		// 发送验证码网络请求
		code, err := login.SmsVerificationCode(phoneNum, uuid, checkNum)
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  err.Error(),
		})
	})
}

func serTest(gp *gin.RouterGroup){

	/**gin是不是并发执行? 结果:并发执行 **/
	gp.GET("/hello", func(c *gin.Context) {
		// time.Sleep(time.Second * 6)
		m := map[string]interface{}{
			"code": "0",
			"msg":  "请求成功",
		}
		var data []map[string]string
		for i := 0; i < 20; i++ {
			t := map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			}
			data = append(data, t)
		}
		m["data"] = data
		c.JSON(http.StatusOK, m)
	})
}

func userLoginChecker() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
	   token :=	c.GetHeader("Authorization")
	   if utf8.RuneCountInString(token) == 0{
		   res := tools.OuterFailtStruct()
		   res.Msg = "用户已退出登录"
		   c.AbortWithStatusJSON(http.StatusOK, res)
	   }
	})
}