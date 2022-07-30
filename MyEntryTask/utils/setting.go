package utils

import (
	"MyEntryTask/utils/logs"
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// section
	expectSectionApp   = "app"
	expectSectionRedis = "redis"
	expectSectionMySQL = "mysql"

	// key
	expectKeyIP         = "ip"
	expectKeyPort       = "port"
	expectKeySalt       = "salt"
	expectKeyUsername   = "username"
	expectKeyPassword   = "password"
	expectKeyDB         = "db"
	expectKeyHost       = "host"
	expectKeyExpireTime = "expiretime"
	expectKeyFilePath   = "imgpath"

	// default value
	defaultHost          = "127.0.0.1"
	defaultMySQLPort     = "3306"
	defaultRedisPort     = "6379"
	defaultSQLDB         = "user"
	defaultMySQLUsername = "root"
	defaultMySQLPassword = "root"
	defaultExpireTime    = "1800"
	defaultAppPort       = "8081"
	defaultSalt          = "seatalk"
	defaultRDB           = "0"
	defaultFilePath      = "../utils/page/user_img/"
)

// AppConf 应用程序配置
type AppConf struct {
	IP      string
	Port    string
	Salt    string
	ImgPath string
	MySQL   *MySQLConf
	Redis   *RedisConf
}

// MySQLConf 数据库mysql配置
type MySQLConf struct {
	Username string
	Password string
	DB       string
	Host     string
	Port     string
}

// RedisConf 缓存redis配置
type RedisConf struct {
	Host       string
	Port       string
	ExpireTime time.Duration
	DB         int
}

var Conf *AppConf
var Logs *logs.Logger

// Init 读取ini文件初始化
func Init(confName string) error {
	var err error
	Logs, err = logs.NewLogger(logs.LevelDebug)
	if err != nil {
		log.Printf("new logger err :[%v]\n", err)
		return err
	}

	mysql := new(MySQLConf)
	redis := new(RedisConf)

	mysql.Username = getValue(confName, expectSectionMySQL, expectKeyUsername, defaultMySQLUsername)
	mysql.Password = getValue(confName, expectSectionMySQL, expectKeyPassword, defaultMySQLPassword)
	mysql.DB = getValue(confName, expectSectionMySQL, expectKeyDB, defaultSQLDB)
	mysql.Host = getValue(confName, expectSectionMySQL, expectKeyHost, defaultHost)
	mysql.Port = getValue(confName, expectSectionMySQL, expectKeyPort, defaultMySQLPort)

	redis.Host = getValue(confName, expectSectionRedis, expectKeyHost, defaultHost)
	redis.Port = getValue(confName, expectSectionRedis, expectKeyPort, defaultRedisPort)
	redisExpireTime, err := strconv.Atoi(getValue(confName, expectSectionRedis, expectKeyExpireTime, defaultExpireTime))
	if err != nil {
		Logs.Warn("redis expire time init err :[%v]\n", err)
		return err
	}
	redis.ExpireTime = time.Duration(redisExpireTime) * time.Second
	redisRDB, err := strconv.Atoi(getValue(confName, expectSectionRedis, expectKeyDB, defaultRDB))
	if err != nil {
		Logs.Warn("redis db init err :[%v]\n", err)
		return err
	}
	redis.DB = redisRDB

	Conf = new(AppConf)
	Conf.MySQL = mysql
	Conf.Redis = redis
	Conf.IP = getValue(confName, expectSectionApp, expectKeyIP, defaultHost)
	Conf.Port = getValue(confName, expectSectionApp, expectKeyPort, defaultAppPort)
	Conf.Salt = getValue(confName, expectSectionApp, expectKeySalt, defaultSalt)
	Conf.ImgPath = getValue(confName, expectSectionApp, expectKeyFilePath, defaultFilePath)
	return nil
}

// GetConfValue 根据文件名，段名，键名，默认值获取ini的值
func getValue(filename, expectSection, expectKey, defaultValue string) string {
	file, err := os.Open(filename)
	if err != nil {
		Logs.Warn("file [%s] open err [%v]\n", filename, err)
		return ""
	}
	defer file.Close()

	// 读取文件
	reader := bufio.NewReader(file)

	var sectionName string
	for {
		lineStr, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			Logs.Warn("read string error: [%v]\n", err)
			break
		}
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}

		// 忽略注释
		if lineStr[0] == ';' {
			continue
		}

		// 行首和尾巴分别是方括号的，说明是段标记的起止符
		if lineStr[0] == '[' && lineStr[len(lineStr)-1] == ']' {
			sectionName = lineStr[1 : len(lineStr)-1]
		} else if sectionName == expectSection {
			pair := strings.Split(lineStr, "=")
			if len(pair) == 2 {
				key := strings.TrimSpace(pair[0])
				if key == expectKey {
					return strings.TrimSpace(pair[1])
				}
			}
		}
	}
	return defaultValue
}
