package tools

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
	"golang.org/x/net/publicsuffix"
	"image/png"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	_ "path"
	"path/filepath"
	"spider-bilibili/app/model"
	"time"
)

func GetCookie() {
	logrus.Info("初始化 Cookie Jar")

	// 初始化 Cookie Jar
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		logrus.Fatal("初始化 Cookie Jar 失败:", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	logrus.Info("请求二维码")

	request, err := http.NewRequest("GET", "https://passport.bilibili.com/x/passport-login/web/qrcode/generate?source=main-fe-header&go_url=https:%2F%2Fwww.bilibili.com%2F", nil)
	if err != nil {
		logrus.Error("请求bilibili数据失败:", err)
		return
	}

	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0")

	response, err := client.Do(request)
	if err != nil {
		logrus.Error("请求 bilibili 数据失败:", err)
		return
	}
	defer response.Body.Close()

	logrus.Info("读取二维码响应内容")

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error("bilibili:读取响应内容失败：", err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logrus.Error("解码失败", err)
	}

	data := result["data"].(map[string]interface{})
	qrData := model.QRData{
		Url:       data["url"].(string),
		QrcodeKey: data["qrcode_key"].(string),
	}

	logrus.Info("生成二维码图像")

	// 定义保存二维码图片的文件夹路径
	qrcodeDir := "qrcodes"
	qrcodeFileName := "qrcode.png"
	qrcodeFilePath := filepath.Join(qrcodeDir, qrcodeFileName)

	// 确保文件夹存在，如果不存在则创建它
	err = os.MkdirAll(qrcodeDir, 0755)
	if err != nil {
		logrus.Fatal("创建文件夹失败:", err)
	}
	// 生成二维码图像
	qr, err := qrcode.New(qrData.Url, qrcode.Medium)
	if err != nil {
		logrus.Error("生成一个新的二维码图像失败", err)
	}
	// 获取二维码的图像表示
	img := qr.Image(256)

	// 创建文件并保存二维码图像
	file, err := os.Create(qrcodeFilePath)
	if err != nil {
		logrus.Error("创建二维码文件失败:", err)
	}
	defer file.Close()

	// 编码二维码图片到文件
	err = png.Encode(file, img)
	if err != nil {
		logrus.Error("编码二维码文件失败:", err)
	}

	// 更新日志信息，反映实际的文件路径
	logrus.Info("二维码保存为", qrcodeFilePath)
	time.Sleep(3 * time.Second)
	SendDingDing()
	// 等待用户扫描二维码并登录
	fmt.Println("请扫描二维码...\n60s后结束...")

	time.Sleep(60 * time.Second)

	// 轮询登录状态
	logrus.Info("轮询登录状态")
	pollUrl := fmt.Sprintf("https://passport.bilibili.com/x/passport-login/web/qrcode/poll?qrcode_key=%s&source=main-fe-header", qrData.QrcodeKey)
	pollRequest, err := http.NewRequest("GET", pollUrl, nil)
	if err != nil {
		logrus.Error("创建轮询请求失败:", err)
		return
	}

	pollRequest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0")

	pollResponse, err := client.Do(pollRequest)
	if err != nil {
		logrus.Error("请求轮询状态失败:", err)
		return
	}
	defer pollResponse.Body.Close()

	pollBody, err := io.ReadAll(pollResponse.Body)
	if err != nil {
		logrus.Error("读取轮询响应内容失败:", err)
		return
	}
	var pollResult map[string]interface{}
	err = json.Unmarshal(pollBody, &pollResult)
	if err != nil {
		logrus.Error("解码轮询响应内容失败:", err)
		return
	}

	code := pollResult["code"].(float64)

	if code == 0 {
		// 获取 cookie
		bilibiliURL, err := url.Parse("https://www.bilibili.com")
		if err != nil {
			logrus.Error("解析 bilibili URL 失败:", err)
			return
		}

		cookies := client.Jar.Cookies(bilibiliURL)

		// 拼接 Cookie 字符串
		cookieMap := make(map[string]string)
		for _, cookie := range cookies {
			cookieMap[cookie.Name] = cookie.Value
		}
		// 使用viper设置配置项
		viper.Set("cookie", cookieMap)
		err = viper.WriteConfigAs("./config.yaml")
		if err != nil {
			logrus.Error("保存cookie到配置文件失败", err)
		} else {
			fmt.Println("cookie已经保存到yaml配置文件")
		}
	} else {
		logrus.Warn("登录失败或超时")
	}
}
