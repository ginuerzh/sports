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
	//"fmt"
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
	Userid     string `json:"userid"`
	LoginCount int    `json:"login_count"`
	LastLog    int64  `json:"last_login_time"`
	TimeLimit  int64  `json:"ban_time,omitempty"`
}

func BindWSPushApi(m *martini.ClassicMartini) {
	m.Get("/1/ws", wsPushHandler)
}

func checkTokenValid(token string) bool {
	if len(token) == 0 {
		return false
	}

	a := strings.Split(token, "-")
	if len(a) != 6 {
		return false
	}

	secs, err := strconv.ParseInt(a[5], 10, 64)
	if err != nil {
		return false
	}

	return time.Unix(secs, 0).After(time.Now())
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
	conn.ReadJSON(&auth)

	//log.Println("check token:", auth.Token)
	if !checkTokenValid(auth.Token) {
		//log.Println("check token valid")
		redisLogger.DelOnlineUser(auth.Token)
		conn.WriteJSON(r)
		return
	}

	uid := redisLogger.OnlineUser(auth.Token)

	user := &models.Account{}
	if find, _ := user.FindByUserid(uid); !find || user.TimeLimit < 0 {
		r.TimeLimit = user.TimeLimit
		conn.WriteJSON(r)
		return
	}

	redisLogger.LogLogin(user.Id)

	days := user.LoginDays
	loginCount := user.LoginCount + 1
	d := nowDate()
	if user.LastLogin.Unix() < d.Unix() { // check wether first time login of one day
		days++
		if user.LastLogin.Unix() < d.Unix()-24*3600 {
			days = 1
		}
		loginCount = 1
	}
	//fmt.Println(uid, "loginCount", loginCount)
	user.SetLastLogin(days, loginCount, time.Now())

	r.Userid = uid
	r.LastLog = user.LastLogin.Unix()
	r.LoginCount = loginCount

	if err := conn.WriteJSON(r); err != nil {
		return
	}

	if len(uid) == 0 {
		return
	}

	redisLogger.LogVisitor(user.Id)
	psc := redisLogger.PubSub(user.Id)

	go func(conn *websocket.Conn) {
		//wg.Add(1)
		//defer log.Println("ws thread closed")
		//defer wg.Done()
		redisLogger.SetOnline(user.Id, true, 0)
		start := time.Now()

		defer psc.Close()

		for {
			event := &models.Event{}
			err := conn.ReadJSON(event)
			if err != nil {
				//log.Println(err)

				dur := int64(time.Since(start) / time.Second)
				redisLogger.SetOnline(user.Id, false, dur)
				user.UpdateStat(models.StatOnlineTime, dur)
				return
			}
			//log.Println("recv msg:", event.Type)
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
				//fmt.Println(user.Id, event.Data.Body)
				switch event.Data.Type {
				case "loc":
					var lat, lng float64
					var locaddr string
					for _, body := range event.Data.Body {
						switch body.Type {
						case "latlng":
							//log.Println("latlng:", body.Content)
							loc := strings.Split(body.Content, ",")
							if len(loc) != 2 {
								break
							}
							lat, _ = strconv.ParseFloat(loc[0], 64)
							lng, _ = strconv.ParseFloat(loc[1], 64)
						case "locaddr":
							//log.Println("locaddr:", body.Content)
							locaddr = body.Content
						}
					}
					user.UpdateLocation(models.Location{Lat: lat, Lng: lng}, locaddr)
				case "device":
					for _, body := range event.Data.Body {
						switch body.Type {
						case "token":
							token := body.Content
							//log.Println("device token:", token)
							user.AddDevice(token)
						}
					}
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
