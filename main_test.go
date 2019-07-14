package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
)

type initConfigTestSuite struct {
	suite.Suite
}

func (suite initConfigTestSuite) SetupSuite() {
	home := filepath.Join(os.TempDir(), uuid.New().String())
	err := os.Mkdir(home, os.ModeDir|0775)
	if err != nil {
		panic(err)
	}

	err = os.Setenv("HOME", home)
	if err != nil {
		panic(err)
	}
}

func (suite initConfigTestSuite) TearDownSuite() {
	home := os.Getenv("HOME")
	if !strings.HasPrefix(home, os.TempDir()) {
		panic("Temp home is not correct, skip clean")
	}

	err := os.RemoveAll(os.Getenv(home))
	if err != nil {
		panic(err)
	}
}

func (suite initConfigTestSuite) TestInitConfigFromFlag() {
	// configPath is from main.go, which will be loaded later.
	configPath = filepath.Join(os.TempDir(), uuid.New().String()+".yaml")
	defer func() {
		configPath = ""
	}()

	config := []byte(`access_key_id: XXX
secret_access_key: YYY
`)
	f, err := os.Create(configPath)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer os.Remove(configPath)
	_, err = f.Write(config)
	if err != nil {
		suite.T().Fatal(err)
	}

	err = initConfig()
	suite.Nil(err)
	suite.NotNil(contexts.Storage)
	suite.Equal(configPath, viper.ConfigFileUsed())
	suite.Equal("XXX", viper.GetString(constants.ConfigAccessKeyID))
	suite.Equal("YYY", viper.GetString(constants.ConfigSecretAccessKey))
}

func (suite initConfigTestSuite) TestInitConfigFromEnv() {
	err := os.Setenv("QSCTL_ACCESS_KEY_ID", "XXX")
	if err != nil {
		suite.T().Fatal(err)
	}

	err = os.Setenv("QSCTL_SECRET_ACCESS_KEY", "YYY")
	if err != nil {
		suite.T().Fatal(err)
	}

	err = initConfig()
	suite.Nil(err)
	suite.NotNil(contexts.Storage)
	suite.Equal("", viper.ConfigFileUsed())
	suite.Equal("XXX", viper.GetString(constants.ConfigAccessKeyID))
	suite.Equal("YYY", viper.GetString(constants.ConfigSecretAccessKey))
}

func (suite initConfigTestSuite) TestInitConfigFromHome() {
	home := os.Getenv("HOME")

	filePath := filepath.Join(home, ".qingstor", "config.yaml")

	config := []byte(`access_key_id: XXX
secret_access_key: YYY
`)
	err := os.MkdirAll(filepath.Join(home, ".qingstor"), os.ModeDir|0775)
	if err != nil {
		suite.T().Fatal(err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer os.Remove(filePath)
	_, err = f.Write(config)
	if err != nil {
		suite.T().Fatal(err)
	}

	err = initConfig()
	suite.Nil(err)
	suite.NotNil(contexts.Storage)
	suite.Equal(filePath, viper.ConfigFileUsed())
	suite.Equal("XXX", viper.GetString(constants.ConfigAccessKeyID))
	suite.Equal("YYY", viper.GetString(constants.ConfigSecretAccessKey))
}
func (suite initConfigTestSuite) TestInitConfigFromHomeConfig() {
	home := os.Getenv("HOME")

	filePath := filepath.Join(home, ".config", "qingstor", "config.yaml")

	config := []byte(`access_key_id: XXX
secret_access_key: YYY
`)
	err := os.MkdirAll(filepath.Join(home, ".config", "qingstor"), os.ModeDir|0775)
	if err != nil {
		suite.T().Fatal(err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer os.Remove(filePath)
	_, err = f.Write(config)
	if err != nil {
		suite.T().Fatal(err)
	}

	err = initConfig()
	suite.Nil(err)
	suite.NotNil(contexts.Storage)
	suite.Equal(filePath, viper.ConfigFileUsed())
	suite.Equal("XXX", viper.GetString(constants.ConfigAccessKeyID))
	suite.Equal("YYY", viper.GetString(constants.ConfigSecretAccessKey))
}

func TestInitConfigTestSuite(t *testing.T) {
	suite.Run(t, new(initConfigTestSuite))
}
