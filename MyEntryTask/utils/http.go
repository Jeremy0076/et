package utils

import (
	"MyEntryTask/models"
	"encoding/json"
	"net/http"
)

// Response http响应封装
func Response(resp http.ResponseWriter, statusCode int, data interface{}, code int) {
	// 响应头
	resp.Header().Set("Content-Type", "application/json")
	// 状态行
	resp.WriteHeader(statusCode)
	// 响应实体
	res := models.Resp{
		Code: code,
		Data: data,
		Msg:  CodeMsg[code],
	}

	body, err := json.Marshal(res)
	if err != nil {
		Logs.Warn("json marshal res :[%v] err :[%v]\n", res, err)
		return
	}

	_, _ = resp.Write(body)
	return
}
