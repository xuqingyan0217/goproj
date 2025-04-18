package friend

import (
	"context"
	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(req *types.FriendPutInHandleReq) (resp *types.FriendPutInHandleResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUId(l.ctx)
	var friend *socialclient.FriendPutInHandleResp
	friend, err = l.svcCtx.Social.FriendPutInHandle(l.ctx, &socialclient.FriendPutInHandleReq{
		FriendReqId:  req.FriendReqId,
		UserId:       userId,
		HandleResult: req.HandleResult,
	})
	// 在上面调用rpc之后，此时就建立了好友关系或者拒绝请求了，我们针对该情况做会话是否建立
	// 调用创建会话的方法
	if constants.HandlerResult(req.HandleResult) == constants.PassHandlerResult {
		_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, &im.SetUpUserConversationReq{
			SendId:   userId,
			RecvId:   friend.Fid,
			ChatType: 2,
		})
	}
	return &types.FriendPutInHandleResp{}, err
}
