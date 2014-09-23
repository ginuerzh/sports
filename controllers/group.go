// group
package controllers

import (
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"net/http"
	"time"
)

func BindGroupApi(m *martini.ClassicMartini) {
	m.Post("/1/user/joinGroup", binding.Json(joinGroupForm{}), ErrorHandler, joinGroupHandler)
	m.Post("/1/user/newGroup", binding.Json(setGroupForm{}), ErrorHandler, setGroupHandler)
	m.Get("/1/user/getGroupInfo", binding.Form(groupInfoForm{}), ErrorHandler, groupInfoHandler)
	m.Get("/1/user/deleteGroup", binding.Json(groupDelForm{}), ErrorHandler, delGroupHandler)
}

type joinGroupForm struct {
	Gid   string `json:"group_id" binding:"required"`
	Leave bool   `json:"leave"`
	Token string `json:"access_token" binding:"required"`
}

func joinGroupHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, form joinGroupForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	redis.JoinGroup(user.Id, form.Gid, !form.Leave)
	writeResponse(request.RequestURI, resp, nil, nil)

	msg := &pushMsg{
		Type: "message",
		Push: pushData{
			From: user.Id,
			To:   form.Gid,
		},
	}
	if !form.Leave {
		msg.Push.Type = "subscribe"
	} else {
		msg.Push.Type = "unsubscribe"
	}

	redis.PubMsg(msg.Push.Type, user.Id, msg.Bytes())
}

type Group struct {
	Id          string   `json:"group_id"`
	Name        string   `json:"group_name"`
	Profile     string   `json:"group_image"`
	Desc        string   `json:"group_description"`
	Creator     string   `json:"group_creater_id"`
	Time        int64    `json:"create_time"`
	MemberCount int      `json:"members_count"`
	Members     []string `json:"member_ids"`
	Level       int      `json:"group_level"`
	models.Address
	models.Location
}

type setGroupForm struct {
	Group
	Token string `json:"access_token" binding:"required"`
}

func setGroupHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, form setGroupForm) {

	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	group := &models.Group{
		Gid:     form.Id,
		Name:    form.Name,
		Profile: form.Profile,
		Desc:    form.Desc,
		Creator: user.Id,
		Time:    time.Now(),
	}

	if form.Address.String() != "" {
		group.Addr = &form.Address
		loc := models.Addr2Loc(form.Address)
		group.Loc = &loc
	}

	var err error
	if len(form.Id) == 0 {
		err = group.Save()
		if err == nil {
			redis.JoinGroup(user.Id, group.Gid, true)
		}
	} else {
		err = group.Update()
	}

	writeResponse(request.RequestURI, resp, map[string]string{"group_id": group.Gid}, err)
}

type groupInfoForm struct {
	Gid   string `form:"group_id" binding:"required"`
	Token string `form:"access_token"`
}

func groupInfoHandler(request *http.Request, resp http.ResponseWriter, form groupInfoForm) {

	group := &models.Group{}
	err := group.FindById(form.Gid)

	grp := &Group{
		Id:          group.Gid,
		Name:        group.Name,
		Profile:     group.Profile,
		Desc:        group.Desc,
		Creator:     group.Creator,
		Time:        group.Time.Unix(),
		MemberCount: len(group.Members),
		Members:     group.Members,
		Level:       group.Level,
		Address:     *group.Addr,
		Location:    *group.Loc,
	}

	writeResponse(request.RequestURI, resp, grp, err)
}

type groupDelForm struct {
	Gid   string `json:"group_id" binding:"required"`
	Token string `json:"access_token" binding:"required"`
}

func delGroupHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, form groupDelForm) {

	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	group := &models.Group{Gid: form.Gid}
	group.Remove(user.Id)

	writeResponse(request.RequestURI, resp, nil, nil)
}
