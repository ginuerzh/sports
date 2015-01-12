// article
package controllers

import (
	"bytes"
	//"encoding/json"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func BindArticleApi(m *martini.ClassicMartini) {
	m.Post("/1/article/new",
		binding.Json(newArticleForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		checkLimitHandler,
		newArticleHandler)
	m.Post("/1/article/delete",
		binding.Json(deleteArticleForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		deleteArticleHandler)
	m.Post("/1/article/thumb",
		binding.Json(articleThumbForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		checkLimitHandler,
		articleThumbHandler)
	m.Get("/1/article/is_thumbed",
		binding.Form(articleIsThumbedForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		articleIsThumbedHandler)
	m.Get("/1/article/timelines",
		binding.Form(articleListForm{}),
		ErrorHandler,
		articleListHandler)
	m.Get("/1/article/get",
		binding.Form(articleInfoForm{}),
		ErrorHandler,
		articleInfoHandler)
	m.Post("/1/article/comments",
		binding.Json(articleCommentsForm{}),
		ErrorHandler,
		articleCommentsHandler)
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
	Content    string           `json:"content"`
	Rewards    int64            `json:"reward_total"`
}

var (
	header = `<!DOCTYPE HTML>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=310, initial-scale=1, maximum-scale=1, user-scalable=no">
		<style>
			body{
				font-size:16px;
				line-height:30px;
				background-color:#f6f6f6;
				text-align:center;
				margin: 0;
			}
			p{
				text-align:left;
				padding-left: 5px;
				padding-right: 5px;
			}
			img{
				max-width:97%;
				height:auto;
				margin:auto;
			}
			div.divimg {
				text-align:center;
			}
		</style>
	</head>
	<body>`

	footer = `
	</body>
</html>`
)

func convertArticle(article *models.Article) *articleJsonStruct {
	jsonStruct := &articleJsonStruct{}
	jsonStruct.Id = article.Id.Hex()
	jsonStruct.Parent = article.Parent
	jsonStruct.Author = article.Author
	//jsonStruct.Contents = article.Contents
	jsonStruct.PubTime = article.PubTime.Unix()
	jsonStruct.Thumbs = len(article.Thumbs)
	jsonStruct.Reviews = len(article.Reviews)
	jsonStruct.Rewards = article.TotalReward

	jsonStruct.Title = article.Title
	jsonStruct.Image = article.Image

	if len(article.Content) > 0 {
		jsonStruct.Content = header + article.Content + footer
	}
	if len(article.Contents) > 0 {
		jsonStruct.Content = header + content2Html(article.Contents) + footer
	}

	return jsonStruct
}

func content2Html(contents []models.Segment) string {
	buf := &bytes.Buffer{}
	for _, content := range contents {
		switch strings.ToUpper(content.ContentType) {
		case "TEXT":
			buf.WriteString("<p>" + content.ContentText + "</p>")
		case "IMAGE":
			buf.WriteString("<div class=\"divimg\"><img src=\"" + content.ContentText + "\" /></div>")
		}
	}

	return buf.String()
}

type newArticleForm struct {
	Parent   string           `json:"parent_article_id"`
	Contents []models.Segment `json:"article_segments" binding:"required"`
	Tags     []string         `json:"article_tag"`
	parameter
}

func articleCover(contents []models.Segment) (text string, images []string) {
	for _, seg := range contents {
		if len(text) == 0 && strings.ToUpper(seg.ContentType) == "TEXT" {
			text = seg.ContentText
		}
		if strings.ToUpper(seg.ContentType) == "IMAGE" {
			images = append(images, seg.ContentText)
		}
	}
	return
}
func newArticleHandler(request *http.Request, resp http.ResponseWriter,
	client *ApnClient, redis *models.RedisLogger, user *models.Account, p Parameter) {
	form := p.(newArticleForm)

	article := &models.Article{
		Author:  user.Id,
		Content: content2Html(form.Contents),
		PubTime: time.Now(),
		Parent:  form.Parent,
		Tags:    form.Tags,
	}
	article.Title, article.Images = articleCover(form.Contents)
	if len(article.Images) > 0 {
		article.Image = article.Images[0]
	}

	if len(article.Tags) == 0 {
		article.Tags = []string{"SPORT_LOG"}
	}

	if err := article.Save(); err != nil {
		log.Println(err)
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	awards := Awards{}
	// only new article
	if len(form.Parent) == 0 {
		awards = Awards{Literal: 10 + user.Props.Level, Wealth: 10 * models.Satoshi, Score: 10 + user.Props.Level}
		if err := GiveAwards(user, awards, redis); err != nil {
			log.Println(err)
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.DbError, err.Error()))
			return
		}
	}

	// comment
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

		//u := &models.User{Id: parent.Author}
		author := &models.Account{}
		author.FindByUserid(parent.Author)

		// ws push
		event := &models.Event{
			Type: models.EventArticle,
			Time: time.Now().Unix(),
			Data: models.EventData{
				Type: models.EventComment,
				Id:   parent.Id.Hex(),
				From: user.Id,
				To:   parent.Author,
				Body: []models.MsgBody{
					{Type: "total_count", Content: strconv.Itoa(parent.CommentCount())},
					{Type: "image", Content: parent.Image},
				},
			},
		}
		redis.PubMsg(models.EventArticle, parent.Author, event.Bytes())
		if err := event.Save(); err == nil {
			redis.IncrEventCount(parent.Author, event.Data.Type, 1)
		}
		// apple push
		if author.Push {
			for _, dev := range author.Devs {
				if err := client.Send(dev, user.Nickname+"评论了你的主题!", 1, ""); err != nil {
					log.Println(err)
				}
			}
		}
	}

	respData := map[string]interface{}{
		"articles_without_content": convertArticle(article),
		"ExpEffect":                awards,
	}
	writeResponse(request.RequestURI, resp, respData, nil)
}

type deleteArticleForm struct {
	Id string `json:"article_id" binding:"required"`
	parameter
}

func deleteArticleHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(deleteArticleForm)

	article := &models.Article{}
	article.Author = user.Id
	article.Id = bson.ObjectIdHex(form.Id)

	err := article.Remove()
	writeResponse(request.RequestURI, resp, nil, err)
}

type articleThumbForm struct {
	Id     string `json:"article_id" binding:"required"`
	Status bool   `json:"thumb_status"`
	parameter
}

func articleThumbHandler(request *http.Request, resp http.ResponseWriter,
	client *ApnClient, redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(articleThumbForm)
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

	awards := Awards{}
	/*
		awards := Awards{Physical: 1, Wealth: 1 * models.Satoshi}
		if err := giveAwards(user, awards); err != nil {
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.DbError, err.Error()))
			return
		}
	*/

	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": awards}, nil)

	if form.Status {
		author := &models.Account{}
		author.FindByUserid(article.Author)

		// ws push
		event := &models.Event{
			Type: models.EventArticle,
			Time: time.Now().Unix(),
			Data: models.EventData{
				Type: models.EventThumb,
				Id:   article.Id.Hex(),
				From: user.Id,
				To:   article.Author,
				Body: []models.MsgBody{
					{Type: "total_count", Content: strconv.Itoa(len(article.Thumbs) + 1)},
					{Type: "image", Content: article.Image},
				},
			},
		}

		redis.PubMsg(models.EventArticle, article.Author, event.Bytes())
		if err := event.Save(); err == nil {
			redis.IncrEventCount(article.Author, event.Data.Type, 1)
		}

		// apple push
		if author.Push {
			for _, dev := range author.Devs {
				if err := client.Send(dev, user.Nickname+"赞了你的主题!", 1, ""); err != nil {
					log.Println(err)
				}
			}
		}
	}
	//user.UpdateAction(ActThumb, nowDate())
}

type articleIsThumbedForm struct {
	Id string `form:"article_id" binding:"required"`
	parameter
}

func articleIsThumbedHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(articleIsThumbedForm)

	article := &models.Article{}
	article.Id = bson.ObjectIdHex(form.Id)
	b, err := article.IsThumbed(user.Id)

	respData := map[string]bool{"is_thumbed": b}
	writeResponse(request.RequestURI, resp, respData, err)
}

type articleListForm struct {
	Token string `form:"access_token"`
	models.Paging
	Tag string `form:"article_tag"`
}

func articleListHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form articleListForm) {
	_, articles, err := models.GetArticles(form.Tag, &form.Paging, true)
	if err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	jsonStructs := make([]*articleJsonStruct, len(articles))
	for i, _ := range articles {
		jsonStructs[i] = convertArticle(&articles[i])
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

func articleInfoHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form articleInfoForm) {
	article := &models.Article{}
	if find, err := article.FindById(form.Id); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError)
		}
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	if uid := redis.OnlineUser(form.Token); uid == article.Author {
		user := &models.Account{Id: uid}
		count := user.ClearEvent(models.EventThumb, article.Id.Hex())
		redis.IncrEventCount(user.Id, models.EventThumb, -count)

		count = user.ClearEvent(models.EventComment, article.Id.Hex())
		redis.IncrEventCount(user.Id, models.EventComment, -count)

		count = user.ClearEvent(models.EventReward, article.Id.Hex())
		redis.IncrEventCount(user.Id, models.EventReward, -count)
	}

	jsonStruct := convertArticle(article)
	writeResponse(request.RequestURI, resp, jsonStruct, nil)
}

type articleCommentsForm struct {
	Id string `json:"article_id"  binding:"required"`
	models.Paging
}

func articleCommentsHandler(request *http.Request, resp http.ResponseWriter,
	form articleCommentsForm) {

	article := &models.Article{Id: bson.ObjectIdHex(form.Id)}
	_, comments, err := article.Comments(&form.Paging, true)

	jsonStructs := make([]*articleJsonStruct, len(comments))
	for i, _ := range comments {
		jsonStructs[i] = convertArticle(&comments[i])
	}

	respData := make(map[string]interface{})
	respData["page_frist_id"] = form.Paging.First
	respData["page_last_id"] = form.Paging.Last
	//respData["page_item_count"] = total
	respData["articles_without_content"] = jsonStructs
	writeResponse(request.RequestURI, resp, respData, err)
}
