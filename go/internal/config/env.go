package config

import (
	"context"
	"fmt"
	"reflect"
	"seolmyeong-tang-server/internal/pkg/logger"

	"github.com/spf13/viper"
)

type env struct {
	APP_ENV string `mapstructure:"APP_ENV"`

	KUBE_CONFIG            string `mapstructure:"KUBE_CONFIG"`
	KUBE_SESSION_NAMESPACE string `mapstructure:"KUBE_SESSION_NAMESPACE" required:"true"`

	CF_TUNNEL_ID string `mapstructure:"CF_TUNNEL_ID" required:"true"`

	AWS_ACCESS_KEY string `mapstructure:"AWS_ACCESS_KEY" required:"true"`
	AWS_SECRET_KEY string `mapstructure:"AWS_SECRET_KEY" required:"true"`
	AWS_REGION     string `mapstructure:"AWS_REGION"`

	DYNAMODB_TABLE string `mapstructure:"DYNAMODB_TABLE" required:"true"`
}

var Env *env

func InitEnv() {
	Env = loadEnv()

	if err := validateEnv(Env); err != nil {
		logger.FatalEvent(context.Background(), "env_validation_failed", "environment validation failed", err)
	}
}

func loadEnv() *env {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		logger.FatalEvent(context.Background(), "env_read_failed", "unable to read env file", err)
	}

	v.AutomaticEnv()

	var e env
	if err := v.Unmarshal(&e); err != nil {
		logger.FatalEvent(context.Background(), "env_unmarshal_failed", "unable to unmarshal env file", err)
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
