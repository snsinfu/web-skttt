package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func randomHex(size int) (string, error) {
	bytes := make([]byte, (size+1)/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:size], nil
}

func (dom *Domain) CreateTopic() (string, string, error) {
	random, err := randomHex(dom.nameLen + dom.keyLen)
	if err != nil {
		return "", "", err
	}
	name := random[:dom.nameLen]
	key := random[dom.nameLen:]

	dom.mutex.Lock()
	defer dom.mutex.Unlock()

	if _, exists := dom.topics[name]; exists {
		return name, "", ErrTopicCollision
	}

	dom.topics[name] = &Topic{
		key:         key,
		subscribers: map[*Subscriber]bool{},
	}

	if dom.expiration > 0 {
		go func() {
			time.Sleep(dom.expiration)
			dom.DeleteTopic(name)
		}()
	}

	return name, key, nil
}

func (dom *Domain) DeleteTopic(name string) error {
	dom.mutex.Lock()
	defer dom.mutex.Unlock()

	topic, ok := dom.topics[name]
	if !ok {
		return ErrTopicNotFound
	}

	delete(dom.topics, name)

	topic.mutex.RLock()
	defer topic.mutex.RUnlock()

	for sub := range topic.subscribers {
		close(sub.message)
	}

	return nil
}
