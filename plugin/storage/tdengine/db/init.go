package db

import (
	"database/sql"
	"fmt"
	_ "github.com/wenj91/taos-driver"
)

type TaosConnConfig struct {
	User       string
	Password   string
	HostName   string
	ServerPort int
	DbName     string
}

func PrepareConnection(config *TaosConnConfig) (*sql.DB, error) {
	url := fmt.Sprintf("%s:%s@/http(%s:%d)/%s", config.User, config.Password, config.HostName, config.ServerPort, config.DbName)
	db, err := sql.Open("taosSql", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("create database if not exists " + config.DbName + " precision 'ns' update 2")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("use " + config.DbName)
	if err != nil {
		return nil, err
	}
	return db, nil

}
