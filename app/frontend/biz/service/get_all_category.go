package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	common "gomall/app/frontend/hertz_gen/frontend/common"
	"gomall/app/frontend/infra/rpc"
	"gomall/rpc_gen/kitex_gen/product"
)

type GetAllCategoryService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetAllCategoryService(Context context.Context, RequestContext *app.RequestContext) *GetAllCategoryService {
	return &GetAllCategoryService{RequestContext: RequestContext, Context: Context}
}

func (h *GetAllCategoryService) Run(req *common.Empty) (resp map[string]interface{}, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()

	// 调用product服务获取所有类别
	categories, err := rpc.ProductClient.GetAllCategory(h.Context, &product.GetAllCategoryReq{})
	if err != nil {
		return nil, err
	}

	return utils.H{
		"categories": categories.Categories,
	}, nil
}
