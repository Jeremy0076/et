package main

import (
	"MyEntryTask/models"
	"MyEntryTask/utils"
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

const ksql = "insert into user (username,password,nickname,picfile) values (?, ?, ?, ?);"
const kpicsql = "update user set picfile = ?  where username = ?;"
const kpwdsql = "update user set password = ?  where username = ?;"

var SqlDB *sql.DB

// 插入1000w数据到mysql
func main() {
	var err error

	cfg := utils.MySQLConf{
		Username: "root",
		Password: "root",
		DB:       "entrytask",
		Host:     "127.0.0.1",
		Port:     "3306",
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	SqlDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("open sql username :[%s] pwd :[%s] host :[%s] port :[%s] db :[%s] err :[%v]\n", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DB, err)
		panic(err.Error())
		return
	}

	if err = SqlDB.Ping(); err != nil {
		log.Printf("mysql ping err: [%v]\n", err)
		return
	}

	//insert1000w()
	//updateAllPic()
	//updateAllPwd()
	updateAllPwdtoRaw()
}

func insert1000w() {
	prepare, err := SqlDB.Prepare(ksql)
	if err != nil {
		log.Printf("prepare err :[%v]\n", err)
		return
	}

	for i := 1; i < 10000000; i++ {
		user := &models.UserInfo{
			Username: "test" + strconv.Itoa(i),
			Password: "test01",
			Nickname: "测试" + strconv.Itoa(i),
			Picfile:  "../utils/page/user_img/20220723232242-436834F8.jpg",
		}
		_, err = prepare.Exec(user.Username, user.Password, user.Nickname, user.Picfile)
		if err != nil {
			log.Printf("exec i:[%d] err:[%v]\n", i, err)
		}
	}
}

func updateAllPwd() {
	prepare, err := SqlDB.Prepare(kpwdsql)
	if err != nil {
		log.Printf("prepare err :[%v]\n", err)
		return
	}

	pwd, err := utils.EncryData("test01")
	if err != nil {
		log.Printf("encry data err")
		return
	}

	for i := 1; i < 10000000; i++ {
		user := &models.UserInfo{
			Username: "test" + strconv.Itoa(i),
			Password: pwd,
		}
		_, err = prepare.Exec(user.Password, user.Username)
		if err != nil {
			log.Printf("exec i:[%d] err:[%v]\n", i, err)
		}
	}
}

func updateAllPwdtoRaw() {
	prepare, err := SqlDB.Prepare(kpwdsql)
	if err != nil {
		log.Printf("prepare err :[%v]\n", err)
		return
	}

	pwd := "test01"
	if err != nil {
		log.Printf("encry data err")
		return
	}

	for i := 1; i < 10000000; i++ {
		user := &models.UserInfo{
			Username: "test" + strconv.Itoa(i),
			Password: pwd,
		}
		_, err = prepare.Exec(user.Password, user.Username)
		if err != nil {
			log.Printf("exec i:[%d] err:[%v]\n", i, err)
		}
	}
}

func updateAllPic() {
	prepare, err := SqlDB.Prepare(kpicsql)
	if err != nil {
		log.Printf("prepare err :[%v]\n", err)
		return
	}

	for i := 1; i < 10000000; i++ {
		user := &models.UserInfo{
			Username: "test" + strconv.Itoa(i),
			Picfile:  "20220723232225-19C24544.jpg",
		}
		_, err = prepare.Exec(user.Picfile, user.Username)
		if err != nil {
			log.Printf("exec i:[%d] err:[%v]\n", i, err)
		}
	}
}
