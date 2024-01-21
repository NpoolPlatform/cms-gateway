package category

import (
	"context"

	category1 "github.com/NpoolPlatform/cms-gateway/pkg/category"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/category"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateCategory(ctx context.Context, in *npool.CreateCategoryRequest) (*npool.CreateCategoryResponse, error) {
	handler, err := category1.NewHandler(
		ctx,
		category1.WithAppID(&in.AppID, true),
		category1.WithParentID(in.ParentID, false),
		category1.WithName(&in.Name, true),
		category1.WithSlug(&in.Slug, true),
		category1.WithEnabled(&in.Enabled, true),
		category1.WithIndex(&in.Index, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateCategory",
			"In", in,
			"Error", err,
		)
		return &npool.CreateCategoryResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCategory(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateCategory",
			"In", in,
			"Error", err,
		)
		return &npool.CreateCategoryResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCategoryResponse{
		Info: info,
	}, nil
}
