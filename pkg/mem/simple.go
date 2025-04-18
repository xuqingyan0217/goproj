package mem

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cloudwego/eino/schema"
)

// GetDefaultMemory 返回默认配置的 SimpleMemory 实例
func GetDefaultMemory() *SimpleMemory {
	return NewSimpleMemory(SimpleMemoryConfig{
		Dir:           "data/memory", // 默认存储目录
		MaxWindowSize: 6,             // 默认最大窗口大小
	})
}

// SimpleMemoryConfig 用于配置 SimpleMemory
// Dir: 存储目录
// MaxWindowSize: 消息窗口最大长度
type SimpleMemoryConfig struct {
	Dir           string
	MaxWindowSize int
}

// NewSimpleMemory 创建一个新的 SimpleMemory 实例
func NewSimpleMemory(cfg SimpleMemoryConfig) *SimpleMemory {
	if cfg.Dir == "" {
		cfg.Dir = "/tmp/eino/memory"
	}
	if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
		return nil
	}

	return &SimpleMemory{
		dir:           cfg.Dir,
		maxWindowSize: cfg.MaxWindowSize,
		conversations: make(map[string]*Conversation),
	}
}

// SimpleMemory 简单内存结构体，可存储每个会话的消息
// 支持多会话并发安全操作
// dir: 存储目录
// maxWindowSize: 消息窗口最大长度
// conversations: 会话ID到会话对象的映射
type SimpleMemory struct {
	mu            sync.Mutex
	dir           string
	maxWindowSize int
	conversations map[string]*Conversation
}

// GetConversation 获取指定ID的会话，若不存在且 createIfNotExist 为 true 则创建
func (m *SimpleMemory) GetConversation(id string, createIfNotExist bool) *Conversation {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.conversations[id]

	filePath := filepath.Join(m.dir, id+".jsonl")
	if !ok {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if createIfNotExist {
				if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
					return nil
				}
				m.conversations[id] = &Conversation{
					ID:            id,
					Messages:      make([]*schema.Message, 0),
					filePath:      filePath,
					maxWindowSize: m.maxWindowSize,
				}
			}
		}

		con := &Conversation{
			ID:            id,
			Messages:      make([]*schema.Message, 0),
			filePath:      filePath,
			maxWindowSize: m.maxWindowSize,
		}
		con.load()
		m.conversations[id] = con
	}

	return m.conversations[id]
}

// ListConversations 列出所有已存在的会话ID
func (m *SimpleMemory) ListConversations() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	files, err := os.ReadDir(m.dir)
	if err != nil {
		return nil
	}

	ids := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ids = append(ids, strings.TrimSuffix(file.Name(), ".jsonl"))
	}

	return ids
}

// DeleteConversation 删除指定ID的会话及其文件
func (m *SimpleMemory) DeleteConversation(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	filePath := filepath.Join(m.dir, id+".jsonl")
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	delete(m.conversations, id)
	return nil
}

// Conversation 表示单个会话，包含消息列表和文件路径
// mu: 并发锁
// ID: 会话ID
// Messages: 消息列表
// filePath: 存储文件路径
// maxWindowSize: 消息窗口最大长度
type Conversation struct {
	mu sync.Mutex

	ID       string            `json:"id"`
	Messages []*schema.Message `json:"messages"`

	filePath string

	maxWindowSize int
}

// Append 向会话追加一条消息，并保存到文件
func (c *Conversation) Append(msg *schema.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Messages = append(c.Messages, msg)

	c.save(msg)
}

// GetFullMessages 获取该会话的全部消息
func (c *Conversation) GetFullMessages() []*schema.Message {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Messages
}

// GetMessages 获取窗口内的最新消息（窗口大小由 maxWindowSize 决定）
func (c *Conversation) GetMessages() []*schema.Message {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.Messages) > c.maxWindowSize {
		return c.Messages[len(c.Messages)-c.maxWindowSize:]
	}

	return c.Messages
}

// load 从文件加载所有消息到内存
func (c *Conversation) load() error {
	reader, err := os.Open(c.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		var msg schema.Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}
		c.Messages = append(c.Messages, &msg)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// save 将单条消息追加写入到文件
func (c *Conversation) save(msg *schema.Message) {
	str, _ := json.Marshal(msg)

	// 追加写入到文件
	f, err := os.OpenFile(c.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.Write(str)
	f.WriteString("\n")
}
