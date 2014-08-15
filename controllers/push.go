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
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
}

type wsAuth struct {
	Lat   string `json:"lat"`
	Lng   string `json:"lng"`
	Token string `json:"token"`
}

type wsAuthResp struct {
	Userid string `json:"userid"`
}

type pushData struct {
	Type string           `json:"type"`
	Id   string           `json:"id"`
	From string           `json:"from"`
	To   string           `json:"to"`
	Body []models.MsgBody `json:"body"`
}

type pushMsg struct {
	Type string   `json:"type"`
	Push pushData `json:"push"`
	Time int64    `json:"time"`
}

func (m *pushMsg) Bytes() []byte {
	b, _ := json.Marshal(m)
	return b
}

func BindWSPushApi(m *martini.ClassicMartini) {
	m.Get("/1/ws", wsPushHandler)
}

func wsPushHandler(request *http.Request, resp http.ResponseWriter, redisLogger *models.RedisLogger) {
	var wg sync.WaitGroup
	defer wg.Wait()

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

	psc := redisLogger.PubSub(user.Id, redisLogger.Groups(user.Id)...)

	go func(conn *websocket.Conn) {
		wg.Add(1)

		defer wg.Done()
		defer psc.Close()

		for {
			msg := &pushMsg{}
			err := conn.ReadJSON(msg)
			if err != nil {
				//log.Println(err)
				return
			}
			switch msg.Type {
			case "message":
				m := &models.Message{
					From: msg.Push.From,
					To:   msg.Push.To,
					Body: msg.Push.Body,
					Time: time.Now(),
				}
				if msg.Push.Type == "chat" || msg.Push.Type == "groupchat" {
					m.Type = msg.Push.Type
					m.Save()
					msg.Push.Id = m.Id.Hex()
					msg.Time = m.Time.Unix()

					redisLogger.PubMsg(m.Type, m.To, msg.Bytes())
				}
			default:
				log.Println("unhandled message type:", msg.Type)
			}
		}
	}(conn)

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			//log.Printf("%s: message: %s\n", v.Channel, v.Data)
			msg := &pushMsg{}
			if err := json.Unmarshal(v.Data, msg); err != nil {
				log.Println("parse push message error:", err)
				continue
			}

			// subscribe group
			if msg.Push.Type == "subscribe" && msg.Push.From == user.Id {
				if err := redisLogger.Subscribe(psc, msg.Push.To); err != nil {
					log.Println(err)
				}
				continue
			}
			// unsubscribe group
			if msg.Push.Type == "unsubscribe" && msg.Push.From == user.Id {
				if err := redisLogger.Unsubscribe(psc, msg.Push.To); err != nil {
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
