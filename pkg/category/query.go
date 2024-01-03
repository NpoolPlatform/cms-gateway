package category

import (
	"context"

	categorymwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/category"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/category"
	categorymwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/category"
	"github.com/google/uuid"
)

type queryHandler struct {
	*Handler
}

func (h *queryHandler) buildCategoryTree(infos []*categorymwpb.Category, parentID string) []*npool.Category {
	result := []*npool.Category{}
	for _, info := range infos {
		if info.ParentID == parentID {
			children := h.buildCategoryTree(infos, info.EntID)
			category := &npool.Category{
				ID:       info.ID,
				EntID:    info.EntID,
				ParentID: info.ParentID,
				Name:     info.Name,
				Slug:     info.Slug,
				Children: children,
			}
			result = append(result, category)
		}
	}
	return result
}

func (h *Handler) GetCategories(ctx context.Context) ([]*npool.Category, error) {
	handler := &queryHandler{
		Handler: h,
	}

	infos, _, err := categorymwcli.GetCategories(ctx, &categorymwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
	if err != nil {
		return nil, err
	}

	nilUUID := uuid.Nil.String()
	categories := handler.buildCategoryTree(infos, nilUUID)

	return categories, nil
}

func (h *Handler) GetCategoryList(ctx context.Context) ([]*categorymwpb.Category, uint32, error) {
	return categorymwcli.GetCategories(ctx, &categorymwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
}
