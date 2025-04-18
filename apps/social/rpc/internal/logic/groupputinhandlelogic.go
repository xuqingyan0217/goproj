package logic

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"easy-chat/apps/social/socialmodels"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/xerr"

	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrGroupReqBeforePass   = xerr.NewMsg("请求已通过")
	ErrGroupReqBeforeRefuse = xerr.NewMsg("请求已拒绝")
)

type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	// todo: add your logic here and delete this line
	// 获取群申请记录
	groupReq, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, int64(in.GroupReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend req err %v req %v", err, in.GroupReqId)
	}
	// 验证该表中记录在创建之后是否已经被处理过（因为已经处理的我们也不打算删，只是状态变了）
	switch constants.HandlerResult(groupReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforeRefuse)
	}
	// 未被处理的话，就只有两种情况了，把当前用户请求期望的状态赋值过去
	groupReq.HandleResult = sql.NullInt64{
		Int64: int64(in.HandleResult),
		Valid: true,
	}
	// 用事务对未处理的申请处理，若同意，需要插入两条好友关系记录
	err = l.svcCtx.GroupRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 更新申请表里面的记录状态，该状态我们也不知道是同意还是不同意，所以往下面又一判断
		if err := l.svcCtx.GroupRequestsModel.Update(l.ctx, session, groupReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend req err %v req %v", err, groupReq)
		}
		// 判断当前的处理方式是不是不通过的，如果不通过，是肯定不能建立群聊关系的
		if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
			return nil
		}
		// 用户在进入群聊后的数据样式初始化
		groupMember := &socialmodels.GroupMembers{
			GroupId:     groupReq.GroupId,
			UserId:      groupReq.ReqId,
			RoleLevel:   int(constants.AtLargeGroupRoleLevel),
			OperatorUid: in.HandleUid,
		}
		// 将用户插入群表中，标志着该用户已进群
		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, groupMember)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert friend err %v req %v", err, groupMember)
		}

		return nil
	})

	if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
		return &social.GroupPutInHandleResp{}, err
	}

	return &social.GroupPutInHandleResp{
		GroupId: groupReq.GroupId,
	}, err
}
