package configserver

import (
	"errors"
	"github.com/zeromicro/go-zero/core/conf"
)

var ErrNotSetConfig = errors.New("未设置配置信息")

type OnChange func([]byte) error

// ConfigServer 定义了一个配置服务器的接口，它提供了从JSON字节中获取配置和错误处理的方法。
type ConfigServer interface {
	Build() error	// 每个程序运行不一样，提供的操作也不一样，所以来个build用于初始化
	SetOnChange(OnChange)    // 用于支持在配置发生改变时执行的方法
	FromJsonBytes() ([]byte, error) // 从JSON字节中加载配置
}

// configServer 实现了ConfigServer接口，添加了配置文件路径的字段。
type configServer struct {
	ConfigServer // 嵌入ConfigServer接口
	configFile   string // 配置文件路径
}

// NewConfigServer 创建一个新的configServer实例。
// 参数：
//   configFile: 配置文件路径
//   s: 实现了ConfigServer接口的服务器配置实例
func NewConfigServer(configFile string, s ConfigServer) *configServer {
	return &configServer{
		ConfigServer: s,
		configFile:   configFile,
	}
}

// MustLoad 加载配置到给定的结构体变量v中。
// 如果ConfigServer有错误，直接返回错误。
// 如果没有配置文件路径且ConfigServer实例为nil，返回ErrNotSetConfig错误。
// 如果ConfigServer实例为nil，使用go-zero默认的方式加载配置。
// 否则，从JSON字节中加载配置。
func (c *configServer) MustLoad(v any, onChange OnChange) error {
	if c.configFile == "" && c.ConfigServer == nil {
		return ErrNotSetConfig
	}
	if c.ConfigServer == nil {
		// 采用go-zero默认的方式
		conf.MustLoad(c.configFile, v)
		return nil
	}
	// 如果不为空，我们就去设置配置更新时要执行的方法
	if onChange != nil {
		c.ConfigServer.SetOnChange(onChange)
	}
	// 然后，再对当前配置中心进行构建
	if err := c.ConfigServer.Build(); err != nil {
		return err
	}

	data, err := c.ConfigServer.FromJsonBytes()
	if err != nil {
		return err
	}
	return LoadFromJsonBytes(data, v)
}

func LoadFromJsonBytes(data []byte, v any) error {
	return conf.LoadFromJsonBytes(data, v)
}