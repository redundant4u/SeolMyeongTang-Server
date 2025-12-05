package main

import (
	"seolmyeong-tang-server/internal/config"
	"seolmyeong-tang-server/internal/db"
	"seolmyeong-tang-server/internal/router"
)

func main() {
	config.InitEnv()

	ddb, err := db.Initddb()
	if err != nil {
		panic(err)
	}

	e := router.New(ddb)

	e.Logger.Fatal(e.Start(":8090"))
}
