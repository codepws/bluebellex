package sqldb

import (
	"bluebell_backend/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var login_db *sqlx.DB

func Init(cfg *settings.DBs) (err error) {

	fmt.Println("初始化数据库")

	// "user:password@tcp(host:port)/dbname"
	//登录数据库
	login_dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.LoginDB.User, cfg.LoginDB.Password, cfg.LoginDB.Host, cfg.LoginDB.Port, cfg.LoginDB.Database)
	//商品数据库
	//shop_dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
	//	cfg.LoginDB.User, cfg.LoginDB.Password, cfg.LoginDB.Host, cfg.ShopDB.Port, cfg.LoginDB.Database)
	//其他数据库

	login_db, err = sqlx.Connect(cfg.LoginDB.Type, login_dsn)
	if err != nil {
		return
	}
	login_db.SetMaxOpenConns(cfg.LoginDB.MaxOpenConns)
	login_db.SetMaxIdleConns(cfg.LoginDB.MaxIdleConns)

	return nil
}

// Close 关闭MySQL连接
func Close() {
	fmt.Println("关闭数据库")
	_ = login_db.Close()
}
