package model

import (
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"go-streaming/broadcast"
)

type Message struct {
	ChannelId string
	Text      string
}

type Listener struct {
	ChannelId string
	Chan      chan interface{}
}

type ChannelManager struct {
	channels      map[string]broadcast.Broadcaster
	open          chan *Listener
	close         chan *Listener
	delete        chan string
	messages      chan *Message
	mutex         sync.RWMutex
	closed        bool
	bufferSize    int
	broadcasterTimeout time.Duration
	enablePing    bool
}

func NewChannelManager(bufferSize int, broadcasterTimeout time.Duration, enablePing bool) *ChannelManager {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	if broadcasterTimeout <= 0 {
		broadcasterTimeout = time.Millisecond * 10
	}
	manager := &ChannelManager{
		channels:      make(map[string]broadcast.Broadcaster),
		open:          make(chan *Listener, bufferSize),
		close:         make(chan *Listener, bufferSize),
		delete:        make(chan string, bufferSize),
		messages:      make(chan *Message, bufferSize),
		bufferSize:    bufferSize,
		broadcasterTimeout: broadcasterTimeout,
		enablePing:    enablePing,
	}
	go manager.run()
	return manager
}

func (m *ChannelManager) run() {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("error", r).Stack().Msg("ChannelManager panicked, restarting")
			time.Sleep(time.Second)
			if !m.closed {
				go m.run()
			}
		}
	}()
	for {
		select {
		case listener, ok := <-m.open:
			if !ok {
				return
			}
			m.register(listener)
		case listener, ok := <-m.close:
			if !ok {
				return
			}
			m.deregister(listener)
		case channelId, ok := <-m.delete:
			if !ok {
				return
			}
			m.deleteBroadcast(channelId)
		case message, ok := <-m.messages:
			if !ok {
				return
			}
			m.getChannel(message.ChannelId).TrySubmit(message.Text)
			// metrics.Increment("channel_manager.messages_processed")
		}
	}
}

func (m *ChannelManager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.closed {
		return nil
	}
	m.closed = true
	close(m.open)
	close(m.close)
	close(m.delete)
	close(m.messages)
	for channelId, b := range m.channels {
		if err := b.Close(); err != nil {
			log.Error().Err(err).Str("channel", channelId).Msg("Failed to close broadcaster")
		}
		delete(m.channels, channelId)
	}
	log.Info().Msg("ChannelManager closed")
	return nil
}

func (m *ChannelManager) register(listener *Listener) {
	m.getChannel(listener.ChannelId).Register(listener.Chan)
	// metrics.Increment("channel_manager.listeners_registered")
}

func (m *ChannelManager) deregister(listener *Listener) {
	b := m.getChannel(listener.ChannelId)
	b.Unregister(listener.Chan)
	if !IsClosed(listener.Chan) {
		close(listener.Chan)
	}
	// metrics.Increment("channel_manager.listeners_deregistered")
}

func (m *ChannelManager) deleteBroadcast(channelId string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if b, ok := m.channels[channelId]; ok {
		if err := b.Close(); err != nil {
			log.Error().Err(err).Str("channel", channelId).Msg("Failed to close broadcaster")
		}
		delete(m.channels, channelId)
		// metrics.Increment("channel_manager.broadcasters_deleted")
	}
}

func (m *ChannelManager) IsExistsChannel(channelId string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, ok := m.channels[channelId]
	return ok
}

func (m *ChannelManager) getChannel(channelId string) broadcast.Broadcaster {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	b, ok := m.channels[channelId]
	if !ok {
		b = broadcast.NewBroadcaster(m.bufferSize, m.broadcasterTimeout)
		m.channels[channelId] = b
		// metrics.Increment("channel_manager.broadcasters_created")
	}
	return b
}

func (m *ChannelManager) OpenListener(prefix string, keys string) chan interface{} {
	if m.closed {
		log.Warn().Msg("Attempt to open listener on closed ChannelManager")
		return nil
	}
	s := strings.Split(keys, ",")
	listener := make(chan interface{}, m.bufferSize)
	for _, key := range s {
		ch := prefix + ":" + key
		m.open <- &Listener{
			ChannelId: ch,
			Chan:      listener,
		}
	}
	if m.enablePing {
		m.open <- &Listener{
			ChannelId: "PING",
			Chan:      listener,
		}
	}
	return listener
}

func (m *ChannelManager) CloseListener(prefix string, keys string, channel chan interface{}) {
	if m.closed {
		log.Warn().Msg("Attempt to close listener on closed ChannelManager")
		return
	}
	s := strings.Split(keys, ",")
	for _, key := range s {
		ch := prefix + ":" + key
		m.close <- &Listener{
			ChannelId: ch,
			Chan:      channel,
		}
	}
	if m.enablePing {
		m.close <- &Listener{
			ChannelId: "PING",
			Chan:      channel,
		}
	}
}

func (m *ChannelManager) DeleteBroadcast(prefix string, keys string) {
	if m.closed {
		log.Warn().Msg("Attempt to delete broadcast on closed ChannelManager")
		return
	}
	s := strings.Split(keys, ",")
	for _, key := range s {
		m.delete <- prefix + ":" + key
	}
}

func (m *ChannelManager) Submit(channelId string, text string) {
	if m.closed {
		log.Warn().Msg("Attempt to submit on closed ChannelManager")
		return
	}
	m.messages <- &Message{
		ChannelId: channelId,
		Text:      text,
	}
}

func IsClosed(ch <-chan interface{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}