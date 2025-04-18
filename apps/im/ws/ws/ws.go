/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package ws

import "easy-chat/pkg/constants"

/* 因为在websocket的message.go里面的消息结构体的data字段是interface{}类型，
   而interface{}类型json序列化后就是map类型，而我们想要的是string类型，
   mapstructure 就是来做这个处理的，把map转换为string，方便数据流通，
   所以在接收如下数据结构时，需要用mapstructure进行序列化转换
   我们目前使用的mapstructure针对的场景就是在服务端获取客户端发来的Message消息结构之后需要使用
   该mapstructure将用户输入消息转到下面的消息结构中，详细看push.go，conversation.go
   而在mq那边的数据结构，就采用json了，因为mq是获取的服务端序列化之后的消息，也就是string类型了
   mq处理后消息也是给服务端，这个过程mq没有重客户端获取消息，所以说是json
   而在最初的user.online那个路由，因为是不需要使用新的数据结构去接收它的，所以那里直接转map了也没事 */
type (
	Msg struct {
		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
		MsgId 		string `mapstructure:"msgId"`
		ReadRecords  map[string]string `mapstructure:"readRecords"`
	}

	Chat struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`
		Msg                `mapstructure:"msg"`
	}

	Push struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		RecvIds			   []string `mapstructure:"recvIds"`
		SendTime           int64  `mapstructure:"sendTime"`

		MsgId 		string `mapstructure:"msgId"`

		ReadRecords  map[string]string `mapstructure:"readRecords"`
		constants.ContentType `mapstructure:"contentType"`

		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}

	MarkRead struct {
		constants.ChatType `mapstructure:"chatType"`
		RecvId             string `mapstructure:"recvId"`
		ConversationId     string `mapstructure:"conversationId"`
		MsgIds			   []string `mapstructure:"msgIds"`
	}
)
