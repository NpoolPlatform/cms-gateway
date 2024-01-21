package media

import (
	"context"
	"fmt"

	mediamwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/media"
	mediamwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/media"

	"github.com/NpoolPlatform/go-service-framework/pkg/oss"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetMedias(ctx context.Context) ([]*mediamwpb.Media, uint32, error) {
	infos, total, err := mediamwcli.GetMedias(ctx, &mediamwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	},
		h.Offset,
		h.Limit,
	)
	if err != nil {
		return nil, 0, err
	}

	return infos, total, nil
}

func (h *Handler) GetMedia(ctx context.Context) ([]byte, error) {
	info, err := mediamwcli.GetMediaOnly(ctx, &mediamwpb.Conds{
		AppID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		MediaURL: &basetypes.StringVal{Op: cruder.EQ, Value: *h.FileName},
	})
	if err != nil {
		return []byte{}, err
	}
	if info == nil {
		return []byte{}, fmt.Errorf("not found")
	}
	key := fmt.Sprintf("media/%v/%v", *h.AppID, *h.FileName)
	content, err := oss.GetObject(ctx, key, true)
	if err != nil {
		return []byte{}, err
	}
	if content == nil {
		return []byte{}, fmt.Errorf("no content")
	}

	return content, nil
}
