package main

import (
	"log"
	"os"
	"strconv"

	"github.com/vborisov12/go_final/pkg/api"
	"github.com/vborisov12/go_final/pkg/server"
	"github.com/vborisov12/go_final/pkg/storage"
)

// Стандартные значения для порта и файла БД
// Задаются из TODO_PORT и TODO_DBFILE для кастомных значений
const (
	defaultPort   = 7540
	defaultDBFile = "scheduler.db"
)

func main() {
	dbFile := getDBFile()

	db, err := storage.Init(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	taskStorage := storage.NewTaskStorage(db)
	taskService := api.NewTaskService(taskStorage)
	taskApi := api.NewApi(taskService)
	taskServer := server.NewTaskServer(getPort(), taskApi)

	if err := taskServer.Start(); err != nil {
		log.Fatal(err)
	}

}

// Функции для получения значений из переменных окружения
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
