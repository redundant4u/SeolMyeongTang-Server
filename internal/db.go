package internal

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var env = os.Getenv("ENV")
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
	var cfg aws.Config
	var err error

	if env == "local" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID {
				return aws.Endpoint{
					URL: "http://localhost:8000",
				}, nil
			}
			// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

		credentials := credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "local",
				SecretAccessKey: "local",
				SessionToken:    "local",
			},
		}

		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(customResolver),
			config.WithCredentialsProvider(credentials),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	}
	errCheck("Unable to load Dynamodb config", err)

	db = *dynamodb.NewFromConfig(cfg)
}

func getPostsFromDB(ctx context.Context) ([]postWithoutContent, error) {
	log.Default().Println(tableName)

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
