package category

import (
	"context"

	"github.com/NpoolPlatform/message/npool/cms/gw/v1/category"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	category.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	category.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return category.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
