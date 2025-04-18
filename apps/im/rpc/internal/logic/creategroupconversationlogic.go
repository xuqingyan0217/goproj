package logic

import (
	"context"

	"github.com/pkg/errors"

	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/svc"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupConversationLogic) CreateGroupConversation(in *im.CreateGroupConversationReq) (*im.CreateGroupConversationResp, error) {
	// todo: add your logic here and delete this line
	// 完成对群会话的验证，先判断群聊会话是否存在
	res := &im.CreateGroupConversationResp{}
	// 返回的错误有三种情况，错误存在和不存在，其中存在又分为了两种：查询错误和查不到
	_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, in.GroupId)
	// 如果为空，则说明查询是没有问题的
	if err == nil {
		return res, nil
	}
	// 如果错误不是群聊还未创建，则属于查询错误
	if err != immodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(),"ConversationModel.FindOne err %v req %v", err,in.GroupId)
	}

	// 到此则说明群聊查询结果是查不到，我们准备新建群聊会话
	err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
		ConversationId: in.GroupId,
		ChatType:       constants.GroupChatType,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(),"ConversationModel.FindOne err %v ", err)
	}

	// 建立会话关系，此处直接引用之前创建会话的方法就行了，只不过要在之前那个方法里面的switch里新增群会话
	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&im.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})

	return res, nil
}
