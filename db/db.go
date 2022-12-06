package db

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

const defConfPath = "./config/database.json"

type config struct {
	Host     string `json:"host"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Port     int64  `json:"port"`
	User     string `json:"user"`
}

func Connect(path string) (*sqlx.DB, error) {
	if path == "" {
		path = defConfPath
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	conf := config{}
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		return nil, err
	}

	// dsn := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	conf.Host, conf.Port, conf.User, conf.Password, conf.Name)
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Name)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Connecting to data store %v\n", conf.Host)

	return db, err
}
