package media

import (
	"context"
	"fmt"

	mediamwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/media"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/cms/mw/v1/media"
)

func (h *Handler) DeleteMedia(ctx context.Context) (*npool.Media, error) {
	exist, err := mediamwcli.ExistMediaConds(ctx, &npool.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid media")
	}

	info, err := mediamwcli.DeleteMedia(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	return info, nil
}
