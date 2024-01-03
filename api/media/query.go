package media

import (
	"context"

	media1 "github.com/NpoolPlatform/cms-gateway/pkg/media"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/media"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetMedia(ctx context.Context, in *npool.GetMediaRequest) (*npool.GetMediaResponse, error) {
	handler, err := media1.NewHandler(
		ctx,
		media1.WithAppID(&in.AppID, true),
		media1.WithFileName(&in.FileName, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMedia",
			"In", in,
			"Error", err,
		)
		return &npool.GetMediaResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetMedia(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMedia",
			"In", in,
			"Error", err,
		)
		return &npool.GetMediaResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetMediaResponse{
		Info: info,
	}, nil
}

func (s *Server) GetMedias(ctx context.Context, in *npool.GetMediasRequest) (*npool.GetMediasResponse, error) {
	handler, err := media1.NewHandler(
		ctx,
		media1.WithAppID(&in.AppID, true),
		media1.WithOffset(in.GetOffset()),
		media1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMedias",
			"In", in,
			"Error", err,
		)
		return &npool.GetMediasResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetMedias(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMedias",
			"In", in,
			"Error", err,
		)
		return &npool.GetMediasResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetMediasResponse{
		Infos: infos,
		Total: total,
	}, nil
}
