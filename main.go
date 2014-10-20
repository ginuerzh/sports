// main
package main

import (
	"flag"
	"github.com/garyburd/redigo/redis"
	"github.com/ginuerzh/sports/controllers"
	"github.com/ginuerzh/sports/models"
	"github.com/zhengying/apns"
	//"github.com/martini-contrib/gzip"
	"gopkg.in/ginuerzh/weedo.v0"
	"gopkg.in/go-martini/martini.v1"
	"log"
	"net/http"
	"os"
	"strings"
	//"strconv"
	"time"
)


var (
	staticDir  string
	listenAddr string
	redisAddr  string
	weedfsAddr string
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.StringVar(&staticDir, "static", "public", "static files directory")
	flag.StringVar(&listenAddr, "l", ":8080", "addr on listen")
	flag.StringVar(&redisAddr, "redis", "localhost:6379", "redis server")
	flag.StringVar(&models.MongoAddr, "mongo", "localhost:27017", "mongodb server")
	flag.StringVar(&controllers.CoinAddr, "cs", "localhost:8087", "coin server")
	flag.StringVar(&weedfsAddr, "weed", "localhost:9334", "weed-fs server")
	flag.Parse()

	if !strings.HasPrefix(controllers.CoinAddr, "http") {
		controllers.CoinAddr = "http://" + controllers.CoinAddr
	}
	controllers.Weedfs = weedo.NewClient(weedfsAddr)
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
	controllers.BindWalletApi(m)
	controllers.BindTaskApi(m)

	//m.Run()
	log.Fatal(http.ListenAndServe(listenAddr, m))
}

func redisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddr)
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
