// event
package controllers

import (
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"log"
	"net/http"
	"strconv"
)

const (
	EventListV1Uri     = "/1/event/timeline"
	NewEventCountV1Uri = "/1/event/news"
	EventReadV1Uri     = "/1/event/change_status_read"
)

func BindEventApi(m *martini.ClassicMartini) {
	m.Get("/1/event/news", binding.Form(eventNewsForm{}), ErrorHandler, eventNewsHandler)
	m.Get("/1/event/news_details", binding.Form(eventNewsForm{}), ErrorHandler, eventDetailHandler)
	m.Post("/1/event/change_status_read", binding.Json(changeEventStatusForm{}), ErrorHandler, changeEventStatusHandler)
}

type eventNewsForm struct {
	Token string `form:"access_token"  binding:"required"`
}

func eventNewsHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form eventNewsForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	u := &models.User{}
	u.FindByUserid(user.Id)

	counts := redis.EventCount(user.Id)
	respData := map[string]int{
		"new_chat_count":      counts[0],
		"new_comment_count":   counts[1],
		"new_thumb_count":     counts[2],
		"new_reward_count":    counts[3] + counts[5],
		"new_attention_count": counts[4],
	}

	writeResponse(request.RequestURI, resp, respData, nil)

	//redis.LogOnlineUser(form.Token, user)
}

type contactStruct struct {
	Id       string         `json:"userid"`
	Profile  string         `json:"user_profile_image"`
	Nickname string         `json:"nikename"`
	Count    int            `json:"new_message_count"`
	Last     *msgJsonStruct `json:"last_message"`
}

func convertContact(contact *models.Contact) *contactStruct {
	return &contactStruct{
		Id:       contact.Id,
		Profile:  contact.Profile,
		Nickname: contact.Nickname,
		Count:    contact.Count,
		Last:     convertMsg(contact.Last),
	}
}

func eventDetailHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form eventNewsForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}
	/*
		u := &models.User{}
		u.FindByUserid(user.Id)


		contacts := []*contactStruct{}
		for _, c := range u.Contacts {
			if c.Count > 0 {
				contacts = append(contacts, convertContact(&c))
			}
		}


		ids := []string{}
		for _, event := range u.Events {
			ids = append(ids, event.Id)
		}
		articles, _ := models.FindArticles(ids...)

		events := []*articleJsonStruct{}
		for i, _ := range articles {
			event := convertArticle(&articles[i])
			for _, e := range u.Events {
				if e.Id == articles[i].Id.Hex() {
					event.NewReviews = len(e.Reviews)
					event.NewThumbs = len(e.Thumbs)
					break
				}
			}
			if event.NewReviews > 0 || event.NewThumbs > 0 {
				events = append(events, event)
			}
		}
	*/
	events, err := models.Events(user.Id)
	if err != nil {
		log.Println(err)
	}

	news := []*models.Event{}
	m := make(map[string]*models.Event)

	for i, event := range events {
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

	for _, v := range m {
		news = append(news, v)
	}

	respData := map[string]interface{}{
		"event_news": news,
	}

	writeResponse(request.RequestURI, resp, respData, nil)
}

type changeEventStatusForm struct {
	Token string `json:"access_token"  binding:"required"`
	Type  string `json:"type" binding:"required"`
	Id    string `json:"id" binding:"required"`
}

func changeEventStatusHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, form changeEventStatusForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	count := user.ClearEvent(form.Type, form.Id)
	if form.Type == models.EventChat {
		u := &models.User{Id: user.Id}
		u.MarkRead(form.Type, form.Id)
	}
	redis.IncrEventCount(user.Id, form.Type, -count)
	writeResponse(request.RequestURI, resp, nil, nil)
}
