package main

import (
	"log"
	"os"
	"strconv"

	"github.com/kigorek/go_final_project/pkg/db"
	"github.com/kigorek/go_final_project/pkg/server"
	"github.com/kigorek/go_final_project/tests"
)

func main() {
	var dbFile string
	if os.Getenv("TODO_DBFILE") != "" {
		dbFile = os.Getenv("TODO_DBFILE")
	} else {
		dbFile = "scheduler.db"
	}

	err := db.Init(dbFile)

	if err != nil {
		log.Fatalf("возникла проблема при инициализации БД: %v", err)

	}

	var port string
	if os.Getenv("TODO_PORT") != "" {
		port = os.Getenv("TODO_PORT")
	} else {
		port = ":" + strconv.Itoa(tests.Port)
	}

	server.Run(port)

}
