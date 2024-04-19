package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/lordofthemind/backendMasterGo/api"
	db "github.com/lordofthemind/backendMasterGo/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:backendMasterGoSecret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "127.0.0.1:9090"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
