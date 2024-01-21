//nolint:dupl
package article

import (
	"context"

	article1 "github.com/NpoolPlatform/cms-gateway/pkg/article"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/article"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteArticle(ctx context.Context, in *npool.DeleteArticleRequest) (*npool.DeleteArticleResponse, error) {
	handler, err := article1.NewHandler(
		ctx,
		article1.WithID(&in.ID, true),
		article1.WithEntID(&in.EntID, true),
		article1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteArticle",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteArticleResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteArticle(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteArticle",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteArticleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteArticleResponse{
		Info: info,
	}, nil
}
