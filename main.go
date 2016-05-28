package main

import (
	"fmt"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

/**存放配置数据key-value形式*/
var CFTablesMap map[string]interface{} // {}为初始化成空
var baseTime time.Time

///////////////////////////////////////////////////
// 使用 channel 自带的阻塞机制实现 连接池，替代cpp的信号量机制
var redis_conn_pool chan redis.Conn

var (
	// 最大连接数
	LimitRedisConns = 1000

	// 当前工作中的连接数
	total_working_conns = 0
)

func InitTestMutexSpinLock() {
	redis_conn_pool = make(chan redis.Conn, LimitRedisConns)
}

func reserve_redis_client_pool(new_count int) {

	for i := 0; i < new_count && total_working_conns < LimitRedisConns; i++ {

		client, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("failed to create redis-client:", err.Error())
			os.Exit(1)
		}

		redis_conn_pool <- client
		total_working_conns++
	}

}

func TestMutexSpinLock(val int) {

	client := pop_redis_client()
	if client == nil {
		fmt.Println("failed to pop new client")
		os.Exit(1)
	}
	defer func() {
		// todo just test max connection :
		push_redis_client(client)
		// deleting_conn.PushBack(client)
	}()

	_, err := client.Do("LPUSH", "my_list", val)
	if err != nil {
		fmt.Println("failed:", err)
		return
	}
}

func pop_redis_client() redis.Conn {

	if len(redis_conn_pool) == 0 && total_working_conns < LimitRedisConns {
		reserve_redis_client_pool(total_working_conns)
	}

	conn := <-redis_conn_pool

	return conn
}

func push_redis_client(conn redis.Conn) {

	redis_conn_pool <- conn

	fmt.Printf("total_redis_conns: %d, redis_conn_pool size: %d, therold: %d\n",
		total_working_conns, len(redis_conn_pool), total_working_conns*3/4)

	// 移除多余连接
	if total_working_conns > 100 && len(redis_conn_pool) >= (total_working_conns*3/4) {
		count := total_working_conns / 2
		for i := 0; i < count && len(redis_conn_pool) > 1; i++ {

			client := <-redis_conn_pool
			client.Close()

			total_working_conns--
		}

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
			fmt.Println("pool size is", total_working_conns)

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
