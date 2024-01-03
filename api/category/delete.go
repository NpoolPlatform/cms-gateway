package category

import (
	"context"

	category1 "github.com/NpoolPlatform/cms-gateway/pkg/category"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/category"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteCategory(ctx context.Context, in *npool.DeleteCategoryRequest) (*npool.DeleteCategoryResponse, error) {
	handler, err := category1.NewHandler(
		ctx,
		category1.WithID(&in.ID, true),
		category1.WithEntID(&in.EntID, true),
		category1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteCategory",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteCategoryResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteCategory(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteCategory",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteCategoryResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteCategoryResponse{
		Info: info,
	}, nil
}
