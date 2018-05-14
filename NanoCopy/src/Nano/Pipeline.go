package Nano

import (
	"Nano/Internal/Message"
	"Nano/Session"
	"fmt"
)

type (
	// Message struct {
	// 	temp *fjdkls
	// }

	// PipelineFunc 消息处理
	PipelineFunc func(me *Session.Session, msg Message.Message) error

	// IPipeline 管道
	IPipeline interface {
		Outbound() IPipelineChannel
		Inbound() IPipelineChannel
	}

	pipeline struct {
		outbound, inbound *pipelineChannel
	}

	// IPipelineChannel 管道频道
	IPipelineChannel interface {
		PushFront(h PipelineFunc)
		PushBack(h PipelineFunc)
		Process(s *Session.Session, msg Message.Message) error
	}

	pipelineChannel struct {
		handlers []PipelineFunc
	}
)

// NewPipeline 新创建管道实例
func NewPipeline() IPipeline {
	return &pipeline{
		outbound: &pipelineChannel{},
		inbound:  &pipelineChannel{},
	}
}

// Outbound 出站
func (me *pipeline) Outbound() IPipelineChannel {
	return me.outbound
}
func (me *pipeline) Inbound() IPipelineChannel {
	return me.inbound
}

// PushFront should not be used after nano running
func (me *pipelineChannel) PushFront(h PipelineFunc) {
	handlers := make([]PipelineFunc, len(me.handlers)+1)
	handlers[0] = h
	copy(handlers[1:], me.handlers)
	me.handlers = handlers
}

// PushBack should not be used after nano running
func (me *pipelineChannel) PushBack(h PipelineFunc) {
	me.handlers = append(me.handlers, h)
}

func (me *pipelineChannel) Process(s *Session.Session, msg Message.Message) error {
	if len(me.handlers) < 1 {
		return nil
	}
	for _, h := range me.handlers {
		if err := h(s, msg); err != nil {
			logger.Println(fmt.Sprintf("nano/handler: broken pipeline: %s", err.Error()))
			return err
		}
	}
	return nil
}
