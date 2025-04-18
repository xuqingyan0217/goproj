package job

import "time"

const (
	DefaultRetryJetLag = time.Second // 重试间隔
	DefaultRetryTimeout = 2 * time.Second // 重试超时时间；毕竟重试也是一个任务，他也会超时
	DefaultRetryNums = 5 // 最大重试次数
)


