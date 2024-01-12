package article

import (
	"context"
	"fmt"

	roleusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role/user"
	articlemwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/article"
	"github.com/NpoolPlatform/go-service-framework/pkg/oss"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	roleusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role/user"
	types "github.com/NpoolPlatform/message/npool/basetypes/cms/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	articlemwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/article"
)

type queryHandler struct {
	*Handler
}

func (h *queryHandler) getContent(ctx context.Context) (string, error) {
	key := fmt.Sprintf("article/%v/%v/%v/%v", *h.Host, *h.ArticleKey, *h.Version, *h.ContentURL)
	content, err := oss.GetObject(ctx, key, true)
	if err != nil {
		return "", err
	}
	if content == nil {
		return "", fmt.Errorf("without content")
	}

	return string(content), nil
}

func (h *Handler) GetArticles(ctx context.Context) ([]*articlemwpb.Article, uint32, error) {
	conds := &articlemwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.CategoryID != nil {
		conds.CategoryID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.CategoryID}
	}
	if h.Latest != nil {
		conds.Latest = &basetypes.BoolVal{Op: cruder.EQ, Value: *h.Latest}
	}
	if h.Status != nil {
		conds.Status = &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.Status)}
	}
	return articlemwcli.GetArticles(ctx, conds, h.Offset, h.Limit)
}

func (h *Handler) GetContent(ctx context.Context) (string, error) {
	handler := &queryHandler{
		Handler: h,
	}
	latest := true

	info, err := articlemwcli.GetArticleOnly(ctx, &articlemwpb.Conds{
		Host:       &basetypes.StringVal{Op: cruder.EQ, Value: *h.Host},
		ContentURL: &basetypes.StringVal{Op: cruder.EQ, Value: *h.ContentURL},
		Latest:     &basetypes.BoolVal{Op: cruder.EQ, Value: latest},
		Status:     &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(types.ArticleStatus_Published)},
	})
	if err != nil {
		return "", err
	}
	if info == nil {
		return "", fmt.Errorf("not found page")
	}
	if info.ACLEnabled && len(info.ACLRoleIDs) == 0 {
		return "", fmt.Errorf("permission define")
	}
	if info.ACLEnabled && len(info.ACLRoleIDs) != 0 {
		if h.UserID == nil || h.AppID == nil {
			return "", fmt.Errorf("permission define")
		}
		exist, err := roleusermwcli.ExistUserConds(ctx, &roleusermwpb.Conds{
			AppID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
			RoleIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: info.ACLRoleIDs},
		})
		if err != nil {
			return "", err
		}
		if !exist {
			return "", fmt.Errorf("permission define")
		}
	}
	h.Host = &info.Host
	h.ArticleKey = &info.ArticleKey
	h.Version = &info.Version
	h.ContentURL = &info.ContentURL
	content, err := handler.getContent(ctx)
	if err != nil {
		return "", nil
	}

	return content, nil
}

func (h *Handler) GetArticleContent(ctx context.Context) (string, error) {
	handler := &queryHandler{
		Handler: h,
	}
	info, err := articlemwcli.GetArticleOnly(ctx, &articlemwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return "", err
	}
	if info == nil {
		return "", fmt.Errorf("article not exist")
	}
	h.Host = &info.Host
	h.ArticleKey = &info.ArticleKey
	h.Version = &info.Version
	h.ContentURL = &info.ContentURL
	content, err := handler.getContent(ctx)
	if err != nil {
		return "", nil
	}

	return content, nil
}
