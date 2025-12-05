package config

import (
	"fmt"
	"reflect"
	"seolmyeong-tang-server/internal/pkg/logger"

	"github.com/spf13/viper"
)

type env struct {
	APP_ENV string `mapstructure:"APP_ENV"`

	KUBE_CONFIG            string `mapstructure:"KUBE_CONFIG" required:"true"`
	KUBE_SESSION_NAMESPACE string `mapstructure:"KUBE_SESSION_NAMESPACE" required:"true"`

	AWS_ACCESS_KEY string `mapstructure:"AWS_ACCESS_KEY"`
	AWS_SECRET_KEY string `mapstructure:"AWS_SECRET_KEY"`
	AWS_REGION     string `mapstructure:"AWS_REGION"`

	DYNAMODB_TABLE string `mapstructure:"DYNAMODB_TABLE"`

	// DB_HOST     string `mapstructure:"DB_HOST"`
	// DB_USER     string `mapstructure:"DB_USER"`
	// DB_NAME     string `mapstructure:"DB_NAME"`
	// DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	// DB_PORT     int    `mapstructure:"DB_PORT"`
}

var Env *env

func InitEnv() {
	Env = loadEnv()

	if err := validateEnv(Env); err != nil {
		logger.Fatal(err, "environment validation failed")
	}
}

func loadEnv() (env *env) {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal(err, "Unable to read env file")
	}

	if err := viper.Unmarshal(&env); err != nil {
		logger.Fatal(err, "Unable to unmarshal env file")
	}

	return
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
