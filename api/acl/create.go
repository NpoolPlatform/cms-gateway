package acl

import (
	"context"
	"fmt"

	acl1 "github.com/NpoolPlatform/cms-gateway/pkg/acl"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/cms/gw/v1/acl"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateACL(ctx context.Context, in *npool.CreateACLRequest) (*npool.CreateACLResponse, error) {
	handler, err := acl1.NewHandler(
		ctx,
		acl1.WithAppID(&in.AppID, true),
		acl1.WithRoleID(&in.RoleID, true),
		acl1.WithArticleKey(&in.ArticleKey, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateACL",
			"In", in,
			"Error", err,
		)
		return &npool.CreateACLResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateACL(ctx)
	fmt.Println("info: ", info)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateACL",
			"In", in,
			"Error", err,
		)
		return &npool.CreateACLResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateACLResponse{
		Info: info,
	}, nil
}
