//nolint:dupl
package acl

import (
	"context"

	acl1 "github.com/NpoolPlatform/cms-gateway/pkg/acl"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/acl"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteACL(ctx context.Context, in *npool.DeleteACLRequest) (*npool.DeleteACLResponse, error) {
	handler, err := acl1.NewHandler(
		ctx,
		acl1.WithID(&in.ID, true),
		acl1.WithEntID(&in.EntID, true),
		acl1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteACL",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteACLResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteACL(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteACL",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteACLResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteACLResponse{
		Info: info,
	}, nil
}
