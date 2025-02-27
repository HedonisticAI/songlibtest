package postgres

import (
	"database/sql"
	"fmt"
	"songlibtest/config"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func NewDB(Config config.Config) (*Postgres, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		Config.Database.DBHost, Config.Database.DBPort, Config.Database.DBUser, Config.Database.DBPwd, Config.Database.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}

func (P *Postgres) Get() *sql.DB {
	return P.DB
}

func (P *Postgres) Close() {
	P.DB.Close()
}
