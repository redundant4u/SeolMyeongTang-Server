package config

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLoadEnvAndRequiredFields(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile(".env", []byte("APP_ENV=development\nKUBE_SESSION_NAMESPACE=sessions\nCF_TUNNEL_ID=test\nAWS_ACCESS_KEY=test\nAWS_SECRET_KEY=test\nAWS_REGION=test-region\nDYNAMODB_TABLE=posts\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("AWS_REGION", "override")
	old := Env
	defer func() { Env = old }()
	InitEnv()
	if Env.APP_ENV != "development" || Env.AWS_REGION != "override" || Env.DYNAMODB_TABLE != "posts" {
		t.Fatalf("config loading/overrides changed: %+v", Env)
	}
	for _, field := range []string{"KUBE_SESSION_NAMESPACE", "CF_TUNNEL_ID", "AWS_ACCESS_KEY", "AWS_SECRET_KEY", "DYNAMODB_TABLE"} {
		copy := *Env
		reflect.ValueOf(&copy).Elem().FieldByName(field).SetString("")
		if err := validateEnv(&copy); err == nil || !strings.Contains(err.Error(), field) {
			t.Fatalf("missing %s accepted: %v", field, err)
		}
	}
	if err := validateEnv(Env); err != nil {
		t.Fatal(err)
	}
}
