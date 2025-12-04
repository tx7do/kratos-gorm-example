package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/http"

	swaggerUI "github.com/tx7do/kratos-swagger-ui"

	"kratos-gorm-example/app/user/service/cmd/server/assets"
	"kratos-gorm-example/app/user/service/internal/service"

	"kratos-gorm-example/api/gen/go/common/conf"
	userV1 "kratos-gorm-example/api/gen/go/user/service/v1"

	"kratos-gorm-example/pkg/bootstrap"
)

// NewRESTServer new an HTTP server.
func NewRESTServer(
	cfg *conf.Bootstrap, logger log.Logger,
	userService *service.UserService,
) *http.Server {
	if cfg == nil || cfg.Server == nil || cfg.Server.Rest == nil {
		return nil
	}

	srv := bootstrap.CreateRestServer(cfg, logging.Server(logger))

	userV1.RegisterUserServiceHTTPServer(srv, userService)

	if cfg.GetServer().GetRest().GetEnableSwagger() {
		swaggerUI.RegisterSwaggerUIServerWithOption(
			srv,
			swaggerUI.WithTitle("Kratos GORM Example User Service API"),
			swaggerUI.WithMemoryData(assets.OpenApiData, "yaml"),
		)
	}

	return srv
}
