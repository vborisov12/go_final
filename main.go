package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/vborisov12/go_final/pkg/db"
	"github.com/vborisov12/go_final/pkg/server"
)

const (
	defaultPort   = 7540
	defaultDBFile = "scheduler.db"
)

type TaskService struct {
	store db.TaskStore
}

func NewTaskService(store db.TaskStore) TaskService {
	return TaskService{store: store}
}

func main() {
	dbFile := getDBFile()

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	dataBase, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	defer dataBase.Close()

	// store := db.NewTaskStore(dataBase)
	// service := NewTaskService(store)

	port := getPort()

	if err := server.Start(port); err != nil {
		log.Fatal(err)
	}

}

func getPort() int {
	port, exists := os.LookupEnv("TODO_PORT")
	if !exists {
		return defaultPort
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Failed to parse TODO_PORT env var: %v", err)
	}
	return p
}

func getDBFile() string {
	dbFile, exists := os.LookupEnv("TODO_DBFILE")
	if !exists {
		return defaultDBFile
	}
	return dbFile
}
