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
	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(
		ctx,
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
		logger.FatalEvent(ctx, "aws_config_load_failed", "failed to load aws config", err)
		return nil, err
	}

	client := dynamodb.NewFromConfig(awsCfg)

	if err := pingddb(ctx, client); err != nil {
		logger.FatalEvent(ctx, "dynamodb_ping_failed", "failed to connect to dynamodb", err)
		return nil, err
	}

	logger.InfoEvent(ctx, "dynamodb_connected", "connected to dynamodb")

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
