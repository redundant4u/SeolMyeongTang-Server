package post

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/labstack/echo/v4"
	"seolmyeong-tang-server/internal/config"
)

func TestPostQueriesAndHTTPContract(t *testing.T) {
	for _, mode := range []string{"success", "empty", "query error", "malformed item"} {
		t.Run(mode, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-Amz-Target") != "DynamoDB_20120810.Query" {
					t.Error("unexpected database operation")
				}
				var input struct {
					TableName                 string
					KeyConditionExpression    string
					ExpressionAttributeValues map[string]struct{ S string }
				}
				if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
					t.Error(err)
				}
				if input.TableName != "posts" || input.ExpressionAttributeValues[":PK"].S != "post" {
					t.Error("wrong partition/table")
				}
				if strings.Contains(input.KeyConditionExpression, "#SK") && input.ExpressionAttributeValues[":SK"].S != "post1" {
					t.Error("wrong post key")
				}
				w.Header().Set("Content-Type", "application/x-amz-json-1.0")
				switch mode {
				case "success":
					_, _ = w.Write([]byte(`{"Items":[{"Title":{"S":"hello"},"SK":{"S":"post1"},"Content":{"S":"body"},"CreatedAt":{"S":"2026-01-01"}}]}`))
				case "empty":
					_, _ = w.Write([]byte(`{"Items":[]}`))
				case "query error":
					w.WriteHeader(400)
					_, _ = w.Write([]byte(`{"__type":"ResourceNotFoundException","message":"missing table"}`))
				case "malformed item":
					_, _ = w.Write([]byte(`{"Items":[{"Title":{"M":{"nested":{"S":"invalid"}}}}]}`))
				}
			}))
			defer server.Close()
			ddb := dynamodb.NewFromConfig(aws.Config{Region: "test", Credentials: credentials.NewStaticCredentialsProvider("test", "test", ""), RetryMaxAttempts: 1}, func(o *dynamodb.Options) { o.BaseEndpoint = aws.String(server.URL) })
			repo := newRepository(ddb, "posts")
			posts, err := repo.getPosts(context.Background())
			p, singleErr := repo.getPost(context.Background(), "post1")
			if mode == "success" {
				if err != nil || singleErr != nil || len(posts) != 1 || p.Title != "hello" {
					t.Fatalf("query: %+v %v %v", posts, err, singleErr)
				}
			}
			if mode == "empty" {
				if err != nil || len(posts) != 0 || singleErr == nil {
					t.Fatal("empty query semantics changed")
				}
			}
			if mode == "query error" || mode == "malformed item" {
				if err == nil || singleErr == nil {
					t.Fatal("provider failure swallowed")
				}
			}
			t.Chdir(t.TempDir())
			if err := os.WriteFile(".env", []byte("APP_ENV=development\nKUBE_SESSION_NAMESPACE=sessions\nCF_TUNNEL_ID=test\nAWS_ACCESS_KEY=test\nAWS_SECRET_KEY=test\nDYNAMODB_TABLE=posts\n"), 0600); err != nil {
				t.Fatal(err)
			}
			old := config.Env
			defer func() { config.Env = old }()
			config.InitEnv()
			e := echo.New()
			if err := Init(e, ddb); err != nil {
				t.Fatal(err)
			}
			for _, path := range []string{"/post", "/post/post1"} {
				w := httptest.NewRecorder()
				e.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
				want := 200
				if mode == "query error" || mode == "malformed item" || (mode == "empty" && path != "/post") {
					want = 400
				}
				if w.Code != want {
					t.Fatalf("%s: %d %s", path, w.Code, w.Body.String())
				}
				if mode == "success" && !strings.Contains(w.Body.String(), `"Title":"hello"`) {
					t.Fatal("wire DTO changed")
				}
			}
			c := e.NewContext(httptest.NewRequest("GET", "/post", nil), httptest.NewRecorder())
			if err := newHandler(repo).getPost(c); err != nil || c.Response().Status != 400 {
				t.Fatal("missing post ID accepted")
			}
		})
	}
}

func TestProductionPostsAreForbidden(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile(".env", []byte("APP_ENV=production\nKUBE_SESSION_NAMESPACE=sessions\nCF_TUNNEL_ID=test\nAWS_ACCESS_KEY=test\nAWS_SECRET_KEY=test\nDYNAMODB_TABLE=posts\n"), 0600); err != nil {
		t.Fatal(err)
	}
	old := config.Env
	defer func() { config.Env = old }()
	config.InitEnv()
	e := echo.New()
	_ = Init(e, nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest("GET", "/post", nil))
	if w.Code != 403 {
		t.Fatal("production post access changed")
	}
}
