package main

import (
	"sync"
)

type Message struct {
	msg       string
	topicName string
}

type Topic struct {
	name   string
	mu     sync.Mutex
	subs   []chan Message
	closed bool
}

type Broker struct {
	mu     sync.RWMutex
	topics map[string]*Topic
}

func NewTopic(name string) *Topic {
	return &Topic{name: name}
}

func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]*Topic),
	}
}

func (b *Broker) getOrCreate(name string) *Topic {
	b.mu.Lock()
	defer b.mu.Unlock()

	if t, ok := b.topics[name]; ok {
		return t
	}

	t := NewTopic(name)
	b.topics[name] = t
	return t
}

type Subscriber struct {
	id    int
	inbox chan Message
}

func (b *Broker) Subscribe(topicName string, subscriber *Subscriber) {
	topic := b.getOrCreate(topicName)
	topic.mu.Lock()
	defer topic.mu.Unlock()
	if topic.closed {
		return
	}
	topic.subs = append(topic.subs, subscriber.inbox)
}

func (b *Broker) Publish(topicName string, msg string) {
	topic := b.getOrCreate(topicName)
	topic.mu.Lock()
	channels := make([]chan Message, len(topic.subs))
	copy(channels, topic.subs)
	topic.mu.Unlock()

	message := Message{
		msg:       msg,
		topicName: topicName,
	}

	for _, ch := range channels {
		ch <- message
	}
}
