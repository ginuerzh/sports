// chat
package controllers

import (
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"github.com/zhengying/apns"
	"gopkg.in/go-martini/martini.v1"
	"log"
	"net/http"
	"time"
)

func BindChatApi(m *martini.ClassicMartini) {
	m.Get("/1/chat/recent_chat_infos", binding.Form(contactsForm{}), ErrorHandler, contactsHandler)
	m.Post("/1/chat/send_message", binding.Json(sendMsgForm{}), ErrorHandler, sendMsgHandler)
	m.Get("/1/chat/get_list", binding.Form(msgListForm{}), ErrorHandler, msgListHandler)
}

type contactsForm struct {
	Token string `form:"access_token" binding:"required"`
}

func contactsHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form contactsForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	u := &models.User{}
	u.FindByUserid(user.Id)

	contacts := make([]*contactStruct, len(u.Contacts))
	for i, _ := range u.Contacts {
		contacts[i] = convertContact(&u.Contacts[i])
	}

	respData := map[string]interface{}{
		"contact_infos": contacts,
	}
	writeResponse(request.RequestURI, resp, respData, nil)
}

type sendMsgForm struct {
	Token   string `json:"access_token" binding:"required"`
	To      string `json:"to_id" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func sendMsgHandler(request *http.Request, resp http.ResponseWriter,
	client *apns.Client, redis *models.RedisLogger, form sendMsgForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	if redis.Relationship(user.Id, form.To) == models.RelBlacklist ||
		redis.Relationship(form.To, user.Id) == models.RelBlacklist {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	touser := &models.Account{}
	if find, err := touser.FindByUserid(form.To); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError, "user '"+form.To+"' not exists")
		}
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	msg := &models.Message{
		From: user.Id,
		To:   form.To,
		Body: []models.MsgBody{models.MsgBody{Type: form.Type, Content: form.Content}},
		Type: "chat",
		Time: time.Now(),
	}
	if err := msg.Save(); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	u := &models.User{Id: user.Id}
	contact := &models.Contact{
		Id:       touser.Id,
		Profile:  touser.Profile,
		Nickname: touser.Nickname,
		Last:     msg,
	}
	if err := u.AddContact(contact); err != nil {
		log.Println(err)
	}

	u.Id = touser.Id
	contact.Id = user.Id
	contact.Profile = user.Profile
	contact.Nickname = user.Nickname
	contact.Count = 1
	if err := u.AddContact(contact); err != nil {
		log.Println(err)
	}

	// ws push
	event := &models.Event{
		Type: models.EventMsg,
		Time: time.Now().Unix(),
		Data: models.EventData{
			Type: models.EventChat,
			Id:   user.Id,
			From: user.Id,
			To:   form.To,
			Body: []models.MsgBody{
				{Type: "msg_type", Content: form.Type},
				{Type: "msg_content", Content: form.Content},
				{Type: "nikename", Content: user.Nickname},
			},
		},
	}

	redis.PubMsg(models.EventMsg, form.To, event.Bytes())
	if err := event.Save(); err == nil {
		redis.IncrEventCount(form.To, event.Data.Type, 1)
	}

	devs, enabled, _ := u.Devices()
	if enabled {
		for _, dev := range devs {
			if err := sendApns(client, dev, user.Nickname+": "+msg.Body[0].Content, 1, ""); err != nil {
				log.Println(err)
			}
		}
	}

	writeResponse(request.RequestURI, resp, map[string]string{"message_id": msg.Id.Hex()}, nil)
}

type msgJsonStruct struct {
	Id      string `json:"message_id"`
	From    string `json:"from_id"`
	To      string `json:"to_id"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

func convertMsg(msg *models.Message) *msgJsonStruct {
	return &msgJsonStruct{
		Id:      msg.Id.Hex(),
		From:    msg.From,
		To:      msg.To,
		Type:    msg.Body[0].Type,
		Content: msg.Body[0].Content,
		Time:    msg.Time.Unix(),
	}
}

type msgListForm struct {
	Token  string `form:"access_token" binding:"required"`
	Userid string `form:"userid" binding:"required"`
	models.Paging
}

func msgListHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form msgListForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	u := &models.User{Id: user.Id}
	_, msgs, err := u.Messages(form.Userid, &form.Paging)
	jsonStructs := make([]*msgJsonStruct, len(msgs))
	for i, _ := range msgs {
		jsonStructs[i] = convertMsg(&msgs[i])
	}

	respData := make(map[string]interface{})
	respData["page_frist_id"] = form.Paging.First
	respData["page_last_id"] = form.Paging.Last
	//respData["page_item_count"] = total
	respData["messages"] = jsonStructs
	writeResponse(request.RequestURI, resp, respData, err)
}
