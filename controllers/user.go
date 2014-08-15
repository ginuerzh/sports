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
	//"time"
)

func BindUserApi(m *martini.ClassicMartini) {
	m.Get("/1/user/articles", binding.Form(userArticlesForm{}), ErrorHandler, userArticlesHandler)
	m.Post("/1/user/send_device_token", binding.Json(sendDevForm{}), ErrorHandler, sendDevHandler)
	m.Post("/1/user/set_push_enable", binding.Json(setPushForm{}), ErrorHandler, setPushHandler)
	m.Get("/1/user/is_push_enabled", binding.Form(pushStatusForm{}), ErrorHandler, pushStatusHandler)
	m.Post("/1/user/enableAttention", binding.Json(followForm{}), ErrorHandler, followHandler)
	m.Get("/1/user/getAttentionFriendsList", binding.Form(getFollowsForm{}), ErrorHandler, getFollowsHandler)
	m.Get("/1/user/getAttentedMembersList", binding.Form(getFollowsForm{}), ErrorHandler, getFollowersHandler)
	m.Get("/1/user/getJoinedGroupsList", binding.Form(getFollowsForm{}), ErrorHandler, getGroupsHandler)
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
		jsonStructs[i] = convertArticle(&articles[i], 0)
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

func sendDevHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, form sendDevForm) {

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

func setPushHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, form setPushForm) {

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

type followForm struct {
	Userid string `json:"userid"`
	Follow bool   `json:"beFriend"`
	Token  string `json:"access_token" binding:"required"`
}

func followHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form followForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}
	//u := &models.User{Id: user.Id}
	//u.SetFollow(form.Userid, !form.Follow)
	redis.SetFollow(user.Id, form.Userid, form.Follow)

	writeResponse(request.RequestURI, resp, nil, nil)
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
	writeResponse(request.RequestURI, resp, redis.Follows(user.Id), nil)
}

func getFollowersHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getFollowsForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	//u := &models.User{Id: user.Id}
	writeResponse(request.RequestURI, resp, redis.Followers(user.Id), nil)
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
