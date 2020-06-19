package config

import (
	"errors"
	"flag"
	"fmt"
	"genids/utils/logs"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

const SERVER_NAME = "genids"

var CONFIG_KEY = fmt.Sprintf("/configs/dm.sp.%s/system", SERVER_NAME)

var Config *config // 静态配置
var ConsulAddr string
var _path string

// 静态配置，程序启动后无法再做更改的参数配置
type config struct {
	BaseConf BaseConf `yaml:"base"`
}
type BaseConf struct {
	// 当前服务监听的端口
	ServerPort int `yaml:"server_port"`
	// 运行日志相关配置
	SystemLogPath string `yaml:"error_log_path"`
	LogMaxAge     int    `yaml:"log_max_age"`
	LogLevel      string `yaml:"log_level"`
	// Sentry相关配置
	SentryDSN string `yaml:"sentry_dsn"`
	// 是否预先生成ID
	PreGen bool `yaml:"pre_gen"`
	// 节点ID
	NodeID int64 `yaml:"node_id"` // Node ID,只能为 0， 1， 2, 3
}

// 初始化解析参数
func init() {
	flag.StringVar(&_path, "c", SERVER_NAME+".yml", "default config path")
	flag.StringVar(&ConsulAddr, "r", os.Getenv("CONSUL"), "default consul address")
}

// 优先从consul中加载配置，没有则从配置文件中加载配置
// consul中的配置文件需为yaml格式
func InitConfig() error {
	var err error
	var content []byte

	if ConsulAddr != "" {
		content, err = fetchConfig(CONFIG_KEY, watchSystemConfig)
	} else {
		content, err = ioutil.ReadFile(_path)
	}

	if err != nil {
		return err
	}

	if len(content) == 0 {
		return errors.New("not found nothing system config")
	}

	Config = &config{}
	if err := yaml.Unmarshal(content, Config); err != nil {
		return err
	}

	level, err := logrus.ParseLevel(Config.BaseConf.LogLevel)
	if err == nil {
		logs.LogSystem.SetLevel(level)
	}

	fmt.Printf("static system config => [%#v]\n", Config)

	return nil
}

// 从consul中获取配置信息
func fetchConfig(configKey string, watchFn func([]byte)) ([]byte, error) {
	config := consulapi.DefaultConfig()
	config.Address = ConsulAddr
	_client, err := consulapi.NewClient(config)
	if err != nil {
		logs.LogSystem.Error("system config consul client error : ", err)
	}
	data, meta, err := _client.KV().Get(configKey, nil)

	if watchFn != nil {
		go func() {
			for {
				lastIndex := meta.LastIndex
				options := &consulapi.QueryOptions{WaitIndex: meta.LastIndex, WaitTime: time.Minute * 5}
				data, meta, err = _client.KV().Get(configKey, options)
				if err == nil {
					if lastIndex != meta.LastIndex {
						lastIndex = meta.LastIndex
						watchFn(data.Value)
					}
				} else {
					for {
						_client, err = consulapi.NewClient(config)
						if err != nil {
							logs.LogSystem.Error("system config consul client error : ", err)
							time.Sleep(time.Second * 10)
						} else {
							data, meta, err = _client.KV().Get(configKey, nil)
							if err == nil {
								break
							}
						}
					}
				}
			}
		}()
	}
	return data.Value, err
}

func watchSystemConfig(val []byte) {
	newConfig := &config{}
	if err := yaml.Unmarshal(val, newConfig); err != nil {
		logs.LogSystem.Error("watchSystemConfig error:", err, " val: ", string(val))
		return
	}

	Config = newConfig
	// 更新运行日志等级
	level, err := logrus.ParseLevel(Config.BaseConf.LogLevel)
	if err == nil {
		logs.LogSystem.SetLevel(level)
	}

	fmt.Printf("Latest System config => [%#v]\n", Config)
}
