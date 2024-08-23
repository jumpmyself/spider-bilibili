package ippool

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type IPInfo struct {
	IP     string
	Client *http.Client // 包含配置了代理的 http.Client
}

// 全局变量
var (
	ProxyInfo  IPInfo
	ProxyMutex sync.Mutex
)

var HttpClient *http.Client

func IpPool() IPInfo {
	tiquApiUrl := "http://proxy.siyetian.com/apis_get.html?token=AesJWLORUUx8EVJdXTqF1dOpWS14EVnFzTR1STqFUeORUQy0karFjTElUMOp3Z04EVJFTTUNme.ANzMTNwUDMycTM&limit=1&type=0&time=&split=1&split_text="
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(tiquApiUrl)
	if err != nil {
		fmt.Println("获取代理IP失败:", err)
		return IPInfo{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应内容失败:", err)
		return IPInfo{}
	}

	// 分割 IP 地址和端口号
	hostPort := strings.SplitN(string(body), ":", 2)
	if len(hostPort) != 2 {
		fmt.Println("无效的代理IP格式")
		return IPInfo{}
	}
	ip := net.ParseIP(hostPort[0])
	if ip == nil {
		fmt.Println("无效的IP地址:", hostPort[0])
		return IPInfo{}
	}

	proxyURL, err := url.Parse("http://" + net.JoinHostPort(hostPort[0], hostPort[1]))
	if err != nil {
		fmt.Println("创建代理URL失败:", err)
		return IPInfo{}
	}

	// 创建并返回配置了代理的 http.Client
	return IPInfo{
		IP: net.JoinHostPort(hostPort[0], hostPort[1]),
		Client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
			Timeout: 10 * time.Second, // 设置一个合理的超时时间
		},
	}
}

func UpdateProxyInfo() {
	ProxyMutex.Lock()
	defer ProxyMutex.Unlock()

	ProxyInfo = IpPool()
	if ProxyInfo.IP == "" || ProxyInfo.Client == nil {
		logrus.Error("没有获取到有效的代理IP或http.Client")
	} else {
		fmt.Println("更新代理 IP 地址:", ProxyInfo.IP)
	}
}
func GetProxyInfo() IPInfo {
	ProxyMutex.Lock()
	defer ProxyMutex.Unlock()
	return ProxyInfo
}
func ValidateProxy(proxy IPInfo) bool {
	request, err := http.NewRequest("GET", "https://www.bilibili.com/", nil)
	if err != nil {
		logrus.Error("验证代理IP时创建请求失败", err)
		return false
	}
	response, err := proxy.Client.Do(request)
	if err != nil || response.StatusCode != http.StatusOK {
		return false
	}
	return true
}
