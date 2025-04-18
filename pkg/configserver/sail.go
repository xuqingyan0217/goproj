package configserver

import (
	"encoding/json"
	"fmt"
	"github.com/HYY-yu/sail-client"
)

type Config struct {
	ETCDEndpoints string `toml:"etcd_endpoints"` // 逗号分隔的ETCD地址，0.0.0.0:2379,0.0.0.0:12379,0.0.0.0:22379

	ProjectKey   string `toml:"project_key"`
	Namespace    string `toml:"namespace"`

	Configs        string `toml:"configs"`          // 逗号分隔的 config_name.config_type，如：mysql.toml,cfg.json,redis.yaml，空代表不下载任何配置
	ConfigFilePath string `toml:"config_file_path"` // 本地配置文件存放路径，空代表不存储本都配置文件
	LogLevel       string `toml:"log_level"`        // 日志级别(DEBUG\INFO\WARN\ERROR)，默认 WARN
}

type Sail struct {
	*sail.Sail
	sail.OnConfigChange
	c *Config
}

func NewSail(cfg *Config) *Sail {
	//s := sail.New(&sail.MetaConfig{
	//	// etcd的连接地址
	//	ETCDEndpoints:  cfg.ETCDEndpoints,
	//	// 项目的key，在sail的ui界面，在项目管理处可以找到
	//	ProjectKey:     cfg.ProjectKey,
	//	// 命名空间和配置文件名
	//	Namespace:      cfg.Namespace,
	//	Configs:        cfg.Configs,
	//	// 本地配置文件路径，尽量为空，因为它会默认先把该目录清空
	//	ConfigFilePath: cfg.ConfigFilePath,
	//	LogLevel:       cfg.LogLevel,
	//})
	return &Sail{c:cfg}
}

// 实现接口

func (s *Sail) FromJsonBytes() ([]byte, error) {
	if err := s.Pull(); err != nil {
		return nil, err
	}

	return s.fromJsonBytes(s.Sail)
}

func (s *Sail) fromJsonBytes(sail *sail.Sail) ([]byte, error) {
	// 将配置文件加载到我们上面的Config对象里
	// MergeVipers 合并Vipers
	// 该方法将多个Viper实例合并为一个，以便统一管理配置
	// 返回值:
	// - *Viper: 合并后的Viper实例
	// - error: 如果合并过程中发生错误，返回错误信息
	v, err := sail.MergeVipers()
	if err != nil {
		return nil, err
	}

	// AllSettings 获取合并后的Viper实例中的所有设置
	// 该方法将配置数据以map形式返回，便于后续处理或导出
	data := v.AllSettings()

	// json.Marshal 将配置数据的map转换为JSON格式的字节切片
	// 这一步通常用于导出配置数据或在网络中传输
	return json.Marshal(data)
}

func (s *Sail) Build() error {
	// 定义一个用于存储配置选项的切片
	var opts []sail.Option
	// 检查是否存在配置变更回调函数
	if s.OnConfigChange != nil {
	    // 如果回调函数存在，则将其添加到配置选项中
	    opts = append(opts, sail.WithOnConfigChange(s.OnConfigChange))
	}
	s.Sail = sail.New(&sail.MetaConfig{
		// etcd的连接地址
		ETCDEndpoints:  s.c.ETCDEndpoints,
		// 项目的key，在sail的ui界面，在项目管理处可以找到
		ProjectKey:     s.c.ProjectKey,
		// 命名空间和配置文件名
		Namespace:      s.c.Namespace,
		Configs:        s.c.Configs,
		// 本地配置文件路径，尽量为空，因为它会默认先把该目录清空
		ConfigFilePath: s.c.ConfigFilePath,
		LogLevel:       s.c.LogLevel,
	}, opts...)
	return s.Sail.Err()
}

// SetOnChange 设置配置变更时的回调函数
// 此函数允许开发者指定一个回调函数f，当Sail对象的配置发生变化时，将调用此函数
// 参数:
//   f: 一个OnChange类型的函数，表示配置变更时需要执行的逻辑
func (s *Sail) SetOnChange(f OnChange) {
    // 将OnConfigChange回调函数设置为一个新的函数
    s.OnConfigChange = func(configFileKey string, sail *sail.Sail) {
        // 从sail对象中反序列化数据
        data, err := s.fromJsonBytes(sail)
        if err != nil {
            // 如果反序列化失败，打印错误信息并返回
            fmt.Println("get config err", err)
            return
        }
        // 调用开发者指定的回调函数f，并传递反序列化后的数据
        if err = f(data); err != nil {
            // 如果回调函数执行过程中发生错误，打印错误信息
            fmt.Println("onchange err", err)
        }
    }
}
