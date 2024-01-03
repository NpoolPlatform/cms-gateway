package acl

import (
	"context"

	aclmwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/acl"
	aclmwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/acl"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetACLs(ctx context.Context) ([]*aclmwpb.ACL, uint32, error) {
	infos, total, err := aclmwcli.GetACLs(ctx, &aclmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit,
	)
	if err != nil {
		return nil, 0, err
	}

	return infos, total, nil
}
