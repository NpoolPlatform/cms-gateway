package article

import (
	"context"

	"github.com/NpoolPlatform/message/npool/cms/gw/v1/article"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	article.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	article.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return article.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
