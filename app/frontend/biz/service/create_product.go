package service

import (
	"context"

	common "gomall/app/frontend/hertz_gen/frontend/common"
	product "gomall/app/frontend/hertz_gen/frontend/product"
	"gomall/app/frontend/infra/rpc"
	rpcproduct "gomall/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type CreateProductService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCreateProductService(Context context.Context, RequestContext *app.RequestContext) *CreateProductService {
	return &CreateProductService{RequestContext: RequestContext, Context: Context}
}

func (h *CreateProductService) Run(req *product.CreateProductReq) (resp *common.Empty, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	_, err = rpc.ProductClient.CreateProduct(h.Context, &rpcproduct.CreateProductReq{
		Product: &rpcproduct.Product{
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
