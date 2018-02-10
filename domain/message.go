package domain

func (dom *Domain) PostToTopic(name, key string, msg Message) error {
	dom.mutex.RLock()
	defer dom.mutex.RUnlock()

	topic, ok := dom.topics[name]
	if !ok {
		return ErrTopicNotFound
	}

	if key != topic.key {
		return ErrInvalidKey
	}

	topic.mutex.RLock()
	defer topic.mutex.RUnlock()

	for sub := range topic.subscribers {
		sub.message <- msg
	}

	return nil
}

func (sub *Subscriber) Read() <-chan Message {
	return sub.message
}
