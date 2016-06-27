package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/Supro/cproxy/lib"
)

var (
	w        sync.WaitGroup
	logFile = flag.String("log", "/var/log/cproxy.log", "Log file location")
)

func main() {
	flag.Parse()

	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer f.Close()

	log.SetOutput(f)

	w.Add(1)

	store := cproxy.RedisStore{cproxy.ConnectRedis(), &sync.Mutex{}}
	defer store.Conn.Close()
	log.Println("Database connected")

	go cproxy.RunServer(store)
	log.Println("Server started")

	cw := cproxy.CacheWorker{store}
	go cw.Run()
	log.Println("Cache worker started")

	w.Wait()
}
