package media

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NpoolPlatform/cms-gateway/common/servermux"

	media1 "github.com/NpoolPlatform/cms-gateway/pkg/media"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/media"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const MaxUploadFileSize = 10 << 20

func init() {
	mux := servermux.AppServerMux()
	mux.HandleFunc("/v1/upload/app/media", Upload)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(MaxUploadFileSize) // 限制最大文件大小为 10MB
	if err != nil {
		fmt.Println(err)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fmt.Println("file: ", file)

	appID := r.FormValue("AppID")
	userID := r.FormValue("UserID")
	name := r.FormValue("Name")
	fmt.Println("appID: ", appID)
	fmt.Println("userID: ", userID)
	fmt.Println("name: ", name)

	ctx := r.Context()
	handler, err := media1.NewHandler(
		ctx,
		media1.WithAppID(&appID, true),
		media1.WithUserID(&userID, true),
		media1.WithFileName(&name, true),
		media1.WithFileData(file),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UploadMedia",
			"Error", err,
		)
		fmt.Fprintf(w, "%v", err.Error())
		return
	}

	info, err := handler.UploadMedia(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UploadMedia",
			"Error", err,
		)
		fmt.Fprintf(w, "%v", err.Error())
		return
	}

	fmt.Fprintf(w, "%v", info)
}

func (s *Server) UploadMedia(ctx context.Context, in *npool.UploadMediaRequest) (*npool.UploadMediaResponse, error) {
	handler, err := media1.NewHandler(
		ctx,
		media1.WithAppID(&in.AppID, true),
		media1.WithUserID(&in.UserID, true),
		media1.WithFileName(&in.Name, true),
		media1.WithMediaData(&in.MediaData, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UploadMedia",
			"In", in,
			"Error", err,
		)
		return &npool.UploadMediaResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UploadMedia(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UploadMedia",
			"In", in,
			"Error", err,
		)
		return &npool.UploadMediaResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UploadMediaResponse{
		Info: info,
	}, nil
}
