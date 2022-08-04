package app

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"log"
	"net/http"
)

func StartApp() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.0.101:26379",
		Password: "livegbs@2019", // no password set
		DB:       0,              // use default DB
	})
	ctx := context.Background()
	val, err := rdb.Get(ctx, "key").Result()
	switch {
	case err == redis.Nil:
		fmt.Println("key does not exist")
	case err != nil:
		fmt.Println("Get failed", err)
	case val == "":
		fmt.Println("value is empty")
	}
	cn := rdb.Conn(ctx)
	defer cn.Close()

	if err := cn.ClientSetName(ctx, "OneGBS").Err(); err != nil {
		panic(err)
	}

	name, err := cn.ClientGetName(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("client name", name)

	// There is no error because go-redis automatically reconnects on error.
	pubsub := rdb.Subscribe(ctx, "alarm")

	// Close the subscription when we are done.
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
		sendAlarm()
	}

	fmt.Println("Done.")
}

func sendAlarm() {
	resp, err := http.Get("http://localhost:58080/async/alarm?method=fireCameraAlarm")
	if err != nil {
		log.Println("请求失败，错误原因：", err)
		return
	}
	defer resp.Body.Close()
	// 200 OK
	fmt.Println("返回码：", resp.Status)
	fmt.Println("返回头：", resp.Header)
	if resp.StatusCode != 200 {
		log.Println("请求失败，返回码：", resp.StatusCode)
		return
	}
	buf := make([]byte, 1024)
	for {
		// 接收服务端信息
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		} else {
			fmt.Println("读取完毕")
			res := string(buf[:n])
			fmt.Println(res)
			break
		}
	}
}
