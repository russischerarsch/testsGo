package main

import (
	"context"
	"fmt"
	"tgtest/conn"
	"tgtest/http"
	"tgtest/repo"
	"tgtest/serv"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	conn, err := conn.CreateConnection(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	repo := repo.CreateRepo(conn)
	serv := serv.CreateServ(repo)
	handler := http.CreateHandler(serv)
	router := gin.Default()
	router.POST("/user", handler.CreateUser)
}
