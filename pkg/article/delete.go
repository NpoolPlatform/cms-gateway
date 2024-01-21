package article

import (
	"context"
	"fmt"

	articlemwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/article"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/cms/mw/v1/article"
)

func (h *Handler) DeleteArticle(ctx context.Context) (*npool.Article, error) {
	exist, err := articlemwcli.ExistArticleConds(ctx, &npool.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid article")
	}

	info, err := articlemwcli.DeleteArticle(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	return info, nil
}
