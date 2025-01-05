package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func GetMd5(pwd string) string  {
	//创建一个哈希
	h := md5.New()
	h.Write([]byte(pwd))
	//返回这个哈希
	return hex.EncodeToString(h.Sum(nil))
}

func main()  {
	md5_pwd := GetMd5("123456")
	fmt.Println(md5_pwd)
}