package config

import (
	"fmt"
	"stylight/internal/services"
	"os"

	"github.com/kelseyhightower/envconfig"

	"github.com/joho/godotenv"
)

//WebServerConfig is used to store the environment and port for the webserver as well as the logging config.
type WebServerConfig struct {
	Port        string `required:"true" split_words:"true" default:"50051"`
	Environment string `required:"true" split_words:"true" default:"dev"`
	Service     *services.ServiceConfig
}

// FromEnv pulls config from the environment
func FromEnv() (cfg *WebServerConfig, err error) {
	fromFileToEnv()

	cfg = &WebServerConfig{}

	err = envconfig.Process("", cfg)
	if err != nil {

		return nil, err
	}

	return cfg, nil
}

func fromFileToEnv() {
	cfgFilename := os.Getenv("ENV_FILE")
	if cfgFilename != "" {
		err := godotenv.Load(cfgFilename)
		if err != nil {
			fmt.Println("ENV_FILE not found. Trying MY_POD_NAMESPACE")
		}
		return
	}

	loc := os.Getenv("ENVIRONMENT_NAME")
	fmt.Printf(loc)
	if loc == "" {
		fmt.Println("ENV_FILE defaulting to dev config")
		loc = "dev"
	}
	cfgFilename = fmt.Sprintf("../../etc/config/config.%s.env", loc)
	fmt.Printf(cfgFilename)
	err := godotenv.Load(cfgFilename)
	if err != nil {
		fmt.Printf("No config files found to load to env. Defaulting to environment. With error: %s", err.Error())
	}

}