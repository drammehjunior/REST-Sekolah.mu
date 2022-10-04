package db

import (
	"database/sql"
	"exampleclean.com/refactor/app/config"
	"exampleclean.com/refactor/app/domain"
	"fmt"
	"github.com/go-gorp/gorp"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDatabase(cfg config.Config) (*gorp.DbMap, error) {
	//psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=''", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort)
	query := fmt.Sprintf("%s:@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, dbErr := sql.Open("mysql", query)
	if dbErr != nil {
		return nil, dbErr
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(domain.Users{}, "user").SetKeys(true, "Id").ColMap("Email").SetUnique(true)
	err := dbmap.CreateTablesIfNotExists()

	if err != nil {
		fmt.Println("failed to make new table")
	}
	return dbmap, nil
}
