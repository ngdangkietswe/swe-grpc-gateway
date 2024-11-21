package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

var (
	GlobalConfig = &Configuration{}
)

func init() {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Can't read config file: %v", err)
		return
	}

	config := &Configuration{}
	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("Can't unmarshal config: %v", err)
		return
	}

	GlobalConfig = config

	fmt.Println("\033[34m================ Loaded Configuration ================\033[0m")
	fmt.Printf("\033[32mPort:           \033[0m%d\n", GlobalConfig.Port)
	fmt.Printf("\033[32mGRPC Task Host:     \033[0m%s\n", GlobalConfig.GrpcTaskHost)
	fmt.Printf("\033[32mGRPC Task Port:     \033[0m%d\n", GlobalConfig.GrpcTaskPort)
	fmt.Println("\033[34m=====================================================\033[0m")
	return
}

type Configuration struct {
	Port         int    `mapstructure:"PORT"`
	GrpcTaskHost string `mapstructure:"GRPC_TASK_HOST"`
	GrpcTaskPort int    `mapstructure:"GRPC_TASK_PORT"`
}
