package mysql

import (
	"CloudDisk/internal/model"
	"CloudDisk/tools/ini"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var DB *xorm.Engine

func init() {
	cfg := ini.GetSectionMap("database")
	drivename := cfg["drivename"]
	// DsName := "root:990219@(127.0.0.1:3306)/db_cloud_disk?charset=utf8"
	DsName := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8", cfg["user"], cfg["password"], cfg["domain"], cfg["port"], cfg["dbname"])
	println(DsName)
	err := errors.New("")
	DB, err = xorm.NewEngine(drivename, DsName)
	if nil != err && "" != err.Error() {
		log.Fatal(err.Error())
	}
	//是否显示SQL语句
	DB.ShowSQL(true)
	//数据库最大打开的连接数
	DB.SetMaxOpenConns(100)

	//自动User
	if err = DB.Sync2(new(model.UserInfo)); err != nil {
		log.Print(err.Error())
	}

	if err = DB.Sync2(new(model.FileData)); err != nil {
		log.Print(err.Error())
	}

	if err = DB.Sync2(new(model.UserFile)); err != nil {
		log.Print(err.Error())
	}

}

func DBConn() *xorm.Engine {
	return DB
}
