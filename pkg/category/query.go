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

func (h *queryHandler) buildCategoryList(infos []*categorymwpb.Category, parentID, parentSlug string) {
	for _, info := range infos {
		if info.ParentID != parentID {
			continue
		}
		fullSlug := fmt.Sprintf("%v/%v", parentSlug, info.Slug)
		if parentSlug == "" {
			fullSlug = info.Slug
		}
		h.buildCategoryList(infos, info.EntID, fullSlug)
		category := &npool.Category{
			ID:       info.ID,
			EntID:    info.EntID,
			AppID:    info.AppID,
			ParentID: info.ParentID,
			Name:     info.Name,
			Slug:     info.Slug,
			Enabled:  info.Enabled,
			FullSlug: fullSlug,
		}
		h.categoryMap[category.EntID] = category
	}
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

func (h *queryHandler) getCategoryFullSlug(ctx context.Context, id string) (string, error) {
	fullSlug := ""
	for {
		category, err := categorymwcli.GetCategory(ctx, id)
		if err != nil {
			return "", err
		}
		if category == nil {
			return "", fmt.Errorf("invalid categoryid")
		}
		if fullSlug == "" {
			fullSlug = category.Slug
		} else {
			fullSlug = fmt.Sprintf("%v/%v", category.Slug, fullSlug)
		}
		nullUUID := uuid.Nil.String()
		if category.ParentID == nullUUID {
			break
		}
		id = category.ParentID
	}
	return fullSlug, nil
}

func (h *Handler) GetCategoryExt(ctx context.Context, row *categorymwpb.Category) (*npool.Category, error) {
	handler := &queryHandler{
		Handler: h,
	}
	fullSlug, err := handler.getCategoryFullSlug(ctx, row.EntID)
	if err != nil {
		return nil, err
	}
	info := &npool.Category{
		ID:       row.ID,
		EntID:    row.EntID,
		AppID:    row.AppID,
		ParentID: row.ParentID,
		Name:     row.Name,
		Slug:     row.Slug,
		FullSlug: fullSlug,
		Enabled:  row.Enabled,
	}

	return info, nil
}

func (h *Handler) GetCategoryList(ctx context.Context) ([]*npool.Category, error) {
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

	return categories, nil
}

func (h *Handler) GetCategories(ctx context.Context) ([]*npool.Category, uint32, error) {
	handler := &queryHandler{
		Handler:     h,
		categoryMap: map[string]*npool.Category{},
	}

	categories, total, err := categorymwcli.GetCategories(ctx, &categorymwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}

	nilUUID := uuid.Nil.String()
	handler.buildCategoryList(categories, nilUUID, "")

	infos := []*npool.Category{}
	for _, info := range handler.categoryMap {
		category := &npool.Category{
			ID:       info.ID,
			EntID:    info.EntID,
			AppID:    info.AppID,
			ParentID: info.ParentID,
			Name:     info.Name,
			Slug:     info.Slug,
			Enabled:  info.Enabled,
			FullSlug: info.FullSlug,
		}
		infos = append(infos, category)
	}

	return infos, total, nil
}
