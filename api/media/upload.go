package media

import (
	"context"

	media1 "github.com/NpoolPlatform/cms-gateway/pkg/media"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/media"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UploadMedia(ctx context.Context, in *npool.UploadMediaRequest) (*npool.UploadMediaResponse, error) {
	handler, err := media1.NewHandler(
		ctx,
		media1.WithAppID(&in.AppID, true),
		media1.WithUserID(&in.UserID, true),
		media1.WithFileName(&in.Name, true),
		media1.WithMediaData(&in.MediaData, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UploadMedia",
			"In", in,
			"Error", err,
		)
		return &npool.UploadMediaResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UploadMedia(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UploadMedia",
			"In", in,
			"Error", err,
		)
		return &npool.UploadMediaResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UploadMediaResponse{
		Info: info,
	}, nil
}
