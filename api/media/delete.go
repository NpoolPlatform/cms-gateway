package media

import (
	"context"

	media1 "github.com/NpoolPlatform/cms-gateway/pkg/media"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/media"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteMedia(ctx context.Context, in *npool.DeleteMediaRequest) (*npool.DeleteMediaResponse, error) {
	handler, err := media1.NewHandler(
		ctx,
		media1.WithID(&in.ID, true),
		media1.WithEntID(&in.EntID, true),
		media1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteMedia",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteMediaResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteMedia(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteMedia",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteMediaResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteMediaResponse{
		Info: info,
	}, nil
}
