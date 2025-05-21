package broadcast

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type broadcaster struct {
	input   chan interface{}
	reg     chan chan<- interface{}
	unreg   chan chan<- interface{}
	outputs map[chan<- interface{}]bool
	mutex   sync.RWMutex
	closed  bool
	timeout time.Duration
}

type Broadcaster interface {
	Register(chan<- interface{})
	Unregister(chan<- interface{})
	Close() error
	Submit(interface{})
	TrySubmit(interface{}) bool
}

func sendMessageToChannel(b *broadcaster, ch chan<- interface{}, m interface{}) {
	defer func() {
		if r := recover(); r != nil {
			b.mutex.Lock()
			delete(b.outputs, ch)
			b.mutex.Unlock()
			log.Warn().Interface("error", r).Msg("Deleted closed channel")
			// metrics.Increment("broadcaster.channel_closed")
		}
	}()
	select {
	case ch <- m:
		// metrics.Increment("broadcaster.messages_sent")
	case <-time.After(b.timeout):
		log.Warn().Msg("Timeout sending message to channel")
		// metrics.Increment("broadcaster.messages_dropped")
	}
}

func (b *broadcaster) broadcast(m interface{}) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for ch := range b.outputs {
		sendMessageToChannel(b, ch, m)
	}
}

func (b *broadcaster) run() {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("error", r).Stack().Msg("Broadcaster panicked, restarting")
			time.Sleep(time.Second) // Backoff
			if !b.closed {
				go b.run()
			}
		}
	}()
	for {
		select {
		case m, ok := <-b.input:
			if !ok {
				return
			}
			b.broadcast(m)
		case ch, ok := <-b.reg:
			if !ok {
				return
			}
			b.mutex.Lock()
			b.outputs[ch] = true
			b.mutex.Unlock()
			// metrics.Gauge("broadcaster.subscribers", float64(len(b.outputs)))
		case ch := <-b.unreg:
			b.mutex.Lock()
			delete(b.outputs, ch)
			b.mutex.Unlock()
			// metrics.Gauge("broadcaster.subscribers", float64(len(b.outputs)))
		}
	}
}

func NewBroadcaster(buflen int, timeout time.Duration) Broadcaster {
	if buflen <= 0 {
		buflen = 100 // Default buffer size
	}
	if timeout <= 0 {
		timeout = time.Millisecond * 10 // Default timeout
	}
	b := &broadcaster{
		input:   make(chan interface{}, buflen),
		reg:     make(chan chan<- interface{}, buflen),
		unreg:   make(chan chan<- interface{}, buflen),
		outputs: make(map[chan<- interface{}]bool),
		timeout: timeout,
	}
	go b.run()
	return b
}

func (b *broadcaster) Register(ch chan<- interface{}) {
	if b.closed {
		log.Warn().Msg("Attempt to register on closed broadcaster")
		return
	}
	b.reg <- ch
}

func (b *broadcaster) Unregister(ch chan<- interface{}) {
	if b.closed {
		log.Warn().Msg("Attempt to unregister on closed broadcaster")
		return
	}
	b.unreg <- ch
}

func (b *broadcaster) Close() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.closed {
		return nil
	}
	b.closed = true
	close(b.reg)
	close(b.unreg)
	close(b.input)
	for ch := range b.outputs {
		delete(b.outputs, ch)
	}
	log.Info().Msg("Broadcaster closed")
	return nil
}

func (b *broadcaster) Submit(m interface{}) {
	if b == nil || b.closed {
		log.Warn().Msg("Attempt to submit on nil or closed broadcaster")
		return
	}
	select {
	case b.input <- m:
	default:
		log.Warn().Msg("Input channel full, dropping message")
		// metrics.Increment("broadcaster.messages_dropped_input")
	}
}

func (b *broadcaster) TrySubmit(m interface{}) bool {
	if b == nil || b.closed {
		log.Warn().Msg("Attempt to try submit on nil or closed broadcaster")
		return false
	}
	select {
	case b.input <- m:
		return true
	default:
		log.Warn().Msg("Input channel full, try submit failed")
		// metrics.Increment("broadcaster.messages_dropped_input")
		return false
	}
}