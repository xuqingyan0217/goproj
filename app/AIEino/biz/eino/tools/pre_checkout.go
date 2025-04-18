package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"gomall/app/AIEino/infra/rpc"
	"gomall/rpc_gen/kitex_gen/checkout"
	"gomall/rpc_gen/kitex_gen/product"
	"log"
)

type PreCheckoutParams struct {
	Query string `json:"query"`
	// 目前实现的是查询一类商品，每类商品数量都设置为同一个值；因为其他类型不知道为什么不支持
	ProductCount uint32 `json:"product_count"`

	UserId        uint32 `json:"user_id"`
	Email         string `json:"email"`
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
}

// ProductInfo 结构体，用于打包返回值
type ProductInfo struct {
	ProductId uint32 `json:"product_id"`
	Quantity  uint32 `json:"quantity"`
}

type PreCheckoutTool struct{}

func (ft *PreCheckoutTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "pre_checkout_tool",
		Desc: "用于依据文本输入实现自动预下单功能",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {
				Desc:     "商品名称，作为查询商品的条件，从文本中获取",
				Type:     schema.String,
				Required: true,
			},
			"product_count": {
				Desc:     "商品购买数量，从文本中获取",
				Type:     schema.Integer,
				Required: true,
			},
			"user_id": {
				Desc:     "预下单用户的用户id，从文本中获取",
				Type:     schema.Integer,
				Required: true,
			},
			"email": {
				Desc:     "预下单用户的邮箱，从文本中获取",
				Type:     schema.String,
				Required: true,
			},
			"street_address": {
				Desc:     "街道信息，从文本中获取",
				Type:     schema.String,
				Required: true,
			},
			"city": {
				Desc:     "城市信息，从文本中获取",
				Type:     schema.String,
				Required: true,
			},
			"state": {
				Desc:     "州信息，从文本中获取",
				Type:     schema.String,
				Required: true,
			},
			"country": {
				Desc:     "国家信息，从文本中获取",
				Type:     schema.String,
				Required: true,
			},
			"zip_code": {
				Desc:     "邮政编码，从文本中获取",
				Type:     schema.String,
				Required: true,
			},
		}),
	}, nil
}

func (ft *PreCheckoutTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	var params PreCheckoutParams
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		log.Printf("Failed to unmarshal arguments: %v", err)
		return "", err
	}
	// 验证参数的有效性
	if params.Query == "" {
		log.Printf("name is required")
		return "", fmt.Errorf("name is required")
	}
	if params.UserId == 0 {
		log.Printf("user_id is required")
		return "", fmt.Errorf("user_id is required")
	}
	if params.Email == "" {
		log.Printf("email is required")
		return "", fmt.Errorf("email is required")
	}
	if params.StreetAddress == "" {
		log.Printf("street_address is required")
		return "", fmt.Errorf("street_address is required")
	}
	if params.City == "" {
		log.Printf("city is required")
		return "", fmt.Errorf("city is required")
	}
	if params.State == "" {
		log.Printf("state is required")
		return "", fmt.Errorf("state is required")
	}
	if params.Country == "" {
		log.Printf("country is required")
		return "", fmt.Errorf("country is required")
	}
	if params.ZipCode == "" {
		log.Printf("zip_code is required")
		return "", fmt.Errorf("zip_code is required")
	}
	// 如果没有从文本中提取到数量，就默认为1件
	if params.ProductCount == 0 {
		params.ProductCount = 1
	}
	log.Println(params)
	// 在此处添加工具的具体逻辑
	var productInfoList []ProductInfo
	productSearchResp, err := rpc.ProductClient.SearchProducts(ctx, &product.SearchProductsReq{
		Query: params.Query,
	})
	if err != nil {
		return "", err
	}
	var productResult = productSearchResp.Results
	for _, pt := range productResult {
		productInfoList = append(productInfoList, ProductInfo{
			ProductId: pt.Id,
			Quantity:  params.ProductCount,
		})
	}
	for _, productItem := range productInfoList {
		log.Printf("product_id: %d, quantity: %d", productItem.ProductId, productItem.Quantity)
	}

	preCheckoutResp, err := ft.Helper(ctx, params, productInfoList)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(preCheckoutResp)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (ft *PreCheckoutTool) Helper(ctx context.Context, params PreCheckoutParams, proList []ProductInfo) (*checkout.PreCheckoutResp, error) {
	// 获取产品信息列表
	var productList []*checkout.ProductInfo
	for _, productItem := range proList {
		productList = append(productList, &checkout.ProductInfo{
			ProductId: productItem.ProductId,
			Quantity:  productItem.Quantity,
		})
	}
	// 调用rpc方法开始ai下单
	preCheckoutResp, err := rpc.CheckOutClient.PreCheckout(ctx, &checkout.PreCheckoutReq{
		UserId: params.UserId,
		Email:  params.Email,
		Address: &checkout.Address{
			StreetAddress: params.StreetAddress,
			City:          params.City,
			State:         params.State,
			Country:       params.Country,
			ZipCode:       params.ZipCode,
		},
		ProductInfoList: productList,
	})
	if err != nil {
		return nil, err
	}
	return &checkout.PreCheckoutResp{
		PreOrderId:  preCheckoutResp.PreOrderId,
		TotalAmount: preCheckoutResp.TotalAmount,
		ValidUntil:  preCheckoutResp.ValidUntil,
	}, nil
}
