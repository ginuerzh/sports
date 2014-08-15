// main
package main

import (
	"flag"
	"github.com/garyburd/redigo/redis"
	"github.com/ginuerzh/sports/controllers"
	//"github.com/martini-contrib/gzip"
	"github.com/zhengying/apns"
	"gopkg.in/go-martini/martini.v1"
	"log"
	"net/http"
	"os"
	//"strconv"
	"time"
)

var (
	staticDir  string
	listenAddr string
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.StringVar(&staticDir, "static", "public", "static files directory")
	flag.StringVar(&listenAddr, "l", ":8080", "addr on listen")
	flag.Parse()
}

func classic() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	//m.Use(gzip.All())
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(martini.Static(staticDir))
	m.Use(controllers.RedisLoggerHandler)
	m.Action(r.Handle)
	return &martini.ClassicMartini{m, r}
}

func main() {
	m := classic()
	m.Map(log.New(os.Stdout, "[sports] ", log.LstdFlags))
	m.Map(redisPool())
	m.Map(apnsClient())

	controllers.BindAccountApi(m)
	controllers.BindUserApi(m)
	controllers.BindArticleApi(m)
	controllers.BindChatApi(m)
	controllers.BindEventApi(m)
	controllers.BindFileApi(m)
	controllers.BindRecordApi(m)
	//controllers.BindStatApi(m)
	controllers.BindWSPushApi(m)
	controllers.BindGroupApi(m)

	//m.Run()
	http.ListenAndServe(listenAddr, m)
}

func redisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Println(err)
				return nil, err
			}
			/*
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			*/
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

var (
	dgw = "gateway.sandbox.push.apple.com:2195"
	gw  = "gateway.push.apple.com:2195"
)

func apnsClient() *apns.Client {
	return apns.ComboPEMClient(dgw, "apns.pem")
}
