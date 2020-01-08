package commons

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sync"
)

type ServerCoordinates struct {
	ApiURL   string `json:"apiUrl"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	Servers map[string]ServerCoordinates `json:"servers"`
}

var cachedConfig *Config
var once sync.Once

func GetCachedApplicationConfig() Config {
	once.Do(func() {
		cfg, err := GetApplicationConfig()
		AbortIfError(err)
		cachedConfig = &cfg
	})
	return *cachedConfig
}

func GetApplicationConfig() (Config, error) {
	var cfg Config
	err := getApplicationConfig(&cfg, ApplicationName)
	if err != nil {
		return cfg, err
	}

	return cfg, err
}

func UpdateApplicationConfig(cfg Config) error {
	return updateApplicationConfig(cfg, ApplicationName)
}

func getApplicationConfig(configStructure interface{}, applicationName string) error {
	cfgPath, err := getAppConfigFilePath(applicationName)
	if err != nil {
		return err
	}

	if _, err := os.Stat(cfgPath); err != nil {
		Debugf("Application config file [%s] does not exists yet", cfgPath)
		return nil
	}

	jsonBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, configStructure)
}

func getAppConfigFilePath(applicationName string) (string, error) {
	appCfgDir, err := GetApplicationConfigDir(applicationName)
	if err != nil {
		return "", err
	}

	cfgPath := filepath.Join(appCfgDir, fmt.Sprintf("%s-config.json", applicationName))
	return cfgPath, nil
}

func GetApplicationConfigDir(applicationName string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, fmt.Sprintf(".%s", applicationName)), nil
}

func updateApplicationConfig(configStructure interface{}, applicationName string) error {
	appCfgDir, err := GetApplicationConfigDir(applicationName)
	if err != nil {
		return err
	}

	err = MakeDir(appCfgDir)
	if err != nil {
		return err
	}

	appCfgPath, err := getAppConfigFilePath(applicationName)
	if err != nil {
		return err
	}

	configJson, err := json.MarshalIndent(configStructure, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(appCfgPath, configJson, 0644)
}
