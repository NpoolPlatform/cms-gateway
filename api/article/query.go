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
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	mux := servermux.AppServerMux()
	mux.HandleFunc("/api/cms/v1/t/", Content)
}

func Content(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	host := r.Host
	fmt.Println("host: ", host)
	parts := strings.Split(path, "/")
	fmt.Println("parts: ", parts)

	for i, item := range parts {
		fmt.Printf("i: %v, item: %v\n", i, item)
	}

	var nonEmptyParts []string
	for _, part := range parts {
		if part != "" {
			nonEmptyParts = append(nonEmptyParts, part)
		}
	}
	fmt.Println("nonEmptyParts: ", nonEmptyParts)
	for i, item := range nonEmptyParts {
		fmt.Printf("i: %v, nonitem: %v\n", i, item)
	}

	minPathLength := 7
	if len(parts) < minPathLength {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	contentURL := nonEmptyParts[4]
	for i := 5; i < len(nonEmptyParts); i++ {
		contentURL = fmt.Sprintf("%v/%v", contentURL, nonEmptyParts[i])
	}
	ctx := r.Context()
	handler, err := article1.NewHandler(
		ctx,
		article1.WithHost(&host, true),
		article1.WithContentURL(&contentURL, true),
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

	fmt.Fprintf(w, "%v", info)
}

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
