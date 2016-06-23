package cproxy

import (
	"reflect"
	"testing"
	"time"
)

func TestRedisStoreAppendGetFirst(t *testing.T) {
	store := getStore()

	_, err := store.Conn.Do("DEL", *table) // delete key from redis before test

	defer store.Conn.Close()

	if err != nil {
		t.Errorf("%v", err)
	}

	m := &Cache{time.Now(), "http://hello.world", nil}

	store.Append(m)

	nm, err := store.GetFirst()

	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(m, nm) {
		t.Errorf("Expected %v to equal %v", m, nm)
	}
}
