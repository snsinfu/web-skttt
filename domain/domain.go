package domain

import (
	"errors"
	"sync"
	"time"
)

type Config struct {
	NameLen    int
	KeyLen     int
	Expiration time.Duration
}

type Domain struct {
	nameLen    int
	keyLen     int
	expiration time.Duration
	topics     map[string]*Topic
	mutex      sync.RWMutex
}

type Topic struct {
	key         string
	subscribers map[*Subscriber]bool
	mutex       sync.RWMutex
}

type Subscriber struct {
	topic   *Topic
	message chan Message
}

type Message struct {
	Data []byte
}

var (
	ErrTopicNotFound  = errors.New("topic not found")
	ErrTopicCollision = errors.New("topic name collision")
	ErrInvalidKey     = errors.New("invalid key")
)

func New(config Config) (*Domain, error) {
	dom := &Domain{
		nameLen:    config.NameLen,
		keyLen:     config.KeyLen,
		expiration: config.Expiration,
		topics:     map[string]*Topic{},
	}
	return dom, nil
}
