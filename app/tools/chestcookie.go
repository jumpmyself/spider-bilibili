package tools

import (
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"spider-bilibili/app/ippool"
)

// CheckCookieValidity 检查当前的 Cookie 是否有效
func CheckCookieValidity() bool {
	// 获取代理信息
	proxyInfo := ippool.GetProxyInfo()
	// 创建新的GET请求
	request, err := http.NewRequest("GET", "https://api.bilibili.com/x/web-interface/nav", nil)
	if err != nil {
		logrus.Error("检查登录态时创建请求失败:", err)
		return false
	}

	// 添加 User-Agent 到请求头
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0")

	// 从配置中获取Cookies并添加到请求头
	cookies := []*http.Cookie{
		{Name: "bili_jct", Value: viper.GetString("cookie.bili_jct")},
		{Name: "DedeUserID", Value: viper.GetString("cookie.DedeUserID")},
		{Name: "DedeUserID__ckMd5", Value: viper.GetString("cookie.DedeUserID__ckMd5")},
		{Name: "SESSDATA", Value: viper.GetString("cookie.SESSDATA")},
		{Name: "sid", Value: viper.GetString("cookie.sid")},
	}
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	// 发送请求并获取响应
	response, err := proxyInfo.Client.Do(request)
	if err != nil {
		logrus.Error("检查登录态时请求失败:", err)
		return false
	}
	defer response.Body.Close()

	// 检查HTTP响应状态码
	if response.StatusCode != http.StatusOK {
		logrus.Error("检查登录态时返回非200状态码:", response.StatusCode)
		return false
	}

	// 解析响应体
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		logrus.Error("解码响应体失败:", err)
		return false
	}

	// 检查响应体中的isLogin字段
	isLogin, ok := result["data"].(map[string]interface{})["isLogin"].(bool)
	if !ok {
		logrus.Error("无法从响应体中获取isLogin字段")
		return false
	}

	// 如果isLogin为true，则返回true，表示Cookie有效
	return isLogin
}
