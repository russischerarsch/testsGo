package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"tgtest/conn"
	"tgtest/http"
	"tgtest/producerkafka"
	"tgtest/repo"
	"tgtest/serv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	connStr := os.Getenv("POSTGRES_CONN")
	conn, err := conn.CreateConnection(connStr, ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	m, err := migrate.New("file://migrations", connStr)
	if err != nil {
		log.Fatalf("failed to find migration file %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to set migrations %v", err)
	}
	m.Close()
	cfg := producerkafka.LoadKafkaConfig()
	writer := producerkafka.CreateWriter(cfg.Brokers, cfg.Topic)
	eventService := serv.CreateEventService(writer)
	err = eventService.PublishEvent(ctx, "user-event", &serv.UserCreatedEvent{
		ID:        "1",
		Name:      "Daria",
		Email:     "daria@example.com",
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Fatalf("failed to send event to kafka, %v", err)
	}
	defer writer.Close()
	repo := repo.CreateRepo(conn)
	serv := serv.CreateServ(repo)
	handler := http.CreateHandler(serv)
	router := gin.Default()
	router.POST("/user", handler.CreateUser)
	if err := router.Run(":8080"); err != nil {
		return
	}
}
