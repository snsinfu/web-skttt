package domain

import (
	"reflect"
	"regexp"
	"testing"
)

func TestDomainUsage(t *testing.T) {
	nameLen := 7
	keyLen := 13
	safePathComponent := `^[0-9a-z]+$`

	dom, err := New(Config{
		NameLen: nameLen,
		KeyLen:  keyLen,
	})
	if err != nil {
		t.Errorf("failed to create domain: %v", err)
	}

	name, key, err := dom.CreateTopic()
	if err != nil {
		t.Errorf("failed to create topic: %v", err)
	}
	if len(name) != nameLen {
		t.Errorf("unexpected name length: %d, expected %d", len(name), nameLen)
	}
	if len(key) != keyLen {
		t.Errorf("unexpected key length: %d, expected %d", len(key), keyLen)
	}
	if ok, _ := regexp.MatchString(safePathComponent, name); !ok {
		t.Errorf("unsafe topic name: %v", name)
	}
	if ok, _ := regexp.MatchString(safePathComponent, key); !ok {
		t.Errorf("unsafe key: %v", key)
	}

	sub, err := dom.Subscribe(name)
	if err != nil {
		t.Errorf("failed to subscribe to topic: %v", err)
	}

	msg := Message{
		Data: []byte("Lorem ipsum dolor sit amet"),
	}
	if err := dom.PostToTopic(name, key, msg); err != nil {
		t.Errorf("failed to post to topic: %v", err)
	}

	select {
	case actual := <-sub.Read():
		if !reflect.DeepEqual(actual, msg) {
			t.Errorf("got unexpected message %v, expected %v", actual, msg)
		}
	default:
		t.Errorf("got no message")
	}

	dom.Unsubscribe(sub)
}
