package service

import (
	"context"
	"testing"

	product "gomall/rpc_gen/kitex_gen/product"
)

func TestCreateProduct_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCreateProductService(ctx)
	// init req and assert value

	req := &product.CreateProductReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
