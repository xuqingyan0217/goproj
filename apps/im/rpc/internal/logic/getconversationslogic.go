package logic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"easy-chat/apps/im/immodels"
	"easy-chat/pkg/xerr"

	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  Get获取会话（可以理解为获取qq上的可聊天对象列表包括好友间的，临时对话的；）
func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
	// todo: add your logic here and delete this line
	// 具体业务:获取到用户的会话列表，然后由该列表进行遍历获取到ids，再由ids去获取每一个具体的会话；然后做已读未读。
	// 根据用户查询用户会话列表conversations
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		// 如下则说明用户的会话列表conversations里没有任何会话conversation，这不算错误
		if err == immodels.ErrNotFound {
			return &im.GetConversationsResp{}, nil
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindByUserId err %v, req %v", err, in.UserId)
	}
	// 定义返回类型
	var res im.GetConversationsResp
	// 将上述获取到的用户会话传到该返回类型中
	copier.Copy(&res, &data)

	// 根据会话列表conversations，存储每个会话id，进而在conversation里找到会话
	ids := make([]string, 0, len(data.ConversationList))
	// 通过for循环获取到conversations里的id集合
	for _, conversation := range data.ConversationList {
		ids = append(ids, conversation.ConversationId)
	}
	// 依据id集合在conversation表中找到具体的会话，进而组成一个集合
	conversations, err := l.svcCtx.ConversationModel.ListByConversationIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.ListByConversationIds err %v, req %v", err, ids)
	}

	// 计算是否存在未读消息
	// 注意，遍历的conversations是一个刚才for循环得到的每个conversation组成的集合，它不是conversations表里面的那个map
	for _, conversation := range conversations {
		// 如果在conversations查询不到，跳过
		if _, ok := res.ConversationList[conversation.ConversationId]; !ok {
			continue
		}
		// conversations表里的map是当作所有已读消息会话的集合，和上面的conversation表中的会话集合区分开，那里面是总消息数
		// 用户读取的消息量（已读）
		total := res.ConversationList[conversation.ConversationId].Total
		// 已读的消息总数如果小于了会话里面的消息，说明有新增消息了
		if total < int32(conversation.Total) {
			// 有新的消息（下面两个主要是前端显示，尤其是已读改成了总消息数，目的是前端可以显示总消息数和未读数）
			// 设置已读会话中，已读消息为新值，这个不能写入数据库，我们只是告诉前端
			res.ConversationList[conversation.ConversationId].Total = int32(conversation.Total)
			// 由二者的差值计算有多少是未读，添加到未读字段中，用于显示前端
			res.ConversationList[conversation.ConversationId].ToRead = int32(conversation.Total) - total
			// 更改当前会话为显示状态，因为用户可能会删除会话（仅是不可见，不是说真删了），当消息来了，我们让它再出来
			res.ConversationList[conversation.ConversationId].IsShow = true
		}
	}
	return &res, nil
}
