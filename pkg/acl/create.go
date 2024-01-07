package acl

import (
	"context"
	"fmt"

	aclmwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/acl"
	articlemwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/article"
	aclmwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/acl"
	articlemwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/article"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/acl"
)

func (h *Handler) CreateACL(ctx context.Context) (*npool.ACL, error) {
	exist, err := articlemwcli.ExistArticleConds(ctx, &articlemwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		ArticleKey: &basetypes.StringVal{Op: cruder.EQ, Value: *h.ArticleKey},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid article")
	}

	info, err := aclmwcli.GetACLOnly(ctx, &aclmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		RoleID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.RoleID},
		ArticleKey: &basetypes.StringVal{Op: cruder.EQ, Value: *h.ArticleKey},
	})
	if err != nil {
		return nil, err
	}
	if info != nil {
		return h.GetACL(ctx, info)
	}

	info, err = aclmwcli.CreateACL(ctx, &aclmwpb.ACLReq{
		AppID:      h.AppID,
		RoleID:     h.RoleID,
		ArticleKey: h.ArticleKey,
	})
	if err != nil {
		return nil, err
	}

	return h.GetACL(ctx, info)
}
