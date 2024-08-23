package tools

import (
	"github.com/CodyGuo/dingtalk"
	"github.com/CodyGuo/glog"
	"github.com/astaxie/beego"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BaseController struct {
	beego.Controller
}

// SendDingMsg 主要代码
func (c *BaseController) SendDingMsg(msg string) {
	//请求地址模板
	webHook := `https://oapi.dingtalk.com/robot/send?access_token=xxxx`
	content := `{"msgtype": "text",
			"text": {"content": "` + msg + `"}
		}`
	//创建一个请求
	req, err := http.NewRequest("POST", webHook, strings.NewReader(content))
	if err != nil {
		// handle error
	}

	client := &http.Client{}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	//发送请求
	resp, err := client.Do(req)
	//关闭请求
	defer resp.Body.Close()

	if err != nil {
		// handle error
	}
}

func SendDingDing() {

	webHook := "https://oapi.dingtalk.com/robot/send?access_token=452d024d9e8383b7d87ef3f2becde90461fd19e61a5a2f083aa0fa89f73ec65e"
	// 密钥，机器人安全设置页面加签勾选之后下面显示的SEC开头的字符串
	secret := "SEC8bd30a47dd729c126ea5f79cdb841cbb32e7642391326e551a4b8de1f0f63de7"
	dt := dingtalk.New(webHook, dingtalk.WithSecret(secret))

	// markdown类型
	markdownTitle := "markdown"
	currentTime := time.Now().Unix() //获取当前时间戳
	markdownText := "该扫描二维码重新登录了@17516118727\n" +
		"> 扫码完成后请稍等片刻，操作失败此条信息将重新发送" +
		"> ![screenshot](http://hzgw8e.natappfree.cc/getimage?imageName=qrcode.png&t=" + strconv.FormatInt(currentTime, 10) + ")\n"
	if err := dt.RobotSendMarkdown(markdownTitle, markdownText); err != nil {
		glog.Fatal(err)
	}
}
