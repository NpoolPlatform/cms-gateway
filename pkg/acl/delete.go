package acl

import (
	"context"
	"fmt"

	aclmwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/acl"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/acl"
	aclmwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/acl"
)

func (h *Handler) DeleteACL(ctx context.Context) (*npool.ACL, error) {
	exist, err := aclmwcli.ExistACLConds(ctx, &aclmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid acl")
	}

	info, err := aclmwcli.DeleteACL(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	return h.GetACL(ctx, info)
}
