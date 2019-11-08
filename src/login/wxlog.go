package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"
)

const (
	//微信请求授权范围
	wxAuthReqScope = "snsapi_userinfo"
	//微信请求授权状态
	wxAuthReqState = "123321"
	//微信APPID
	wxAppID = "wx6c96d67ffae969f1"
	//微信APPSecret
	wxAppSecret = "f205d59b8ecb9218721f4fe9e5feee70"
)


// 利用code 获取acces_token  openID
func fetchWxAccesTockenOpenID(code string) (accToken, openID ,refreshToken string, err error) {
	if utf8.RuneCountInString(code) == 0 {
		return "", "","", errors.New("Can`t fetch Wx acessToken and openID, reason:code = 0")
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", wxAppID, wxAppSecret, code)
	resp, err := http.Get(url)
	if err != nil {
		return "","","", err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", "", "",err
	}
	if _, ok := result["errcode"]; ok {
		//存在错误
		return "","","",fmt.Errorf("Can`t fetch Wx acessToken and openID, reason:errcode=%v,errmsg=%v", result["errcode"], result["errmsg"])
	}
	return result["access_token"].(string), result["openid"].(string),result["refresh_token"].(string),nil

}

// 利用refresh_token 刷新 accces_token openID  
func refreshToken(tk string) (accesToken,openID string, err error) {
    if utf8.RuneCountInString(tk) == 0 {
		return "", "", errors.New("Can`t refresh acessToken, reason:token = 0")
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s", wxAppID,tk)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", "", err
	}
	if _, ok := result["errcode"]; ok {
		return "","",fmt.Errorf("Can`t refresh Wx acessToken, reason:errcode=%v,errmsg=%v", result["errcode"], result["errmsg"])
	}
	return result["access_token"].(string), result["openid"].(string), nil
}


// 利用accessTocken ,openID获取个人信息
func fetchUserInfo(accesTocken,openID string) (map[string]interface{},error) {
	if utf8.RuneCountInString(accesTocken) == 0 || utf8.RuneCountInString(openID) == 0  {
		return nil, errors.New("Can`t fetch user info, reason:accessTocken = 0 | openID = 0")
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s",accesTocken,openID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil,err
	}
	if _, ok := result["errcode"]; ok {
		return nil,fmt.Errorf("Can`t refresh Wx acessToken, reason:errcode=%v,errmsg=%v", result["errcode"], result["errmsg"])
	}
	return result,nil
}