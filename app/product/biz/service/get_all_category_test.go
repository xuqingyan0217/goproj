package service

import (
	"context"
	"testing"
	product "gomall/rpc_gen/kitex_gen/product"
)

func TestGetAllCategory_Run(t *testing.T) {
	ctx := context.Background()
	s := NewGetAllCategoryService(ctx)
	// init req and assert value

	req := &product.GetAllCategoryReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
