package tools

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/CodyGuo/dingtalk"
//	"github.com/CodyGuo/glog"
//	"github.com/skip2/go-qrcode"
//	"github.com/spf13/viper"
//	"image/png"
//	"io/ioutil"
//	"net/http"
//	"net/http/cookiejar"
//	"os"
//)
//
//type QRData struct {
//	Url       string `json:"url"`
//	QrcodeKey string `json:"qrcode_key"`
//}
//
//func GetQRUrl() (QRData, error) {
//	client := &http.Client{}
//	req, err := http.NewRequest("GET", "https://passport.bilibili.com/x/passport-login/web/qrcode/generate?source=main-fe-header&go_url=https:%2F%2Fwww.bilibili.com%2F", nil)
//	if err != nil {
//		return QRData{}, err
//	}
//
//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
//
//	resp, err := client.Do(req)
//	if err != nil {
//		return QRData{}, err
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return QRData{}, err
//	}
//
//	var result map[string]interface{}
//	err = json.Unmarshal(body, &result)
//	if err != nil {
//		return QRData{}, err
//	}
//
//	data := result["data"].(map[string]interface{})
//	qrData := QRData{
//		Url:       data["url"].(string),
//		QrcodeKey: data["qrcode_key"].(string),
//	}
//
//	return qrData, nil
//}
//
//func MakeQRCode(url string) error {
//	qr, err := qrcode.New(url, qrcode.Medium)
//	if err != nil {
//		return err
//	}
//
//	img := qr.Image(256)
//	file, err := os.Create("QRPic/qrcode.png")
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	err = png.Encode(file, img)
//	if err != nil {
//		return err
//	}
//
//	fmt.Println("QR Code saved as qrcode.png")
//	return nil
//}
//
//func saveCookie(data map[string]string, id string) error {
//	dir := "./bilibili_login/cookie"
//	err := os.MkdirAll(dir, 0755)
//	if err != nil {
//		return err
//	}
//
//	filePath := fmt.Sprintf("%s/%s.json", dir, id)
//	file, err := os.Create(filePath)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	jsonData, err := json.Marshal(data)
//	if err != nil {
//		return err
//	}
//
//	_, err = file.Write(jsonData)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func GetResponse(data string) error {
//	jar, _ := cookiejar.New(nil)
//	client := &http.Client{
//		Jar: jar,
//	}
//	url := fmt.Sprintf("https://passport.bilibili.com/x/passport-login/web/qrcode/poll?qrcode_key=%s&source=main-fe-header", data)
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		return err
//	}
//
//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36 QIHU 360SE")
//
//	resp, err := client.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return err
//	}
//
//	var result map[string]interface{}
//	err = json.Unmarshal(body, &result)
//	if err != nil {
//		return err
//	}
//
//	dataLogin := result["data"].(map[string]interface{})
//	code := int(dataLogin["code"].(float64))
//	if code == 0 {
//		cookies := client.Jar.Cookies(req.URL)
//		fmt.Println("cookies", cookies)
//		cookieData := make(map[string]string)
//		for _, cookie := range cookies {
//			cookieData[cookie.Name] = cookie.Value
//			key := "cookie.bilibili." + cookie.Name
//			viper.Set(key, cookie.Value)
//			err = viper.WriteConfigAs("./config.yaml")
//			if err != nil {
//				fmt.Println("写入配置文件时出错：", err)
//			}
//		}
//
//	}
//
//	return nil
//}
//
//func loadCookie() (map[string]string, error) {
//	filePath := "./bilibili_login/cookie/test.json"
//	file, err := os.Open(filePath)
//	if err != nil {
//		fmt.Println("a")
//		return nil, err
//	}
//	defer file.Close()
//
//	var cookieData map[string]string
//	err = json.NewDecoder(file).Decode(&cookieData)
//	if err != nil {
//		return nil, err
//	}
//
//	return cookieData, nil
//}
//
//func IfLogin() bool {
//
//	client := &http.Client{}
//	url := "https://api.bilibili.com/x/web-interface/nav"
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		return false
//	}
//
//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36 QIHU 360SE")
//
//	req.AddCookie(&http.Cookie{Name: "bili_jct", Value: viper.GetString("cookie.bilibili.bili_jct")})
//	req.AddCookie(&http.Cookie{Name: "DedeUserID", Value: viper.GetString("cookie.bilibili.DedeUserID")})
//	req.AddCookie(&http.Cookie{Name: "DedeUserID__ckMd5", Value: viper.GetString("cookie.bilibili.DedeUserID__ckMd5")})
//	req.AddCookie(&http.Cookie{Name: "SESSDATA", Value: viper.GetString("cookie.bilibili.SESSDATA")})
//	req.AddCookie(&http.Cookie{Name: "sid", Value: viper.GetString("cookie.bilibili.sid")})
//	fmt.Println(req)
//	resp, err := client.Do(req)
//	if err != nil {
//		return false
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return false
//	}
//
//	var result map[string]interface{}
//	err = json.Unmarshal(body, &result)
//	if err != nil {
//		return false
//	}
//
//	if result["code"].(float64) == 0 {
//		return true
//	}
//
//	//personData := result["data"].(map[string]interface{})
//	//userName := personData["uname"].(string)
//	//coinNum := personData["money"]
//	//level := personData["level_info"].(map[string]interface{})["current_level"].(float64)
//	//face := personData["face"].(string)
//	//
//	//fmt.Printf("用户名: %s\n硬币数量: %.2f\n等级: %d\n头像链接: %s\n", userName, coinNum, int(level), face)
//	return false
//
//}
//func SendDingDing() {
//
//	webHook := "https://oapi.dingtalk.com/robot/send?access_token=d4ea4aadf0ed52873ca00f7e604209d9b619fe212fd6f3aabcf46c77f1736df7"
//	// 密钥，机器人安全设置页面，加签一栏勾选之后下面显示的SEC开头的字符串
//	secret := "SEC1254b4ea889fac96c100a6c99f9ef85d40d83b7cf28bfba875a1f5cbe3714390"
//	dt := dingtalk.New(webHook, dingtalk.WithSecret(secret))
//
//	// markdown类型
//	markdownTitle := "markdown"
//	markdownText := "请您扫码保持Bilibili登录态 @17516118727 \n" +
//		"> 扫码完成后请稍等片刻，操作失败此条信息将重新发送" +
//		"> ![screenshot](http://52rreq.natappfree.cc/getImage?imageName=./QRPic/qrcode.png)\n"
//	if err := dt.RobotSendMarkdown(markdownTitle, markdownText); err != nil {
//		glog.Fatal(err)
//	}
//
//}
