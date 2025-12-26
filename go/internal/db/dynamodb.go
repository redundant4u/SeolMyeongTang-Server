package db

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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
		logger.Fatal(err, "failed to load aws config")
		return nil, err
	}

	client := dynamodb.NewFromConfig(awsCfg)

	err = pingddb(context.Background(), client)
	if err != nil {
		logger.Fatal(err, "failed to connect to dynamodb")
		return nil, err
	}

	logger.Info("connect to dynamodb")

	return client, nil
}

func pingddb(ctx context.Context, client *dynamodb.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.ListTables(ctx, &dynamodb.ListTablesInput{
		Limit: aws.Int32(1),
	})

	return err
}
