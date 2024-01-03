//nolint:dupl
package category

import (
	"context"
	"fmt"

	categorymwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/category"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	categorymwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/category"
)

type updateHandler struct {
	*Handler
	oldInfo *categorymwpb.Category
}

func (h *updateHandler) checkParentExist(ctx context.Context) error {
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

func (h *Handler) UpdateCategory(ctx context.Context) (*categorymwpb.Category, error) {
	handler := &updateHandler{
		Handler: h,
	}

	if err := handler.checkParentExist(ctx); err != nil {
		return nil, err
	}

	if err := handler.checkCagegoryExist(ctx); err != nil {
		return nil, err
	}

	if err := handler.checkCagegoryName(ctx); err != nil {
		return nil, err
	}

	return categorymwcli.UpdateCategory(ctx, &categorymwpb.CategoryReq{
		ID:       h.ID,
		ParentID: h.ParentID,
		Name:     h.Name,
		Slug:     h.Slug,
		Enabled:  h.Enabled,
	})
}
