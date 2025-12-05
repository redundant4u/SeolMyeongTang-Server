package post

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type repository struct {
	db        *dynamodb.Client
	tableName string
}

func newRepository(db *dynamodb.Client, tableName string) *repository {
	return &repository{
		db:        db,
		tableName: tableName,
	}
}

func (r *repository) getPosts(ctx context.Context) ([]post, error) {
	q, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("#PK = :PK"),
		ExpressionAttributeNames: map[string]string{
			"#PK": "PK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":PK": &types.AttributeValueMemberS{Value: "post"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("dynamodb query failed %w", err)
	}

	var posts []post

	err = attributevalue.UnmarshalListOfMaps(q.Items, &posts)
	if err != nil {
		return nil, fmt.Errorf("unmarshall failed %w", err)
	}

	return posts, nil
}

func (r *repository) getPost(ctx context.Context, postId string) (post, error) {
	q, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("#PK = :PK AND #SK = :SK"),
		ExpressionAttributeNames: map[string]string{
			"#PK": "PK",
			"#SK": "SK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":PK": &types.AttributeValueMemberS{Value: "post"},
			":SK": &types.AttributeValueMemberS{Value: postId},
		},
	})
	if err != nil {
		return post{}, fmt.Errorf("dynamodb query failed: %w", err)
	}

	if len(q.Items) == 0 {
		return post{}, fmt.Errorf("post not found: %s", postId)
	}

	var p post
	if err := attributevalue.UnmarshalMap(q.Items[0], &p); err != nil {
		return post{}, fmt.Errorf("unmarshall failed: %w", err)
	}

	return p, nil
}
