package main

import (
	"log"
	"net/http"

	"github.com/Popolzen/go_final_project/internal/server"
	"github.com/Popolzen/go_final_project/internal/storage"
)

func main() {
	dsn := "host=localhost port=5432 user=postgres password=123456 dbname=gophkeeper sslmode=disable"

	repo, err := storage.NewPostgresRepository(dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer repo.Close()

	log.Println("Подключились к БД")

	jwtSecret := []byte("my-secret-key")
	encKey := make([]byte, 32)
	handler := server.NewHandler(repo, jwtSecret, encKey)

	addr := ":8080"
	log.Printf("Сервер запущен на http://localhost%s", addr)

	if err := http.ListenAndServe(addr, handler.Routes()); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
