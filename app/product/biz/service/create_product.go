package service

import (
	"context"
	"fmt"
	"net/url"

	product "gomall/rpc_gen/kitex_gen/product"

	"gomall/app/product/biz/dal/mysql"
	"gomall/app/product/biz/dal/redis"
	"gomall/app/product/biz/model"
)

type CreateProductService struct {
	ctx context.Context
} // NewCreateProductService new CreateProductService
func NewCreateProductService(ctx context.Context) *CreateProductService {
	return &CreateProductService{ctx: ctx}
}

// Run create note info
func (s *CreateProductService) Run(req *product.CreateProductReq) (resp *product.CreateProductResp, err error) {
	// Finish your business logic.
	productMutation := model.NewProductMutation(s.ctx, mysql.DB, redis.RedisClient)

	categories := make([]model.Category, 0, len(req.Product.Categories))
	for _, categoryName := range req.Product.Categories {
		categories = append(categories, model.Category{
			Name: categoryName,
		})
	}
	if req.Product.Picture == "" {
		return nil, fmt.Errorf("商品图片链接不能为空")
	}
	_, err = url.Parse(req.Product.Picture)
	if err != nil {
		return nil, fmt.Errorf("无效的图片URL")
	}

	row, err := productMutation.CreateProduct(model.Product{
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Price:       req.Product.Price,
		Picture:     req.Product.Picture,
		Categories:  categories,
	})
	if err != nil {
		return nil, err
	}
	resp = &product.CreateProductResp{Id: uint32(row.ID)}
	return resp, nil
}
