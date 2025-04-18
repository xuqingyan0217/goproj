package service

import (
	"context"
	"github.com/cloudwego/eino/components/tool"
	mytools "gomall/app/AIEino/biz/eino/tools"
	"gomall/rpc_gen/kitex_gen/AIEino"
)

type AIWithPreCheckoutService struct {
	ctx context.Context
} // NewAIWithPreCheckoutService new AIWithPreCheckoutService
func NewAIWithPreCheckoutService(ctx context.Context) *AIWithPreCheckoutService {
	return &AIWithPreCheckoutService{ctx: ctx}
}

// Run create note info
func (s *AIWithPreCheckoutService) Run(req *AIEino.AIWithPreCheckoutReq) (res *AIEino.AIWithPreCheckoutResp, err error) {
	// Finish your business logic.
	preCheckoutTool := &mytools.PreCheckoutTool{}
	tools := []tool.BaseTool{
		preCheckoutTool,
	}
	// 调用ai大模型逻辑，传入需要的参数
	resp, err := mytools.LogicAI(s.ctx, tools, req.UserInput)
	if err != nil {
		return nil, err
	}
	// 组装合理的返回值
	var temp []string
	for _, msg := range resp {
		temp = append(temp, msg.Content)
	}
	return &AIEino.AIWithPreCheckoutResp{
		PreCheckoutRes: temp,
	}, nil
}
