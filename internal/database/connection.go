package connection

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "financial_control"
	cfg.ParseTime = true
	cfg.Loc = time.Local

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
