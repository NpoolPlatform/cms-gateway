//nolint:dupl
package media

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	mediamwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/media"
	"github.com/NpoolPlatform/go-service-framework/pkg/oss"
	mediamwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/media"
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) checkFileExt() error {
	if h.FileName == nil {
		return fmt.Errorf("invalid name")
	}
	ext := filepath.Ext(*h.FileName)
	h.Ext = &ext
	return nil
}

func (h *createHandler) uploadFile(ctx context.Context) (string, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	mediaURL := fmt.Sprintf("%v%v", *h.EntID, *h.Ext)
	key := fmt.Sprintf("media/%v/%v", *h.AppID, mediaURL)
	fileBytes, err := io.ReadAll(h.FileData)
	if err != nil {
		return "", err
	}

	if err := oss.PutObject(ctx, key, fileBytes, true); err != nil {
		return "", err
	}

	return mediaURL, nil
}

func (h *createHandler) uploadMedia(ctx context.Context) (string, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	mediaURL := fmt.Sprintf("%v%v", *h.EntID, *h.Ext)
	key := fmt.Sprintf("media/%v/%v", *h.AppID, mediaURL)
	content := h.MediaData
	if content == nil || *content == "" {
		return "", fmt.Errorf("invalid content")
	}

	if err := oss.PutObject(ctx, key, []byte(*content), true); err != nil {
		return "", err
	}

	return mediaURL, nil
}

func (h *Handler) UploadMedia(ctx context.Context) (*mediamwpb.Media, error) {
	handler := &createHandler{
		Handler: h,
	}

	if err := handler.checkFileExt(); err != nil {
		return nil, err
	}

	mediaURL, err := handler.uploadMedia(ctx)
	if err != nil {
		return nil, err
	}

	return mediamwcli.CreateMedia(ctx, &mediamwpb.MediaReq{
		EntID:     h.EntID,
		AppID:     h.AppID,
		Name:      h.FileName,
		MediaURL:  &mediaURL,
		Ext:       h.Ext,
		CreatedBy: h.UserID,
	})
}

func (h *Handler) UploadFile(ctx context.Context) (*mediamwpb.Media, error) {
	handler := &createHandler{
		Handler: h,
	}

	if err := handler.checkFileExt(); err != nil {
		return nil, err
	}

	mediaURL, err := handler.uploadFile(ctx)
	if err != nil {
		return nil, err
	}

	return mediamwcli.CreateMedia(ctx, &mediamwpb.MediaReq{
		EntID:     h.EntID,
		AppID:     h.AppID,
		Name:      h.FileName,
		MediaURL:  &mediaURL,
		Ext:       h.Ext,
		CreatedBy: h.UserID,
	})
}
