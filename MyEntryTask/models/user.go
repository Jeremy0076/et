package models

type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Picfile  string `json:"picfile"`
}

type Session struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Token    string `json:"token"`
	UserId   int64  `json:"userid"`
	IsLogin  bool   `json:"islogin"`
}
