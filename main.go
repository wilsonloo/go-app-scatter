package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

/**存放配置数据key-value形式*/
var CFTablesMap map[string]interface{} // {}为初始化成空
var baseTime time.Time

///////////////////////////////////////////////////
// 使用 channel 自带的阻塞机制实现 连接池，替代cpp的信号量机制
var redis_conn_pool chan redis.Conn

// 进行扩充连接数的互斥锁
var reserve_mutex sync.Mutex

var (
	// 最大连接数
	LimitRedisConns = 1000
)

func InitTestMutexSpinLock() {
	redis_conn_pool = make(chan redis.Conn, LimitRedisConns)

	reserve_redis_client_pool(LimitRedisConns)
}

func reserve_redis_client_pool(new_count int) {

	for i := 0; i < new_count; i++ {

		client, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("failed to create redis-client:", err.Error())
			os.Exit(1)
		}

		redis_conn_pool <- client
	}

}

func TestMutexSpinLock(val int) {

	conn := <-redis_conn_pool
	if conn == nil {
		fmt.Println("failed to pop new client")
		os.Exit(1)
	}
	defer func() {
		redis_conn_pool <- conn
	}()

	_, err := conn.Do("LPUSH", "my_list", val)
	if err != nil {
		fmt.Println("failed:", err)
		return
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

	var N = 10000
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
			fmt.Println("index: ", index, "consumed: ", tmpint)

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
