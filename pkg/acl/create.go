package acl

import (
	"context"
	"fmt"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	aclmwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/acl"
	articlemwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/article"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	aclmwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/acl"
	articlemwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/article"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/acl"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) checkRoleExist(ctx context.Context) error {
	exist, err := rolemwcli.ExistRoleConds(ctx, &rolemwpb.Conds{
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.RoleID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid role")
	}

	return nil
}

func (h *createHandler) checkArticleExist(ctx context.Context) error {
	exist, err := articlemwcli.ExistArticleConds(ctx, &articlemwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		ArticleKey: &basetypes.StringVal{Op: cruder.EQ, Value: *h.ArticleKey},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid article")
	}

	return nil
}

func (h *createHandler) checkACLExist(ctx context.Context) error {
	exist, err := aclmwcli.ExistACLConds(ctx, &aclmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		RoleID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.RoleID},
		ArticleKey: &basetypes.StringVal{Op: cruder.EQ, Value: *h.ArticleKey},
	})
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("acl exist")
	}

	return nil
}

func (h *Handler) CreateACL(ctx context.Context) (*npool.ACL, error) {
	handler := &createHandler{
		Handler: h,
	}

	if err := handler.checkRoleExist(ctx); err != nil {
		return nil, err
	}
	if err := handler.checkArticleExist(ctx); err != nil {
		return nil, err
	}
	if err := handler.checkACLExist(ctx); err != nil {
		return nil, err
	}

	info, err := aclmwcli.CreateACL(ctx, &aclmwpb.ACLReq{
		AppID:      h.AppID,
		RoleID:     h.RoleID,
		ArticleKey: h.ArticleKey,
	})
	if err != nil {
		return nil, err
	}

	return h.GetACL(ctx, info)
}
