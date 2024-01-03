package acl

import (
	"context"

	"github.com/NpoolPlatform/message/npool/cms/gw/v1/acl"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	acl.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	acl.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return acl.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
