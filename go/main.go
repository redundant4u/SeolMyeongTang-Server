package main

import (
	"seolmyeong-tang-server/internal/config"
	"seolmyeong-tang-server/internal/router"
)

func main() {
	config.InitEnv()

	e, err := router.New()
	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":8090"))
}
