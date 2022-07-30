package dao

import (
	"MyEntryTask/models"
	"MyEntryTask/utils"
	"github.com/go-redis/redis"
)

// AddToken 添加token-username映射到redis，并设置过期时间
func AddToken(sess *models.Session) error {
	// Set 如果key存在，刷新token
	_, err := utils.Rdb.Set(sess.Token, sess.Username, utils.Conf.Redis.ExpireTime).Result()

	if err != nil {
		utils.Logs.Warn("redis setex token err: [%v]\n", err)
		return err
	}
	return nil
}

// DeleteToken 删除token
func DeleteToken(token string) error {
	_, err := GetUsernameByToken(token)
	// token已经过期
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		utils.Logs.Warn("get username by token err :[%v]\n", err)
		return err
	}

	_, err = utils.Rdb.Del(token).Result()
	if err != nil {
		utils.Logs.Warn("redis delete token err :[%v]\n", err)
		return err
	}
	return nil
}

// GetUsernameByToken 通过token获取username
func GetUsernameByToken(token string) (string, error) {
	value, err := utils.Rdb.Get(token).Result()
	if err != nil {
		if err != redis.Nil {
			utils.Logs.Warn("redis get token err :[%v]\n", err)
			return "", err
		}
	}
	return value, err
}
