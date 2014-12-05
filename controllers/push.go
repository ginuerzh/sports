// push
package controllers

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/gorilla/websocket"
	"gopkg.in/go-martini/martini.v1"
	"log"
	"net/http"
	"strconv"
	"strings"
	//"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
}

type wsAuth struct {
	Token string `json:"token"`
}

type wsAuthResp struct {
	Userid string `json:"userid"`
}

func BindWSPushApi(m *martini.ClassicMartini) {
	m.Get("/1/ws", wsPushHandler)
}

func wsPushHandler(request *http.Request, resp http.ResponseWriter, redisLogger *models.RedisLogger) {
	conn, err := upgrader.Upgrade(resp, request, nil)
	if err != nil {
		conn.WriteJSON(errors.NewError(errors.HttpError, err.Error()))
		return
	}
	defer conn.Close()

	r := wsAuthResp{}
	var auth wsAuth
	if err := conn.ReadJSON(&auth); err != nil {
		conn.WriteJSON(r)
		log.Println(auth.Token)
		return
	}
	//log.Println(auth.Token)
	user := redisLogger.OnlineUser(auth.Token)
	if user != nil {
		r.Userid = user.Id
	}
	if err := conn.WriteJSON(r); err != nil {
		return
	}

	if user == nil {
		return
	}

	redisLogger.SetOnline(user.Id)
	redisLogger.LogVisitor(user.Id)

	psc := redisLogger.PubSub(user.Id, redisLogger.Groups(user.Id)...)

	go func(conn *websocket.Conn) {
		//wg.Add(1)
		defer log.Println("ws thread closed")
		//defer wg.Done()
		defer psc.Close()

		for {
			event := &models.Event{}
			err := conn.ReadJSON(event)
			if err != nil {
				//log.Println(err)
				return
			}
			log.Println("recv msg:", event.Type)
			switch event.Type {
			case models.EventMsg:
				m := &models.Message{
					From: event.Data.From,
					To:   event.Data.To,
					Body: event.Data.Body,
					Time: time.Now(),
				}
				if event.Data.Type == models.EventChat || event.Data.Type == models.EventGChat {
					m.Type = event.Data.Type
					m.Save()
					event.Data.Id = m.Id.Hex()
					event.Time = m.Time.Unix()

					redisLogger.PubMsg(m.Type, m.To, event.Bytes())
				}
			case "status":
				if event.Data.Type == "loc" && len(event.Data.Body) > 0 {
					log.Println("loc", event.Data.Body[0].Content)
					loc := strings.Split(event.Data.Body[0].Content, ",")
					if len(loc) != 2 {
						break
					}
					lat, _ := strconv.ParseFloat(loc[0], 64)
					lng, _ := strconv.ParseFloat(loc[1], 64)

					user.UpdateLocation(models.Location{Lat: lat, Lng: lng})
				}
			default:
				log.Println("unhandled message type:", event.Type)
			}
		}
	}(conn)

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			//log.Printf("%s: message: %s\n", v.Channel, v.Data)
			event := &models.Event{}
			if err := json.Unmarshal(v.Data, event); err != nil {
				log.Println("parse push message error:", err)
				continue
			}

			// subscribe group
			if event.Data.Type == models.EventSub && event.Data.From == user.Id {
				if err := redisLogger.Subscribe(psc, event.Data.To); err != nil {
					log.Println(err)
				}
				continue
			}
			// unsubscribe group
			if event.Data.Type == models.EventUnsub && event.Data.From == user.Id {
				if err := redisLogger.Unsubscribe(psc, event.Data.To); err != nil {
					log.Println(err)
				}
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, v.Data); err != nil {
				log.Println(err)
				return
			}
		case redis.Subscription:
			//log.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			//log.Println(v)
			return
		}
	}
}
