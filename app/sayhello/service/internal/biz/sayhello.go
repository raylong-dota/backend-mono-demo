package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"

	v1 "github.com/ray-dota/backend-mono/api/sayhello/service/v1"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Sayhello is a Sayhello model.
type Sayhello struct {
	Hello string
}

// SayhelloRepo is a Greater repo.
type SayhelloRepo interface {
	Save(context.Context, *Sayhello) (*Sayhello, error)
	Update(context.Context, *Sayhello) (*Sayhello, error)
	FindByID(context.Context, int64) (*Sayhello, error)
	ListByHello(context.Context, string) ([]*Sayhello, error)
	ListAll(context.Context) ([]*Sayhello, error)
}

// SayhelloUsecase is a Sayhello usecase.
type SayhelloUsecase struct {
	repo SayhelloRepo
	log  *log.Helper
}

// NewSayhelloUsecase new a Sayhello usecase.
func NewSayhelloUsecase(repo SayhelloRepo, logger log.Logger) *SayhelloUsecase {
	return &SayhelloUsecase{repo: repo, log: log.NewHelper(log.With(logger, "module", "sayhello"))}
}

// CreateSayhello creates a Sayhello, and returns the new Sayhello.
func (uc *SayhelloUsecase) CreateSayhello(ctx context.Context, g *Sayhello) (*Sayhello, error) {
	log.Infof("CreateSayhello: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}
