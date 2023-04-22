package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var region = os.Getenv("AWS_REGION")
var tableName = os.Getenv("DYNAMODB_TABLE")
var pk = os.Getenv("DYNAMODB_PK")

var db dynamodb.Client

type post struct {
	Title     string
	SK        string
	Content   string
	CreatedAt string
}

type postWithoutContent struct {
	Title     string
	SK        string
	CreatedAt string
}

func init() {
	log.Default().Println(region, tableName, pk)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	errCheck("Unable to load SDK config", err)

	db = *dynamodb.NewFromConfig(cfg)
}

func getPostsFromDB(ctx context.Context) ([]postWithoutContent, error) {
	res, err := db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("#PK = :PK"),
		ExpressionAttributeNames: map[string]string{
			"#PK": "PK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":PK": &types.AttributeValueMemberS{Value: pk},
		},
	})
	errCheck("Couldn't query", err)

	var posts []postWithoutContent
	err = attributevalue.UnmarshalListOfMaps(res.Items, &posts)
	errCheck("Couldn't unmarshal query response", err)

	return posts, nil
}

func getPostFromDB(ctx context.Context, postId string) ([]post, error) {
	res, err := db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("#PK = :PK AND #SK = :SK"),
		ExpressionAttributeNames: map[string]string{
			"#PK": "PK",
			"#SK": "SK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":PK": &types.AttributeValueMemberS{Value: pk},
			":SK": &types.AttributeValueMemberS{Value: postId},
		},
	})
	errCheck("Couldn't query", err)

	var posts []post
	err = attributevalue.UnmarshalListOfMaps(res.Items, &posts)
	errCheck("Couldn't unmarshal query response", err)

	return posts, nil
}
