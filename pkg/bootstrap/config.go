package bootstrap

import (
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	// file
	fileKratos "github.com/go-kratos/kratos/v2/config/file"

	"kratos-gorm-example/api/gen/go/common/conf"
)

const remoteConfigSourceConfigFile = "remote.yaml"

// NewConfigProvider 创建一个配置
func NewConfigProvider(configPath string) config.Config {
	err, rc := LoadRemoteConfigSourceConfigs(configPath)
	if err != nil {
		log.Error("LoadRemoteConfigSourceConfigs: ", err.Error())
	}
	if rc != nil {
		return config.New(
			config.WithSource(
				NewFileConfigSource(configPath),
				NewRemoteConfigSource(rc),
			),
		)
	} else {
		return config.New(
			config.WithSource(
				NewFileConfigSource(configPath),
			),
		)
	}
}

// LoadBootstrapConfig 加载程序引导配置
func LoadBootstrapConfig(configPath string) *conf.Bootstrap {
	cfg := NewConfigProvider(configPath)
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := cfg.Scan(&bc); err != nil {
		panic(err)
	}

	if bc.Server == nil {
		bc.Server = &conf.Server{}
		_ = cfg.Scan(&bc.Server)
	}

	if bc.Client == nil {
		bc.Client = &conf.Client{}
		_ = cfg.Scan(&bc.Client)
	}

	if bc.Data == nil {
		bc.Data = &conf.Data{}
		_ = cfg.Scan(&bc.Data)
	}

	if bc.Trace == nil {
		bc.Trace = &conf.Tracer{}
		_ = cfg.Scan(&bc.Trace)
	}

	if bc.Logger == nil {
		bc.Logger = &conf.Logger{}
		_ = cfg.Scan(&bc.Logger)
	}

	if bc.Registry == nil {
		bc.Registry = &conf.Registry{}
		_ = cfg.Scan(&bc.Registry)
	}

	if bc.Oss == nil {
		bc.Oss = &conf.OSS{}
		_ = cfg.Scan(&bc.Oss)
	}

	return &bc
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// LoadRemoteConfigSourceConfigs 加载远程配置源的本地配置
func LoadRemoteConfigSourceConfigs(configPath string) (error, *conf.RemoteConfig) {
	configPath = configPath + "/" + remoteConfigSourceConfigFile
	if !pathExists(configPath) {
		return nil, nil
	}

	cfg := config.New(
		config.WithSource(
			NewFileConfigSource(configPath),
		),
	)
	defer func(cfg config.Config) {
		err := cfg.Close()
		if err != nil {
			panic(err)
		}
	}(cfg)

	var err error

	if err = cfg.Load(); err != nil {
		return err, nil
	}

	var rc conf.Bootstrap
	if err = cfg.Scan(&rc); err != nil {
		return err, nil
	}

	return nil, rc.Config
}

type ConfigType string

const (
	ConfigTypeLocalFile  ConfigType = "file"
	ConfigTypeNacos      ConfigType = "nacos"
	ConfigTypeConsul     ConfigType = "consul"
	ConfigTypeEtcd       ConfigType = "etcd"
	ConfigTypeApollo     ConfigType = "apollo"
	ConfigTypeKubernetes ConfigType = "kubernetes"
	ConfigTypePolaris    ConfigType = "polaris"
)

// NewRemoteConfigSource 创建一个远程配置源
func NewRemoteConfigSource(c *conf.RemoteConfig) config.Source {
	switch ConfigType(c.Type) {
	default:
		fallthrough
	case ConfigTypeLocalFile:
		return nil
	}
}

// getConfigKey 获取合法的配置名
func getConfigKey(configKey string, useBackslash bool) string {
	if useBackslash {
		return strings.Replace(configKey, `.`, `/`, -1)
	} else {
		return configKey
	}
}

// NewFileConfigSource 创建一个本地文件配置源
func NewFileConfigSource(filePath string) config.Source {
	return fileKratos.NewSource(filePath)
}
