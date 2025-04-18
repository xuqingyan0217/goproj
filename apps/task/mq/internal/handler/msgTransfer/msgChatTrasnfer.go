/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/apps/task/mq/mq"
	"easy-chat/pkg/bitmap"
)

type MsgChatTransfer struct {
	*baseMsgTransfer
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svc),
	}
}

func (m *MsgChatTransfer) Consume(key, value string) error {
	fmt.Println("key : ", key, " value : ", value)

	var (
		data mq.MsgChatTransfer
		ctx  = context.Background()
		msgId = primitive.NewObjectID()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 记录数据
	if err := m.addChatLog(ctx, msgId, &data); err != nil {
		return err
	}

	// 推送消息

	return m.Transfer(ctx,&ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MType:          data.MType,
		MsgId:          msgId.Hex(),
		// MsgId: 			data.MsgId,
		Content:        data.Content,
	})
}

/*// 私聊
func (m *MsgChatTransfer) single(data *mq.MsgChatTransfer) error {
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

// 群聊
func (m *MsgChatTransfer) group(ctx context.Context, data *mq.MsgChatTransfer) error {
	// 获取到群聊中的所有用户
	users, err := m.svc.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
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
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}*/

func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	// 记录消息
	chatLog := immodels.ChatLog{
		ID: 			msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}
	// 这里在记录消息后，把该消息对应的位图里的发送方用户改为已读，发送方都未读，那就离谱了
	// 群聊消息和私聊消息都是在一个数据库里，
	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(chatLog.SendId)
	chatLog.ReadRecords = readRecords.Export()

	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}

	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
}
