package cproxy

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/garyburd/redigo/redis"
)

var tw sync.WaitGroup

func randNumber(nums *[]int) int {
	i := rand.Int()

	for _, num := range *nums {
		if num == i {
			return randNumber(nums)
		}
	}

	*nums = append(*nums, i)
	return i
}

func TestConcurrency(t *testing.T) {
	pr := &testProxy{Code: 200, Body: "Prerender content"}
	ps := mockPrerenderServer(pr)

	def := *prerenderUrl
	prerenderUrl = &ps.URL

	defer func() {
		prerenderUrl = &def
	}()

	defer ps.Close()

	store := getStore()
	ser := Server{store}

	defer store.Conn.Close()

	s := httptest.NewServer(ser)

	_, err := store.Conn.Do("DEL", *table)

	if err != nil {
		t.Error(err)
	}

	defer s.Close()

	dir := resetCacheDir(t)
	defer os.Remove(dir)

	var nums = []int{}

	expect := 300

	for i := 1; i <= expect; i++ {
		tw.Add(1)
		go func() {
			w := httptest.NewRecorder()
			r, err := http.NewRequest("GET", s.URL+"/"+strconv.Itoa(randNumber(&nums)), nil)

			if err != nil {
				t.Error(err)
			}

			ser.ServeHTTP(w, r)

			c, err := store.GetFirst()

			if err != nil {
				t.Error(err)
			}

			store.Append(c)

			tw.Done()
		}()
	}

	tw.Wait()

	got, err := redis.Int(store.Conn.Do("LLEN", *table))

	if err != nil {
		t.Error(err)
	}

	if got != expect {
		t.Errorf("Expected request to create record in redis, expect %v, got %v", expect, got)
	}
}
