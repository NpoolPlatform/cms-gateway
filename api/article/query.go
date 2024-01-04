//nolint:dupl
package article

import (
	"context"

	article1 "github.com/NpoolPlatform/cms-gateway/pkg/article"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/cms/v1"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/article"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetContent(ctx context.Context, in *npool.GetContentRequest) (*npool.GetContentResponse, error) {
	host := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Sugar().Errorw(
			"GetContent",
			"In", in,
			"Error", "invalid host",
		)
		return &npool.GetContentResponse{}, status.Error(codes.InvalidArgument, "invalid host")
	}
	headHost, ok := md["x-forwarded-host"]
	if !ok {
		logger.Sugar().Errorw(
			"GetContent",
			"In", in,
			"Error", "invalid host",
		)
		return &npool.GetContentResponse{}, status.Error(codes.InvalidArgument, "invalid host")
	}
	host = headHost[0]
	handler, err := article1.NewHandler(
		ctx,
		article1.WithHost(&host, true),
		article1.WithContentURL(&in.ContentURL, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetContent",
			"In", in,
			"Error", err,
		)
		return &npool.GetContentResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetContent(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetContent",
			"In", in,
			"Error", err,
		)
		return &npool.GetContentResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetContentResponse{
		Info: info,
	}, nil
}

func (s *Server) GetContents(ctx context.Context, in *npool.GetContentsRequest) (*npool.GetContentsResponse, error) {
	latest := true
	articleStatus := basetypes.ArticleStatus_Published
	handler, err := article1.NewHandler(
		ctx,
		article1.WithAppID(&in.AppID, true),
		article1.WithCategoryID(in.CategoryID, false),
		article1.WithLatest(&latest, true),
		article1.WithStatus(&articleStatus, true),
		article1.WithOffset(in.GetOffset()),
		article1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetContents",
			"In", in,
			"Error", err,
		)
		return &npool.GetContentsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetArticles(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetContents",
			"In", in,
			"Error", err,
		)
		return &npool.GetContentsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetContentsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppArticles(ctx context.Context, in *npool.GetAppArticlesRequest) (*npool.GetAppArticlesResponse, error) {
	handler, err := article1.NewHandler(
		ctx,
		article1.WithAppID(&in.AppID, true),
		article1.WithOffset(in.GetOffset()),
		article1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppArticles",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppArticlesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetArticles(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppArticles",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppArticlesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppArticlesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppArticleContent(ctx context.Context, in *npool.GetAppArticleContentRequest) (*npool.GetAppArticleContentResponse, error) {
	handler, err := article1.NewHandler(
		ctx,
		article1.WithID(&in.ID, true),
		article1.WithEntID(&in.EntID, true),
		article1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppArticleContent",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppArticleContentResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetArticleContent(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppArticleContent",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppArticleContentResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppArticleContentResponse{
		Info: info,
	}, nil
}
