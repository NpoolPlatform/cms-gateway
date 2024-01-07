//nolint:dupl
package article

import (
	"context"
	"fmt"
	"strings"

	articlemwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/article"
	categorymwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/category"
	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	"github.com/NpoolPlatform/go-service-framework/pkg/oss"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	articlemwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/article"
	categorymwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/category"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) checkCategory(ctx context.Context) error {
	exist, err := categorymwcli.ExistCategory(ctx, *h.CategoryID)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid category")
	}
	return nil
}

func (h *createHandler) checkTitle(ctx context.Context) error {
	latest := true
	exist, err := articlemwcli.ExistArticleConds(ctx, &articlemwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		Title:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.Title},
		Latest: &basetypes.BoolVal{Op: cruder.EQ, Value: latest},
	})
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("repeated title")
	}
	return nil
}

func (h *createHandler) getCategories(ctx context.Context, id string, level int32) ([]*categorymwpb.Category, error) {
	categories := []*categorymwpb.Category{}
	category, err := categorymwcli.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, fmt.Errorf("invalid categoryid")
	}
	categories = append(categories, category)

	nullUUID := uuid.Nil.String()
	if category.ParentID != nullUUID {
		parentCategories, err := h.getCategories(ctx, category.ParentID, level+1)
		if err != nil {
			return nil, err
		}
		categories = append(categories, parentCategories...)
	}
	return categories, nil
}

func (h *createHandler) getGategorySlugs(ctx context.Context) (string, error) {
	categories, err := h.getCategories(ctx, *h.CategoryID, 0)
	if err != nil {
		return "", err
	}
	categorySlugs := ""

	for index, info := range categories {
		if index == 0 {
			categorySlugs = info.Slug
			if *h.CategoryID == info.EntID {
				break
			}
			continue
		}
		categorySlugs = fmt.Sprintf("%v/%v", categorySlugs, info.Slug)
		if *h.CategoryID == info.EntID {
			break
		}
	}
	return categorySlugs, nil
}

func (h *createHandler) getISO(ctx context.Context) error {
	info, err := applangmwcli.GetLangOnly(ctx, &applangmwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		LangID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.LangID},
	})
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("invalid langID")
	}
	h.ISO = &info.Lang
	return nil
}

func (h *createHandler) getPageName() string {
	title := strings.ToLower(*h.Title)
	title = strings.ReplaceAll(title, " ", "-")
	title = fmt.Sprintf("%v.html", title)
	return title
}

func (h *createHandler) uploadContent(ctx context.Context) error {
	articleKey := uuid.NewString()
	h.ArticleKey = &articleKey

	version := uint32(1)
	h.Version = &version

	cateGorySlugs, err := h.getGategorySlugs(ctx)
	if err != nil {
		return err
	}

	pageName := h.getPageName()

	contentURL := fmt.Sprintf("%v/%v/%v", cateGorySlugs, *h.ISO, pageName)
	h.ContentURL = &contentURL

	key := fmt.Sprintf("article/%v/%v/%v/%v", *h.Host, *h.ArticleKey, *h.Version, *h.ContentURL)

	content := h.Content
	if content == nil || *content == "" {
		return fmt.Errorf("invalid content")
	}
	if err := oss.PutObject(ctx, key, []byte(*content), true); err != nil {
		return err
	}
	return nil
}

func (h *createHandler) createArticle(ctx context.Context) (*articlemwpb.Article, error) {
	info, err := articlemwcli.CreateArticle(ctx, &articlemwpb.ArticleReq{
		AppID:      h.AppID,
		CategoryID: h.CategoryID,
		AuthorID:   h.UserID,
		ArticleKey: h.ArticleKey,
		Title:      h.Title,
		Subtitle:   h.Subtitle,
		Digest:     h.Digest,
		Status:     h.Status,
		Version:    h.Version,
		Host:       h.Host,
		ISO:        h.ISO,
		ContentURL: h.ContentURL,
	})
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (h *Handler) CreateArticle(ctx context.Context) (*articlemwpb.Article, error) {
	handler := &createHandler{
		Handler: h,
	}

	if err := handler.checkCategory(ctx); err != nil {
		return nil, err
	}
	if err := handler.checkTitle(ctx); err != nil {
		return nil, err
	}
	if err := handler.getISO(ctx); err != nil {
		return nil, err
	}
	if err := handler.uploadContent(ctx); err != nil {
		return nil, err
	}

	info, err := handler.createArticle(ctx)
	if err != nil {
		return nil, err
	}

	return info, nil
}
