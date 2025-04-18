package logic

import (
	"context"
	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/user/rpc/user"
	"easy-chat/pkg/bitmap"
	"easy-chat/pkg/constants"

	"easy-chat/apps/im/api/internal/svc"
	"easy-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogReadRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChatLogReadRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogReadRecordsLogic {
	return &GetChatLogReadRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetChatLogReadRecords 获取聊天记录的阅读状态
// 该方法接收一个GetChatLogReadRecordsReq请求对象，返回一个GetChatLogReadRecordsResp响应对象和一个错误信息
func (l *GetChatLogReadRecordsLogic) GetChatLogReadRecords(req *types.GetChatLogReadRecordsReq) (resp *types.GetChatLogReadRecordsResp, err error) {
    // todo: add your logic here and delete this line
    // 调用IM服务获取聊天记录，但依据消息id获取，只能获取一条消息的记录，该消息可能是群聊消息可能是私聊消息
	// 这两种消息虽然在同一数据库，并且结构都一样，但其中的bitmap会略有不同。
    chatlogs, err := l.svcCtx.Im.GetChatLog(l.ctx, &im.GetChatLogReq{
        MsgId:          req.MsgId,
    })
	// 调用出现错误，或者聊天记录列表为空
    if err != nil || len(chatlogs.List) == 0 {
        return nil, err
    }

    var (
        // chatlog 代表聊天记录列表中的第一个元素，因为只能获取一条
        chatlog = chatlogs.List[0]
        // 消息的发送方肯定是已读的，所以我们直接放进去，chatlog的ReadRecord里面也有记录
        reads   = []string{chatlog.SendId}
        // unreads 用于存储未读消息的接收方ID
        unreads  []string
        // ids 用于存储业务涉及到的用户id，便于后续由它作为调用user-rpc的参数
        ids      []string
    )

    // 根据聊天类型分别处理已读和未读状态
    switch constants.ChatType(chatlog.ChatType) {
    case constants.SingleChatType:
        // 私聊处理逻辑
		// 我们知道，对于私聊，只有两方，所以字节数值只需要一个字节的数据即可，通过判断该数据是0是1即可
		// 但是，我们在记录消息时，无论私聊群聊，都已经添加过一个bitmap，而私聊消息如果已读了的话，那么该bitmap就是只有1个字节的1
		// 而如果消息是未读，那么此时该消息的bitmap还是我们之前设置的那个长的，里面只记录了消息发送者的id所在位为1
		// 所以说，对于私聊的处理，我们需要先直接判断一手当前bitmap的len是否为1，为1说明bitmap已经被覆盖了，也就是[]byte{1}
		// 因为我们在消息记录时，设置的size大小是0，那么就会默认变为250，所以可以由此判断私聊是否已读。
		if len(chatlog.ReadRecords) == 1 && chatlog.ReadRecords[0] == 1 {
            // 如果bitmap长度为1，并且第一个字节为1，则认为是已读消息
            reads = append(reads, chatlog.RecvId)
        } else {
			unreads = []string{chatlog.RecvId}
		}
		// 下面的是视频里面的，我认为不合理
        /*if len(chatlog.ReadRecords) == 0 || chatlog.ReadRecords[0] == 0 {
            // 如果没有读取记录，或者第一个读取记录为0，则认为是未读消息
            unreads = []string{chatlog.RecvId}
        } else {
            // 否则，将接收者ID添加到已读消息列表中
            reads = append(reads, chatlog.RecvId)
        }*/
        // 构建消息交互的双方ID列表，包括发送者和接收者
        ids = []string{chatlog.RecvId, chatlog.SendId}
    case constants.GroupChatType:
        // 群聊处理逻辑
		// 群聊的处理就很简单了，因为它至始至终都是一个bitmap，只会增加1的个数，不会长度变化
        groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
            GroupId: chatlog.RecvId,
        })
        if err != nil {
            return nil, err
        }
        bitmaps := bitmap.Load(chatlog.ReadRecords)
        for _, members := range groupUsers.List {
            ids = append(ids, members.UserId)
			// 如果是消息发送者的话，跳过，因为在消息记录时已经做过了
            if members.UserId == chatlog.SendId {
                continue
            }
			// 否则，调用bitmap里面的依据用户id判断是否有过记录，有则已读，无则未读
            if bitmaps.IsSet(members.UserId) {
                reads = append(reads, members.UserId)
            } else {
                unreads = append(unreads, members.UserId)
            }
        }
    }

    // 获取用户信息，为了得到phone 字段
    userEntitys, err := l.svcCtx.User.FindUser(l.ctx, &user.FindUserReq{
        Ids:   ids,
    })
    if err != nil {
        return nil, err
    }
	// 将用户信息的集合以map的方法进行存储，其中key就是用户id，value就是用户信息
    userEntitySet := make(map[string]*user.UserEntity, len(userEntitys.User))
    for i, entity := range userEntitys.User {
        userEntitySet[entity.Id] = userEntitys.User[i]
    }

    // 设置手机号码
	// 后续就很简单了，有了key为用户id的用户map后，再加上read和unread里面的用户id信息，即可实现设置手机号。
    for i, read := range reads {
        if u := userEntitySet[read]; u != nil {
            reads[i] = u.Phone
        }
    }
    for i, unread := range unreads {
        if u := userEntitySet[unread]; u != nil {
            unreads[i] = u.Phone
        }
    }
    // 返回消息已读未读的用户电话
    return &types.GetChatLogReadRecordsResp{
        Reads:   reads,
        UnReads: unreads,
    }, nil
}

