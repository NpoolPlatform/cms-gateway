package category

import (
	"context"

	category1 "github.com/NpoolPlatform/cms-gateway/pkg/category"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/category"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateCategory(ctx context.Context, in *npool.UpdateCategoryRequest) (*npool.UpdateCategoryResponse, error) {
	handler, err := category1.NewHandler(
		ctx,
		category1.WithID(&in.ID, true),
		category1.WithEntID(&in.EntID, true),
		category1.WithAppID(&in.AppID, true),
		category1.WithParentID(in.ParentID, false),
		category1.WithName(in.Name, false),
		category1.WithEnabled(in.Enabled, false),
		category1.WithIndex(in.Index, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateCategory",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateCategoryResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateCategory(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateCategory",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateCategoryResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateCategoryResponse{
		Info: info,
	}, nil
}
