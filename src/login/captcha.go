package login

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/color"
	"image/png"

	// "io/ioutil"
	"path"

	"stephanie.io/tools"

	"unicode/utf8"

	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"
)

// Capthca 图片验证码
func Capthca(c *gin.Context) (user map[string]interface{}, err error) {

	type replyJSON struct {
		UUID       string `json:"uuid"`
		UserClient string `json:"user_client"`
	}
	var re replyJSON
	err = c.BindJSON(&re)
	if err != nil {
		panic(err)
	}

	// uuid 位数不对
	if utf8.RuneCountInString(re.UUID) == 0 {
		u := map[string]interface{}{
			"code": 1,
			"msg":  "uuid 不能为空",
		}
		// 不写err,上层如果捕捉到err会停止服务
		// err = errors.New("手机号位数不对")
		// return u,nil
		return u, nil
	}

	cap := captcha.New()

	//这里使用绝对路径
	// dir := filepath.Dir("comic")
	//v, err := filepath.unixAbs("comic.ttf")
	//fmt.Printf("%s",v)
	pth, err := tools.GetCurrentFileDir()
	// 部署的时候可能会出错.文件资源的路径不对.
	path := path.Join(pth, "comic.ttf")
	if err := cap.SetFont(path); err != nil {
		panic(err.Error())
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

	img, str := cap.Create(6, captcha.NUM)
	// var bfs []byte
	buffer := bytes.NewBuffer(make([]byte, 0))
	err = png.Encode(buffer, img)
	if err != nil {
		return nil, errors.New("无法获取图片校验码")
	}

	u := map[string]interface{}{
		"code": 0,
		"msg":  "获取成功",
		"img":  base64.StdEncoding.EncodeToString(buffer.Bytes()), //base64  编码
		"num":  str,
	}
	return u, nil

}
