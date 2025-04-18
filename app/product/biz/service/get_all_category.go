package service

import (
	"context"
	"gomall/app/product/biz/model"
	"gomall/app/product/biz/dal/mysql"
	product "gomall/rpc_gen/kitex_gen/product"
)

type GetAllCategoryService struct {
	ctx context.Context
} // NewGetAllCategoryService new GetAllCategoryService
func NewGetAllCategoryService(ctx context.Context) *GetAllCategoryService {
	return &GetAllCategoryService{ctx: ctx}
}

// Run create note info
func (s *GetAllCategoryService) Run(req *product.GetAllCategoryReq) (resp *product.GetAllCategoryResp, err error) {
	// Finish your business logic.
	categoryQuery := model.NewCategoryQuery(s.ctx, mysql.DB)
	categoryNames, err := categoryQuery.GetAllCategoryNames()
	if err != nil {
		return nil, err
	}
	return &product.GetAllCategoryResp{Categories: categoryNames}, nil
}
