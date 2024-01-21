package category

import (
	"context"

	category1 "github.com/NpoolPlatform/cms-gateway/pkg/category"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/category"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCategoryList(ctx context.Context, in *npool.GetCategoryListRequest) (*npool.GetCategoryListResponse, error) {
	handler, err := category1.NewHandler(
		ctx,
		category1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCategoryList",
			"In", in,
			"Error", err,
		)
		return &npool.GetCategoryListResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, err := handler.GetCategoryList(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCategoryList",
			"In", in,
			"Error", err,
		)
		return &npool.GetCategoryListResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCategoryListResponse{
		Infos: infos,
	}, nil
}

func (s *Server) GetCategories(ctx context.Context, in *npool.GetCategoriesRequest) (*npool.GetCategoriesResponse, error) {
	handler, err := category1.NewHandler(
		ctx,
		category1.WithAppID(&in.AppID, true),
		category1.WithOffset(in.GetOffset()),
		category1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCategories",
			"In", in,
			"Error", err,
		)
		return &npool.GetCategoriesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCategories(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCategories",
			"In", in,
			"Error", err,
		)
		return &npool.GetCategoriesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCategoriesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
