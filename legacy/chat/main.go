package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redundant4u/SeolMyeongTang-Server/chat/configs"
	"github.com/redundant4u/SeolMyeongTang-Server/chat/models"
	"github.com/redundant4u/SeolMyeongTang-Server/chat/router"
)

func main() {
	configs.InitEnv()

	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3001"},
		AllowMethods: []string{"GET", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type"},
	}))

	models.DBConnect()
	router.InitRouter(engine)

	url := fmt.Sprint(configs.Env.Host, ":", configs.Env.Port)

	engine.Run(url)
}
