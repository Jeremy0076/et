package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	typeMySQL = "mysql"
)

var (
	SqlDB *sql.DB
)

func InitMySQL(cfg *MySQLConf) error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	SqlDB, err = sql.Open(typeMySQL, dsn)
	if err != nil {
		Logs.Error("open sql username :[%s] pwd :[%s] host :[%s] port :[%s] db :[%s] err :[%v]\n", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DB, err)
		panic(err.Error())
		return err
	}

	// 最大连接数
	SqlDB.SetMaxOpenConns(1000)
	// 闲置连接数
	SqlDB.SetMaxIdleConns(200)
	// 最大存活时间
	SqlDB.SetConnMaxLifetime(120 * time.Second)

	if err = SqlDB.Ping(); err != nil {
		Logs.Error("mysql ping err: [%v]\n", err)
		return err
	}
	return nil
}

func CloseMySQL() {
	var err error
	err = SqlDB.Close()
	if err != nil {
		panic(err.Error())
		return
	}
}
