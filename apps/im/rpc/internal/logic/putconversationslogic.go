package logic

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/svc"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/xerr"
)

type PutConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新会话
func (l *PutConversationsLogic) PutConversations(in *im.PutConversationsReq) (*im.PutConversationsResp, error) {
	// todo: add your logic here and delete this line
	// 先查询到该会话，然后才能更新；该会话一定是存在的，但是会话里面有没有聊天会话就不好说了，所以说判断一手
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		// 我们会在用户建立会话的时候，两者之间都会有对应的会话记录列表，所以不会为空
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindByUserId err %v, req %v", err, in.UserId)
	}
	// 存在会话，但会话列表可能为空
	if data.ConversationList == nil {
		data.ConversationList = make(map[string]*immodels.Conversation)
	}
	// 遍历会话map
	for s, conversation := range in.ConversationList {
		// 记录历史已读消息
		var oldTotal int
		if data.ConversationList[s] != nil {
			// 不为空说明已读会话中有值，修改oldTotal为该值作为历史已读记录
			oldTotal = data.ConversationList[s].Total
		}
		// 更新conversations里的会话列表为最新值
		data.ConversationList[s] = &immodels.Conversation{
			ConversationId: conversation.ConversationId,
			ChatType:       constants.ChatType(conversation.ChatType),
			IsShow:         conversation.IsShow,
			// 已读的消息量和原本的消息量，加起来就是实实在在的消息总理
			Total:          int(conversation.Read) + oldTotal,
			Seq:            conversation.Seq,
		}
	}
	// 更新，此处如果没有，这个update就相当于执行了一个插入操作
	err = l.svcCtx.ConversationsModel.Update(l.ctx, data)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Update err %v, req %v", err, data)
	}
	
	return &im.PutConversationsResp{}, nil
}
