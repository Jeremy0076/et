package main

import (
	"MyEntryTask/httpserver/controller"
	"MyEntryTask/httpserver/routers"
	"MyEntryTask/utils"
	"fmt"
	"net/http"
	"os"
)

// http server
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage：../conf/conf.ini")
		return
	}

	// 加载配置文件
	if err := utils.Init(os.Args[1]); err != nil {
		fmt.Printf("load config file %s failed, err:%v\n", os.Args[1], err)
		return
	}

	defer utils.RecoverPanic()
	// 初始化rpc pool
	err := controller.PoolInit()
	if err != nil {
		panic(err)
	}

	//初始化相关路由
	routers.InitRouter()

	utils.Logs.Info("http server start at %s:%s", utils.Conf.IP, utils.Conf.Port)
	//监听端口
	err = http.ListenAndServe(":"+utils.Conf.Port, nil)
	if err != nil {
		utils.Logs.Info("http server listen and serve err %v\n", err)
		panic(err)
	}
}
