package data

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/logger"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"

	"gorm.io/gorm"

	"kratos-gorm-example/api/gen/go/common/conf"
	"kratos-gorm-example/app/user/service/internal/data/models"
)

type gormLoggerWriter struct {
	helper *log.Helper
}

func (w gormLoggerWriter) Printf(format string, args ...interface{}) {
	w.helper.Infof(format, args...)
}

func NewGormLogger(l log.Logger) logger.Interface {
	w := gormLoggerWriter{helper: log.NewHelper(log.With(l, "module", "gorm/data/user-service"))}

	return logger.New(
		w,
		logger.Config{
			SlowThreshold: time.Millisecond * 100, // 慢 SQL 阈值（超过 100ms 标为慢 SQL）
			LogLevel:      logger.Info,            // 核心：Info 级别会打印所有 SQL
			Colorful:      true,                   // 终端彩色输出（文件输出需关闭）
		},
	)
}

// NewGormClient 创建数据库客户端
func NewGormClient(cfg *conf.Bootstrap, logger log.Logger) *gorm.DB {
	l := log.NewHelper(log.With(logger, "module", "gorm/data/user-service"))

	var driver gorm.Dialector
	switch cfg.Data.Database.Driver {
	default:
		fallthrough
	case "mysql":
		driver = mysql.Open(cfg.Data.Database.Source)
		break
	case "postgres":
		driver = postgres.Open(cfg.Data.Database.Source)
		break
	case "clickhouse":
		driver = clickhouse.Open(cfg.Data.Database.Source)
		break
	case "sqlite":
		driver = sqlite.Open(cfg.Data.Database.Source)
		break
	case "sqlserver":
		driver = sqlserver.Open(cfg.Data.Database.Source)
		break
	}

	client, err := gorm.Open(driver, &gorm.Config{
		Logger: NewGormLogger(logger),
	})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
		return nil
	}

	// 运行数据库迁移工具
	if cfg.Data.Database.Migrate {
		if err = client.AutoMigrate(
			models.GetMigrates()...,
		); err != nil {
			l.Fatalf("failed creating schema resources: %v", err)
			return nil
		}
	}
	return client
}
