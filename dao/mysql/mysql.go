package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"web_app/config"
)

var Db *sqlx.DB

func Init(app *config.AppConfig) (err error) {
	//设置连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		app.Mysql.Username,
		app.Mysql.Password,
		app.Mysql.Host,
		app.Mysql.Port,
		app.Mysql.Dbname,
		app.Mysql.Charset,
	)
	Db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	Db.SetMaxOpenConns(app.Mysql.MaxOpenConns)
	Db.SetMaxIdleConns(app.Mysql.MaxIdleConns)
	return
}
func Close() {
	_ = Db.Close()
}
