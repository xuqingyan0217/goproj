package service

import (
	"context"

	product "gomall/rpc_gen/kitex_gen/product"

	"gomall/app/product/biz/dal/mysql"
	"gomall/app/product/biz/dal/redis"
	"gomall/app/product/biz/model"
)

type DeleteProductService struct {
	ctx context.Context
} // NewDeleteProductService new DeleteProductService
func NewDeleteProductService(ctx context.Context) *DeleteProductService {
	return &DeleteProductService{ctx: ctx}
}

// Run create note info
func (s *DeleteProductService) Run(req *product.DeleteProductReq) (resp *product.DeleteProductResp, err error) {
	// Finish your business logic.
	productMutation := model.NewProductMutation(s.ctx, mysql.DB, redis.RedisClient)
	err = productMutation.DeleteProduct(int(req.Id))
	if err != nil {
		return nil, err
	}
	resp = &product.DeleteProductResp{Id: req.Id}
	return resp, nil
}
