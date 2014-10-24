// article
package admin

import (
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
	"time"
)

func BindArticleApi(m *martini.ClassicMartini) {
	m.Get("/admin/article/info", binding.Form(articleInfoForm{}), adminErrorHandler, articleInfoHandler)
	m.Get("/admin/article/timeline", binding.Form(articleListForm{}), adminErrorHandler, articleListHandler)
	m.Get("/admin/article/comments", binding.Form(articleListForm{}), adminErrorHandler, articleCommentsHandler)
	m.Post("/admin/article/post", binding.Json(postForm{}), adminErrorHandler, articlePostHandler)
	m.Post("/admin/article/delete", binding.Json(delArticleForm{}), adminErrorHandler, delArticleHandler)
	m.Get("/admin/article/search", binding.Form(articleSearchForm{}), adminErrorHandler, articleSearchHandler)
	m.Post("/admin/article/update", binding.Json(articleUpdateForm{}), adminErrorHandler, articleUpdateHandler)
}

type articleInfo struct {
	Id          string           `json:"article_id"`
	Author      string           `json:"author"`
	Image       string           `json:"cover_image"`
	Title       string           `json:"cover_text"`
	Time        int64            `json:"time"`
	Thumbs      int              `json:"thumbs_count"`
	Comments    int              `json:"comments_count"`
	Rewards     int64            `json:"rewards_value"`
	RewardUsers []string         `json:"rewards_users"`
	Tags        []string         `json:"tags"`
	Contents    []models.MsgBody `json:"contents"`
}

func convertArticle(article *models.Article) *articleInfo {
	info := &articleInfo{}
	info.Id = article.Id.Hex()
	info.Author = article.Author
	info.Time = article.PubTime.Unix()
	info.Thumbs = len(article.Thumbs)
	info.Comments = len(article.Reviews)
	info.Rewards = article.TotalReward
	info.RewardUsers = article.Rewards
	info.Tags = article.Tags
	info.Title, info.Image = article.Cover()
	info.Contents = []models.MsgBody{}

	for _, content := range article.Contents {
		info.Contents = append(info.Contents,
			models.MsgBody{Type: strings.ToLower(content.ContentType), Content: content.ContentText})
	}
	return info
}

type articleInfoForm struct {
	Id    string `form:"article_id"`
	Token string `form:"access_token" binding:"required"`
}

func articleInfoHandler(w http.ResponseWriter, redis *models.RedisLogger, form articleInfoForm) {
	/*
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(w, errors.NewError(errors.AccessError))
			return
		}
	*/

	article := &models.Article{}
	if find, err := article.FindById(form.Id); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError)
		}
		writeResponse(w, err)
		return
	}

	writeResponse(w, convertArticle(article))
}

type articleListForm struct {
	Userid string `form:"userid"`
	Id     string `form:"article_id"`
	AdminPaging
	Token string `form:"access_token" binding:"required"`
}

func articleListHandler(w http.ResponseWriter, redis *models.RedisLogger, form articleListForm) {
	/*
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(w, errors.NewError(errors.AccessError))
			return
		}
	*/

	u := &models.User{Id: form.Userid}
	paging := &models.Paging{First: form.Pre, Last: form.Next, Count: form.Count}
	total, articles, _ := u.Articles("ARTICLES", paging)

	list := make([]*articleInfo, len(articles))
	for i, _ := range articles {
		list[i] = convertArticle(&articles[i])
	}

	resp := map[string]interface{}{
		"articles":     list,
		"prev_cursor":  paging.First,
		"next_cursor":  paging.Last,
		"total_number": total,
	}

	writeResponse(w, resp)
}

func articleCommentsHandler(w http.ResponseWriter, redis *models.RedisLogger, form articleListForm) {
	/*
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(w, errors.NewError(errors.AccessError))
			return
		}
	*/

	paging := &models.Paging{First: form.Pre, Last: form.Next, Count: form.Count}
	article := &models.Article{Id: bson.ObjectIdHex(form.Id)}
	total, comments, _ := article.Comments(paging)

	list := make([]*articleInfo, len(comments))
	for i, _ := range comments {
		list[i] = convertArticle(&comments[i])
	}

	resp := map[string]interface{}{
		"comments":     list,
		"next_cursor":  paging.Last,
		"prev_cursor":  paging.First,
		"total_number": total,
	}
	writeResponse(w, resp)
}

type postForm struct {
	Id       string           `json:"article_id"`
	Author   string           `json:"author"`
	Contents []models.MsgBody `json:"contents"`
	Tags     []string         `json:"tags"`
	Token    string           `json:"access_token" binding:"required"`
}

func articlePostHandler(w http.ResponseWriter, redis *models.RedisLogger, form postForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(w, errors.NewError(errors.AccessError))
		return
	}

	article := &models.Article{
		Parent:  form.Id,
		Author:  form.Author,
		PubTime: time.Now(),
		Tags:    form.Tags,
	}

	if len(form.Author) == 0 {
		article.Author = user.Id
	}
	for _, content := range form.Contents {
		article.Contents = append(article.Contents,
			models.Segment{ContentType: strings.ToUpper(content.Type), ContentText: content.Type})
	}

	if err := article.Save(); err != nil {
		writeResponse(w, err)
		return
	}

	writeResponse(w, map[string]string{"article_id": article.Id.Hex()})
}

type delArticleForm struct {
	Id    string `json:"article_id"`
	Token string `json:"access_token" binding:"required"`
}

func delArticleHandler(w http.ResponseWriter, redis *models.RedisLogger, form delArticleForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(w, errors.NewError(errors.AccessError))
		return
	}

	article := &models.Article{}

	if bson.IsObjectIdHex(form.Id) {
		article.Id = bson.ObjectIdHex(form.Id)
	}

	if err := article.RemoveId(); err != nil {
		writeResponse(w, err)
		return
	}

	writeResponse(w, map[string]interface{}{})
}

type articleSearchForm struct {
	Keyword string `form:"keyword" binding:"required"`
	AdminPaging
	Token string `form:"access_token"`
}

func articleSearchHandler(w http.ResponseWriter, redis *models.RedisLogger, form articleSearchForm) {
	/*
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(w, errors.NewError(errors.AccessError))
			return
		}
	*/

	paging := &models.Paging{First: form.Pre, Last: form.Next, Count: form.Count}
	total, articles, err := models.SearchArticle(form.Keyword, paging)

	if err != nil {
		writeResponse(w, err)
	}

	list := make([]*articleInfo, len(articles))
	for i, _ := range articles {
		list[i] = convertArticle(&articles[i])
	}

	resp := map[string]interface{}{
		"articles":     list,
		"next_cursor":  paging.Last,
		"prev_cursor":  paging.First,
		"total_number": total,
	}
	writeResponse(w, resp)
}

type articleUpdateForm struct {
	Id       string           `json:"article_id" binding:"required"`
	Author   string           `json:"author,omitempty"`
	Time     int64            `json:"time,omitempty"`
	Tags     []string         `json:"tags,omitempty"`
	Contents []models.MsgBody `json:"contents,omitempty"`
	Token    string           `json:"access_token"`
}

func articleUpdateHandler(w http.ResponseWriter, redis *models.RedisLogger, form articleUpdateForm) {
	/*
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(w, errors.NewError(errors.AccessError))
			return
		}
	*/

	var contents []models.Segment
	for _, content := range form.Contents {
		contents = append(contents, models.Segment{
			ContentType: content.Type,
			ContentText: content.Content,
		})
	}

	article := &models.Article{
		Id:       bson.ObjectIdHex(form.Id),
		Author:   form.Author,
		Contents: contents,
		PubTime:  time.Unix(form.Time, 0),
		Tags:     form.Tags,
	}

	if err := article.Update(); err != nil {
		writeResponse(w, err)
		return
	}

	writeResponse(w, map[string]interface{}{})
}
