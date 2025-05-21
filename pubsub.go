package model

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"sync"
)

type PubsubManager struct {
	channels       map[string]int
	subscribeRedis *redis.PubSub
	open           chan []string
	close          chan []string
	mutex          sync.RWMutex
	closed         bool
}

func NewPubsubManager(rdb *redis.PubSub, bufferSize int) *PubsubManager {
	if bufferSize <= 0 {
		bufferSize = 100 // Default buffer size
	}
	manager := &PubsubManager{
		channels:       make(map[string]int),
		subscribeRedis: rdb,
		open:           make(chan []string, bufferSize),
		close:          make(chan []string, bufferSize),
	}
	go manager.run()
	return manager
}

func (m *PubsubManager) run() {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("error", r).Stack().Msg("PubsubManager panicked, restarting")
			time.Sleep(time.Second) // Backoff
			if !m.closed {
				go m.run()
			}
		}
	}()
	for {
		select {
		case arr, ok := <-m.open:
			if !ok {
				return
			}
			m.openRedisPubsub(arr)
		case arr, ok := <-m.close:
			if !ok {
				return
			}
			m.closeRedisPubsub(arr)
		}
	}
}

func (m *PubsubManager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.closed {
		return nil
	}
	m.closed = true
	close(m.open)
	close(m.close)
	if err := m.subscribeRedis.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close Redis PubSub")
		return err
	}
	log.Info().Msg("PubsubManager closed")
	return nil
}

func (m *PubsubManager) changeCount(ch string, isInc bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if isInc {
		m.channels[ch]++
	} else {
		m.channels[ch]--
		if m.channels[ch] < 0 {
			m.channels[ch] = 0
		}
	}
	// metrics.Gauge("pubsub.channels", float64(len(m.channels)))
}

func (m *PubsubManager) getChannelCount(ch string) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.channels[ch]
}

func (m *PubsubManager) openRedisPubsub(array []string) {
	for _, ch := range array {
		if m.getChannelCount(ch) <= 0 {
			log.Info().Str("channel", ch).Msg("Subscribing to Redis channel")
			start := time.Now()
			if err := m.subscribeRedis.Subscribe(context.Background(), ch); err != nil {
				log.Error().Err(err).Str("channel", ch).Msg("Failed to subscribe")
				continue
			}
			// metrics.Observe("pubsub.subscribe_latency", time.Since(start).Seconds())
		}
		m.changeCount(ch, true)
	}
}

func (m *PubsubManager) closeRedisPubsub(array []string) {
	for _, ch := range array {
		m.changeCount(ch, false)
		if m.getChannelCount(ch) > 0 {
			continue
		}
		log.Info().Str("channel", ch).Msg("Unsubscribing from Redis channel")
		start := time.Now()
		if err := m.subscribeRedis.Unsubscribe(context.Background(), ch); err != nil {
			log.Error().Err(err).Str("channel", ch).Msg("Failed to unsubscribe")
			continue
		}
		// metrics.Observe("pubsub.unsubscribe_latency", time.Since(start).Seconds())
	}
}

func (m *PubsubManager) RegisterRedisPubsub(array []string) {
	if m.closed {
		log.Warn().Msg("Attempt to register on closed PubsubManager")
		return
	}
	m.open <- array
}

func (m *PubsubManager) UnregisterRedisPubsub(array []string) {
	if m.closed {
		log.Warn().Msg("Attempt to unregister on closed PubsubManager")
		return
	}
	m.close <- array
}