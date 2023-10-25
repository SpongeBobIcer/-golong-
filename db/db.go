package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB // 全局数据库连接池
