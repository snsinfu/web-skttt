package domain

func (dom *Domain) Subscribe(name string) (*Subscriber, error) {
	dom.mutex.RLock()
	defer dom.mutex.RUnlock()

	topic, ok := dom.topics[name]
	if !ok {
		return nil, ErrTopicNotFound
	}

	topic.mutex.Lock()
	defer topic.mutex.Unlock()

	sub := &Subscriber{
		topic:   topic,
		message: make(chan Message, 1),
	}
	topic.subscribers[sub] = true

	return sub, nil
}

func (dom *Domain) Unsubscribe(sub *Subscriber) {
	topic := sub.topic

	topic.mutex.Lock()
	defer topic.mutex.Unlock()

	delete(topic.subscribers, sub)
}
