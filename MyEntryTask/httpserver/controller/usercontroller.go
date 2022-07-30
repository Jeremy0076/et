package controller

import (
	"MyEntryTask/models"
	"MyEntryTask/utils"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	HeaderToken = "Authorization"
	PIC         = "PicFile"
	kUserName   = "username"
	kNickName   = "nickname"
	layOut      = "20060102150405"
	symbolLink  = "-"
	kReqId      = "requestID"
)

// UserAuth token鉴权
func UserAuth(requestId, token string, username string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	authArgs := &models.ServiceUserAuthReq{Username: username, Token: token, RequestId: requestId}
	authReply := &models.ServiceUserAuthResp{}
	code, err := CallRPC(ctx, "UserService.UserAuth", authArgs, authReply)
	if err != nil {
		utils.Logs.Warn("requestId :[%s] call rpc err :[%v]\n", requestId, err)
		return false
	}
	if code != utils.CodeSucc {
		utils.Logs.Warn("requestId :[%s] user auth code incorrect, code:[%d]\n", requestId, code)
		return false
	}
	return authReply.Islogin
}

func Login(resp http.ResponseWriter, req *http.Request) {
	defer utils.RecoverPanic()
	if req.Method != http.MethodPost {
		utils.Logs.Warn("request method wrong, method is %s\n", req.Method)
		utils.Response(resp, http.StatusMethodNotAllowed, nil, utils.CodeFormFileFailed)
		return
	}
	starttime := time.Now()
	// 解析请求
	token := req.Header.Get(HeaderToken)
	loginReq := &models.LoginReq{}
	loginResp := &models.LoginResp{}
	loginReq.RequestID = req.URL.Query().Get(kReqId)
	err := json.NewDecoder(req.Body).Decode(loginReq)
	if err != nil {
		utils.Logs.Warn("requestId :[%s] json decode err :[%v]\n", loginReq.RequestID, err)
		utils.Response(resp, http.StatusBadRequest, nil, utils.CodeFormFileFailed)
		return
	}
	defer req.Body.Close()

	////校验签名和时间戳
	//ok, err := utils.CheckSignAndTimestamp(req, token, loginReq.Username, loginReq.Password)
	//if err != nil {
	//	utils.Logs.Warn("requestId :[%s] check sign and timestamp err :[%v]\n", loginReq.RequestID, err)
	//	utils.Response(resp, http.StatusInternalServerError, nil, utils.CodeInternalErr)
	//	return
	//}
	//if !ok {
	//	utils.Logs.Info("requestId :[%s] sign wrong, request refuse, username: [%s]\n", loginReq.RequestID, loginReq.Username)
	//	utils.Response(resp, http.StatusForbidden, nil, utils.CodeFormFileFailed)
	//	return
	//}

	// 校验格式
	if !utils.CheckUsername(loginReq.Username) || !utils.CheckPassword(loginReq.Password) {
		utils.Logs.Warn("requestId :[%s] username :[%s] password incorrect\n", loginReq.RequestID, loginReq.Username)
		utils.Response(resp, http.StatusOK, nil, utils.CodeInvalidPasswd)
		return
	}
	//fmt.Printf("requestId %s use time check %v \n", loginReq.RequestID, time.Since(starttime))

	// token鉴权
	isLogin := UserAuth(loginReq.RequestID, token, loginReq.Username)
	// 已经登陆
	if isLogin {
		loginResp.Token = token
		utils.Logs.Info("API:[/api/login] requestId:[%s] usertime :[%v] username :[%s] login sussess\n", loginReq.RequestID, time.Since(starttime), loginReq.Username)
		utils.Response(resp, http.StatusOK, loginResp, utils.CodeSucc)
		return
	}
	//fmt.Printf("requestId %s use time auth %v\n", loginReq.RequestID, time.Since(starttime))

	// 未登录，调用rpc生成token
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	args := &models.ServiceLoginReq{Username: loginReq.Username, Pwd: loginReq.Password, RequestId: loginReq.RequestID}
	reply := &models.ServiceLoginResp{}
	code, err := CallRPC(ctx, "UserService.Login", args, reply)
	if err != nil {
		utils.Logs.Warn("requestId [%s] call rpc err :[%v]\n", loginReq.RequestID, err)
		utils.Response(resp, http.StatusInternalServerError, nil, utils.CodeTCPRpcServiceErr)
		return
	}
	//fmt.Printf("requestId %s use time login %v\n", loginReq.RequestID, time.Since(starttime))
	utils.Logs.Info("API:[/api/login] requestId:[%s] usertime :[%v] username :[%s] login sussess\n", loginReq.RequestID, time.Since(starttime), loginReq.Username)
	// 组装resp返回
	loginResp.Token = reply.Token
	utils.Response(resp, http.StatusOK, loginResp, int(code))
}

func SignOut(resp http.ResponseWriter, req *http.Request) {
	defer utils.RecoverPanic()
	if req.Method != http.MethodPost {
		utils.Logs.Warn("request method wrong, method is %s\n", req.Method)
		utils.Response(resp, http.StatusMethodNotAllowed, nil, utils.CodeFormFileFailed)
		return
	}

	// 解析请求
	token := req.Header.Get(HeaderToken)
	signoutReq := &models.SignOutReq{}
	signoutResp := &models.SignOutResp{}
	signoutReq.RequestID = req.URL.Query().Get(kReqId)
	err := json.NewDecoder(req.Body).Decode(signoutReq)
	if err != nil {
		utils.Logs.Warn("requestId :[%s] json decode err :[%v]\n", signoutReq.RequestID, err)
		utils.Response(resp, http.StatusBadRequest, nil, utils.CodeFormFileFailed)
		return
	}
	defer req.Body.Close()

	//// 校验签名和时间戳
	//ok, err := utils.CheckSignAndTimestamp(req, token, signoutReq.Username)
	//if err != nil {
	//	utils.Logs.Warn("requestId :[%s] check sign and timestamp err :[%v]\n", signoutReq.RequestID, err)
	//	utils.Response(resp, http.StatusInternalServerError, nil, utils.CodeInternalErr)
	//	return
	//}
	//if !ok {
	//	utils.Logs.Info("requestId :[%s] sign wrong, request refuse, username: [%s]\n", signoutReq.RequestID, signoutReq.Username)
	//	utils.Response(resp, http.StatusForbidden, nil, utils.CodeFormFileFailed)
	//	return
	//}

	// token鉴权
	isLogin := UserAuth(signoutReq.RequestID, token, signoutReq.Username)
	// 未登陆
	if !isLogin {
		utils.Logs.Warn("requestId :[%s] user :[%s] not login or token err\n", signoutReq.RequestID, signoutReq.Username)
		utils.Response(resp, http.StatusUnauthorized, nil, utils.CodeTokenNotFound)
		return
	}

	// 调用rpc
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	args := &models.ServiceSignOutReq{
		Token: token, RequestId: signoutReq.RequestID,
	}
	reply := &models.ServiceSignOutResp{}
	code, err := CallRPC(ctx, "UserService.SignOut", args, reply)
	if err != nil {
		utils.Logs.Warn("requestId :[%s] call rpc err :[%v]\n", signoutReq.RequestID, err)
		utils.Response(resp, http.StatusInternalServerError, nil, utils.CodeTCPRpcServiceErr)
		return
	}

	utils.Logs.Info("API:[/api/signout] requestId:[%s] username :[%s] signout sussess\n", signoutReq.RequestID, signoutReq.Username)
	utils.Response(resp, http.StatusOK, signoutResp, int(code))
}

func GetProfile(resp http.ResponseWriter, req *http.Request) {
	defer utils.RecoverPanic()
	if req.Method != http.MethodGet {
		utils.Logs.Warn("request method wrong, method is %s\n", req.Method)
		utils.Response(resp, http.StatusMethodNotAllowed, nil, utils.CodeFormFileFailed)
		return
	}

	// 解析请求
	token := req.Header.Get(HeaderToken)
	getprofileReq := &models.GetUserReq{}
	getprofileResp := &models.GetUserResp{}
	getprofileReq.Username = req.URL.Query().Get("username")
	getprofileReq.RequestID = req.URL.Query().Get(kReqId)
	if !utils.CheckUsername(getprofileReq.Username) {
		utils.Logs.Warn("requestId :[%s] username :[%s] incorrect\n", getprofileReq.RequestID, getprofileReq.Username)
		utils.Response(resp, http.StatusOK, nil, utils.CodeFormFileFailed)
		return
	}

	//// 校验签名和时间戳
	//ok, err := utils.CheckSignAndTimestamp(req, token, getprofileReq.Username)
	//if err != nil {
	//	utils.Logs.Warn("requestId :[%s] check sign and timestamp err :[%v]\n", getprofileReq.RequestID, err)
	//	utils.Response(resp, http.StatusInternalServerError, nil, utils.CodeInternalErr)
	//	return
	//}
	//if !ok {
	//	utils.Logs.Info("requestId :[%s] sign wrong, request refuse, username: [%s]\n", getprofileReq.RequestID, getprofileReq.Username)
	//	utils.Response(resp, http.StatusForbidden, nil, utils.CodeFormFileFailed)
	//	return
	//}

	// token鉴权
	isLogin := UserAuth(getprofileReq.RequestID, token, getprofileReq.Username)
	// 未登陆
	if !isLogin {
		utils.Logs.Warn("requestId :[%s] user :[%s] not login or token err\n", getprofileReq.RequestID, getprofileReq.Username)
		utils.Response(resp, http.StatusUnauthorized, nil, utils.CodeTokenNotFound)
		return
	}

	// 调用rpc
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	args := &models.ServiceGetUserReq{
		Username: getprofileReq.Username, RequestId: getprofileReq.RequestID,
	}
	reply := &models.ServiceGetUserResp{}
	code, err := CallRPC(ctx, "UserService.GetUserInfo", args, reply)
	if err != nil {
		utils.Logs.Warn("requestId :[%s] call rpc err :[%v]\n", getprofileReq.RequestID, err)
		utils.Response(resp, http.StatusInternalServerError, nil, utils.CodeTCPRpcServiceErr)
		return
	}

	utils.Logs.Info("API:[/api/profile] requestId :[%s] username :[%s] get profile sussess\n", getprofileReq.RequestID, getprofileReq.Username)
	// 组装resp返回
	getprofileResp.Username = reply.Username
	getprofileResp.Nickname = reply.Nickname
	getprofileResp.Picfile = reply.Picfile
	utils.Response(resp, http.StatusOK, getprofileResp, int(code))
}

func UpdateProfile(resp http.ResponseWriter, req *http.Request) {
	defer utils.RecoverPanic()
	if req.Method != http.MethodPost {
		utils.Logs.Warn("request method wrong, method is %s\n", req.Method)
		utils.Response(resp, http.StatusMethodNotAllowed, nil, utils.CodeFormFileFailed)
		return
	}

	// 解析请求
	token := req.Header.Get(HeaderToken)
	updateProfileReq := &models.UpdateReq{}
	updateProfileResp := &models.UpdateResp{}
	updateProfileReq.RequestID = req.URL.Query().Get(kReqId)
	updateProfileReq.Username = req.FormValue(kUserName)
	updateProfileReq.Nickname = req.FormValue(kNickName)
	if !utils.CheckUsername(updateProfileReq.Username) {
		utils.Logs.Warn("requestId :[%s] username :[%s] incorrect\n", updateProfileReq.RequestID, updateProfileReq.Username)
		utils.Response(resp, http.StatusOK, nil, utils.CodeFormFileFailed)
		return
	}

	if updateProfileReq.Nickname == "" {
		utils.Logs.Warn("requestId :[%s] nickname :[%s] incorrect\n", updateProfileReq.RequestID, updateProfileReq.Username)
		utils.Response(resp, http.StatusOK, nil, utils.CodeFormFileFailed)
		return
	}

	//// 校验签名和时间戳
	//ok, err := utils.CheckSignAndTimestamp(req, token, updateProfileReq.Username, updateProfileReq.Nickname)
	//if err != nil {
	//	utils.Logs.Warn("requestId :[%s] check sign and timestamp err :[%v]\n", updateProfileReq.RequestID, err)
	//	utils.Response(resp, http.StatusInternalServerError, nil, utils.CodeInternalErr)
	//	return
	//}
	//if !ok {
	//	utils.Logs.Info("requestId :[%s] sign wrong, request refuse, username: [%s]\n", updateProfileReq.RequestID, updateProfileReq.Username)
	//	utils.Response(resp, http.StatusForbidden, nil, utils.CodeFormFileFailed)
	//	return
	//}

	// token鉴权
	isLogin := UserAuth(updateProfileReq.RequestID, token, updateProfileReq.Username)
	// 未登陆
	if !isLogin {
		utils.Logs.Warn("requestId :[%s] user :[%s] not login or token err\n", updateProfileReq.RequestID, updateProfileReq.Username)
		utils.Response(resp, http.StatusUnauthorized, nil, utils.CodeTokenNotFound)
		return
	}

	var picPath string
	var picBytes []byte
	var filename string
	// 解析图片文件
	picFile, _, _ := req.FormFile(PIC)

	// 如果有图像
	if picFile != nil {
		defer picFile.Close()
		// 校验文件后缀 png jpeg jpg gif
		picBytes, _ = io.ReadAll(picFile)
		ok, ext := utils.CheckImageCode(&picBytes)
		if !ok {
			utils.Logs.Info("requestId :[%s] pic ext wrong, request refuse, username: [%s]\n", updateProfileReq.RequestID, updateProfileReq.Username)
			utils.Response(resp, http.StatusForbidden, nil, utils.CodeFormFileFailed)
			return
		}

		uuid, err := utils.GenUUid()
		if err != nil {
			utils.Logs.Warn("requestId :[%s] gen uuid err :[%v]\n", updateProfileReq.RequestID, err)
			utils.Response(resp, http.StatusInternalServerError, nil, utils.CodeInternalErr)
			return
		}

		filename = time.Now().Format(layOut) + symbolLink + strings.Split(uuid, "-")[0] + ext
		picPath = utils.Conf.ImgPath + filename

		err = ioutil.WriteFile(picPath, picBytes, 0777)
		if err != nil {
			utils.Logs.Warn("requestId :[%s] write file err :[%v]\n", updateProfileReq.RequestID, err)
			utils.Response(resp, http.StatusInternalServerError, nil, utils.CodeInternalErr)
			return
		}
	}

	// 调用rpc
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	args := &models.ServiceEditReq{
		Username: updateProfileReq.Username, Nickname: updateProfileReq.Nickname, Picpath: filename, RequestId: updateProfileReq.RequestID,
	}
	reply := &models.ServiceEditResp{}
	code, err := CallRPC(ctx, "UserService.EditUser", args, reply)
	if err != nil {
		utils.Logs.Warn("requestId :[%s] call rpc err :[%v]\n", updateProfileReq.RequestID, err)
		utils.Response(resp, http.StatusInternalServerError, nil, utils.CodeTCPRpcServiceErr)
		return
	}

	utils.Logs.Info("API:[/api/updateProfile] requestId :[%s] username :[%s] update profile sussess\n", updateProfileReq.RequestID, updateProfileReq.Username)
	// 组装resp返回
	updateProfileResp.Username = reply.Username
	updateProfileResp.Nickname = reply.Nickname
	updateProfileResp.Picfile = reply.Picpath
	utils.Response(resp, http.StatusOK, updateProfileResp, int(code))
}

func HandleCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 允许访问所有域
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许header值
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		// 允许携带cookie
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		// 允许请求方法
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// 允许返回任意格式数据
		w.Header().Set("content-type", "*")

		// 跨域第一次OPTIONS请求，直接放行
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	}
}
