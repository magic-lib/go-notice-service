package msgbuild

import (
	"context"
	"fmt"
	"github.com/magic-lib/go-notice-service/msg"
	"log"
)

// ChannelAdapter 渠道适配器接口
type ChannelAdapter interface {
	SupportedChannels() []msg.ChannelType
	Send(ctx context.Context, msg msg.Message) (string, error) // 发送消息
}

// ChannelAdapterRegistry 渠道适配器注册表
type ChannelAdapterRegistry interface {
	Register(channel msg.ChannelType, adapter ChannelAdapter) error // 注册适配器
	GetAdapter(channel msg.ChannelType) ChannelAdapter              // 获取适配器
}

// ChannelAdapterManager 渠道适配器管理器实现
type ChannelAdapterManager struct {
	adapters map[msg.ChannelType]ChannelAdapter
}

// NewChannelAdapterManager 创建适配器管理器
func NewChannelAdapterManager() ChannelAdapterRegistry {
	return &ChannelAdapterManager{
		adapters: make(map[msg.ChannelType]ChannelAdapter),
	}
}

// Register 注册渠道适配器
func (m *ChannelAdapterManager) Register(channel msg.ChannelType, adapter ChannelAdapter) error {
	if existing, ok := m.adapters[channel]; ok {
		return fmt.Errorf("渠道 %s 已注册适配器 %T，无法重复注册", channel, existing)
	}
	m.adapters[channel] = adapter
	return nil
}

func (m *ChannelAdapterManager) RegisterAll(adapter ChannelAdapter) {
	for _, channel := range adapter.SupportedChannels() {
		// 检查是否已存在，避免意外覆盖
		if existing, ok := m.adapters[channel]; ok {
			log.Printf("警告: 渠道 %s 的适配器已存在 (%T)，将被新适配器 (%T) 覆盖",
				channel, existing, adapter)
		}
		m.adapters[channel] = adapter
	}
}

// GetAdapter 获取渠道适配器
func (m *ChannelAdapterManager) GetAdapter(channel msg.ChannelType) ChannelAdapter {
	if adapter, ok := m.adapters[channel]; ok {
		return adapter
	}
	return nil
}

// MessageSender 消息发送器
type MessageSender struct {
	registry ChannelAdapterRegistry
}

// NewMessageSender 创建消息发送器
func NewMessageSender(registry ChannelAdapterRegistry) *MessageSender {
	return &MessageSender{registry: registry}
}

// Send 发送消息
func (s *MessageSender) Send(ctx context.Context, msg msg.Message) (string, error) {
	// 验证消息是否正确
	if err := msg.Validate(); err != nil {
		return "", err
	}

	// 获取对应渠道的适配器
	adapter := s.registry.GetAdapter(msg.Channel())
	if adapter == nil {
		return "", fmt.Errorf("unsupported channel: %s", msg.Channel())
	}

	// 委托给适配器发送
	return adapter.Send(ctx, msg)
}
