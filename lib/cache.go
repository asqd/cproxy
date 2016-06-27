package cproxy

import (
	"encoding/hex"
	"flag"
	"io/ioutil"
	"os"
	"time"
)

var (
	timeout  = flag.Duration("timeout", 36*time.Hour, "Timeout unless clear file cache")
	cacheDir = flag.String("dir", "/home/deploy/apps/cache/", "Dir where page content cache will stored")
)

type Cache struct {
	CreatedAt time.Time `json:"created_at"`
	URL       string    `json:"url"`
	Store     Store     `json:"-"`
}

// Path to cache file
func (c *Cache) Path() string {
	return *cacheDir + c.Hex()
}

// return hex from URL
func (c *Cache) Hex() string {
	result := hex.EncodeToString([]byte(c.URL))

	return result
}

// Checks if Cache is need to be
// recached by it's CreatedAt field
func (c *Cache) ShouldRecache() bool {
	tnow := time.Now()

	tcreated := c.CreatedAt.Add(*timeout)

	return tcreated.Before(tnow)
}

// Checks if cache file exists
func (c *Cache) Exists() bool {
	if _, err := os.Stat(c.Path()); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Read Cache from file
func (c *Cache) Read() ([]byte, error) {
	return ioutil.ReadFile(c.Path())
}

// Write Cache to file
func (c *Cache) Write(body []byte) error {
	return ioutil.WriteFile(c.Path(), body, 0644)
}

// Remove cache file
func (c *Cache) Remove() error {
	return os.Remove(c.Path())
}

// Updates cache by removing and writing it
// and appending record to store
func (c *Cache) Update(body []byte) error {
	var err error

	if c.Exists() {
		err = c.Remove()

		if err != nil {
			return err
		}
	}

	err = c.Write(body)
	c.CreatedAt = time.Now()
	
	c.Store.Append(c)

	if err != nil {
		return err
	}

	return nil
}
