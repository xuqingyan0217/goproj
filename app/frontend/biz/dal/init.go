package dal

import (
	"gomall/app/frontend/biz/dal/mysql"
	"gomall/app/frontend/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
