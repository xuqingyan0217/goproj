package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"
	"easy-chat/apps/social/socialmodels"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/xerr"
)

type GroupPutinLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinLogic {
	return &GroupPutinLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutin 此处目的是申请加群，如果有就直接返回，没有就添加一条申请记录，同时判断是否需要直接同意
func (l *GroupPutinLogic) GroupPutin(in *social.GroupPutinReq) (*social.GroupPutinResp, error) {
	// todo: add your logic here and delete this line

	//  1. 普通用户申请 ： 如果群无验证直接进入
	//  2. 群成员邀请： 如果群无验证直接进入
	//  3. 群管理员/群创建者邀请：直接进入群

	var (
		inviteGroupMember *socialmodels.GroupMembers
		userGroupMember   *socialmodels.GroupMembers
		groupInfo         *socialmodels.Groups

		err error
	)
	// 判断用户是否已经在群中
	userGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.ReqId, in.GroupId)
	// 查询出现错误并且错误不是不存在用户
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by groud id and  req id err %v, req %v, %v", err,
			in.GroupId, in.ReqId)
	}
	// 查询结果不为空则说明用户在群里，无需申请入群，返回
	if userGroupMember != nil {
		return &social.GroupPutinResp{}, nil
	}
	// 否则说明用户还未进群，下面开始入群申请
	// 首先判断用户是否有过申请
	groupReq, err := l.svcCtx.GroupRequestsModel.FindByGroupIdAndReqId(l.ctx, in.GroupId, in.ReqId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group req by groud id and user id err %v, req %v, %v", err,
			in.GroupId, in.ReqId)
	}
	// 不为空，则说明用户已经申请过了，只是还没有同意
	if groupReq != nil {
		return &social.GroupPutinResp{}, nil
	}

	// 否则，说明用户还没有申请
	// 申请进群的数据结果初始化
	groupReq = &socialmodels.GroupRequests{
		ReqId:   in.ReqId,
		GroupId: in.GroupId,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: sql.NullTime{
			Time:  time.Unix(in.ReqTime, 0),
			Valid: true,
		},
		JoinSource: sql.NullInt64{
			Int64: int64(in.JoinSource),
			Valid: true,
		},
		InviterUserId: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	}
	// 一个匿名函数，用于在无需验证时将用户添加到群里
	createGroupMember := func() {
		if err != nil {
			return
		}
		err = l.createGroupMember(in)
	}
	// 查找当前群信息，看是否需要验证
	groupInfo, err = l.svcCtx.GroupsModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group by groud id err %v, req %v", err, in.GroupId)
	}

	// 验证是否要验证
	if !groupInfo.IsVerify {
		// 不需要，则执行入群
		defer createGroupMember()

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}
		// 返回群号
		return l.createGroupReq(groupReq, true)
	}

	// 验证进群方式
	if constants.GroupJoinSource(in.JoinSource) == constants.PutInGroupJoinSource {
		// 申请
		// 返回空
		return l.createGroupReq(groupReq, false)
	}
	// 查询邀请新用户的那个人的信息
	inviteGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.InviterUid, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by groud id and user id err %v, req %v",
			in.InviterUid, in.GroupId)
	}

	if constants.GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.CreatorGroupRoleLevel || constants.
		GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.ManagerGroupRoleLevel {
		// 是管理者或创建者邀请
		defer createGroupMember()

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}
		groupReq.HandleUserId = sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		}
		return l.createGroupReq(groupReq, true)
	}
	return l.createGroupReq(groupReq, false)

}

// 往群申请表里添加一条数据，符合直接同意的，就返回群号，否则返回空，表示待同意（请求表有值）
func (l *GroupPutinLogic) createGroupReq(groupReq *socialmodels.GroupRequests, isPass bool) (*social.GroupPutinResp, error) {

	_, err := l.svcCtx.GroupRequestsModel.Insert(l.ctx, groupReq)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert group req err %v req %v", err, groupReq)
	}

	if isPass {
		return &social.GroupPutinResp{GroupId: groupReq.GroupId}, nil
	}

	return &social.GroupPutinResp{}, nil
}

// 申请通过，把用户添加到群里
func (l *GroupPutinLogic) createGroupMember(in *social.GroupPutinReq) error {
	groupMember := &socialmodels.GroupMembers{
		GroupId:     in.GroupId,
		UserId:      in.ReqId,
		RoleLevel:   int(constants.AtLargeGroupRoleLevel),
		OperatorUid: in.InviterUid,
	}
	_, err := l.svcCtx.GroupMembersModel.Insert(l.ctx, nil, groupMember)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "insert friend err %v req %v", err, groupMember)
	}

	return nil
}
