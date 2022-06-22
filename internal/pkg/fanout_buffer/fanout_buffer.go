package fanout_buffer

import (
	"kafka-forwarder/internal/pkg/time_buffer"
	"sync"
	"time"
)

type FBuffer struct {
	cache     sync.Map
	size      uint
	topicSize uint
	OutC      chan []time_buffer.Message
	timeout   time.Duration
}

type FanoutBuffer interface {
	Put(elem string)
}

//New
func New(outChannel chan []time_buffer.Message, topicSize uint, timeout time.Duration) *FBuffer {
	return &FBuffer{size: 0,
		topicSize: topicSize,
		OutC:      outChannel,
		timeout:   timeout,
	}
}

func (fbuf *FBuffer) Put(key string, elem string) {
	value, exists := fbuf.cache.Load(key)
	var tbuf *time_buffer.Buffer
	if exists {
		tbuf = value.(*time_buffer.Buffer)
	} else {
		tbuf = time_buffer.New(key, fbuf.timeout, fbuf.OutC)
		tbuf.SetMaxLength(fbuf.topicSize)
		fbuf.cache.Store(key, tbuf)
		tbuf.StartTimer()
		tbuf.Run()
	}

	message := time_buffer.Message{Topic: key, Content: elem}
	tbuf.InsertElement(message)
}

func (fbuf *FBuffer) Get(key string) (*time_buffer.Buffer, bool) {
	value, exists := fbuf.cache.Load(key)
	var tbuf *time_buffer.Buffer
	if !exists {
		return &time_buffer.Buffer{}, false
	}

	tbuf = value.(*time_buffer.Buffer)
	return tbuf, true
}
