package acl

import (
	"context"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	aclmwcli "github.com/NpoolPlatform/cms-middleware/pkg/client/acl"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/acl"
	aclmwpb "github.com/NpoolPlatform/message/npool/cms/mw/v1/acl"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetACLs(ctx context.Context) ([]*npool.ACL, uint32, error) {
	rows, total, err := aclmwcli.GetACLs(ctx, &aclmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		ArticleKey: &basetypes.StringVal{Op: cruder.EQ, Value: *h.ArticleKey},
	}, h.Offset, h.Limit,
	)
	if err != nil {
		return nil, 0, err
	}

	roleIDs := []string{}
	for _, val := range rows {
		if val.RoleID != "" {
			roleIDs = append(roleIDs, val.RoleID)
		}
	}

	roleMap := map[string]*rolemwpb.Role{}
	if len(roleIDs) > 0 {
		roleInfos, _, err := rolemwcli.GetRoles(ctx, &rolemwpb.Conds{
			EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: roleIDs},
		}, 0, int32(len(roleIDs)))
		if err != nil {
			return nil, 0, err
		}

		for _, val := range roleInfos {
			roleMap[val.EntID] = val
		}
	}

	infos := []*npool.ACL{}
	for _, val := range rows {
		role, ok := roleMap[val.RoleID]
		if !ok {
			continue
		}
		infos = append(infos, &npool.ACL{
			ID:         val.ID,
			EntID:      val.EntID,
			AppID:      val.AppID,
			Role:       role.Role,
			RoleID:     val.RoleID,
			ArticleKey: val.ArticleKey,
			CreatedAt:  val.CreatedAt,
			UpdatedAt:  val.UpdatedAt,
		})
	}
	return infos, total, nil
}

func (h *Handler) GetACL(ctx context.Context, row *aclmwpb.ACL) (*npool.ACL, error) {
	info := &npool.ACL{
		ID:         row.ID,
		EntID:      row.EntID,
		AppID:      row.AppID,
		RoleID:     row.RoleID,
		ArticleKey: row.ArticleKey,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}

	role, err := rolemwcli.GetRoleOnly(ctx, &rolemwpb.Conds{
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: row.RoleID},
	})
	if err != nil {
		return nil, err
	}
	if role != nil {
		return info, nil
	}
	info.Role = role.Role

	return info, nil
}
