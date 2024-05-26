package adapters

import (
	"backend/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConnect struct {
	Connection *sql.DB
	Error      error
}

func Connect() (*sql.DB, error) {
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, err
	}

	Dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.Database.User, config.Database.Pwd, config.Database.Host, config.Database.Port, config.Database.Dbname)
	db, err := sql.Open("mysql", Dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewSql() (Myc *MysqlConnect, err error) {
	con, err := Connect()
	if err != nil {
		return nil, err
	}

	return &MysqlConnect{
		Connection: con,
	}, nil
}
