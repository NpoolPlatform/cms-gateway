package media

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/NpoolPlatform/cms-gateway/common/servermux"
	media1 "github.com/NpoolPlatform/cms-gateway/pkg/media"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/media"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	mux := servermux.AppServerMux()
	mux.HandleFunc("/v1/m/", Content)
}

func Content(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
	}

	var nonEmptyParts []string
	for _, part := range parts {
		if part != "" {
			nonEmptyParts = append(nonEmptyParts, part)
		}
	}

	minPathLength := 4
	if len(parts) < minPathLength {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	appID := nonEmptyParts[2]
	fileName := nonEmptyParts[3]

	ctx := r.Context()
	handler, err := media1.NewHandler(
		ctx,
		media1.WithAppID(&appID, true),
		media1.WithFileName(&fileName, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMedia",
			"Error", err,
		)
		fmt.Fprintf(w, "%v", err.Error())
		return
	}

	info, err := handler.GetMedia(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMedia",
			"Error", err,
		)
		fmt.Fprintf(w, "%v", err.Error())
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(string(info))
	if err != nil {
		http.Error(w, "Failed to decode base64 data", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", fmt.Sprint(len((decoded))))
	w.Header().Set("Content-Disposition", "inline")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(decoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetMedias(ctx context.Context, in *npool.GetMediasRequest) (*npool.GetMediasResponse, error) {
	handler, err := media1.NewHandler(
		ctx,
		media1.WithAppID(&in.AppID, true),
		media1.WithOffset(in.GetOffset()),
		media1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMedias",
			"In", in,
			"Error", err,
		)
		return &npool.GetMediasResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetMedias(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMedias",
			"In", in,
			"Error", err,
		)
		return &npool.GetMediasResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetMediasResponse{
		Infos: infos,
		Total: total,
	}, nil
}
