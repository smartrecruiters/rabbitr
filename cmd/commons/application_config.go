package commons

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"sync"

	"github.com/zalando/go-keyring"
)

// ServerCoordinates describes server configuration parameters
type ServerCoordinates struct {
	APIURL   string `json:"apiUrl"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Config holds configuration data uses by the application
type Config struct {
	Servers map[string]*ServerCoordinates `json:"servers"`
}

// GetServerNames returns slice of sorted server names taken from configuration
func (cfg *Config) GetServerNames() []string {
	i, serverNames := 0, make([]string, len(cfg.Servers))
	for k := range cfg.Servers {
		serverNames[i] = k
		i++
	}
	sort.Strings(serverNames)
	return serverNames
}

var cachedConfig *Config
var once sync.Once

// GetCachedApplicationConfig returns cached application config
func GetCachedApplicationConfig() Config {
	once.Do(func() {
		cfg, err := GetApplicationConfig()
		AbortIfError(err)
		cachedConfig = &cfg
	})
	return *cachedConfig
}

// GetApplicationConfig returns fresh application config read from a file
func GetApplicationConfig() (Config, error) {
	var cfg Config
	err := getApplicationConfig(&cfg, ApplicationName)
	if err != nil {
		return cfg, err
	}

	return cfg, err
}

// UpdateApplicationConfig writes application config to a file
func UpdateApplicationConfig(cfg Config) error {
	if IsOSX() {
		err := storePasswordsInKeyChain(&cfg)
		if err != nil {
			return err
		}
	}

	return updateApplicationConfig(cfg, ApplicationName)
}

func getApplicationConfig(configStructure *Config, applicationName string) error {
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

	err = json.Unmarshal(jsonBytes, configStructure)
	if err != nil {
		return err
	}

	// for MacOS try to obtain passwords from keychain
	if IsOSX() {
		return fillPasswordsFromKeyChain(configStructure)
	}

	return nil
}

func getAppConfigFilePath(applicationName string) (string, error) {
	appCfgDir, err := getApplicationConfigDir(applicationName)
	if err != nil {
		return "", err
	}

	cfgPath := filepath.Join(appCfgDir, fmt.Sprintf("%s-config.json", applicationName))
	return cfgPath, nil
}

func getApplicationConfigDir(applicationName string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, fmt.Sprintf(".%s", applicationName)), nil
}

func updateApplicationConfig(configStructure interface{}, applicationName string) error {
	appCfgDir, err := getApplicationConfigDir(applicationName)
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

	configJSON, err := json.MarshalIndent(configStructure, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(appCfgPath, configJSON, 0644)
}

func fillPasswordsFromKeyChain(configStructure *Config) error {
	for serverName, coordinates := range configStructure.Servers {
		srvPass, err := keyring.Get(ApplicationName, serverName)
		if err != nil && err != keyring.ErrNotFound {
			return err
		}
		if len(srvPass) > 0 && err != keyring.ErrNotFound {
			Debugf("Pass for %s obtained from keychain successfully", serverName)
			coordinates.Password = srvPass
		}
	}
	return nil
}

func storePasswordsInKeyChain(configStructure *Config) error {
	for serverName, coordinates := range configStructure.Servers {
		err := keyring.Set(ApplicationName, serverName, coordinates.Password)
		if err != nil {
			return err
		}
		Debugf("Password for %s stored in keychain successfully", serverName)
		coordinates.Password = ""
	}
	return nil
}
