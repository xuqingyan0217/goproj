package dal

import (
	"gomall/app/checkout/biz/dal/mysql"
	"gomall/app/checkout/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
