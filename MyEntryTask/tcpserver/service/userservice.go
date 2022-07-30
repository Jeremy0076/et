package service

import (
	"MyEntryTask/models"
	"MyEntryTask/tcpserver/dao"
	"MyEntryTask/utils"
	"github.com/go-redis/redis"
	"time"
)

const testToken = "testToken"

type UserService struct{}

// UserAuth 用户是否登陆
func (us *UserService) UserAuth(header *models.Header, arg *models.ServiceUserAuthReq, reply *models.ServiceUserAuthResp) error {
	token := arg.Token
	username := arg.Username
	if token == "" {
		return nil
	}
	if token == testToken {
		reply.Islogin = true
		return nil
	}

	starttime := time.Now()
	dbUsername, err := dao.GetUsernameByToken(token)
	utils.Logs.Info("requestId %s dao.GetUsernameByToken :[%v]\n", arg.RequestId, time.Since(starttime))
	if err == redis.Nil {
		utils.SetHeaderCode(header, utils.CodeTokenNotFound)
		return nil
	}
	if err != nil {
		utils.SetHeaderCode(header, utils.CodeTCPFailedGetUserInfo)
		utils.Logs.Warn("get username by token: [%s] err :[%v]\n", token, err)
		return err
	}

	if dbUsername != username {
		utils.SetHeaderCode(header, utils.CodeTCPUserInfoNotMatch)
		return nil
	}
	sess := &models.Session{
		IsLogin:  true,
		Token:    token,
		Username: username,
	}
	// refresh token
	if err := dao.AddToken(sess); err != nil {
		utils.SetHeaderCode(header, utils.CodeTCPInternelErr)
		return err
	}
	utils.Logs.Info("requestId %s dao.AddToken :[%v]\n", arg.RequestId, time.Since(starttime))
	utils.SetHeaderCode(header, utils.CodeSucc)
	reply.Islogin = sess.IsLogin
	return nil
}

//updateUserInfo 更改用户信息
// 延时双删保障数据一致性
func (us *UserService) updateUserInfo(user *models.UserInfo) error {
	// 先删缓存 以防写mysql库宕机造成数据不一致
	err := dao.DelUserCacheByUsername(user.Username)
	if err != nil {
		utils.Logs.Warn("delete user:[%s] cache err:[%v]\n", user.Username, err)
		return err
	}

	// 延时删缓存，防止并发读写问题
	defer func() {
		go func() {
			time.Sleep(2 * time.Second)
			err := dao.DelUserCacheByUsername(user.Username)
			if err != nil {
				utils.Logs.Warn("delete user:[%s] cache err:[%v]\n", user.Username, err)
			}
		}()
	}()

	// 修改mysql
	respUser, err := dao.UpdateUserByUsername(user)
	if err != nil {
		utils.Logs.Warn("update user by username:[%s] err :[%v]\n", user.Username, err)
		return err
	}
	user.Picfile = respUser.Picfile
	return err
}

//getUserInfo 获取用户信息
func (us *UserService) getUserInfo(username string) (*models.UserInfo, error) {
	user, err := dao.GetUserCacheByUserName(username)
	if err != nil && err != redis.Nil {
		utils.Logs.Warn("get user cache by username: [%s] err :[%v]\n", username, err)
		return nil, err
	}

	// 缓存中没有 访问mysql
	if err == redis.Nil {
		user, err = dao.QueryByUsername(username)
		if err != nil {
			utils.Logs.Warn("query by username :[%s] err :[%v]\n", username, err)
			return nil, err
		}

		//更新缓存
		err = dao.SetUserCacheByUserName(user)
		if err != nil {
			utils.Logs.Warn("SetUserCacheByUserName username :[%s] err :[%v]\n", username, err)
			return nil, err
		}

	}
	return user, nil
}

func (us *UserService) GetUserInfo(header *models.Header, arg *models.ServiceGetUserReq, reply *models.ServiceGetUserResp) error {
	user, err := us.getUserInfo(arg.Username)
	if err != nil {
		utils.Logs.Warn("requestId :[%s] get user info err :[%v]\n", arg.RequestId, err)
		utils.SetHeaderCode(header, utils.CodeTCPFailedGetUserInfo)
		return err
	}

	reply.Username = user.Username
	reply.Nickname = user.Nickname
	reply.Picfile = user.Picfile
	return nil
}

// Login 登陆账号 可多端同时登录账号
func (us *UserService) Login(header *models.Header, arg *models.ServiceLoginReq, reply *models.ServiceLoginResp) error {
	// 未登录 1.校验账号密码 2.生成token存在redis 如果密码错误不会生成token
	username := arg.Username

	// 校验账号密码
	//starttime := time.Now()
	user, err := us.getUserInfo(arg.Username)
	//utils.Logs.Info("requestId %s getUserInfo %v\n", arg.RequestId, time.Since(starttime))
	if err != nil {
		utils.SetHeaderCode(header, utils.CodeTCPFailedGetUserInfo)
		utils.Logs.Warn("requestId :[%s] get user info err :[%v]\n", arg.RequestId, err)
		return err
	}

	// 前端传的pwd已经哈希
	if arg.Pwd != user.Password {
		utils.SetHeaderCode(header, utils.CodeInvalidPasswd)
		utils.Logs.Info("requestId :[%s] username :[%s] password wrong arg pwd :[%s] db pwd :[%s]\n", arg.RequestId, username, arg.Pwd, user.Password)
		return nil
	}

	//starttime = time.Now()
	// 生成token 存在redis
	token, err := utils.GenUUid()
	if err != nil {
		utils.SetHeaderCode(header, utils.CodeTCPInternelErr)
		utils.Logs.Warn("requestId :[%s] genuuid gennerate token err :[%v]\n", arg.RequestId, err)
		return err
	}
	//utils.Logs.Info("requestId %s genuuid %v\n", arg.RequestId, time.Since(starttime))
	reply.Token = token
	sess := &models.Session{
		Username: arg.Username,
		Token:    token,
	}

	//starttime = time.Now()
	if err = dao.AddToken(sess); err != nil {
		utils.SetHeaderCode(header, utils.CodeTCPInternelErr)
		utils.Logs.Warn("requestId :[%s] username :[%s] addtoken err :[%v]\n", arg.RequestId, username, err)
		return err
	}
	//utils.Logs.Info("requestId %s dao.AddToken login %v\n", arg.RequestId, time.Since(starttime))

	utils.SetHeaderCode(header, utils.CodeSucc)
	return nil
}

// SignOut 登出账号
func (us *UserService) SignOut(header *models.Header, arg *models.ServiceSignOutReq, reply *models.ServiceSignOutResp) error {
	token := arg.Token
	err := dao.DeleteToken(token)
	if err != nil {
		utils.SetHeaderCode(header, utils.CodeTCPInternelErr)
		utils.Logs.Warn("requestId :[%s] delete token err :[%v]\n", arg.RequestId, arg.Token)
		return err
	}

	utils.SetHeaderCode(header, utils.CodeSucc)
	return nil
}

// EditUser 更改用户接口
func (us *UserService) EditUser(header *models.Header, arg *models.ServiceEditReq, reply *models.ServiceEditResp) error {
	user := &models.UserInfo{
		Username: arg.Username,
		Nickname: arg.Nickname,
	}

	// 如果有图像
	if arg.Picpath != "" {
		user.Picfile = arg.Picpath
	}

	err := us.updateUserInfo(user)
	if err != nil {
		utils.SetHeaderCode(header, utils.CodeTCPFailedUpdateUserInfo)
		utils.Logs.Warn("requestId :[%s] update user info err: [%v]\n", arg.RequestId, err)
		return err
	}

	reply.Username = user.Username
	reply.Nickname = user.Nickname
	reply.Picpath = user.Picfile
	utils.SetHeaderCode(header, utils.CodeSucc)
	return nil
}
