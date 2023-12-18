package custconfigs

import (
	"context"
	"encoding/json"
	"github.com/bytedance/sonic"
	"gopkg.in/yaml.v3"
	custerror "labs/service-mesh/helper/error"
	"log"
	"os"
)

type ServiceConfigs interface {
	GetPublic() *HttpConfigs
	GetPrivate() *HttpConfigs
	GetLogger() *LoggerConfigs
	GetEventStore() *EventStoreConfigs
	GetMqttStore() *EventStoreConfigs
	GetSqlite() *DatabaseConfigs

	String() string
}

type Configs struct {
	Public     HttpConfigs       `json:"public,omitempty" yaml:"public,omitempty"`
	Private    HttpConfigs       `json:"private,omitempty" yaml:"private,omitempty"`
	Logger     LoggerConfigs     `json:"logger,omitempty" yaml:"logger,omitempty"`
	EventStore EventStoreConfigs `json:"eventStore,omitempty" yaml:"eventStore,omitempty"`
	MqttStore  EventStoreConfigs `json:"mqttStore,omitempty" yaml:"mqttStore,omitempty"`
	Sqlite     DatabaseConfigs   `json:"sqlite,omitempty" yaml:"sqlite,omitempty"`
}

func (c Configs) String() string {
	configBytes, _ := sonic.Marshal(c)
	return string(configBytes)
}

func (c Configs) GetPublic() *HttpConfigs {
	return &c.Public
}

func (c Configs) GetPrivate() *HttpConfigs {
	return &c.Private
}

func (c Configs) GetLogger() *LoggerConfigs {
	return &c.Logger
}

func (c Configs) GetEventStore() *EventStoreConfigs {
	return &c.EventStore
}

func (c Configs) GetMqttStore() *EventStoreConfigs {
	return &c.MqttStore
}

func (c Configs) GetSqlite() *DatabaseConfigs {
	return &c.Sqlite
}

func Init(ctx context.Context, destConfig ServiceConfigs) {
	err := readConfig(destConfig)
	if err != nil {
		log.Fatal(err)
		return
	}
}

type HttpConfigs struct {
	Name string           `json:"name,omitempty" yaml:"name,omitempty"`
	Port int              `json:"port,omitempty" yaml:"port,omitempty"`
	Tls  TlsConfig        `json:"tls,omitempty" yaml:"tls,omitempty"`
	Auth BasicAuthConfigs `json:"auth,omitempty" yaml:"auth,omitempty"`
}

type TlsConfig struct {
	Cert      string `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key       string `json:"key,omitempty" yaml:"key,omitempty"`
	Authority string `json:"authority,omitempty" yaml:"authority,omitempty"`
}

func (c TlsConfig) Enabled() bool {
	if len(c.Cert) > 0 && len(c.Key) > 0 {
		return true
	}
	return false
}

type LoggerConfigs struct {
	Level    string `json:"level,omitempty" yaml:"level,omitempty"`
	Encoding string `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

type BasicAuthConfigs struct {
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Token    string `json:"token,omitempty" yaml:"token,omitempty"`
}

type EventStoreConfigs struct {
	Tls      TlsConfig `json:"tls,omitempty" yaml:"tls,omitempty"`
	Host     string    `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int       `json:"port,omitempty" yaml:"port,omitempty"`
	Name     string    `json:"name,omitempty" yaml:"name,omitempty"`
	Enabled  bool      `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Username string    `json:"username,omitempty" yaml:"username,omitempty"`
	Password string    `json:"password,omitempty" yaml:"password,omitempty"`
	Level    string    `json:"level,omitempty" yaml:"level,omitempty"`
}

type DatabaseConfigs struct {
	Connection string `json:"connection,omitempty" yaml:"connection,omitempty"`
}

func (c *EventStoreConfigs) HasAuth() bool {
	return len(c.Username) > 0 && len(c.Password) > 0
}

func readConfig(dest ServiceConfigs) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	configFile, err := readConfigFile(path)
	if err != nil {
		return err
	}

	err = parseConfig(configFile, dest)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	path := os.Getenv(ENV_CONFIG_FILE_PATH)
	if len(path) == 0 {
		return "", custerror.FormatNotFound("ENV_CONFIG_FILE_PATH not found, unable to read configurations")
	}
	return path, nil
}

func readConfigFile(path string) ([]byte, error) {
	fs, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, custerror.FormatNotFound("readConfigFile: file not found")
		}
		return nil, custerror.FormatInternalError("readConfigFile: err = %s", err)
	}

	contents, err := os.ReadFile(fs.Name())
	if err != nil {
		return nil, custerror.FormatInternalError("readConfigFile: err = %s", err)
	}

	return contents, nil
}

func parseConfig(contents []byte, configs ServiceConfigs) error {
	if jsonErr := json.Unmarshal(contents, configs); jsonErr != nil {
		if yamlErr := yaml.Unmarshal(contents, configs); yamlErr != nil {
			return custerror.FormatInvalidArgument("parseConfig: config parse JSON err = %s YAML err = %s", jsonErr, yamlErr)
		}
	}
	return nil
}
