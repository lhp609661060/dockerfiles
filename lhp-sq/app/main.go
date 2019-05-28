package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"strconv"
	"time"
)

var redisAddr = flag.String("addr", "", "redis address")
var redisPassword = flag.String("password", "", "redis password")

func redisConn(add string, password string) (redis.Conn, error) {
	if password == "" {
		return redis.Dial("tcp", add)
	}
	return redis.Dial("tcp", add, redis.DialPassword(password))
}

func redispool(add string, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 20, // adjust to your needs
		IdleTimeout: 240 * time.Second, // adjust to your needs
		Dial: func() (redis.Conn, error) {
			fmt.Println("这里执行了")
			return redisConn(add, password)
		},
	}
}


/**
 * 启动函数
 */
func main() {

	flag.Parse()

	r := gin.Default()

	pool := redispool(*redisAddr, *redisPassword)

    defer pool.Close()

	//
	r.GET("/:name", func(c *gin.Context) {

		// 获取缓存连接
		conn := pool.Get()
		defer conn.Close()

		// 序列号对应的名称
		name := c.Param("name")

		// 位数
		n, ne := strconv.Atoi(c.DefaultQuery("n", "4"))
		if ne != nil {
			n = 4
		}

		now := time.Now()

		key := name + now.Format("20060102")

		i, _ := redis.Int(conn.Do("INCR", key))

		var buffer bytes.Buffer
		buffer.WriteString("%s")
		buffer.WriteString("%0")
		buffer.WriteString(strconv.Itoa(n))
		buffer.WriteString("d")

		c.String(http.StatusOK, buffer.String(), key, i)
	})

	_ = r.Run(":80")
}