package service

import (
	"context"
	"gomall/app/product/biz/dal/mysql"
	"gomall/app/product/biz/model"
	product "gomall/rpc_gen/kitex_gen/product"
)

type ListProductsService struct {
	ctx context.Context
} // NewListProductsService new ListProductsService
func NewListProductsService(ctx context.Context) *ListProductsService {
	return &ListProductsService{ctx: ctx}
}

// Run create note info
func (s *ListProductsService) Run(req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	// Finish your business logic.
	categoryQuery := model.NewCategoryQuery(s.ctx, mysql.DB)
	c, err := categoryQuery.GetProductsByCategoryName(req.CategoryName)
	if err != nil {
		return nil, err
	}
	resp = &product.ListProductsResp{}
	// 使用map来存储已处理过的商品ID，避免重复
	processedProducts := make(map[int]bool)

	for _, v1 := range c {
		for _, v := range v1.Products {
			// 检查商品是否已经被处理过
			if _, exists := processedProducts[v.ID]; exists {
				continue
			}
			// 标记商品为已处理
			processedProducts[v.ID] = true

			categories, err := categoryQuery.GetCategoriesByProductId(v.ID)
			if err != nil {
				return nil, err
			}
			resp.Products = append(resp.Products, &product.Product{
				Id:          uint32(v.ID),
				Name:        v.Name,
				Description: v.Description,
				Picture:     v.Picture,
				Price:       v.Price,
				Categories:  categories,
			})
		}
	}
	return resp, nil
}
