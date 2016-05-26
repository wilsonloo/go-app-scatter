package main

import (
	"container/list"
	"fmt"
	"os"
	"sync"
)

var reIDsMutex sync.Mutex
var reIDsClients list.List

var (
	redisCount = 4
)

func reserve_redis_client_pool() {

	/*
		for i := 0; i < redisCount; i++ {
			client := redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("127.0.0.1:6379"),
				Password: "",
				DB:       redisDb,
			})

			if client != nil {
				fmt.Println("failed to create redis-client")
				os.Exit(1)
			}

			reIDsClients.PushBack(client)
		}
	*/
}

func TestMutexSpinLock() {

	client := pop_redis_client()
	if !client {
		fmt.Println("failed to pop new client")
		os.Exit(1)
	}

	cal := 1
	cal++

	push_redis_client()
}

func pop_redis_client() bool {
	reIDsMutex.Lock()
	defer reIDsMutex.Unlock()

	for i := 0; i < 100; i++ {
		k := 3
		k++
	}
	/*
		if reIDsClients.Len() == 0 {
			reserve_redis_client_pool()
		}

		client := reIDsClients.Front().Value.(*redis.Client)
		reIDsClients.Remove(reIDsClients.Front())
	*/

	return true
}

func push_redis_client() {
	reIDsMutex.Lock()
	defer reIDsMutex.Unlock()

	for i := 0; i < 100; i++ {
		k := 3
		k++
	}

	/*
		reIDsClients.PushBack(client)

		if redisCount > 100 && reIDsClients.Len() >= (redisCount*3/4) {
			count := redisCount / 2
			for i := 0; i < count && reIDsClients.Len() > 1; i++ {
				client := reIDsClients.Front().Value.(*redis.Client)
				reIDsClients.Remove(reIDsClients.Front())
				client.Close()
			}

			redisCount -= count
		}*/
}
