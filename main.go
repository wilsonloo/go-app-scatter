package main

import (
	"sync"
	"time"

	"fmt"

	"container/list"
	"os"

	"github.com/garyburd/redigo/redis"
)

/**存放配置数据key-value形式*/
var CFTablesMap map[string]interface{} // {}为初始化成空
var baseTime time.Time

///////////////////////////////////////////////////

var reIDsMutex sync.Mutex
var redis_conn_pool list.List

var (
	total_redis_conns = 0
)

///////////////////////////////////////////////////

func reserve_redis_client_pool(new_count int) {

	for i := 0; i < new_count; i++ {

		client, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("failed to create redis-client:", err.Error())
			os.Exit(1)
		}

		redis_conn_pool.PushBack(client)
	}

	total_redis_conns += new_count
}

func TestMutexSpinLock(val int) {

	client := pop_redis_client()
	if client == nil {
		fmt.Println("failed to pop new client")
		os.Exit(1)
	}
	defer func() {
		push_redis_client(client)
	}()

	_, err := client.Do("LPUSH", "my_list", val)
	if err != nil {
		fmt.Println("failed:", err)
		return
	}
}

func InitTestMutexSpinLock() {
	reserve_redis_client_pool(4)
}

func pop_redis_client() redis.Conn {
	reIDsMutex.Lock()
	defer reIDsMutex.Unlock()

	if redis_conn_pool.Len() == 0 {
		reserve_redis_client_pool(total_redis_conns)
	}

	conn := redis_conn_pool.Front().Value.(redis.Conn)
	redis_conn_pool.Remove(redis_conn_pool.Front())

	return conn
}

func push_redis_client(client redis.Conn) {
	reIDsMutex.Lock()
	defer reIDsMutex.Unlock()

	redis_conn_pool.PushBack(client)

	// 移除多余连接
	if total_redis_conns > 100 && redis_conn_pool.Len() >= (total_redis_conns*3/4) {
		count := total_redis_conns / 2
		for i := 0; i < count && redis_conn_pool.Len() > 1; i++ {
			client := redis_conn_pool.Front().Value.(redis.Conn)
			redis_conn_pool.Remove(redis_conn_pool.Front())
			client.Close()
		}

		total_redis_conns -= count
	}
}

func main() {

	tnow := time.Now()
	tmptime, _ := time.Parse("2006-01-02 15:04:05", "2015-06-19 17:01:47.590")
	subdur := tnow.Sub(tmptime)
	fmt.Printf("subdur: %f,%s ", subdur.Nanoseconds(), tmptime, " ")
	datastr := fmt.Sprintf("%d%d%d%d%d%d", tnow.Year(), tnow.Month(), tnow.Day(), tnow.Hour(), tnow.Minute(), tnow.Second())

	fmt.Println(time.Now().Format("2006-01-02_15:04:05"), "   ", datastr)
	CFTablesMap = loadConfig()
	baseTime = time.Now()

	//////////////////////////////////////
	InitTestMutexSpinLock()
	//////////////////////////////////////

	var N = 9000
	sem := make(chan string, N)
	for i := 0; i < N; i++ {
		go func(index int) {

			tnow := time.Now()
			starttimeint := (int)(tnow.Sub(baseTime).Seconds() * 1000000)

			///////////////////////////////////////
			TestMutexSpinLock(index)
			///////////////////////////////////////

			tmptime := time.Now()
			subdur := tmptime.Sub(tnow)
			tmpint := (int)(subdur.Seconds() * 1000)
			fmt.Println("subdur:  ", tmpint)
			sem <- (fmt.Sprintf("%d,%d", starttimeint, tmpint))
		}(i)
	}

	//	var max = 0
	outString := "["
	for m := 0; m < N; m++ {
		tmp := <-sem
		outString += fmt.Sprintf("[%s]", tmp)
		if m == N-1 {

		} else {
			outString += ","
		}
	}
	outString += "]"
	writeFileWithData("./config/output.html", outString, N)

}
