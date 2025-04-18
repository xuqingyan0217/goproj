/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"

	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/task/mq/internal/config"
	"easy-chat/pkg/constants"
	"net/http"
)

type ServiceContext struct {
	config.Config

	WsClient websocket.Client
	*redis.Redis

	immodels.ChatLogModel
	immodels.ConversationModel

	socialclient.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ChatLogModel:      immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel: immodels.MustConversationModel(c.Mongo.Url, c.Mongo.Db),

		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}

	token, err := svc.GetSystemToken()
	if err != nil {
		panic(any(err))
	}

	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = websocket.NewClient(c.Ws.Host,
		websocket.WithClientHeader(header),
		websocket.WithClientDiscover(websocket.NewRedisDiscover(header, constants.REDIS_DISCOVER_SRV, c.Redisx)),
	)
	return svc
}

func (svc *ServiceContext) GetSystemToken() (string, error) {
	return svc.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
}
