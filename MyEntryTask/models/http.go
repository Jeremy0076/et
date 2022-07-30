package models

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type LoginReq struct {
	RequestID string `json:"requestID"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type LoginResp struct {
	Token string `json:"token"`
}

type SignOutReq struct {
	RequestID string `json:"requestID"`
	Username  string `json:"username"`
}

type SignOutResp struct{}

type GetUserReq struct {
	RequestID string `json:"requestID"`
	Username  string `json:"username"`
}

type GetUserResp struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Picfile  string `json:"picfile"`
}

type UpdateReq struct {
	RequestID string `json:"requestID"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
}

type UpdateResp struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Picfile  string `json:"picfile"`
}
