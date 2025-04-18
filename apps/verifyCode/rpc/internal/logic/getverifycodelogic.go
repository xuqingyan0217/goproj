package logic

import (
	"context"
	"easy-chat/apps/verifyCode/rpc/internal/svc"
	"easy-chat/apps/verifyCode/rpc/verifyCode"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
)

type GetVerifyCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVerifyCodeLogic {
	return &GetVerifyCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVerifyCodeLogic) GetVerifyCode(in *verifyCode.GetVerifyCodeRequest) (*verifyCode.GetVerifyCodeReply, error) {
	// todo: add your logic here and delete this line
	var code = RandCode(int(in.Length), in.Type)
	l.svcCtx.Redis.Setex(in.Phone, code, 60)
	return &verifyCode.GetVerifyCodeReply{
		Code: code,
	}, nil
}

// RandCode 依据请求参数做不同的处理
func RandCode(l int, t verifyCode.TYPE) string {
	switch t {
	case verifyCode.TYPE_DEFAULT: //默认就是type=1，意思是请求参数没有表明type的话，就按1的来
		fallthrough
	case verifyCode.TYPE_DIGIT:
		return randCode("0123456789", l, 4) //4是1-9的二进制个数
	case verifyCode.TYPE_LETTER:
		return randCode("ABCDEFGHIJKLMNOPQRSTUVWXYZ", l, 5) //5是a-z的26个字母的二进制个数
	case verifyCode.TYPE_MIXED:
		return randCode("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ", l, 6) //6是36个二进制个数
	default:

	}
	return ""
}

//randCode 做随机的核心方法(最简单的实现)
/*func randCode(chars string, l int) string  {
	charsLen := len(chars)
	result := make([]byte, l)
	for i:=0; i<l; i++ {
		//这里是随机一个下标，然后再去相应的chars里面对应
		randIndex := rand.Intn(charsLen)
		result[i] = chars[randIndex]
	}
	return string(result)
}*/

// randCode 做随机数的核心方法（优化实现）
// 一次随机多个随机位，分部分多次使用
func randCode(chars string, l, idxBits int) string {
	// 依据我们的优化思路，我们是从一次随机63位里面选择其中的4，5，6位
	// 推荐是在实参里面写死的，如果不写死，我们还需要如下计算出需要的二进制的个数
	// idxBits = len(fmt.Sprintf("%b", len(chars)))

	// 怎么切割呢，通过形成掩码 mask
	// 例如，使用低六位：00000000000111111
	idxMask := 1<<idxBits - 1 //00001000000 - 1 = 00000111111
	// 63位可以用多少次
	idxMax := 63 / idxBits
	// 结果
	result := make([]byte, l)
	// 生成随机字符
	// cache随机位缓存
	// remain当前还可以用多少次
	for i, cache, remain := 0, rand.Int63(), idxMax; i < l; {
		//如果使用的次数剩余为0次，则重新获取随机
		if 0 == remain {
			cache, remain = rand.Int63(), idxMax
		}
		//利用掩码获取有效部位的随机数位
		// 00110100 10110100 10110100 10110100 10110100 10110100 10110100 10110100
		// &
		// 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00111111
		// =
		// 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00110100
		if randIndex := int(cache & int64(idxMask)); randIndex < len(chars) {
			result[i] = chars[randIndex]
			i++
		}

		//使用下一组随机数
		// 00110100 10110100 10110100 10110100 10110100 10110100 10110100 10110100
		// cache >> 6
		// 000000 00110100 10110100 10110100 10110100 10110100 10110100 10110100 10
		cache >>= idxBits
		// 减少一次使用次数
		remain--
	}
	// 将获取到的随机数存储到设置过期时间的redis里面

	return string(result)
}
