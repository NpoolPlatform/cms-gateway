package article

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/cms-gateway/pkg/const"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/cms/v1"
	articlemw "github.com/NpoolPlatform/message/npool/cms/mw/v1/article"

	"github.com/google/uuid"
)

type Handler struct {
	ID            *uint32
	EntID         *string
	Host          *string
	ArticleKey    *string
	AppID         *string
	UserID        *string
	LangID        *string
	CategoryID    *string
	Title         *string
	Subtitle      *string
	Digest        *string
	Content       *string
	ISO           *string
	UpdateContent *bool
	ContentURL    *string
	Version       *uint32
	Latest        *bool
	Status        *basetypes.ArticleStatus
	Reqs          []*articlemw.ArticleReq
	Offset        int32
	Limit         int32
}

func NewHandler(ctx context.Context, options ...func(context.Context, *Handler) error) (*Handler, error) {
	handler := &Handler{}
	for _, opt := range options {
		if err := opt(ctx, handler); err != nil {
			return nil, err
		}
	}
	return handler, nil
}

func WithID(id *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid id")
			}
			return nil
		}
		h.ID = id
		return nil
	}
}

func WithEntID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid entid")
			}
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.EntID = id
		return nil
	}
}

func WithHost(host *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if host == nil {
			if must {
				return fmt.Errorf("invalid host")
			}
			return nil
		}
		if *host == "" {
			return fmt.Errorf("invalid host")
		}
		h.Host = host
		return nil
	}
}

func WithAppID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid appid")
			}
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.AppID = id
		return nil
	}
}

func WithLangID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid langid")
			}
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.LangID = id
		return nil
	}
}

func WithCategoryID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid categoryid")
			}
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.CategoryID = id
		return nil
	}
}

func WithUserID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid userid")
			}
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.UserID = id
		return nil
	}
}

func WithArticleKey(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid articlekey")
			}
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.ArticleKey = id
		return nil
	}
}

func WithTitle(title *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if title == nil {
			if must {
				return fmt.Errorf("invalid title")
			}
			return nil
		}
		if *title == "" {
			return fmt.Errorf("invalid title")
		}
		h.Title = title
		return nil
	}
}

func WithSubtitle(subtitle *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if subtitle == nil {
			if must {
				return fmt.Errorf("invalid subtitle")
			}
			return nil
		}
		if *subtitle == "" {
			return fmt.Errorf("invalid subtitle")
		}
		h.Subtitle = subtitle
		return nil
	}
}

func WithDigest(digest *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if digest == nil {
			if must {
				return fmt.Errorf("invalid digest")
			}
			return nil
		}
		if *digest == "" {
			return fmt.Errorf("invalid digest")
		}
		h.Digest = digest
		return nil
	}
}

func WithContent(content *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if content == nil {
			if must {
				return fmt.Errorf("invalid content")
			}
			return nil
		}
		if *content == "" {
			return fmt.Errorf("invalid content")
		}
		h.Content = content
		return nil
	}
}

func WithUpdateContent(updatecontent *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if updatecontent == nil {
			if must {
				return fmt.Errorf("invalid updatecontent")
			}
			return nil
		}
		h.UpdateContent = updatecontent
		return nil
	}
}

func WithLatest(latest *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if latest == nil {
			if must {
				return fmt.Errorf("invalid latest")
			}
			return nil
		}
		h.Latest = latest
		return nil
	}
}

func WithStatus(status *basetypes.ArticleStatus, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if status == nil {
			if must {
				return fmt.Errorf("invalid status")
			}
			return nil
		}
		switch *status {
		case basetypes.ArticleStatus_Draft:
		case basetypes.ArticleStatus_Published:
		default:
			return fmt.Errorf("invalid status")
		}

		h.Status = status
		return nil
	}
}

func WithContentURL(url *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if url == nil {
			if must {
				return fmt.Errorf("invalid url")
			}
			return nil
		}
		if *url == "" {
			return fmt.Errorf("invalid url")
		}
		h.ContentURL = url
		return nil
	}
}

func WithOffset(offset int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Offset = offset
		return nil
	}
}

func WithLimit(limit int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if limit == 0 {
			limit = constant.DefaultRowLimit
		}
		h.Limit = limit
		return nil
	}
}
