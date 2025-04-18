package service

import (
	"context"

	product "gomall/rpc_gen/kitex_gen/product"

	"gomall/app/product/biz/dal/mysql"
	"gomall/app/product/biz/dal/redis"
	"gomall/app/product/biz/model"
)

type UpdateProductService struct {
	ctx context.Context
} // NewUpdateProductService new UpdateProductService
func NewUpdateProductService(ctx context.Context) *UpdateProductService {
	return &UpdateProductService{ctx: ctx}
}

// Run create note info
func (s *UpdateProductService) Run(req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	// Finish your business logic.
	productMutation := model.NewProductMutation(s.ctx, mysql.DB, redis.RedisClient)
	categories := make([]model.Category, 0, len(req.Product.Categories))
	for _, categoryName := range req.Product.Categories {
		categories = append(categories, model.Category{
			Name: categoryName,
		})
	}
	err = productMutation.UpdateProduct(int(req.Product.Id), model.Product{
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Price:       req.Product.Price,
		Picture:     req.Product.Picture,
		Categories:  categories,
	})
	if err != nil {
		return nil, err
	}
	resp = &product.UpdateProductResp{Id: req.Product.Id}
	return resp, nil
}
