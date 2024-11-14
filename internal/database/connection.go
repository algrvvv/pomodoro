package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/algrvvv/pomodoro/internal/config"
)

var C *sql.DB

func GetConnectionString() string {
	// host=%s port=%d user=%s password=%s dbname=%s sslmode=disable
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.Config.DB.User, config.Config.DB.Passwd,
		config.Config.DB.Addr, config.Config.DB.DBName,
	)
}

func Connect() error {
	var err error

	C, err = sql.Open("postgres", GetConnectionString())
	if err != nil {
		return err
	}

	return C.Ping()
}

func Close() error {
	return C.Close()
}
