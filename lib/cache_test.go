package cproxy

import (
	"reflect"
	"testing"
	"time"
)

func TestExistsReadWriteRemoveUpdate(t *testing.T) {
	resetCacheDir(t)

	c := Cache{time.Now(), "http://google.com", getStore()}

	if c.Exists() {
		t.Error("Expected just created cache to be false got true")
	}

	err := c.Write([]byte{})

	if err != nil {
		t.Error(err)
	}

	if !c.Exists() {
		t.Error("Expected Exists to be true after creatring file")
	}

	expect := []byte("Cache")

	err = c.Update(expect)

	if err != nil {
		t.Error(err)
	}

	got, err := c.Read()

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expect, got) {
		t.Errorf("After rewriting cache expected content to be %v got %v", string(expect), string(got))
	}

	nc, err := c.Store.GetFirst()
	
	if err != nil {
		t.Error(err)
	}

	if nc.URL != c.URL {
		t.Errorf("After update in store should be created record with same url expect %v, got %v", c.URL, nc.URL)
	}

	err = c.Remove()

	if err != nil {
		t.Error(err)
	}

        if c.Exists() {
                t.Error("Expected cache to be false after removing got true")
        }	
}
