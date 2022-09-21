package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"go-entry-task/tcpserver/config"
)

var db *sql.DB

func init() {
	conf := config.DefaultDBConfig()
	db, _ = sql.Open("mysql", conf.Host)
	//设置数据库最大连接数
	db.SetMaxOpenConns(conf.MaxConn)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(conf.IdleConn)
	//db.SetConnMaxIdleTime(time.Duration(conf.MysqlIdleTimeout) * time.Second)

	//验证连接
	if err := db.Ping(); err != nil {
		panic(err.Error())
	}
}

func Query(user string) (string, string, string, error) {
	var nick, pass string
	var url sql.NullString
	err := db.QueryRow("SELECT nickname,password,picture FROM userinfo WHERE  username = ?", user).Scan(&nick, &pass, &url)
	if err != nil {
		return "", "", "", err
	}
	return nick, pass, url.String, nil
}

func UpdateNic(nick, user string) error {
	_, err := db.Exec("UPDATE userinfo SET nickname=? where username=?", nick, user)
	return err
}

func UpdatePic(pic, user string) error {
	_, err := db.Exec("UPDATE userinfo SET picture=? where username=?", pic, user)
	return err
}
