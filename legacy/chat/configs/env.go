package configs

import (
	"log"

	"github.com/spf13/viper"
)

type env struct {
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBNAME     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`
}

var Env *env

func InitEnv() {
	Env = loadEnv()
}

func loadEnv() (env *env) {
	viper.AddConfigPath(".")
	viper.SetConfigFile("local.env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&env); err != nil {
		log.Fatal(err)
	}

	return
}
