package config_test

import (
	"os"
	"testing"

	"github.com/pocockn/shouts-api/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigCreation(t *testing.T) {

	expectedConfigStruct := config.Config{
		Database: config.Database{
			Host:         "127.0.0.1",
			DatabaseName: "docker_pocockn",
			Port:         "3306",
			Password:     "pocockn",
			Username:     "pocockn",
		},
		SNS: config.SNSConfig{
			Arn: "arn:aws:sns:eu-west-1:314572533291:image-similarity-analysis",
		},
	}

	err := os.Setenv("ENV", "test")
	assert.NoError(t, err)

	config := config.NewConfig()

	assert.Equal(t, expectedConfigStruct.Database.Host, config.Database.Host)
	assert.Equal(t, expectedConfigStruct.Database.DatabaseName, config.Database.DatabaseName)
	assert.Equal(t, expectedConfigStruct.Database.Port, config.Database.Port)
	assert.Equal(t, expectedConfigStruct.Database.Password, config.Database.Password)
	assert.Equal(t, expectedConfigStruct.Database.Username, config.Database.Username)
	assert.Equal(t, expectedConfigStruct.SNS.Arn, config.SNS.Arn)
}
