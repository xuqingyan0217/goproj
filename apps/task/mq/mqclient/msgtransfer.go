/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package mqclient

import (
	"encoding/json"

	"github.com/zeromicro/go-queue/kq"

	"easy-chat/apps/task/mq/mq"
)

// 消息转发客户端

type MsgChatTransferClient interface {
	Push(msg *mq.MsgChatTransfer) error
}


type msgChatTransferClient struct {
	pusher *kq.Pusher
}

func NewMsgChatTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgChatTransferClient {
	return &msgChatTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

func (c *msgChatTransferClient) Push(msg *mq.MsgChatTransfer) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.pusher.Push(string(body))
}

// 消息已读未读处理客户端

type MsgReadTransferClient interface {
	Push(msg *mq.MsgMarkRead) error
}

type msgReadTransferClient struct {
	pusher *kq.Pusher
}

func NewMsgReadTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgReadTransferClient {
	return &msgReadTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

func (c *msgReadTransferClient) Push(msg *mq.MsgMarkRead) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.pusher.Push(string(body))
}
