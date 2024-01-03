package acl

import (
	"context"

	acl1 "github.com/NpoolPlatform/cms-gateway/pkg/acl"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/acl"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetACLs(ctx context.Context, in *npool.GetACLsRequest) (*npool.GetACLsResponse, error) {
	handler, err := acl1.NewHandler(
		ctx,
		acl1.WithAppID(&in.AppID, true),
		acl1.WithOffset(in.GetOffset()),
		acl1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetACLs",
			"In", in,
			"Error", err,
		)
		return &npool.GetACLsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetACLs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetACLs",
			"In", in,
			"Error", err,
		)
		return &npool.GetACLsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetACLsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
