package time_buffer

import (
	"fmt"
	"sync"
	"time"
)

const MAX = 200

type Message struct {
	Topic   string
	Content string
}

func (m Message) String() string {
	return fmt.Sprintf("{topic: %v, msg: %v}", m.Topic, m.Content)
}

type Buffer struct {
	id        string
	length    uint
	maxLength uint
	messages  []Message
	timer     *time.Timer
	timeout   time.Duration
	lock      sync.RWMutex
	OutC      chan []Message
	insertC   chan Message
	innerC    chan bool
}

type timerFunc func()

type TimedBuffer interface {
	insertElement(element Message)
	InsertElement(element Message)
	StartTimer()
	StopTimer()
	SetMaxLength(max uint)
	Size() uint
	Id() string
	Cancel()
}

func New(id string, timeout time.Duration, outChannel chan []Message) *Buffer {
	buffer := &Buffer{
		id:        id,
		length:    0,
		maxLength: MAX,
		messages:  make([]Message, 0),
		timeout:   timeout,
		timer:     nil,
		OutC:      outChannel,
		insertC:   make(chan Message),
		innerC:    make(chan bool),
	}

	return buffer
}

func (b *Buffer) InsertElement(element Message) {
	b.insertC <- element
}

func (b *Buffer) insertElement(element Message) bool {
	if b.length >= b.maxLength {
		return false
	}

	b.length++
	b.messages = append(b.messages, element)

	if b.length == b.maxLength {
		b.flushToChannel("SizeLimit")
	}
	return true
}

func (b *Buffer) Run() {

	go func() {
		for {
			select {
			case element := <-b.insertC:
				fmt.Printf("insert: %v\n", element)
				b.insertElement(element)
			case <-b.innerC:
				fmt.Printf("Buffer: %v canceled\n", b.id)
				return
			}
		}
	}()
}

func (b *Buffer) reset() {
	b.messages = make([]Message, 0)
	b.length = 0
}

func (b *Buffer) Size() uint {
	return b.length
}

func (b *Buffer) StartTimer() {
	b.timer = time.NewTimer(b.timeout)
	go func() {
		<-b.timer.C
		b.flushToChannel("Timer")
	}()
}

func (b *Buffer) flushToChannel(originMessage string) {
	b.StopTimer()
	b.StartTimer()
	var messagesCopy []Message = nil
	if b.length > 0 {
		messagesCopy = b.messages
		fmt.Printf("{%v} Copied: %v\n", originMessage, messagesCopy)
		b.messages = make([]Message, 0)
		b.length = 0
	}

	if messagesCopy != nil {
		fmt.Println("Message sent ...")
		b.OutC <- messagesCopy
	}
}

func (b *Buffer) StopTimer() {
	if !b.timer.Stop() {
		select {
		case <-b.timer.C:
		case <-time.After(1 * time.Millisecond):
		}
	}
}

func (b *Buffer) Cancel() {
	b.timer.Stop()
	// b.innerC <- true
}

func (b *Buffer) SetMaxLength(max uint) {
	b.maxLength = max
}

func (b *Buffer) Id() string {
	return b.id
}
