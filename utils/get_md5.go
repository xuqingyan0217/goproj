package utils

import (
	"crypto/md5"
	"encoding/hex"
)
//GetMd5 外部使用时调用即可
func GetMd5(pwd string) string  {
	//创建一个哈希
	h := md5.New()
	h.Write([]byte(pwd))
	//返回这个哈希
	return hex.EncodeToString(h.Sum(nil))
}
