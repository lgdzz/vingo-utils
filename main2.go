package main

import (
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo"
	"github.com/lgdzz/vingo-utils/vingo/queue"
	"time"
)

func main() {

	vingo.InitRedisService(&vingo.RedisConfig{Host: "127.0.0.1", Port: "6379", Select: 0})
	//db.InitMysqlService(&db.MysqlConfig{Host: "127.0.0.1", Port: "3306", Dbname: "shs", Username: "root", Password: "123456789"})

	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println(queue.Redis.PushDelay("test", vingo.FileInfoSimple{
				Name:     fmt.Sprintf("名称_%v", i),
				Realpath: fmt.Sprintf("路径_%v", i),
			}, 10))
		}(i)
	}

	fmt.Println("Success")

	time.Sleep(time.Second * 1000)
}
