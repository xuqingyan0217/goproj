/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package push

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/pkg/constants"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	fmt.Println("000++++++++++++++++++++++++++++++++++++++++++++++++++")
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		fmt.Println("001++++++++++++++++++++++++++++++++++++++++++++++++++")
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		fmt.Println("002++++++++++++++++++++++++++++++++++++++++++++++++++")

		// 发送的目标
		switch data.ChatType {
		case constants.SingleChatType:
			fmt.Println("私聊推送---")
			single(srv, &data, data.RecvId)
			fmt.Println("004++++++++++++++++++++++++++++++++++++++++++++++++++")
		case constants.GroupChatType:
			fmt.Println("群聊推送---")
			group(srv, &data)
			fmt.Println("005++++++++++++++++++++++++++++++++++++++++++++++++++")
		}
	}
}

// 封装一个函数，存放上面的私聊push
func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	// 发送的目标
	rconn := srv.GetConn(recvId)
	if rconn == nil {
		// todo: 目标离线
		fmt.Println("0000++++++++++++++++++++++++++++++++++++++++++++++++++")
		return nil
	}

	fmt.Println("003++++++++++++++++++++++++++++++++++++++++++++++++++")

	srv.Infof("push msg %v", data)

	return srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Msg: ws.Msg{
			MType:   data.MType,
			Content: data.Content,
			MsgId: 	 data.MsgId,
			ReadRecords: data.ReadRecords,
		},
	}), rconn)
}

// 封装一个函数，用于处理群聊
func group(srv *websocket.Server, data *ws.Push) error {
	for _, id := range data.RecvIds {
		func(id string) {
			single(srv, data, id)
		}(id)
	}
	return nil
}