package database

import (
	"b2match_api/utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() error {
	var err error
	DB, err = sql.Open("mysql", buildConnectionString())
	return err
}

func buildConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(mysql:3306)/%s?parseTime=true", utils.GetEnv("DB_USER", "root"), utils.GetEnv("DB_PASSWORD", ""), utils.GetEnv("DB_NAME", "api"))
}
