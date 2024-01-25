package internal

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func Router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return get(ctx, req)
	default:
		return httpException(http.StatusMethodNotAllowed, "NOT_ALLOWED_METHOD")
	}
}

func get(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	postId, ok := req.PathParameters["postId"]

	if !ok {
		return getPosts(ctx)
	}

	return getPost(ctx, postId)
}

func getPosts(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	posts, err := getPostsFromDB(ctx)
	errCheck("Couldn't get posts from db", err)

	body, err := json.Marshal(posts)
	errCheck("Couldn't marshal", err)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS,GET",
			"Access-Control-Allow-Headers": "Content-Type",
			"Content-Type":                 "application/json",
		},
		Body: string(body),
	}, nil
}

func getPost(ctx context.Context, postId string) (events.APIGatewayProxyResponse, error) {
	posts, err := getPostFromDB(ctx, postId)
	errCheck("Couldn't get post from db", err)

	if len(posts) == 0 {
		return httpException(http.StatusNotFound, "NOT_EXISTING_POST")
	}

	body, err := json.Marshal(posts[0])
	errCheck("Couldn't marshal", err)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS,GET",
			"Access-Control-Allow-Headers": "Content-Type",
			"Content-Type":                 "application/json",
		},
		Body: string(body),
	}, nil
}
