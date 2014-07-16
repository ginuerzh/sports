// article
package controllers

import (
	//"bytes"
	//"encoding/json"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"github.com/zhengying/apns"
	"gopkg.in/go-martini/martini.v1"
	//"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	//"strconv"
	//"strings"
	"time"
)

func BindArticleApi(m *martini.ClassicMartini) {
	m.Post("/1/article/new", binding.Json(newArticleForm{}), ErrorHandler, newArticleHandler)
	m.Post("/1/article/delete", binding.Json(deleteArticleForm{}), ErrorHandler, deleteArticleHandler)
	m.Post("/1/article/thumb", binding.Json(articleThumbForm{}), ErrorHandler, articleThumbHandler)
	m.Get("/1/article/is_thumbed", binding.Form(articleIsThumbedForm{}), ErrorHandler, articleIsThumbedHandler)
	m.Get("/1/article/timelines", binding.Form(articleListForm{}), ErrorHandler, articleListHandler)
	m.Get("/1/article/get", binding.Form(articleInfoForm{}), ErrorHandler, articleInfoHandler)
	m.Post("/1/article/comments", binding.Json(articleCommentsForm{}), ErrorHandler, articleCommentsHandler)
}

type articleJsonStruct struct {
	Id         string           `json:"article_id"`
	Parent     string           `json:"parent_article_id"`
	Author     string           `json:"author"`
	Title      string           `json:"cover_text"`
	Image      string           `json:"cover_image"`
	PubTime    int64            `json:"time"`
	Thumbs     int              `json:"thumb_count"`
	NewThumbs  int              `json:"new_thumb_count"`
	Reviews    int              `json:"sub_article_count"`
	NewReviews int              `json:"new_sub_article_count"`
	Contents   []models.Segment `json:"article_segments"`
	Score      int              `json:"exp_effect,omitempty"`
}

func convertArticle(article *models.Article, score int) *articleJsonStruct {
	jsonStruct := &articleJsonStruct{}
	jsonStruct.Id = article.Id.Hex()
	jsonStruct.Parent = article.Parent
	jsonStruct.Author = article.Author
	jsonStruct.Contents = article.Contents
	jsonStruct.PubTime = article.PubTime.Unix()
	jsonStruct.Thumbs = len(article.Thumbs)
	jsonStruct.Reviews = len(article.Reviews)
	jsonStruct.Score = score
	for _, seg := range jsonStruct.Contents {
		if seg.ContentType == "IMAGE" {
			jsonStruct.Image = seg.ContentText
		}
		if seg.ContentType == "TEXT" {
			jsonStruct.Title = seg.ContentText
		}
		if len(jsonStruct.Image) > 0 && len(jsonStruct.Title) > 0 {
			break
		}
	}

	if len(jsonStruct.Contents) == 0 {
		jsonStruct.Contents = []models.Segment{}
	}

	return jsonStruct
}

type newArticleForm struct {
	Parent   string           `json:"parent_article_id"`
	Contents []models.Segment `json:"article_segments" binding:"required"`
	Token    string           `json:"access_token" binding:"required"`
}

func newArticleHandler(request *http.Request, resp http.ResponseWriter,
	client *apns.Client, redis *models.RedisLogger, form newArticleForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	article := &models.Article{
		Author:   user.Id,
		Contents: form.Contents,
		PubTime:  time.Now(),
		Parent:   form.Parent,
	}

	if err := article.Save(); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	score := 0
	if len(form.Parent) == 0 {
		score = actionExps[ActPost]
	} else {
		score = actionExps[ActComment]
	}
	redis.AddScore(user.Id, score)
	writeResponse(request.RequestURI, resp, convertArticle(article, score), nil)

	if len(form.Parent) > 0 {
		parent := &models.Article{}
		if find, err := parent.FindById(form.Parent); !find {
			e := errors.NewError(errors.NotExistsError)
			if err != nil {
				e = errors.NewError(errors.DbError, err.Error())
			}
			writeResponse(request.RequestURI, resp, nil, e)
			return
		}

		u := &models.User{Id: parent.Author}
		event := &models.Event{
			Id:      form.Parent,
			Reviews: []string{article.Id.Hex()},
		}
		if err := u.AddEvent(event); err != nil {
			log.Println(err)
		}

		devs, enabled, _ := u.Devices()
		if enabled {
			for _, dev := range devs {
				if err := sendApns(client, dev, user.Nickname+"评论了你的主题!", 1, ""); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

type deleteArticleForm struct {
	Id    string `json:"article_id" binding:"required"`
	Token string `json:"access_token" binding:"required"`
}

func deleteArticleHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form deleteArticleForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	article := &models.Article{}
	article.Author = user.Id
	article.Id = bson.ObjectIdHex(form.Id)

	err := article.Remove()
	writeResponse(request.RequestURI, resp, nil, err)
}

type articleThumbForm struct {
	Id     string `json:"article_id" binding:"required"`
	Status bool   `json:"thumb_status"`
	Token  string `json:"access_token" binding:"required"`
}

func articleThumbHandler(request *http.Request, resp http.ResponseWriter,
	client *apns.Client, redis *models.RedisLogger, form articleThumbForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	article := &models.Article{}
	if find, err := article.FindById(form.Id); !find {
		e := errors.NewError(errors.NotExistsError)
		if err != nil {
			e = errors.NewError(errors.DbError, err.Error())
		}
		writeResponse(request.RequestURI, resp, nil, e)
		return
	}

	if err := article.SetThumb(user.Id, form.Status); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	writeResponse(request.RequestURI, resp, nil, nil)

	if form.Status {
		u := &models.User{Id: article.Author}
		event := &models.Event{
			Id:     form.Id,
			Thumbs: []string{user.Id},
		}
		//log.Println("add thumb event", article.Author, event)
		if err := u.AddEvent(event); err != nil {
			log.Println(err)
		}

		devs, enabled, _ := u.Devices()
		if enabled {
			for _, dev := range devs {
				if err := sendApns(client, dev, user.Nickname+"赞了你的主题!", 1, ""); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

type articleIsThumbedForm struct {
	Id    string `form:"article_id" binding:"required"`
	Token string `form:"access_token" binding:"required"`
}

func articleIsThumbedHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form articleIsThumbedForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	article := &models.Article{}
	article.Id = bson.ObjectIdHex(form.Id)
	b, err := article.IsThumbed(user.Id)

	respData := map[string]bool{"is_thumbed": b}
	writeResponse(request.RequestURI, resp, respData, err)
}

type articleListForm struct {
	Token string `form:"access_token"`
	models.Paging
}

func articleListHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form articleListForm) {
	_, articles, err := models.GetArticles(&form.Paging)
	if err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	jsonStructs := make([]*articleJsonStruct, len(articles))
	for i, _ := range articles {
		jsonStructs[i] = convertArticle(&articles[i], 0)
	}

	respData := make(map[string]interface{})
	respData["page_frist_id"] = form.Paging.First
	respData["page_last_id"] = form.Paging.Last
	//respData["page_item_count"] = total
	respData["articles_without_content"] = jsonStructs
	writeResponse(request.RequestURI, resp, respData, err)
}

type articleInfoForm struct {
	Id    string `form:"article_id" binding:"required"`
	Token string `form:"access_token"`
}

func articleInfoHandler(request *http.Request, resp http.ResponseWriter, form articleInfoForm) {
	article := &models.Article{}
	if find, err := article.FindById(form.Id); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError)
		}
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	jsonStruct := convertArticle(article, 0)
	writeResponse(request.RequestURI, resp, jsonStruct, nil)
}

type articleCommentsForm struct {
	Id string `json:"article_id"  binding:"required"`
	models.Paging
}

func articleCommentsHandler(request *http.Request, resp http.ResponseWriter, form articleCommentsForm) {
	article := &models.Article{Id: bson.ObjectIdHex(form.Id)}
	_, comments, err := article.Comments(&form.Paging)

	jsonStructs := make([]*articleJsonStruct, len(comments))
	for i, _ := range comments {
		jsonStructs[i] = convertArticle(&comments[i], 0)
	}

	respData := make(map[string]interface{})
	respData["page_frist_id"] = form.Paging.First
	respData["page_last_id"] = form.Paging.Last
	//respData["page_item_count"] = total
	respData["articles_without_content"] = jsonStructs
	writeResponse(request.RequestURI, resp, respData, err)
}

/*
type articleThumbForm struct {
	ArticleId   string `form:"article_id" json:"article_id" binding:"required"`
	Status      bool   `form:"thumb_status" json:"thumb_status"`
	AccessToken string `form:"access_token" json:"access_token" binding:"required"`
}

func (form *articleThumbForm) Validate(e *binding.Errors, req *http.Request) {
}

func articleSetThumbHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form articleThumbForm) {
	user := redis.OnlineUser(form.AccessToken)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.AccessError)
		return
	}

	//var article models.Article
	//article.Id = bson.ObjectIdHex(form.ArticleId)
	//err := article.SetThumb(userid, form.Status)

	if form.Status {
		user.RateArticle(form.ArticleId, models.ThumbRate, false)
	} else {
		user.RateArticle(form.ArticleId, models.ThumbRateMask, true)
	}

	redis.LogArticleThumb(user.Userid, form.ArticleId, form.Status)

	writeResponse(request.RequestURI, resp, nil, errors.NoError)
}

func checkArticleThumbHandler(request *http.Request, resp http.ResponseWriter, form articleThumbForm, redis *models.RedisLogger) {
	user := redis.OnlineUser(form.AccessToken)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.AccessError)
		return
	}

	respData := map[string]bool{"is_thumbed": redis.ArticleThumbed(user.Userid, form.ArticleId)}
	writeResponse(request.RequestURI, resp, respData, errors.NoError)
}

type relatedArticleForm struct {
	ArticleId   string `form:"article_id" json:"article_id"`
	AccessToken string `form:"access_token" json:"access_token" binding:"required"`
}

func relatedArticleHandler(request *http.Request, resp http.ResponseWriter, form relatedArticleForm, redis *RedisLogger) {
	userid := redis.OnlineUser(form.AccessToken)
	if len(userid) == 0 {
		writeResponse(request.RequestURI, resp, nil, errors.AccessError)
		return
	}
	articleIds := redis.RelatedArticles(form.ArticleId, 3)

	articles, err := models.GetArticles(articleIds...)

	jsonStructs := make([]articleJsonStruct, len(articles))

	for i, _ := range articles {
		jsonStructs[i].Id = articles[i].Id.Hex()
		jsonStructs[i].Title = articles[i].Title
		jsonStructs[i].Source = articles[i].Source
		jsonStructs[i].Url = articles[i].Url
		jsonStructs[i].PubTime = articles[i].PubTime.Format(TimeFormat)
		jsonStructs[i].Image = imageUrl(articles[i].Image, ImageThumbnail)
		//jsonStructs[i].Thumbs = redis.ArticleThumbCount(articles[i].Id.Hex())
		//jsonStructs[i].Reviews = redis.ArticleReviewCount(articles[i].Id.Hex())
		//jsonStructs[i].Read = reads[i]
	}

	respData := make(map[string]interface{})
	respData["related_articles"] = jsonStructs
	writeResponse(request.RequestURI, resp, respData, err)
}

func relatedArticleHandler(request *http.Request, resp http.ResponseWriter, form relatedArticleForm, redis *models.RedisLogger) {
	user := redis.OnlineUser(form.AccessToken)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.AccessError)
		return
	}
	mRate := make(map[string]int)

	if userRate, err := user.ArticleRate(); err == errors.NoError {
		for _, rate := range userRate.Rates {
			mRate[rate.Article] = rate.Rate
		}
	}
	//log.Println(mRate)
	data, err := json.Marshal(&mRate)
	if err != nil {
		log.Println(err)
	}
	r, err := http.Post(SlopeOneUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		writeResponse(request.RequestURI, resp, nil, errors.DbError)
		return
	}
	defer r.Body.Close()

	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		writeResponse(request.RequestURI, resp, nil, errors.DbError)
		return
	}

	var ids []string
	if err := json.Unmarshal(data, &ids); err != nil {
		log.Println(err)
	}
	//log.Println(ids)
	if len(ids) > 3 {
		ids = ids[:3]
	}
	articles, e := models.GetArticles(ids...)

	jsonStructs := make([]articleJsonStruct, len(articles))
	for i, _ := range articles {
		jsonStructs[i].Id = articles[i].Id.Hex()
		jsonStructs[i].Title = articles[i].Title
		//jsonStructs[i].Source = articles[i].Source
		jsonStructs[i].Url = articles[i].Url
		//jsonStructs[i].PubTime = articles[i].PubTime.Format(TimeFormat)
		//jsonStructs[i].Image = imageUrl(articles[i].Image, ImageThumbnail)
		//jsonStructs[i].Thumbs = redis.ArticleThumbCount(articles[i].Id.Hex())
		//jsonStructs[i].Reviews = redis.ArticleReviewCount(articles[i].Id.Hex())
		//jsonStructs[i].Read = reads[i]
	}

	respData := make(map[string]interface{})
	respData["related_articles"] = jsonStructs
	writeResponse(request.RequestURI, resp, respData, e)
}

func articleViewersHandler(request *http.Request, resp http.ResponseWriter, params martini.Params, redis *models.RedisLogger) {
	aid := params["id"]

	viewers := redis.ArticleViewers(aid)

	users, err := models.FindUsers(viewers)

	jsonStructs := make([]userJsonStruct, len(users))
	for i, _ := range users {
		//view, thumb, review, _ := users[i].ArticleCount()

		jsonStructs[i].Userid = users[i].Userid
		jsonStructs[i].Nickname = users[i].Nickname
		jsonStructs[i].Type = users[i].Role
		jsonStructs[i].Profile = users[i].Profile
		jsonStructs[i].Phone = users[i].Phone
		jsonStructs[i].Location = users[i].Location
		jsonStructs[i].About = users[i].About
		jsonStructs[i].RegTime = users[i].RegTime.Format(TimeFormat)
		//jsonStructs[i].Views = view
		//jsonStructs[i].Thumbs = thumb
		//jsonStructs[i].Reviews = review
		//jsonStructs[i].Online = redis.IsOnline(users[i].Userid)
	}

	writeResponse(request.RequestURI, resp, jsonStructs, err)
}

func articleThumbsHandler(request *http.Request, resp http.ResponseWriter, params martini.Params, redis *models.RedisLogger) {
	aid := params["id"]

	viewers := redis.ArticleThumbers(aid)

	users, err := models.FindUsers(viewers)

	jsonStructs := make([]userJsonStruct, len(users))
	for i, _ := range users {
		//view, thumb, review, _ := users[i].ArticleCount()

		jsonStructs[i].Userid = users[i].Userid
		jsonStructs[i].Nickname = users[i].Nickname
		jsonStructs[i].Type = users[i].Role
		jsonStructs[i].Profile = users[i].Profile
		jsonStructs[i].Phone = users[i].Phone
		jsonStructs[i].Location = users[i].Location
		jsonStructs[i].About = users[i].About
		jsonStructs[i].RegTime = users[i].RegTime.Format(TimeFormat)
		//jsonStructs[i].Views = view
		//jsonStructs[i].Thumbs = thumb
		//jsonStructs[i].Reviews = review
		jsonStructs[i].Online = redis.IsOnline(users[i].Userid)
	}

	writeResponse(request.RequestURI, resp, jsonStructs, err)
}

func articleWeightHandler(request *http.Request, resp http.ResponseWriter, params martini.Params) {
	id := params["id"]
	log.Println(id, params["weight"])
	if !bson.IsObjectIdHex(id) {
		writeResponse(request.RequestURI, resp, nil, errors.JsonError)
		return
	}

	weight, _ := strconv.Atoi(params["weight"])

	article := models.Article{Id: bson.ObjectIdHex(id)}
	err := article.SetWeight(weight)

	writeResponse(request.RequestURI, resp, nil, err)
}
*/
