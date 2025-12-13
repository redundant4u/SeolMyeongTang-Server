package main

import (
	"context"
	"seolmyeong-tang-server/internal/api/session"
	"seolmyeong-tang-server/internal/config"
	"seolmyeong-tang-server/internal/db"
	"seolmyeong-tang-server/internal/pkg/k8s"
	"seolmyeong-tang-server/internal/router"
)

func main() {
	config.InitEnv()

	ddb, err := db.Initddb()
	if err != nil {
		panic(err)
	}

	kubeClient, err := k8s.NewClient(config.Env.KUBE_CONFIG)
	if err != nil {
		panic(err)
	}
	kube := session.NewKube(kubeClient, config.Env.KUBE_SESSION_NAMESPACE)

	ctx := context.Background()
	go kube.Gc.Run(ctx)

	e := router.New(ddb, kube)

	e.Logger.Fatal(e.Start(":8090"))
}
