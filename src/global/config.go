package global

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config 全局配置
var Config struct {
	Service struct {
		IP              string `yaml:"ip"`
		Port            int    `yaml:"port"`
		QuitWaitTimeout int    `yaml:"quitWaitTimeout"`
		Limiter         int    `yaml:"limiter"`
		Debug           bool   `yaml:"debug"`
	} `yaml:"service"`
	Logger struct {
		Level        string `yaml:"level"`
		Outputs      string `yaml:"outputs"`
		Encode       string `yaml:"encode"`
		ColorLevel   bool   `yaml:"colorLevel"`
		EnableTrace  bool   `yaml:"enableTrace"`
		EnableCaller bool   `yaml:"enableCaller"`
	}
	Snowflake struct {
		Epoch int64 `yaml:"epoch"`
		Node  int64 `yaml:"node"`
	} `yaml:"snowflake"`
	Database struct {
		Addr      string `yaml:"addr"`
		User      string `yaml:"user"`
		Password  string `yaml:"password"`
		Name      string `yaml:"name"`
		EnableLog bool   `yaml:"enableLog"`
		Timeout   int    `yaml:"timeout"`
	} `yaml:"database"`
	Session struct {
		Key            string `yaml:"key"`
		CookieName     string `yaml:"cookieName"`
		HTTPOnly       bool   `yaml:"httpOnly"`
		Secure         bool   `yaml:"secure"`
		MaxAge         int    `yaml:"maxAge"`
		IdleTime       int    `yaml:"idleTime"`
		RedisAddr      string `yaml:"redisAddr"`
		RedisDB        int    `yaml:"redisDb"`
		RedisKeyPrefix string `yaml:"redisKeyPrefix"`
	} `yaml:"session"`
	Token struct {
		Key    string `yaml:"key"`
		Expire int    `yaml:"expire"`
	} `yaml:"token"`
}

// 解析配置文件路径
func parseFilePath(configFilePath string) (configName, configType string) {
	if configFilePath == "" {
		return
	}
	fullExt := filepath.Ext(configFilePath)
	if fullExt == "" {
		return
	}
	extArr := strings.Split(fullExt, ".")
	if len(extArr) > 1 {
		configType = extArr[1]
	}
	extPos := strings.LastIndex(configFilePath, fullExt)
	configName = configFilePath[:extPos]
	return
}

// SetConfig 设置全局配置变量
func SetConfig() error {
	configFilePath := flag.String("c", "./config.toml", "配置文件路径")
	flag.Parse()

	configName, configType := parseFilePath(*configFilePath)
	if configName == "" || configType == "" {
		return errors.New("无法读取配置文件")
	}
	viper.SetConfigName(configName) // 配置文件名，不需要扩展名
	viper.SetConfigType(configType) // 文件类型
	viper.AddConfigPath(".")        // 文件路径

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return viper.Unmarshal(&Config)
}

func SetConfigYaml() error {
	configFilePath := flag.String("c", "./config.yml", "配置文件路径")
	flag.Parse()

	cfgFile, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		return fmt.Errorf("读取配置文件错误：%w", err)
	}

	err = yaml.Unmarshal(cfgFile, &Config)
	if err != nil {
		return fmt.Errorf("读取配置文件错误：%w", err)
	}

	return nil
}
