package config

import (
	"fmt"
	"reflect"
	"seolmyeong-tang-server/internal/pkg/logger"

	"github.com/spf13/viper"
)

type env struct {
	APP_ENV string `mapstructure:"APP_ENV"`

	KUBE_CONFIG            string `mapstructure:"KUBE_CONFIG"`
	KUBE_SESSION_NAMESPACE string `mapstructure:"KUBE_SESSION_NAMESPACE" required:"true"`

	AWS_ACCESS_KEY string `mapstructure:"AWS_ACCESS_KEY" required:"true"`
	AWS_SECRET_KEY string `mapstructure:"AWS_SECRET_KEY" required:"true"`
	AWS_REGION     string `mapstructure:"AWS_REGION"`

	DYNAMODB_TABLE string `mapstructure:"DYNAMODB_TABLE" required:"true"`
}

var Env *env

func InitEnv() {
	Env = loadEnv()

	if err := validateEnv(Env); err != nil {
		logger.Fatal(err, "environment validation failed")
	}
}

func loadEnv() *env {
	v := viper.New() // 전역 viper보다는 지역 인스턴스 사용 권장
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		logger.Fatal(err, "Unable to read env file")
	}

	v.AutomaticEnv()

	var e env
	if err := v.Unmarshal(&e); err != nil {
		logger.Fatal(err, "Unable to unmarshal env file")
	}

	return &e
}

func validateEnv(e *env) error {
	val := reflect.ValueOf(*e)
	typ := reflect.TypeOf(*e)

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		if fieldType.Tag.Get("required") == "true" {
			if fieldVal.Kind() == reflect.String && fieldVal.String() == "" {
				return fmt.Errorf("required environment variable %s is missing", fieldType.Tag.Get("mapstructure"))
			}
		}
	}

	return nil
}
