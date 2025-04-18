package service

import (
	"context"

	common "gomall/app/frontend/hertz_gen/frontend/common"
	product "gomall/app/frontend/hertz_gen/frontend/product"
	"gomall/app/frontend/infra/rpc"
	rpcproduct "gomall/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateService(Context context.Context, RequestContext *app.RequestContext) *UpdateService {
	return &UpdateService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateService) Run(req *product.UpdateProductReq) (resp *common.Empty, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	_, err = rpc.ProductClient.UpdateProduct(h.Context, &rpcproduct.UpdateProductReq{
		Id: req.Id,
		Product: &rpcproduct.Product{
			Id:          req.Id,
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Picture:     req.Picture,
			Categories:  req.Categories,
		},
	})
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
