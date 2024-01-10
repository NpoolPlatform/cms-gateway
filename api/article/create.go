package article

import (
	"context"

	"google.golang.org/grpc/metadata"

	article1 "github.com/NpoolPlatform/cms-gateway/pkg/article"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/article"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateArticle(ctx context.Context, in *npool.CreateArticleRequest) (*npool.CreateArticleResponse, error) {
	host := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Sugar().Errorw(
			"CreateArticle",
			"In", in,
			"Error", "invalid host",
		)
		return &npool.CreateArticleResponse{}, status.Error(codes.InvalidArgument, "invalid host")
	}
	headHost, ok := md["x-forwarded-host"]
	if !ok {
		logger.Sugar().Errorw(
			"CreateArticle",
			"In", in,
			"Error", "invalid host",
		)
		return &npool.CreateArticleResponse{}, status.Error(codes.InvalidArgument, "invalid host")
	}
	host = headHost[0]

	handler, err := article1.NewHandler(
		ctx,
		article1.WithAppID(&in.AppID, true),
		article1.WithLangID(&in.TargetLangID, true),
		article1.WithCategoryID(&in.CategoryID, true),
		article1.WithUserID(&in.UserID, true),
		article1.WithTitle(&in.Title, true),
		article1.WithSubtitle(&in.Subtitle, false),
		article1.WithDigest(&in.Digest, false),
		article1.WithContent(&in.Content, true),
		article1.WithHost(&host, true),
		article1.WithACLEnabled(in.ACLEnabled, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateArticle",
			"In", in,
			"Error", err,
		)
		return &npool.CreateArticleResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateArticle(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateArticle",
			"In", in,
			"Error", err,
		)
		return &npool.CreateArticleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateArticleResponse{
		Info: info,
	}, nil
}
