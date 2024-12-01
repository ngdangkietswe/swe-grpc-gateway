package configs

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

var (
	GlobalConfig = &Configuration{}
)

func init() {
	env := os.Getenv("K8S_ENV")
	log.Printf("K8S_ENV is set to %s", env)
	if strings.ToLower(env) == "prod" {
		log.Println("Using production config")
		viper.AutomaticEnv()
	} else {
		log.Println("Using local config")
		viper.AddConfigPath("./configs")
		viper.SetConfigName("config")
		viper.SetConfigType("env")

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Can't read config file: %v", err)
			return
		}
	}

	err := viper.Unmarshal(&GlobalConfig)
	if err != nil {
		log.Fatalf("Can't unmarshal config: %v", err)
	}
}

type Configuration struct {
	Port            int           `mapstructure:"PORT"`
	GrpcTaskHost    string        `mapstructure:"GRPC_TASK_HOST"`
	GrpcTaskPort    int           `mapstructure:"GRPC_TASK_PORT"`
	GrpcAuthHost    string        `mapstructure:"GRPC_AUTH_HOST"`
	GrpcAuthPort    int           `mapstructure:"GRPC_AUTH_PORT"`
	GrpcStorageHost string        `mapstructure:"GRPC_STORAGE_HOST"`
	GrpcStoragePort int           `mapstructure:"GRPC_STORAGE_PORT"`
	JwtSecret       string        `mapstructure:"JWT_SECRET"`
	JwtIssuer       string        `mapstructure:"JWT_ISSUER"`
	JwtExp          time.Duration `mapstructure:"JWT_EXPIRATION"`
}
