package article

import (
	"context"

	article1 "github.com/NpoolPlatform/cms-gateway/pkg/article"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/article"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateArticle(ctx context.Context, in *npool.UpdateArticleRequest) (*npool.UpdateArticleResponse, error) {
	handler, err := article1.NewHandler(
		ctx,
		article1.WithID(&in.ID, true),
		article1.WithEntID(&in.EntID, true),
		article1.WithAppID(&in.AppID, true),
		article1.WithCategoryID(in.CategoryID, false),
		article1.WithUserID(in.UserID, false),
		article1.WithTitle(in.Title, false),
		article1.WithSubtitle(in.Subtitle, false),
		article1.WithDigest(in.Digest, false),
		article1.WithContent(in.Content, false),
		article1.WithUpdateContent(in.UpdateContent, false),
		article1.WithStatus(in.Status, false),
		article1.WithACLEnabled(in.ACLEnabled, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateArticle",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateArticleResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateArticle(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateArticle",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateArticleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateArticleResponse{
		Info: info,
	}, nil
}
