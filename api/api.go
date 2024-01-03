package api

import (
	"context"

	"github.com/NpoolPlatform/cms-gateway/api/acl"
	"github.com/NpoolPlatform/cms-gateway/api/article"
	"github.com/NpoolPlatform/cms-gateway/api/category"
	"github.com/NpoolPlatform/cms-gateway/api/media"
	cms "github.com/NpoolPlatform/message/npool/cms/gw/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	cms.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	cms.RegisterGatewayServer(server, &Server{})
	acl.Register(server)
	article.Register(server)
	category.Register(server)
	media.Register(server)
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := cms.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	if err := acl.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := article.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := category.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := media.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}

	return nil
}
