package category

import (
	"context"
	"fmt"

	articlemwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/article"
	categorymwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/category"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	articlemwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/article"
	npool "github.com/NpoolPlatform/message/npool/cms/mw/v1/category"
)

func (h *Handler) DeleteCategory(ctx context.Context) (*npool.Category, error) {
	exist, err := categorymwcli.ExistCategoryConds(ctx, &npool.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid category")
	}

	exist, err = articlemwcli.ExistArticleConds(ctx, &articlemwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CategoryID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("exist article")
	}

	return categorymwcli.DeleteCategory(ctx, *h.ID)
}
