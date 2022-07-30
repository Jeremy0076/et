package routers

import (
	"MyEntryTask/httpserver/controller"
	"net/http"
)

// InitRouter 初始化路由
func InitRouter() {
	// 登陆接口
	http.HandleFunc("/api/login", controller.HandleCors(controller.Login))
	// 登出接口
	http.HandleFunc("/api/signout", controller.HandleCors(controller.SignOut))
	// 获取简介
	http.HandleFunc("/api/profile", controller.HandleCors(controller.GetProfile))
	// 修改简介
	http.HandleFunc("/api/updateProfile", controller.HandleCors(controller.UpdateProfile))
}
