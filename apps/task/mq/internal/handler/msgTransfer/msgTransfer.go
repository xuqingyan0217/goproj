package msgTransfer

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/pkg/constants"
)

type baseMsgTransfer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewBaseMsgTransfer 初始化
func NewBaseMsgTransfer(svc *svc.ServiceContext) *baseMsgTransfer {
	return &baseMsgTransfer{
		svcCtx: svc,
		Logger: logx.WithContext(context.Background()),
	}
}

// Transfer 该方法专门用于消息转发
func (b *baseMsgTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	var err error
	switch data.ChatType {
	case constants.GroupChatType:
		err = b.group(ctx, data)
	case constants.SingleChatType:
		err = b.single(ctx, data)
	}
	return err
}

// single 私聊的转发
func (b *baseMsgTransfer) single(ctx context.Context ,data *ws.Push) error {
	return b.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

// group 群聊
func (b *baseMsgTransfer) group(ctx context.Context, data *ws.Push) error {
	// 获取到群聊中的所有用户
	users, err := b.svcCtx.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})
	if err != nil {
		return err
	}

	// 初始化接收用户集
	data.RecvIds = make([]string, 0, len(users.List))
	// 迭代用户，获取到用户id，注意，自己的要过滤到
	for _, members := range users.List {
		if members.UserId == data.SendId{
			continue
		}
		// 将其余用户放到刚初始化的集合里面
		data.RecvIds = append(data.RecvIds, members.UserId)
	}
	// 获取到了id集合后，准备推送
	return b.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}