package dal

import (
	"gomall/app/AIEino/biz/dal/mysql"
	"gomall/app/AIEino/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
