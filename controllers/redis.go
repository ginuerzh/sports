// redis
package controllers

import (
	"github.com/garyburd/redigo/redis"
	"github.com/ginuerzh/sports/models"
	"gopkg.in/go-martini/martini.v1"
	"net/http"
	//"strings"
	//"fmt"
)

func RedisLoggerHandler(request *http.Request, c martini.Context, pool *redis.Pool) {
	logger := models.NewRedisLogger(pool, pool.Get())
	defer logger.Close()

	/*
		s := strings.Split(request.RemoteAddr, ":")
		if len(s) > 0 {
			logger.LogVisitor(s[0])
		}
	*/
	//logger.LogPV(request.URL.Path)

	c.Map(logger)
	c.Next()
}
