package main

import (
	"log"
	"src/global"
	"src/service"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// 读取配置文件
	if err := global.SetConfig(); err != nil {
		panic(err)
	}

	// 设置日志记录器
	//if err := global.SetAdvLogger(); err != nil {
	if err := global.SetLogger(); err != nil {
		panic(err)
	}

	// 设置ID节点
	if err := global.SetIDnode(); err != nil {
		global.Logger.Fatal(err.Error())
		return
	}

	// 设置数据库
	if err := global.SetDatabase(); err != nil {
		global.Logger.Fatal(err.Error())
		return
	}

	// 设置sessions
	//if err := global.SetSessions(); err != nil {
	//	global.Logger.Fatal(err.Error())
	//	return
	//}

	service.Start()
}
