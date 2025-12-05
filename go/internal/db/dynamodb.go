package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	appConfig "seolmyeong-tang-server/internal/config"
	"seolmyeong-tang-server/internal/pkg/logger"
)

func Initddb() (*dynamodb.Client, error) {
	awsCfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(appConfig.Env.AWS_REGION),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				appConfig.Env.AWS_ACCESS_KEY,
				appConfig.Env.AWS_SECRET_KEY,
				"",
			),
		),
	)
	if err != nil {
		logger.Fatal(err, "failed to connect dynamodb")
		return nil, err
	}

	client := dynamodb.NewFromConfig(awsCfg)

	logger.Info("connect to dynamodb")

	return client, nil
}
