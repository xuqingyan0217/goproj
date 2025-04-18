package service

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"gomall/app/product/biz/dal/mysql"
	"gomall/app/product/biz/dal/redis"
	"gomall/app/product/biz/model"
	product "gomall/rpc_gen/kitex_gen/product"
)

type GetProductService struct {
	ctx context.Context
} // NewGetProductService new GetProductService
func NewGetProductService(ctx context.Context) *GetProductService {
	return &GetProductService{ctx: ctx}
}

// Run create note info
func (s *GetProductService) Run(req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	// Finish your business logic.
	if req.Id == 0 {
		return nil, kerrors.NewGRPCBizStatusError(2004001, "product_id param error")
	}
	productQuery := model.NewCachedProductQuery(s.ctx, mysql.DB, redis.RedisClient)
	p, err := productQuery.GetById(int(req.Id))
	if err != nil {
		return nil, err
	}
	return &product.GetProductResp{
		Product: &product.Product{
			Id:          uint32(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Picture:     p.Picture,
		},
	}, nil
}
