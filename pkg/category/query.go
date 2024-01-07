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

type queryHandler struct {
	*Handler
	categoryMap map[string]*npool.Category
}

func (h *queryHandler) buildCategoryTree(infos []*categorymwpb.Category, parentID, parentSlug string) []*npool.Category {
	result := []*npool.Category{}
	for _, info := range infos {
		if info.ParentID != parentID {
			continue
		}
		fullSlug := fmt.Sprintf("%v/%v", parentSlug, info.Slug)
		if parentSlug == "" {
			fullSlug = info.Slug
		}
		children := h.buildCategoryTree(infos, info.EntID, fullSlug)
		category := &npool.Category{
			ID:       info.ID,
			EntID:    info.EntID,
			ParentID: info.ParentID,
			Name:     info.Name,
			Slug:     info.Slug,
			FullSlug: fullSlug,
			Children: children,
		}
		result = append(result, category)
	}
	return result
}

func (h *Handler) GetCategories(ctx context.Context) ([]*npool.Category, error) {
	handler := &queryHandler{
		Handler:     h,
		categoryMap: map[string]*npool.Category{},
	}

	infos, _, err := categorymwcli.GetCategories(ctx, &categorymwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
	if err != nil {
		return nil, err
	}

	nilUUID := uuid.Nil.String()
	categories := handler.buildCategoryTree(infos, nilUUID, "")
	for _, item := range handler.categoryMap {
		fmt.Println("ITEM: ", item)
	}

	return categories, nil
}

func (h *Handler) GetCategoryList(ctx context.Context) ([]*categorymwpb.Category, uint32, error) {
	return categorymwcli.GetCategories(ctx, &categorymwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
}
