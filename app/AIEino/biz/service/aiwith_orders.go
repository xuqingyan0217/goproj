package service

import (
	"context"
	"github.com/cloudwego/eino/components/tool"
	mytools "gomall/app/AIEino/biz/eino/tools"
	"gomall/rpc_gen/kitex_gen/AIEino"
)

type AIWithOrdersService struct {
	ctx context.Context
} // NewAIWithOrdersService new AIWithOrdersService
func NewAIWithOrdersService(ctx context.Context) *AIWithOrdersService {
	return &AIWithOrdersService{ctx: ctx}
}

// Run create note info
func (s *AIWithOrdersService) Run(req *AIEino.AIWithOrdersReq) (res *AIEino.AIWithOrdersResp, err error) {
	// Finish your business logic.
	// 组装tool工具
	getOrderTool := &mytools.GetOrderTool{}
	tools := []tool.BaseTool{
		getOrderTool,
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
	return &AIEino.AIWithOrdersResp{
		Orders: temp,
	}, nil
}
