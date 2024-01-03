package media

import (
	"context"

	"github.com/NpoolPlatform/message/npool/cms/gw/v1/media"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	media.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	media.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return media.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
