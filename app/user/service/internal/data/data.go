package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"

	"gorm.io/gorm"

	"kratos-gorm-example/api/gen/go/common/conf"

	"kratos-gorm-example/pkg/bootstrap"
)

// Data .
type Data struct {
	log *log.Helper
	db  *gorm.DB
	rdb *redis.Client
}

// NewData .
func NewData(db *gorm.DB, rdb *redis.Client, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "data/user-service"))

	d := &Data{
		db:  db,
		rdb: rdb,
		log: l,
	}

	return d, func() {
		l.Info("message", "closing the data resources")
		if err := d.rdb.Close(); err != nil {
			l.Error(err)
		}
	}, nil
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(cfg *conf.Bootstrap, logger log.Logger) *redis.Client {
	l := log.NewHelper(log.With(logger, "module", "redis/data/user-service"))
	return bootstrap.NewRedisClient(cfg, l)
}
