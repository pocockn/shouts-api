package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pocockn/awswrappers/s3"
)

type (
	// Config contains config for the application.
	Config struct {
		Database Database
		S3       S3Config
		SNS      SNSConfig
	}

	// Database holds database values in our config.
	Database struct {
		Host           string
		DatabaseName   string
		Port           string
		Password       string
		MaxConnections int
		Username       string
		URL            string
	}

	// S3Config holds information to connect to S3.
	S3Config struct {
		Client *s3.Client
		Domain string
		Bucket string
	}

	// SNSConfig holds the topic ARN for the image similarity SNS topic.
	SNSConfig struct {
		Arn string
	}
)

// NewConfig creates a new config struct.
func NewConfig() Config {
	var config Config
	if _, err := toml.DecodeFile(config.generatePath(), &config); err != nil {
		fmt.Println(err)
	}

	s3ClientConfig := &s3.ClientConfig{
		Endpoint: config.S3.Domain,
	}
	config.S3.Client = s3.NewClient(s3ClientConfig, config.isDev(), nil)

	return config
}

func (c Config) environment() string {
	environment := "development"

	if os.Getenv("ENV") != "" {
		environment = os.Getenv("ENV")
	}

	return environment
}

func (c Config) generatePath() string {
	if os.Getenv("ENV") == "test" {
		return "development.toml"
	}

	return fmt.Sprintf("config/%s.toml", c.environment())
}

func (c Config) isDev() bool {
	return os.Getenv("ENV") == "development"
}
