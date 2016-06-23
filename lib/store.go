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

	*sync.Mutex
}

// Redis implementation of appending to list
func (s RedisStore) Append(c *Cache) {
	cont, err := json.Marshal(c)

	if err != nil {
		return
	}

	_, err = s.Conn.Do("RPUSH", *table, string(cont))

	if err != nil {
		return
	}
}

// Redis implementation of getting first record from list
// and removing it from it
func (s RedisStore) GetFirst() (*Cache, error) {
	m := &Cache{Store: s}

	first, err := redis.String(s.Conn.Do("LPOP", *table))

	if err != nil {
		return m, err
	}

	err = json.Unmarshal([]byte(first), m)

	if err != nil {
		return m, err
	}

	return m, nil
}

// Creates connection to Redis
func ConnectRedis() redis.Conn {
	conn, err := redis.Dial("tcp", *redisPort)

	if err != nil {
		log.Fatalln(err)
	}

	return conn
}
