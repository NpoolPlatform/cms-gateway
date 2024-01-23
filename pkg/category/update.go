package category

import (
	"context"
	"fmt"

	articlemwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/article"
	categorymwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/category"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/category"
	articlemwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/article"
	categorymwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/category"
)

type updateHandler struct {
	*Handler
	oldInfo *categorymwpb.Category
}

func (h *updateHandler) checkParentID(ctx context.Context) error {
	if h.ParentID == nil {
		return nil
	}
	exist, err := categorymwcli.ExistCategoryConds(ctx, &categorymwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.ParentID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid parentid")
	}

	latest := true
	exist, err = articlemwcli.ExistArticleConds(ctx, &articlemwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CategoryID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		Latest:     &basetypes.BoolVal{Op: cruder.EQ, Value: latest},
	})
	if err != nil {
		return err
	}
	if exist && (*h.ParentID != h.oldInfo.ParentID) {
		return fmt.Errorf("invalid category parentid")
	}

	return nil
}

func (h *updateHandler) checkCagegoryExist(ctx context.Context) error {
	info, err := categorymwcli.GetCategoryOnly(ctx, &categorymwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("invalid category")
	}
	h.oldInfo = info
	return nil
}

func (h *updateHandler) checkCagegoryName(ctx context.Context) error {
	if h.Name == nil {
		return nil
	}
	if h.ParentID == nil {
		h.ParentID = &h.oldInfo.ParentID
	}
	exist, err := categorymwcli.ExistCategoryConds(ctx, &categorymwpb.Conds{
		EntID:    &basetypes.StringVal{Op: cruder.NEQ, Value: *h.EntID},
		AppID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		ParentID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.ParentID},
		Name:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.Name},
	})
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("invalid name")
	}
	return nil
}

func (h *Handler) UpdateCategory(ctx context.Context) (*npool.Category, error) {
	handler := &updateHandler{
		Handler: h,
	}

	if err := handler.checkCagegoryExist(ctx); err != nil {
		return nil, err
	}

	if err := handler.checkParentID(ctx); err != nil {
		return nil, err
	}

	if err := handler.checkCagegoryName(ctx); err != nil {
		return nil, err
	}

	info, err := categorymwcli.UpdateCategory(ctx, &categorymwpb.CategoryReq{
		ID:       h.ID,
		ParentID: h.ParentID,
		Name:     h.Name,
		Enabled:  h.Enabled,
		Index:    h.Index,
	})

	if err != nil {
		return nil, err
	}

	return h.GetCategoryExt(ctx, info)
}
