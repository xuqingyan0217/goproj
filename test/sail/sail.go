package main

import (
	"fmt"
	"github.com/HYY-yu/sail-client"
	"time"
)

type Config struct {
	Name string
	Host string
	Port string
	Mode string

	Database string

	UserRpc struct{
		Etcd struct{
			Hosts []string
			Key   string
		}
	}
	Redisx struct{
		Host string
		Pass string
	}
	JwtAuth struct{
		AccessSecret string
	}
}

func main() {
	var cfg Config
	// 我们下面所书写的配置，会组成一个请求key
	s := sail.New(&sail.MetaConfig{
		// etcd的连接地址
		ETCDEndpoints:  "10.0.0.17:3379",
		// 项目的key，在sail的ui界面，在项目管理处可以找到
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		// 命名空间和配置文件名
		Namespace:      "user",
		Configs:        "user-api.yaml",
		// 本地配置文件路径，尽量为空，因为它会默认先把该目录清空
		ConfigFilePath: "./conf",
		LogLevel:       "DEBUG",
	}, sail.WithOnConfigChange(func(configFileKey string, s *sail.Sail) {
		// 判断异常
		if s.Err() != nil {
			fmt.Println(s.Err())
			return
		}
		// Pull加载配置文件
		fmt.Println(s.Pull())
		// 将配置文件加载到我们上面的Config对象里
		v, err := s.MergeVipers()
		if err != nil {
			fmt.Println(err)
			return
		}
		v.Unmarshal(&cfg)
		fmt.Println(cfg, "\n" ,cfg.Database)
	}))
	// 判断异常
	if s.Err() != nil {
		fmt.Println(s.Err())
		return
	}
	// Pull加载配置文件
	fmt.Println(s.Pull())
	// 将配置文件加载到我们上面的Config对象里
	v, err := s.MergeVipers()
	if err != nil {
		fmt.Println(err)
		return
	}
	v.Unmarshal(&cfg)
	fmt.Println(cfg, "\n" ,cfg.Database)

	for {
		time.Sleep(time.Second)
	}
}
