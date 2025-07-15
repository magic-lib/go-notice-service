package msgbuild

import (
	"context"
	"fmt"
	"github.com/magic-lib/go-notice-service/msg"
	"github.com/samber/lo"
)

// ChannelAdapter 渠道适配器接口
type ChannelAdapter interface {
	SupportedChannels() []msg.ChannelType //在一个渠道中唯一,可注册到多个渠道中
	ChannelKey() msg.ChannelKey
	Send(ctx context.Context, message msg.MessageTemplate) (string, error) // 发送消息
}

// ChannelAdapterRegistry 渠道适配器注册表
type ChannelAdapterRegistry interface {
	Register(channel msg.ChannelType, key msg.ChannelKey, adapter ChannelAdapter) error // 注册适配器
	RegisterAll(adapter ...ChannelAdapter)                                              // 注册适配器
	GetChannels() map[msg.ChannelType][]msg.ChannelKey                                  // 获取所有渠道适配器
	GetAdapter(channel msg.ChannelType, key msg.ChannelKey) ChannelAdapter              // 获取适配器
}

// ChannelAdapterManager 渠道适配器管理器实现
type ChannelAdapterManager struct {
	adapters map[msg.ChannelType]map[msg.ChannelKey]ChannelAdapter
}

// NewChannelAdapterManager 创建适配器管理器
func NewChannelAdapterManager() ChannelAdapterRegistry {
	return &ChannelAdapterManager{
		adapters: make(map[msg.ChannelType]map[msg.ChannelKey]ChannelAdapter),
	}
}

// Register 注册渠道适配器
func (m *ChannelAdapterManager) Register(channel msg.ChannelType, key msg.ChannelKey, adapter ChannelAdapter) error {
	if m.adapters[channel] == nil {
		m.adapters[channel] = make(map[msg.ChannelKey]ChannelAdapter)
	}
	if key == "" {
		key = msg.ChannelKeyDefault
	}
	if existing, ok := m.adapters[channel]; ok {
		if oneAdapter, ok := existing[key]; ok {
			return fmt.Errorf("渠道 %s 的 %s 已注册适配器 %T，无法重复注册", channel, key, oneAdapter)
		}
	}
	m.adapters[channel][key] = adapter
	return nil
}

func (m *ChannelAdapterManager) RegisterAll(adapters ...ChannelAdapter) {
	if len(adapters) == 0 {
		return
	}
	lo.ForEach(adapters, func(adapter ChannelAdapter, _ int) {
		channels := adapter.SupportedChannels()
		if channels == nil {
			return
		}
		lo.ForEach(channels, func(channel msg.ChannelType, _ int) {
			_ = m.Register(channel, adapter.ChannelKey(), adapter)
		})
	})
}
func (m *ChannelAdapterManager) GetChannels() map[msg.ChannelType][]msg.ChannelKey {
	allChannels := make(map[msg.ChannelType][]msg.ChannelKey)
	for channel, adapters := range m.adapters {
		if _, ok := allChannels[channel]; !ok {
			allChannels[channel] = make([]msg.ChannelKey, 0, len(adapters))
		}
		for key, _ := range adapters {
			allChannels[channel] = append(allChannels[channel], key)
		}
	}
	return allChannels
}

// GetAdapter 获取渠道适配器
func (m *ChannelAdapterManager) GetAdapter(channel msg.ChannelType, key msg.ChannelKey) ChannelAdapter {
	if adapter, ok := m.adapters[channel]; ok {
		if key != "" {
			if oneAdapter, ok := adapter[key]; ok {
				return oneAdapter
			}
		}
		if len(adapter) == 1 { // 只有一个适配器，直接返回
			for _, oneAdapter := range adapter {
				return oneAdapter
			}
		}
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
func (s *MessageSender) Send(ctx context.Context, msg msg.MessageTemplate) (string, error) {
	// 验证消息是否正确
	if err := msg.Validate(); err != nil {
		return "", fmt.Errorf("invalid Validate message: %w", err)
	}
	msgChannel, msgChannelKey := msg.Channel()
	if msgChannel == "" {
		return "", fmt.Errorf("message channel is empty")
	}
	adapter := s.registry.GetAdapter(msgChannel, msgChannelKey)
	if adapter == nil {
		return "", fmt.Errorf("unsupported channel: %s, %s, use WithChannelAdapter or WithChannel", msgChannel, msgChannelKey)
	}
	return adapter.Send(ctx, msg)
}
