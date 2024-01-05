package main

import (
	"context"
	"fmt"
	"net/http"

	apicli "github.com/NpoolPlatform/basal-middleware/pkg/client/api"
	"github.com/NpoolPlatform/cms-gateway/api"
	"github.com/NpoolPlatform/cms-gateway/common/servermux"
	"github.com/NpoolPlatform/cms-gateway/config"
	"github.com/NpoolPlatform/cms-gateway/pkg/migrator"
	"github.com/NpoolPlatform/go-service-framework/pkg/action"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/go-service-framework/pkg/oss"
	ossconst "github.com/NpoolPlatform/go-service-framework/pkg/oss/const"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	cli "github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const BukectKey = "cms_bucket"

var runCmd = &cli.Command{
	Name:    "run",
	Aliases: []string{"s"},
	Usage:   "Run the daemon",
	Before: func(ctx *cli.Context) error {
		return logger.Init(logger.DebugLevel, config.GetConfig().CMS.LogFile)
	},
	Action: func(c *cli.Context) error {
		go runHTTPServer(config.GetConfig().CMS.HTTPPort, config.GetConfig().CMS.GrpcPort)
		err := action.Run(
			c.Context,
			run,
			rpcRegister,
			rpcGatewayRegister,
			watch,
		)

		return err
	},
}

func run(ctx context.Context) error {
	if err := oss.Init(ossconst.SecretStoreKey, BukectKey); err != nil {
		return err
	}
	if err := migrator.Migrate(ctx); err != nil {
		return err
	}
	return nil
}

func shutdown(ctx context.Context) {
	<-ctx.Done()
	logger.Sugar().Infow(
		"Watch",
		"State", "Done",
		"Error", ctx.Err(),
	)
}

func watch(ctx context.Context, cancel context.CancelFunc) error {
	go shutdown(ctx)
	return nil
}

func rpcRegister(server grpc.ServiceRegistrar) error {
	api.Register(server)

	apicli.RegisterGRPC(server)

	return nil
}

func rpcGatewayRegister(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	err := api.RegisterGateway(mux, endpoint, opts)
	if err != nil {
		return err
	}

	_ = apicli.Register(mux)

	return nil
}

func runHTTPServer(httpPort, grpcPort int) {
	httpEndpoint := fmt.Sprintf(":%v", httpPort)
	grpcEndpoint := fmt.Sprintf(":%v", grpcPort)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := api.RegisterGateway(mux, grpcEndpoint, opts)
	if err != nil {
		logger.Sugar().Infow(
			"Watch",
			"State", "Done",
			"Error", fmt.Sprintf("Fail to register gRPC gateway service endpoint: %v", err),
		)
	}

	http.Handle("/v1/", mux)
	servermux.AppServerMux().Handle("/v1/", mux)
	err = http.ListenAndServe(httpEndpoint, servermux.AppServerMux())
	if err != nil {
		logger.Sugar().Infow(
			"Watch",
			"State", "Done",
			"Error", fmt.Sprintf("failed to listen: %v", err),
		)
	}
}
