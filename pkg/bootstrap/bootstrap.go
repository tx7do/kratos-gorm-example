package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"

	"kratos-gorm-example/api/gen/go/common/conf"
)

// Bootstrap 应用引导启动
func Bootstrap(serviceInfo *ServiceInfo) (*conf.Bootstrap, log.Logger) {
	// inject command flags
	Flags := NewCommandFlags()
	Flags.Init()

	// load configs
	cfg := LoadBootstrapConfig(Flags.Conf)
	if cfg == nil {
		panic("load config failed")
	}

	// init logger
	ll := NewLoggerProvider(cfg.Logger, serviceInfo)

	// init tracer
	err := NewTracerProvider(cfg.Trace, serviceInfo)
	if err != nil {
		panic(err)
	}

	return cfg, ll
}
