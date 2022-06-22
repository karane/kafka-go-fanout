package time_buffer

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {

	var a string = "Hello"
	var b string = "Hello"

	assert.Equal(t, a, b, "The two words should be the same.")

}

func TestSimpleInsertion(t *testing.T) {

	outChannel := make(chan []Message)
	tbuf := New("my-buffer", 10*time.Second, outChannel)

	tbuf.SetMaxLength(5)
	tbuf.StartTimer()
	tbuf.Run()
	tbuf.InsertElement(Message{Topic: "A", Content: "w1"})
	tbuf.InsertElement(Message{Topic: "A", Content: "w2"})

	time.Sleep(1 * time.Second)

	assert.Equal(t, "my-buffer", tbuf.Id())
	assert.Equal(t, uint(2), tbuf.Size())
	assert.EqualValues(t, []Message{{Topic: "A", Content: "w1"}, {Topic: "A", Content: "w2"}}, <-tbuf.OutC)

}

func TestRacingInsertion(t *testing.T) {

	outChannel := make(chan []Message)
	tbuf := New("my-buffer", 10*time.Second, outChannel)

	tbuf.StartTimer()
	tbuf.SetMaxLength(201)
	tbuf.Run()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			msg := Message{Topic: "A", Content: fmt.Sprint("w", i)}
			tbuf.InsertElement(msg)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			msg := Message{Topic: "A", Content: fmt.Sprint("p", i)}
			tbuf.InsertElement(msg)
		}
	}()

	wg.Wait()
	time.Sleep(3 * time.Second)
	assert.Equal(t, "my-buffer", tbuf.Id())
	assert.EqualValues(t, uint(200), tbuf.Size())
}

func TestStop(t *testing.T) {
	outChannel := make(chan []Message)
	tbuf := New("my-buffer", 1*time.Second, outChannel)

	tbuf.StartTimer()
	tbuf.Run()
	tbuf.InsertElement(Message{Topic: "A", Content: "w"})
	tbuf.StopTimer()

	select {
	case <-tbuf.OutC:
		assert.Fail(t, "Channel should be empty")
	default:
	}

}

func TestTimeout(t *testing.T) {
	outChannel := make(chan []Message)
	tbuf := New("my-buffer", 1*time.Second, outChannel)

	tbuf.StartTimer()
	tbuf.Run()
	tbuf.InsertElement(Message{Topic: "A", Content: "w"})
	time.Sleep(2 * time.Second)

	select {
	case message := <-tbuf.OutC:
		assert.Equal(t, []Message{{Topic: "A", Content: "w"}}, message)
	default:
		assert.Fail(t, "Channel should have data")
	}

}

func TestMaxSize(t *testing.T) {
	outChannel := make(chan []Message)
	tbuf := New("my-buffer", 10*time.Second, outChannel)
	tbuf.SetMaxLength(2)

	tbuf.StartTimer()
	tbuf.Run()
	go func() {
		tbuf.InsertElement(Message{Topic: "A", Content: "a"})
		tbuf.InsertElement(Message{Topic: "A", Content: "b"})
		tbuf.InsertElement(Message{Topic: "A", Content: "c"})
	}()

	select {
	case message := <-tbuf.OutC:
		assert.Equal(t, []Message{{Topic: "A", Content: "a"}, {Topic: "A", Content: "b"}}, message)
	case <-time.After(2 * time.Second):
		assert.Fail(t, "Channel should have data")
	}

}

func TestMaxSize2(t *testing.T) {
	outChannel := make(chan []Message)
	tbuf := New("my-buffer", 10*time.Second, outChannel)
	tbuf.SetMaxLength(2)

	tbuf.StartTimer()
	tbuf.Run()
	go func() {

		tbuf.InsertElement(Message{Topic: "A", Content: "a"})
		tbuf.InsertElement(Message{Topic: "A", Content: "b"})
		tbuf.InsertElement(Message{Topic: "A", Content: "c"})
	}()

	select {
	case message := <-tbuf.OutC:
		assert.Equal(t, []Message{{Topic: "A", Content: "a"}, {Topic: "A", Content: "b"}}, message)
	case <-time.After(2 * time.Second):
		assert.Fail(t, "Channel should have data")
	}

	time.Sleep(6 * time.Second)
	select {
	case message := <-tbuf.OutC:
		assert.Equal(t, []Message{{Topic: "A", Content: "c"}}, message)
	case <-time.After(6 * time.Second):
		assert.Fail(t, "Channel should have data")
	}

}
