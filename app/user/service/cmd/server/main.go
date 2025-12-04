package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"kratos-gorm-example/pkg/bootstrap"
	"kratos-gorm-example/pkg/service"
)

// go build -ldflags "-X main.Service.Version=x.y.z"

var (
	Service = bootstrap.NewServiceInfo(
		service.UserService,
		"1.0.0",
		"",
	)
)

func newApp(ll log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(Service.GetInstanceId()),
		kratos.Name(Service.Name),
		kratos.Version(Service.Version),
		kratos.Metadata(Service.Metadata),
		kratos.Logger(ll),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	// bootstrap
	cfg, ll := bootstrap.Bootstrap(Service)

	app, cleanup, err := initApp(ll, cfg)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err = app.Run(); err != nil {
		panic(err)
	}
}
