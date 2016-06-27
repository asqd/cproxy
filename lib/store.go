package cproxy

import (
	"encoding/json"
	"flag"
	"log"
	"sync"

	"github.com/garyburd/redigo/redis"
)

var (
	redisPort = flag.String("redis-port", ":6379", "Redis tcp port")
	table     = flag.String("table", "recache_urls", "Table where urls json stored")
)

// Interface for working with different databases
// wich is storing Cache elements
type Store interface {
	// Should return first record
	// and remove it from list
	GetFirst() (*Cache, error)

	// Should append element to end of list
	Append(*Cache)

	// Mutex lock
	Lock()

	// Mutex unlock
	Unlock()
}

// Store interface impelementation for Redis DB
type RedisStore struct {
	// Redis connection
	Conn redis.Conn

	// Redis exec lock
	*sync.Mutex
}

// Redis implementation of appending to list
func (s *RedisStore) Append(c *Cache) {
	cont, err := json.Marshal(c)

	if err != nil {
		log.Println(err)
		return
	}

	s.Lock()
	_, err = s.Conn.Do("RPUSH", *table, string(cont))
	s.Unlock()

	if err != nil {
		s.TryRestore()
		s.Append(c)
	}
}

// Redis implementation of getting first record from list
// and removing it from it
func (s *RedisStore) GetFirst() (*Cache, error) {
	m := &Cache{Store: s}

	s.Lock()
	first, err := redis.String(s.Conn.Do("LPOP", *table))
	s.Unlock()

	if err != nil {
		s.TryRestore()
		return s.GetFirst()	
	}

	err = json.Unmarshal([]byte(first), m)

	if err != nil {
		return m, err
	}

	return m, nil
}

// Try to restore connection if it got error
func (s *RedisStore) TryRestore() {
	if err := s.Conn.Err(); if != nil {
		log.Println("Reconnectin redis")
		s.Conn.Close()
		s.Conn = ConnectRedis()
	}
}


// Creates connection to Redis
func ConnectRedis() redis.Conn {
	conn, err := redis.Dial("tcp", *redisPort)

	if err != nil {
		log.Fatalln(err)
	}

	return conn
}

