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
}

type eventNewsForm struct {
	parameter
}

func eventNewsHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

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

func eventDetailHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	events, err := models.Events(user.Id)
	if err != nil {
		log.Println(err)
	}

	news := []*models.Event{}
	m := make(map[string]*models.Event) // TODO: don't use map

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

type changeEventStatusForm struct {
	Type string `json:"type" binding:"required"`
	Id   string `json:"id" binding:"required"`
	parameter
}

func changeEventStatusHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(changeEventStatusForm)

	count := user.ClearEvent(form.Type, form.Id)
	if form.Type == models.EventChat {
		//u := &models.User{Id: user.Id}
		user.MarkRead(form.Type, form.Id)
	}
	redis.IncrEventCount(user.Id, form.Type, -count)
	writeResponse(request.RequestURI, resp, nil, nil)
}
