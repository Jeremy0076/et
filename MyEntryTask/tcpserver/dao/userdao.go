package dao

import (
	"MyEntryTask/models"
	"MyEntryTask/utils"
	"encoding/json"
)

const (
	sqlQueryUser         = "select id,username,password,nickname,picfile from user where username = ?"
	sqlUpdateUser        = "update user set nickname = ?  where username = ? "
	sqlUpdateUserProfile = "update user set nickname = ? ,picfile = ? where username = ? "
)

// QueryByUsername 通过username查询用户信息
func QueryByUsername(username string) (*models.UserInfo, error) {
	// 预编译 防止sql注入
	prepare, err := utils.SqlDB.Prepare(sqlQueryUser)
	if err != nil {
		utils.Logs.Warn("mysql query by username prepare err :[%v]\n", err)
		return nil, err
	}
	defer prepare.Close()
	queryRow := prepare.QueryRow(username)
	user := &models.UserInfo{}
	err = queryRow.Scan(&user.ID, &user.Username, &user.Password, &user.Nickname, &user.Picfile)

	return user, err
}

// UpdateUserByUsername 修改用户信息 返回修改后的结果
func UpdateUserByUsername(user *models.UserInfo) (*models.UserInfo, error) {
	var sqlStr string
	if len(user.Picfile) == 0 {
		sqlStr = sqlUpdateUser
	} else {
		sqlStr = sqlUpdateUserProfile
	}
	// 预编译
	prepare, err := utils.SqlDB.Prepare(sqlStr)
	if err != nil {
		utils.Logs.Warn("mysql update user by username prepare err :[%v]\n", err)
		return nil, err
	}
	defer prepare.Close()

	if len(user.Picfile) != 0 {
		_, err = prepare.Exec(user.Nickname, user.Picfile, user.Username)
		if err != nil {
			utils.Logs.Warn("UpdateUserByUserName exec err :[%v]", err)
			return nil, err
		}
	} else {
		_, err = prepare.Exec(user.Nickname, user.Username)
		if err != nil {
			utils.Logs.Warn("UpdateUserByUserName exec err :[%v]", err)
			return nil, err
		}
	}

	respUser, err := QueryByUsername(user.Username)
	if err != nil {
		utils.Logs.Warn("QueryByUserName exec err :[%v]", err)
		return nil, err
	}
	return respUser, nil
}

//GetUserCacheByUserName  查询用户缓存数据
func GetUserCacheByUserName(username string) (*models.UserInfo, error) {
	result, err := utils.Rdb.Get(username).Result()

	var user = &models.UserInfo{}
	if len(result) != 0 {
		err = json.Unmarshal([]byte(result), user)
		if err != nil {
			utils.Logs.Warn("json unmarshal result :[%v] err:[%v]\n", result, err)
			return nil, err
		}
		return user, err
	}
	return nil, err
}

//SetUserCacheByUserName 更改用户缓存
func SetUserCacheByUserName(user *models.UserInfo) error {
	value, err := json.Marshal(user)
	if err != nil {
		utils.Logs.Warn("json marshal user :[%v] err :[%v]\n", user, err)
		return err
	}

	// 如果存在则刷新过期时间
	err = utils.Rdb.Set(user.Username, value, 2*utils.Conf.Redis.ExpireTime).Err()
	return err
}

func DelUserCacheByUsername(username string) error {
	err := utils.Rdb.Del(username).Err()
	return err
}
