package login

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"image/png"

	// "io/ioutil"
	"path"

	"gopkg.in/go-playground/validator.v9"
	"stephanie.io/constv"
	"stephanie.io/tools"

	"github.com/gomodule/redigo/redis"

	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"
)

// Capthca 图片验证码
func Capthca(c *gin.Context) interface{} {

	f := func(c redis.Conn){
		fmt.Println(c)
	}
	tools.HandleRedisErr(f)

	type request struct {
		UUID   string `form:"uuid" json:"uuid" binding:"required"`
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
			for _, err := range err.(validator.ValidationErrors) {
				switch err.Field() {
				case "Client":
					outerRes.Msg = "client 不能为空"
				case "UUID":
					outerRes.Msg = "uuid 不能为空"
				default:
					outerRes.Msg = fmt.Sprintf("%s 该字段存在问题", err.Field())
				}
			}
		}
		return outerRes
	}

	//不管是否存在 都尝试从redis中移除
	tools.RedisHelperDel(re.UUID)

	cap := captcha.New()

	//这里使用绝对路径
	//dir := filepath.Dir("comic")
	//v, err := filepath.unixAbs("comic.ttf")
	//fmt.Printf("%s",v)
	pth, err := tools.GetCurrentFileDir()
	//部署的时候可能会出错.文件资源的路径不对.
	path := path.Join(pth, "comic.ttf")
	if err := cap.SetFont(path); err != nil {
		//panic(err.Error())
		outerRes.Code = constv.QTFetchFailtCode
		outerRes.Msg = "生产图片验证时出错"
		return outerRes
	}

	/*
	   //We can load font not only from localfile, but also from any []byte slice
	   	fontContenrs, err := ioutil.ReadFile("comic.ttf")
	   	if err != nil {
	   		panic(err.Error())
	   	}
	   	err = cap.AddFontFromBytes(fontContenrs)
	   	if err != nil {
	   		panic(err.Error())
	   	}
	*/

	cap.SetSize(128, 64)
	cap.SetDisturbance(captcha.MEDIUM)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})

	//设置图片验证码长度
	img, imgCode := cap.Create(4, captcha.NUM)
	buffer := bytes.NewBuffer(make([]byte, 0))
	err = png.Encode(buffer, img)
	if err != nil {
		outerRes.Code = constv.QTFetchFailtCode
		outerRes.Msg = "无法生成图片验证码"
		return outerRes
	}

	type inerRes struct {
		Img       string `form:"img" json:"img"`
		ImageCode string `form:"imageCode" json:"imageCode"`
	}
	inerStrc := new(inerRes)
	inerStrc.Img = base64.StdEncoding.EncodeToString(buffer.Bytes()) //base64编码
	inerStrc.ImageCode = imgCode
	outerRes.Data = inerStrc

	//uuid 存放在redis中
	expireTime := 60 * 30 //过期时间 半个小时
	_, err = tools.RedisHelperSet(re.UUID, imgCode, expireTime)
	if err != nil {
		outerRes.Code = constv.QTFetchFailtCode
		outerRes.Msg = "reids中存放图片验证码时出错"
		return outerRes
	}

	return outerRes

}
