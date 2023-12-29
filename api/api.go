package api

import (
	"context"

	cms "github.com/NpoolPlatform/message/npool/cms/gw/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	cms.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	cms.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := cms.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}

	return nil
}
