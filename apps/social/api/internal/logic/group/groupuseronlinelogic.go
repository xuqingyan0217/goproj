package group

import (
	"context"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/pkg/constants"

	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群在线用户
func NewGroupUserOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserOnlineLogic {
	return &GroupUserOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserOnlineLogic) GroupUserOnline(req *types.GroupUserOnlineReq) (resp *types.GroupUserOnlineResp, err error) {
	// todo: add your logic here and delete this line
	// 获取当前群的所有用户列表
	groupUserList, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
		GroupId: req.GroupId,
	})
	if err != nil || len(groupUserList.List) == 0{
		return &types.GroupUserOnlineResp{}, err
	}
	// 依据群用户列表获取到用户id
	uids := make([]string, len(groupUserList.List))

	for _, v := range groupUserList.List {
		uids = append(uids, v.UserId)
	}
	// 通过公共key来获取整个hash
	onlineList, err := l.svcCtx.Redis.Hgetall(constants.REDIS_OLINE_USER)
	if err != nil {
		return &types.GroupUserOnlineResp{}, err
	}
	// 通过遍历该群用户列表
	resOnlineList := make(map[string]bool, len(uids))
	for _, v := range uids {
		// 判断群用户是否在线
		if _, ok := onlineList[v]; ok {
			resOnlineList[v] = true
		} else {
			resOnlineList[v] = false
		}
	}

	return &types.GroupUserOnlineResp{
		OnlineList: resOnlineList,
	}, nil
}
