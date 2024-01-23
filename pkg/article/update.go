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
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"
	"github.com/google/uuid"
)

type updateHandler struct {
	*Handler
	article *articlemwpb.Article
}

func (h *updateHandler) checkCategory(ctx context.Context) error {
	if h.CategoryID == nil {
		h.CategoryID = &h.article.CategoryID
		return nil
	}
	exist, err := categorymwcli.ExistCategory(ctx, *h.CategoryID)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid category")
	}
	return nil
}

func (h *updateHandler) checkArticleExist(ctx context.Context) error {
	info, err := articlemwcli.GetArticleOnly(ctx, &articlemwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("invalid article")
	}
	if !info.Latest {
		return fmt.Errorf("invalid old version article")
	}
	h.article = info
	return nil
}

func (h *updateHandler) checkTitle(ctx context.Context) error {
	if h.Title == nil {
		return nil
	}
	latest := true
	exist, err := articlemwcli.ExistArticleConds(ctx, &articlemwpb.Conds{
		EntID:  &basetypes.StringVal{Op: cruder.NEQ, Value: *h.EntID},
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		ISO:    &basetypes.StringVal{Op: cruder.EQ, Value: h.article.ISO},
		Title:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.Title},
		Latest: &basetypes.BoolVal{Op: cruder.EQ, Value: latest},
		Host:   &basetypes.StringVal{Op: cruder.EQ, Value: h.article.Host},
	})
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("repeated title")
	}
	return nil
}

func (h *updateHandler) getCategoryFullSlug(ctx context.Context, id string) (string, error) {
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

func (h *updateHandler) getISO(ctx context.Context) error {
	if h.LangID == nil {
		h.ISO = &h.article.ISO
		return nil
	}
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

func (h *updateHandler) getPageName() string {
	title := strings.ToLower(*h.Title)
	title = strings.ReplaceAll(title, " ", "-")
	title = fmt.Sprintf("%v.html", title)
	return title
}

func (h *updateHandler) checkVersion() {
	if h.UpdateContent == nil {
		return
	}
	if !*h.UpdateContent {
		return
	}
	if h.Content == nil {
		return
	}
	newVersion := h.article.Version + 1
	h.Version = &newVersion
	h.ArticleKey = &h.article.ArticleKey
	h.Host = &h.article.Host
}

func (h *updateHandler) checkContent(ctx context.Context) error {
	if h.Version == nil {
		return nil
	}

	cateGorySlugs, err := h.getCategoryFullSlug(ctx, *h.CategoryID)
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

func (h *Handler) UpdateArticle(ctx context.Context) (*articlemwpb.Article, error) {
	handler := &updateHandler{
		Handler: h,
	}

	if err := handler.checkArticleExist(ctx); err != nil {
		return nil, err
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
	handler.checkVersion()
	if err := handler.checkContent(ctx); err != nil {
		return nil, err
	}
	if h.Version != nil {
		return articlemwcli.CreateArticle(ctx, &articlemwpb.ArticleReq{
			AppID:      h.AppID,
			ISO:        h.ISO,
			CategoryID: h.CategoryID,
			Host:       h.Host,
			ArticleKey: h.ArticleKey,
			AuthorID:   h.UserID,
			Title:      h.Title,
			Subtitle:   h.Subtitle,
			Digest:     h.Digest,
			Status:     h.Status,
			Version:    h.Version,
			ContentURL: h.ContentURL,
			ACLEnabled: h.ACLEnabled,
		})
	}

	return articlemwcli.UpdateArticle(ctx, &articlemwpb.ArticleReq{
		ID:         h.ID,
		ISO:        h.ISO,
		AuthorID:   h.UserID,
		Title:      h.Title,
		Subtitle:   h.Subtitle,
		Digest:     h.Digest,
		Status:     h.Status,
		Version:    h.Version,
		ContentURL: h.ContentURL,
		ACLEnabled: h.ACLEnabled,
	})
}
