package dal

import (
	"gomall/app/email/biz/dal/mysql"
	"gomall/app/email/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
