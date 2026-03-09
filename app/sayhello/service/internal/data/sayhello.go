package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/ray-dota/backend-mono/app/sayhello/service/internal/biz"
)

type sayhelloRepo struct {
	data *Data
	log  *log.Helper
}

// NewSayhelloRepo .
func NewSayhelloRepo(data *Data, logger log.Logger) biz.SayhelloRepo {
	return &sayhelloRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *sayhelloRepo) Save(ctx context.Context, g *biz.Sayhello) (*biz.Sayhello, error) {
	return g, nil
}

func (r *sayhelloRepo) Update(ctx context.Context, g *biz.Sayhello) (*biz.Sayhello, error) {
	return g, nil
}

func (r *sayhelloRepo) FindByID(_ context.Context, _ int64) (*biz.Sayhello, error) {
	return nil, biz.ErrUserNotFound
}

func (r *sayhelloRepo) ListByHello(context.Context, string) ([]*biz.Sayhello, error) {
	return nil, nil
}

func (r *sayhelloRepo) ListAll(context.Context) ([]*biz.Sayhello, error) {
	return nil, nil
}
