// event
package controllers

import (
	//"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	EventListV1Uri     = "/1/event/timeline"
	NewEventCountV1Uri = "/1/event/news"
	EventReadV1Uri     = "/1/event/change_status_read"
)

func BindEventApi(m *martini.ClassicMartini) {
	m.Get("/1/event/news",
		binding.Form(eventNewsForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		eventNewsHandler)
	m.Get("/1/event/news_details",
		binding.Form(eventNewsForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		eventDetailHandler)
	m.Post("/1/event/change_status_read",
		binding.Json(changeEventStatusForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		changeEventStatusHandler)
	m.Get("/1/event/notices",
		binding.Form(eventNoticesForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		eventNoticesHandler)
}

type eventNewsForm struct {
	parameter
}

func eventNewsHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	//counts := redis.EventCount(user.Id)
	respData := map[string]int{
		"new_chat_count": user.EventCount("", models.EventChat),
		"new_comment_count": user.EventCount("", models.EventComment) +
			user.EventCount("", models.EventCoach) + user.EventCount("", models.EventCoachPass) + user.EventCount("", models.EventCoachNPass),
		"new_thumb_count": user.EventCount("", models.EventThumb),
		//"new_reward_count":    user.EventCount(models.EventReward) + user.EventCount(models.EventTx),
		"new_reward_count":    user.EventCount(models.EventTx, ""),
		"new_attention_count": user.EventCount("", models.EventSub) + user.EventCount(models.EventSystem, ""),
	}

	writeResponse(request.RequestURI, resp, respData, nil)

	//redis.LogOnlineUser(form.Token, user)
}

func eventDetailHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	events, err := models.Events(user.Id)
	if err != nil {
		log.Println(err)
	}

	news := []*models.Event{}
	m := make(map[string]*models.Event) // TODO: don't use map

	for i, event := range events {
		if event.Data.Type == models.EventChat {
			continue
		}
		key := event.Data.Type + "_" + event.Data.Id
		if e, ok := m[key]; ok {
			count, err := strconv.Atoi(e.Data.Body[len(e.Data.Body)-1].Content)
			if err != nil {
				log.Println(err)
			}
			e.Data.Body[len(e.Data.Body)-1].Content = strconv.Itoa(count + 1)
		} else {
			events[i].Data.Body = append(events[i].Data.Body, models.MsgBody{Type: "new_count", Content: "1"})
			m[key] = &events[i]
		}
	}

	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		news = append(news, m[key])
	}

	respData := map[string]interface{}{
		"event_news": news,
	}

	writeResponse(request.RequestURI, resp, respData, nil)
}

type eventNoticesForm struct {
	parameter
}

func eventNoticesHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {
	notices, _ := models.NoticeList(user.Id)
	writeResponse(r.RequestURI, w, map[string]interface{}{"notices": notices}, nil)
	var ids []interface{}
	for i, _ := range notices {
		ids = append(ids, notices[i].Id)
	}
	models.RemoveEvents(ids...)
}

type changeEventStatusForm struct {
	Type string `json:"type" binding:"required"`
	Id   string `json:"id" binding:"required"`
	parameter
}

func changeEventStatusHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(changeEventStatusForm)

	event := &models.Event{}
	event.Data.Type = strings.ToLower(form.Type)
	event.Data.Id = form.Id
	event.Data.To = user.Id
	event.Clear()
	//count := user.ClearEvent(form.Type, form.Id)
	/*
		if form.Type == models.EventChat { //TODO
			//u := &models.User{Id: user.Id}
			user.MarkRead(form.Type, form.Id)
		}
	*/
	//redis.IncrEventCount(user.Id, form.Type, -count)
	writeResponse(request.RequestURI, resp, nil, nil)
}
