//nolint:dupl
package category

import (
	"context"
	"fmt"

	categorymwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/category"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/category"
	categorymwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/category"
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) checkParentExist(ctx context.Context) error {
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

func (h *createHandler) checkCagegoryName(ctx context.Context) error {
	if h.Name == nil {
		return fmt.Errorf("invalid name")
	}
	parentID := h.ParentID
	if h.ParentID == nil {
		nilUUID := uuid.Nil.String()
		parentID = &nilUUID
	}
	exist, err := categorymwcli.ExistCategoryConds(ctx, &categorymwpb.Conds{
		AppID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		ParentID: &basetypes.StringVal{Op: cruder.EQ, Value: *parentID},
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

func (h *createHandler) checkCagegorySlug(ctx context.Context) error {
	if h.Slug == nil {
		return fmt.Errorf("invalid slug")
	}
	parentID := h.ParentID
	if h.ParentID == nil {
		nilUUID := uuid.Nil.String()
		parentID = &nilUUID
	}
	exist, err := categorymwcli.ExistCategoryConds(ctx, &categorymwpb.Conds{
		AppID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		ParentID: &basetypes.StringVal{Op: cruder.EQ, Value: *parentID},
		Slug:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.Slug},
	})
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("invalid slug")
	}
	return nil
}

func (h *Handler) CreateCategory(ctx context.Context) (*npool.Category, error) {
	handler := &createHandler{
		Handler: h,
	}

	if err := handler.checkParentExist(ctx); err != nil {
		return nil, err
	}

	if err := handler.checkCagegoryName(ctx); err != nil {
		return nil, err
	}

	if err := handler.checkCagegorySlug(ctx); err != nil {
		return nil, err
	}

	info, err := categorymwcli.CreateCategory(ctx, &categorymwpb.CategoryReq{
		AppID:    h.AppID,
		ParentID: h.ParentID,
		Name:     h.Name,
		Slug:     h.Slug,
		Enabled:  h.Enabled,
		Index:    h.Index,
	})

	if err != nil {
		return nil, err
	}

	return h.GetCategoryExt(ctx, info)
}
