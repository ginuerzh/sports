// user
package controllers

import (
	//"encoding/json"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"io/ioutil"
	//"log"
	//"math/rand"
	"net/http"
	//"net/url"
	//"strconv"
	//"strings"
	"time"
)

func BindUserApi(m *martini.ClassicMartini) {
	m.Get("/1/user/articles", binding.Form(userArticlesForm{}), ErrorHandler, userArticlesHandler)
	m.Post("/1/user/send_device_token", binding.Json(sendDevForm{}, (*GetToken)(nil)), ErrorHandler, CheckHandler, sendDevHandler)
	m.Post("/1/user/set_push_enable", binding.Json(setPushForm{}, (*GetToken)(nil)), ErrorHandler, CheckHandler, setPushHandler)
	m.Get("/1/user/is_push_enabled", binding.Form(pushStatusForm{}), ErrorHandler, pushStatusHandler)
	m.Post("/1/user/enableAttention", binding.Json(relationshipForm{}, (*GetToken)(nil)), ErrorHandler, CheckHandler, followHandler)
	m.Post("/1/user/enableDefriend", binding.Json(relationshipForm{}, (*GetToken)(nil)), ErrorHandler, CheckHandler, blacklistHandler)
	m.Get("/1/user/getAttentionFriendsList", binding.Form(getFollowsForm{}), ErrorHandler, getFollowsHandler)
	m.Get("/1/user/getAttentedMembersList", binding.Form(getFollowsForm{}), ErrorHandler, getFollowersHandler)
	m.Get("/1/user/getJoinedGroupsList", binding.Form(getFollowsForm{}), ErrorHandler, getGroupsHandler)
	m.Get("/1/user/getRelatedMembersList", binding.Form(socialListForm{}), ErrorHandler, socialListHandler)
}

type userArticlesForm struct {
	Id    string `form:"userid" binding:"required"`
	Type  string `form:"article_type"`
	Token string `form:"access_token"`
	models.Paging
}

func userArticlesHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, form userArticlesForm) {

	user := &models.User{Id: form.Id}
	_, articles, err := user.Articles(form.Type, &form.Paging)

	jsonStructs := make([]*articleJsonStruct, len(articles))
	for i, _ := range articles {
		jsonStructs[i] = convertArticle(&articles[i])
	}

	respData := make(map[string]interface{})
	if len(articles) > 0 {
		respData["page_frist_id"] = form.Paging.First
		respData["page_last_id"] = form.Paging.Last
		//respData["page_item_count"] = total
	}
	respData["articles_without_content"] = jsonStructs

	writeResponse(request.RequestURI, resp, respData, err)
}

type sendDevForm struct {
	Token string `json:"access_token" binding:"required"`
	Dev   string `json:"device_token" binding:"required"`
}

func (this sendDevForm) getTokenId() string {
	return this.Token
}

func sendDevHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, getT GetToken) {

	form := getT.(sendDevForm)
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	u := &models.User{Id: user.Id}
	err := u.AddDevice(form.Token)
	writeResponse(request.RequestURI, resp, nil, err)
}

type setPushForm struct {
	Token   string `json:"access_token" binding:"required"`
	Enabled bool   `json:"is_enabled"`
}

func (this setPushForm) getTokenId() string {
	return this.Token
}

func setPushHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, getT GetToken) {

	form := getT.(setPushForm)
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	u := &models.User{Id: user.Id}
	err := u.SetPush(form.Enabled)
	writeResponse(request.RequestURI, resp, nil, err)
}

type pushStatusForm struct {
	Token string `json:"access_token" binding:"required"`
}

func pushStatusHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, form pushStatusForm) {

	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	u := &models.User{Id: user.Id}
	enabled, err := u.PushEnabled()
	writeResponse(request.RequestURI, resp, map[string]bool{"is_enabled": enabled}, err)
}

type relationshipForm struct {
	Userid    string `json:"userid"`
	Follow    bool   `json:"bAttention"`
	Blacklist bool   `json:"bDefriend"`
	Token     string `json:"access_token" binding:"required"`
}

func (this relationshipForm) getTokenId() string {
	return this.Token
}

func followHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, getT GetToken) {
	form := getT.(relationshipForm)
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	u := &models.Account{}
	if find, err := u.FindByUserid(form.Userid); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError)
		}
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	redis.SetRelationship(user.Id, form.Userid, models.RelFollowing, form.Follow)

	if form.Follow {
		event := &models.Event{
			Type: models.EventMsg,
			Time: time.Now().Unix(),
			Data: models.EventData{
				Type: models.EventSub,
				Id:   user.Id,
				From: user.Id,
				To:   u.Id,
				Body: []models.MsgBody{
					{Type: "nikename", Content: user.Nickname},
					{Type: "image", Content: user.Profile},
				},
			},
		}
		redis.PubMsg(models.EventMsg, u.Id, event.Bytes())
		if err := event.Save(); err == nil {
			redis.IncrEventCount(u.Id, event.Data.Type, 1)
		}
	} else {
		count := u.ClearEvent(models.EventSub, user.Id)
		redis.IncrEventCount(u.Id, models.EventSub, -count)
	}

	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": Awards{}}, nil)
}

func blacklistHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, getT GetToken) {
	form := getT.(relationshipForm)
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	u := &models.Account{}
	if find, err := u.FindByUserid(form.Userid); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError)
		}
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	redis.SetRelationship(user.Id, form.Userid, models.RelBlacklist, form.Blacklist)

	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": Awards{}}, nil)
}

type getFollowsForm struct {
	Token string `form:"access_token" binding:"required"`
}

func getFollowsHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getFollowsForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	//u := &models.User{Id: user.Id}
	writeResponse(request.RequestURI, resp, redis.Friends(models.RelFollowing, user.Id), nil)
}

func getFollowersHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getFollowsForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	//u := &models.User{Id: user.Id}
	writeResponse(request.RequestURI, resp, redis.Friends(models.RelFollower, user.Id), nil)
}

func getFriendsHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getFollowsForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}
	writeResponse(request.RequestURI, resp, redis.Friends(models.RelFriend, user.Id), nil)
}

func getGroupsHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getFollowsForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	//u := &models.User{Id: user.Id}
	writeResponse(request.RequestURI, resp, redis.Groups(user.Id), nil)
}

type socialListForm struct {
	Token string `form:"access_token" binding:"required"`
	Type  string `form:"member_type" binding:"required"`
	models.Paging
}

func socialListHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form socialListForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	var ids []string
	switch form.Type {
	case "FRIENDS":
		ids = redis.Friends(models.RelFriend, user.Id)
	case "ATTENTION":
		ids = redis.Friends(models.RelFollowing, user.Id)
	case "FANS":
		ids = redis.Friends(models.RelFollower, user.Id)
	case "DEFRIEND":
		ids = redis.Friends(models.RelBlacklist, user.Id)
	case "WEIBO":
		ids = redis.Friends("weibo", user.Id)
	}
	users, err := models.Users(ids, &form.Paging)
	if err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	lb := make([]leaderboardResp, len(users))
	for i, _ := range users {
		lb[i].Userid = users[i].Id
		lb[i].Score = users[i].Props.Score
		lb[i].Level = users[i].Props.Level
		lb[i].Profile = users[i].Profile
		lb[i].Nickname = users[i].Nickname
		lb[i].Gender = users[i].Gender
		lb[i].LastLog = users[i].LastLogin.Unix()
		lb[i].Birth = users[i].Birth
		if users[i].Loc != nil {
			lb[i].Location = *users[i].Loc
		}
	}

	respData := map[string]interface{}{
		"members_list":  lb,
		"page_frist_id": form.Paging.First,
		"page_last_id":  form.Paging.Last,
	}
	writeResponse(request.RequestURI, resp, respData, nil)

	return
}
