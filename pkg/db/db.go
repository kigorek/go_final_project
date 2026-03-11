package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `CREATE TABLE scheduler (
   		 		id INTEGER PRIMARY KEY AUTOINCREMENT,
    			date CHAR(8) NOT NULL DEFAULT '',
				title VARCHAR(256) NOT NULL DEFAULT '',
				comment TEXT ,
				repeat VARCHAR(128) DEFAULT '' );

			CREATE INDEX idx_scheduler_date ON scheduler(date); `

// если install равен true, после открытия БД требуется выполнить
// sql-запрос с CREATE TABLE и CREATE INDEX
func Init(dbFile string) error {
	var install bool

	_, err := os.Stat(dbFile)
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("проблема при открытии БД: %w", err)
	}

	if install {
		_, err := db.Exec(schema)
		if err != nil {
			return fmt.Errorf("при выполнении запроса произошла ошибка: %w", err)
		} else {
			fmt.Println("схема БД создана")
		}
	}

	return nil

}
