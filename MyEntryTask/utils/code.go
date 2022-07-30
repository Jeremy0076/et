package utils

import "MyEntryTask/models"

// 1.全局错误码

const (
	// CodeSucc          succ code
	CodeSucc = 0

	// tcp 1000 ~ 2000
	// CodeTCPFailedGetUserInfo code succ
	CodeTCPFailedGetUserInfo = 1101
	// CodeTCPPasswdErr password error
	CodeTCPPasswdErr = 1102
	// CodeTCPInvalidToken invalid token
	CodeTCPInvalidToken = 1200
	// CodeTCPTokenExpired token expired
	CodeTCPTokenExpired = 1201
	// CodeTCPUserInfoNotMatch token info not match userinfo
	CodeTCPUserInfoNotMatch = 1202
	// CodeTCPFailedUpdateUserInfo update userinfo failed
	CodeTCPFailedUpdateUserInfo = 1301
	// CodeTCPInternelErr internel error
	CodeTCPInternelErr = 1401
	// CodeTCPRpcTransportErr rpc transport err
	CodeTCPRpcTransportErr = 1405
	// CodeTCPRpcNotFindSvcOrMethod rpc can't find svc or method
	CodeTCPRpcNotFindSvcOrMethod = 1406
	// CodeTCPRpcServiceErr rpc service or method err
	CodeTCPRpcServiceErr = 1407
	// CodeTCPRpcTimeout rpc time out
	CodeTCPRpcTimeout = 1408

	// HTTP 2000 ~ 3000
	// CodeInternalErr   internel err
	CodeInternalErr = 2101
	// CodeTokenNotFound missing token
	CodeTokenNotFound = 2102
	// CodeInvalidToken  token format is invalid
	CodeInvalidToken = 2103
	// CodeErrBackend    failed to comm with backend server
	CodeErrBackend = 2201
	// CodeInvalidPasswd passwd format isn't right
	CodeInvalidPasswd = 2301
	// CodeFormFileFailed formFile get error
	CodeFormFileFailed = 2401
	// CodeFileSizeErr file size not match (too small or too large)
	CodeFileSizeErr = 2402
)

// CodeMsg code to msg description
var CodeMsg = map[int]string{
	// http
	CodeSucc:           "succ",
	CodeInternalErr:    "please try again!",
	CodeTokenNotFound:  "param error: token not found",
	CodeInvalidToken:   "invalid token",
	CodeErrBackend:     "Error found!please try again!",
	CodeInvalidPasswd:  "username/passwd error!",
	CodeFormFileFailed: "Incorrect request parameters",
	CodeFileSizeErr:    "File size err (should less than 5MB)!",

	// tcp
	CodeTCPFailedGetUserInfo:     "tcp server: failed to get userinfo",
	CodeTCPPasswdErr:             "tcp server: wrong passwd",
	CodeTCPInvalidToken:          "tcp server: invalid token format",
	CodeTCPTokenExpired:          "tcp server: token expired",
	CodeTCPUserInfoNotMatch:      "tcp server: token cache info not match",
	CodeTCPFailedUpdateUserInfo:  "tcp server: failed to update userinfo",
	CodeTCPInternelErr:           "tcp server: internel error",
	CodeTCPRpcTransportErr:       "tcp server: rpc transport err",
	CodeTCPRpcNotFindSvcOrMethod: "tcp server: rpc can't find svc or method",
	CodeTCPRpcServiceErr:         "tcp server: rpc service or method err",
	CodeTCPRpcTimeout:            "rpc time out",
}

func SetRespCode(reply *models.RPCMessage, code uint32) {
	reply.Header.RespHeader.Code = code
	reply.Header.RespHeader.Msg = CodeMsg[int(code)]
}

func SetHeaderCode(header *models.Header, code uint32) {
	header.RespHeader.Code = code
	header.RespHeader.Msg = CodeMsg[int(code)]
}

// RecoverPanic 全局panic处理
func RecoverPanic() {
	if err := recover(); err != nil {
		Logs.Error("Panic :[%v]\n", err)
	}
}
