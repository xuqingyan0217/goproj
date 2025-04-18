package friend

import (
	"context"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/ctxdata"

	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendsOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友在线情况
func NewFriendsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendsOnlineLogic {
	return &FriendsOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendsOnlineLogic) FriendsOnline(req *types.FriendOnlineReq) (resp *types.FriendOnlineResp, err error) {
	// todo: add your logic here and delete this line

	// 得到当前用户id
	uid := ctxdata.GetUId(l.ctx)
	// 获取到当前用户的好友列表
	friendsList, err := l.svcCtx.Social.FriendList(l.ctx, &socialclient.FriendListReq{
		UserId: uid,
	})
	if err != nil || len(friendsList.List) == 0 {
		return &types.FriendOnlineResp{}, err
	}
	// 依据好友列表获取到好友id
	uids := make([]string, len(friendsList.List))

	for _, v := range friendsList.List {
		uids = append(uids, v.UserId)
	}
	// 通过公共key来获取整个hash
	onlineList, err := l.svcCtx.Redis.Hgetall(constants.REDIS_OLINE_USER)
	if err != nil {
		return nil, err
	}
	// 通过遍历该用户的好友id列表, 判断是否在线,并把信息放入map中
	resOnlineList := make(map[string]bool, len(uids))
	for _, v := range uids {
		// 判断好友是否在线
		if _, ok := onlineList[v]; ok {
			resOnlineList[v] = true
		} else {
			resOnlineList[v] = false
		}
	}

	return &types.FriendOnlineResp{
		OnlineList: resOnlineList,
	}, nil
}
