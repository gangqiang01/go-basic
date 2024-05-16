package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

const (
	EnvironmentalConfigPath = "QENV_CONFIG_PATH"
)

/*
* Config for yaml.
 */
type Config struct {
	ConfigFile string //configure file name
	ConfigPath string //configure file path
	Config     *viper.Viper
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		klog.Errorf("err: %v", err)
		return "./"
	}

	return strings.Replace(dir, "\\", "/", -1)
}

func getConfigDirectory() string {
	//get env
	config := os.Getenv(EnvironmentalConfigPath)
	if config != "" {
		return config
	}

	return GetCurrentDirectory()
}

// New yaml configuration for app.
func NewYamlConfig(fileName string) *Config {
	config := viper.New()
	config.SetConfigType("yaml")
	name := strings.TrimSuffix(fileName, ".yaml")
	config.SetConfigName(name)
	//config.SetConfigFile(fileName)

	confLocation := getConfigDirectory() + "/conf"
	_, err := os.Stat(confLocation)
	if !os.IsExist(err) {
		os.MkdirAll(confLocation, os.ModePerm)
	}
	config.AddConfigPath(confLocation)

	err = config.ReadInConfig()
	if err != nil {
		klog.Errorf("err: %v", err)
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return &Config{
		ConfigFile: fileName,
		ConfigPath: confLocation,
		Config:     config,
	}
}

func (c *Config) GetString(key string) string {
	return c.Config.GetString(key)
}

func (c *Config) GetBool(key string) bool {
	return c.Config.GetBool(key)
}

func (c *Config) GetInt(key string) int {
	return c.Config.GetInt(key)
}

func (c *Config) GetInt64(key string) int64 {
	return c.Config.GetInt64(key)
}

func (c *Config) GetUint(key string) uint {
	return c.Config.GetUint(key)
}

func (c *Config) GetUint64(key string) uint64 {
	return c.Config.GetUint64(key)
}

func (c *Config) GetFloat64(key string) float64 {
	return c.Config.GetFloat64(key)
}

func (c *Config) GetIntSlice(key string) []int {
	return c.Config.GetIntSlice(key)
}

func (c *Config) GetStringSlice(key string) []string {
	return c.Config.GetStringSlice(key)
}

// set key-value.
func (c *Config) Set(key string, value interface{}) {
	c.Config.Set(key, value)
}

func (c *Config) SetString(key, value string) {
	c.Config.Set(key, value)
}

func (c *Config) SetBool(key string, value bool) {
	c.Config.Set(key, value)
}

func (c *Config) SetInt(key string, value int) {
	c.Config.Set(key, value)
}

func (c *Config) SetInt64(key string, value int64) {
	c.Config.Set(key, value)
}

func (c *Config) SetUint(key string, value uint) {
	c.Config.Set(key, value)
}

func (c *Config) SetUint64(key string, value uint64) {
	c.Config.Set(key, value)
}

func (c *Config) SetFloat64(key string, value float64) {
	c.Config.Set(key, value)
}

func (c *Config) SetIntSlice(key string, value []int) {
	c.Config.Set(key, value)
}

func (c *Config) SetStringSlice(key string, value []string) {
	c.Config.Set(key, value)
}

// save config to conf/xx.yaml
func (c *Config) SaveConfig() error {
	fileName := c.ConfigPath + "/" + c.ConfigFile
	return c.Config.WriteConfigAs(fileName)
}
