// article
package admin

import (
	"bytes"
	"fmt"
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
	m.Get("/admin/article/list", binding.Form(articleListForm{}), adminErrorHandler, articleListHandler)
	m.Get("/admin/article/timeline", binding.Form(articleListForm{}), adminErrorHandler, articleTimelineHandler)
	m.Get("/admin/article/comments", binding.Form(articleListForm{}), adminErrorHandler, articleCommentsHandler)
	m.Options("/admin/article/post", articlePostOptionsHandler)
	m.Post("/admin/article/post", binding.Json(postForm{}), adminErrorHandler, articlePostHandler)
	m.Post("/admin/article/delete", binding.Json(delArticleForm{}), adminErrorHandler, delArticleHandler)
	m.Get("/admin/article/search", binding.Form(articleSearchForm{}), adminErrorHandler, articleSearchHandler)
	m.Post("/admin/article/update", binding.Json(articleUpdateForm{}), adminErrorHandler, articleUpdateHandler)
}

type articleInfo struct {
	Id     string              `json:"article_id"`
	Parent string              `json:"parent"`
	Author *userInfoJsonStruct `json:"author"`
	Image  string              `json:"cover_image"`
	Title  string              `json:"cover_text"`
	Time   int64               `json:"time"`
	//TimeStr      string         `json:"time_str"`
	Thumbs       int            `json:"thumbs_count"`
	CommentCount int            `json:"comments_count"`
	Comments     []*articleInfo `json:"comments"`
	Rewards      int64          `json:"rewards_value"`
	RewardUsers  []string       `json:"rewards_users"`
	Tags         []string       `json:"tags"`
	Contents     string         `json:"contents"`
}

func convertArticle(article *models.Article, redis *models.RedisLogger) *articleInfo {
	info := &articleInfo{}
	info.Id = article.Id.Hex()
	info.Parent = article.Parent

	user := &models.Account{}
	user.FindByUserid(article.Author)
	info.Author = convertUser(user, redis)
	info.Time = article.PubTime.Unix()
	//info.TimeStr = article.PubTime.Format("2006-01-02 15:04:05")
	info.Thumbs = len(article.Thumbs)
	info.CommentCount = len(article.Reviews)
	info.Rewards = article.TotalReward
	info.RewardUsers = article.Rewards
	info.Tags = article.Tags
	info.Title = article.Title
	info.Image = article.Image
	info.Contents = article.Content
	if len(article.Contents) > 0 {
		info.Contents = formatArticleContent(article.Contents)
	}
	return info
}

func formatArticleContent(contents []models.Segment) string {
	buffer := &bytes.Buffer{}
	images := &bytes.Buffer{}
	j := 1
	for _, seg := range contents {
		switch strings.ToUpper(seg.ContentType) {
		case "TEXT":
			buffer.WriteString(seg.ContentText + "\n\n")
		case "IMAGE":
			fmt.Fprintf(buffer, "![pic%d][%d]\n\n", j, j)
			fmt.Fprintf(buffer, "[%d]: %s\n", j, seg.ContentText)
			j++
		}
	}
	if images.Len() > 0 {
		buffer.WriteString("\n\n")
		buffer.WriteString(images.String())
	}
	return buffer.String()
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

	a := convertArticle(article, redis)
	a.Comments = comments(redis, article, 0, 0)
	writeResponse(w, a)
}

func comments(redis *models.RedisLogger, article *models.Article, pageIndex, pageCount int) (a []*articleInfo) {
	_, list, _ := article.AdminComments(pageIndex, pageCount)
	for _, c := range list {
		art := convertArticle(&c, redis)
		if len(c.Reviews) > 0 {
			art.Comments = comments(redis, &c, 0, 0)
		}
		a = append(a, art)
	}
	return a
}

type articleListForm struct {
	Userid string `form:"userid"`
	Id     string `form:"article_id"`
	Sort   string `form:"sort"`
	AdminPaging
	Token string `form:"access_token" binding:"required"`
}

func articleListHandler(w http.ResponseWriter, redis *models.RedisLogger, form articleListForm) {
	if form.PageCount == 0 {
		form.PageCount = 50
	}
	total, articles, _ := models.ArticleList(form.Sort, form.PageIndex, form.PageCount)
	list := make([]*articleInfo, len(articles))
	for i, _ := range articles {
		list[i] = convertArticle(&articles[i], redis)
	}
	pages := total / form.PageCount
	if total%form.PageCount > 0 {
		pages++
	}
	resp := map[string]interface{}{
		"articles":     list,
		"page_index":   form.PageIndex,
		"page_total":   pages,
		"total_number": total,
	}

	writeResponse(w, resp)
}

func hotArticleHandler(w http.ResponseWriter, redis *models.RedisLogger, form articleListForm) {

}

func articleTimelineHandler(w http.ResponseWriter, redis *models.RedisLogger, form articleListForm) {
	/*
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(w, errors.NewError(errors.AccessError))
			return
		}
	*/
	if form.PageCount == 0 {
		form.PageCount = 50
	}
	u := &models.Account{Id: form.Userid}
	total, articles, _ := u.ArticleTimeline(form.PageIndex, form.PageCount)

	list := make([]*articleInfo, len(articles))
	for i, _ := range articles {
		list[i] = convertArticle(&articles[i], redis)
	}

	pages := total / form.PageCount
	if total%form.PageCount > 0 {
		pages++
	}
	resp := map[string]interface{}{
		"articles":     list,
		"page_index":   form.PageIndex,
		"page_total":   pages,
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
	if form.PageCount == 0 {
		form.PageCount = 50
	}
	article := &models.Article{Id: bson.ObjectIdHex(form.Id)}
	list := comments(redis, article, form.PageIndex, form.PageCount)

	total := article.CommentCount()
	pages := total / form.PageCount
	if total%form.PageCount > 0 {
		pages++
	}
	resp := map[string]interface{}{
		"comments":     list,
		"page_index":   form.PageIndex,
		"page_total":   pages,
		"total_number": total,
	}
	writeResponse(w, resp)
}

type postForm struct {
	Id       string      `json:"article_id"`
	Author   string      `json:"author"`
	Contents string      `json:"contents"`
	Title    string      `json:"title"`
	Image    []string    `json:"image"`
	Tags     interface{} `json:"tags"`
	Token    string      `json:"access_token" binding:"required"`
}

func articlePostOptionsHandler(w http.ResponseWriter) {
	writeResponse(w, nil)
}

func articlePostHandler(w http.ResponseWriter, redis *models.RedisLogger, form postForm) {
	uid := redis.OnlineUser(form.Token)
	if len(uid) == 0 {
		writeResponse(w, errors.NewError(errors.AccessError))
		return
	}
	user := &models.Account{Id: uid}

	article := &models.Article{
		Parent:  form.Id,
		Author:  form.Author,
		PubTime: time.Now(),
		//Tags:    []string{form.Tags},
	}

	switch v := form.Tags.(type) {
	case string:
		article.Tags = []string{v}
	case []string:
		article.Tags = v
	}

	if len(article.Tags) == 0 {
		article.Tags = []string{"SPORT_LOG"}
	}

	if len(form.Author) == 0 {
		article.Author = user.Id
	}

	article.Content = form.Contents
	article.Title = form.Title
	article.Images = form.Image
	if len(article.Images) > 0 {
		article.Image = article.Images[0]
	}
	/*
		article.Contents = append(article.Contents,
			models.Segment{ContentType: "TEXT", ContentText: form.Contents})
	*/
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
	Keyword string `form:"keyword"`
	Tag     string `form:"tag"`
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
	if form.PageCount == 0 {
		form.PageCount = 50
	}
	total, articles, _ := models.AdminSearchArticle(form.Keyword, form.Tag, form.PageIndex, form.PageCount)

	list := make([]*articleInfo, len(articles))
	for i, _ := range articles {
		list[i] = convertArticle(&articles[i], redis)
	}

	pages := total / form.PageCount
	if total%form.PageCount > 0 {
		pages++
	}
	resp := map[string]interface{}{
		"articles":     list,
		"page_index":   form.PageIndex,
		"page_total":   pages,
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
