//nolint:dupl
package article

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/NpoolPlatform/cms-gateway/common/servermux"
	article1 "github.com/NpoolPlatform/cms-gateway/pkg/article"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/cms/v1"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/article"
	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	mux := servermux.AppServerMux()
	mux.HandleFunc("/v1/c/", Content)
}

func Content(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	host := r.Host
	parts := strings.Split(path, "/")

	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
	}

	nilUUID := uuid.Nil.String()
	appID := r.Header.Get("X-App-Id")
	userID := r.Header.Get("X-User-Id")
	if appID == "" {
		appID = nilUUID
	}
	if userID == "" {
		userID = nilUUID
	}

	var nonEmptyParts []string
	for _, part := range parts {
		if part != "" {
			nonEmptyParts = append(nonEmptyParts, part)
		}
	}

	minPathLength := 3
	if len(nonEmptyParts) < minPathLength {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	contentURL := nonEmptyParts[2]
	for i := 3; i < len(nonEmptyParts); i++ {
		contentURL = fmt.Sprintf("%v/%v", contentURL, nonEmptyParts[i])
	}
	ctx := r.Context()
	handler, err := article1.NewHandler(
		ctx,
		article1.WithHost(&host, true),
		article1.WithContentURL(&contentURL, true),
		article1.WithAppID(&appID, false),
		article1.WithUserID(&userID, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetContent",
			"Error", err,
		)
		fmt.Fprintf(w, "%v", err.Error())
		return
	}

	info, err := handler.GetContent(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetContent",
			"Error", err,
		)
		fmt.Fprintf(w, "%v", err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "%v", info)
}

func (s *Server) GetContentList(ctx context.Context, in *npool.GetContentListRequest) (*npool.GetContentListResponse, error) {
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
			"GetContentList",
			"In", in,
			"Error", err,
		)
		return &npool.GetContentListResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetArticles(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetContentList",
			"In", in,
			"Error", err,
		)
		return &npool.GetContentListResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetContentListResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetArticle(ctx context.Context, in *npool.GetArticleRequest) (*npool.GetArticleResponse, error) {
	handler, err := article1.NewHandler(
		ctx,
		article1.WithEntID(&in.EntID, true),
		article1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetArticle",
			"In", in,
			"Error", err,
		)
		return &npool.GetArticleResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetArticle(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetArticle",
			"In", in,
			"Error", err,
		)
		return &npool.GetArticleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetArticleResponse{
		Info: info,
	}, nil
}

func (s *Server) GetArticles(ctx context.Context, in *npool.GetArticlesRequest) (*npool.GetArticlesResponse, error) {
	handler, err := article1.NewHandler(
		ctx,
		article1.WithAppID(&in.AppID, true),
		article1.WithOffset(in.GetOffset()),
		article1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetArticles",
			"In", in,
			"Error", err,
		)
		return &npool.GetArticlesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetArticles(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetArticles",
			"In", in,
			"Error", err,
		)
		return &npool.GetArticlesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetArticlesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetArticleContent(ctx context.Context, in *npool.GetArticleContentRequest) (*npool.GetArticleContentResponse, error) {
	handler, err := article1.NewHandler(
		ctx,
		article1.WithID(&in.ID, true),
		article1.WithEntID(&in.EntID, true),
		article1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetArticleContent",
			"In", in,
			"Error", err,
		)
		return &npool.GetArticleContentResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetArticleContent(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetArticleContent",
			"In", in,
			"Error", err,
		)
		return &npool.GetArticleContentResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetArticleContentResponse{
		Info: info,
	}, nil
}
