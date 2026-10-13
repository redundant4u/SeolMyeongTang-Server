package db

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"seolmyeong-tang-server/internal/config"
)

func TestDatabaseInitializationAndPing(t *testing.T) {
	fail := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Amz-Target") != "DynamoDB_20120810.ListTables" {
			t.Error("unexpected database operation")
		}
		var input struct{ Limit int }
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			t.Error(err)
		}
		if input.Limit != 1 {
			t.Error("ping must not scan tables")
		}
		w.Header().Set("Content-Type", "application/x-amz-json-1.0")
		if fail {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(`{"__type":"ResourceNotFoundException","message":"missing"}`))
		} else {
			_, _ = w.Write([]byte(`{"TableNames":[]}`))
		}
	}))
	defer server.Close()
	t.Chdir(t.TempDir())
	if err := os.WriteFile(".env", []byte("KUBE_SESSION_NAMESPACE=sessions\nCF_TUNNEL_ID=test\nAWS_ACCESS_KEY=test\nAWS_SECRET_KEY=test\nAWS_REGION=us-east-1\nDYNAMODB_TABLE=posts\n"), 0600); err != nil {
		t.Fatal(err)
	}
	old := config.Env
	defer func() { config.Env = old }()
	config.InitEnv()
	t.Setenv("AWS_ENDPOINT_URL_DYNAMODB", server.URL)
	t.Setenv("AWS_CONFIG_FILE", t.TempDir()+"/missing")
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", t.TempDir()+"/missing")
	client, err := Initddb()
	if err != nil || client == nil {
		t.Fatalf("initialization failed: %v", err)
	}
	fail = true
	client = dynamodb.NewFromConfig(aws.Config{Region: "us-east-1", Credentials: credentials.NewStaticCredentialsProvider("test", "test", ""), RetryMaxAttempts: 1}, func(o *dynamodb.Options) { o.BaseEndpoint = aws.String(server.URL) })
	if err := pingddb(context.Background(), client); err == nil {
		t.Fatal("database ping failure swallowed")
	}
}
