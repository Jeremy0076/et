package main

import (
	"MyEntryTask/tcpserver/rpc/server"
	"MyEntryTask/tcpserver/service"
	"MyEntryTask/utils"
	"log"
	"os"
)

// tcp server
func main() {
	if len(os.Args) < 2 {
		log.Println("Usage：../conf/conf.ini")
		return
	}

	var err error
	// 加载配置文件
	if err = utils.Init(os.Args[1]); err != nil {
		log.Printf("load config file %s failed, err:%v\n", os.Args[1], err)
		return
	}

	defer utils.RecoverPanic()

	// 初始化mysql数据库连接
	err = utils.InitMySQL(utils.Conf.MySQL)
	if err != nil {
		utils.Logs.Error("init mysql failed, err:%v\n", err)
		return
	}
	// 回收SQLdb资源
	defer utils.CloseMySQL()

	//初始化redis连接
	err = utils.InitRedis(utils.Conf.Redis)
	if err != nil {
		utils.Logs.Error("init redis failed, err :%v\n", err)
		return
	}
	// 回收redis资源
	defer utils.CloseRedis()

	// rpc服务初始化
	rpcSvrConf := server.NewRpcServerConf("127.0.0.1", "8889")
	rpcSvr := server.NewRpcServer(rpcSvrConf)

	// 服务注册
	us := &service.UserService{}
	err = rpcSvr.Register(us)
	if err != nil {
		utils.Logs.Error("register UserService err: [%v]\n", err)
		panic(err)
	}

	rpcSvr.Run()
}
