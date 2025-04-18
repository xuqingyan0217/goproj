package mq

import (
	"github.com/nats-io/nats.go"
	"os"
)

var (
	Nc  *nats.Conn
	err error
)

func Init() {
	natsUrl := os.Getenv("NATS_URL")
	Nc, err = nats.Connect(natsUrl)
	if err != nil {
		// 连接如果失败，后面一切都是没有意义的，直接panic
		panic(err)
	}
}
