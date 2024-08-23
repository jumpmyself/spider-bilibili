package app

import (
	"spider-bilibili/app/ippool"
	"spider-bilibili/app/model"
	"spider-bilibili/app/router"
	"spider-bilibili/app/tools"
)

func Start() {
	tools.InitFile("app/log/", "")

	model.NewMySql()
	model.Redis()
	defer func() {
		model.RedisClose()
		model.Close()
	}()

	tools.LoadConfig()
	ippool.UpdateProxyInfo() // 初始化时更新一次代理信息

	router.Router()

}
