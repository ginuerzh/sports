// article
package jsgen

import (
	"bytes"
	"fmt"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"io"
	//"log"
	"net/http"
	"strconv"
	//"strings"
)

func BindArticleApi(m *martini.ClassicMartini) {
	m.Get("/api/article/latest", binding.Form(pagination{}), latestArticlesHandler)
	m.Get("/api/article/hots", binding.Form(pagination{}), latestArticlesHandler)
	m.Get("/api/article/comment", binding.Form(pagination{}), latestArticlesHandler)
	m.Get("/api/article/update", binding.Form(pagination{}), latestArticlesHandler)
	m.Get("/api/article/A:id", articleInfoHandler)
}

type author struct {
	Id     string `json:"_id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Score  string `json:"score"`
}

type refer struct {
	Id  interface{} `json:"_id"`
	Url string      `json:"url"`
}

type tag struct {
	Id       string `json:"_id"`
	Tag      string `json:"tag"`
	Articles int    `json:"articles"`
	Users    int    `json:"users"`
}

type Article struct {
	Id          string     `json:"_id"`
	Author      author     `json:"author"`
	Comment     bool       `json:"comment"`
	Comments    int        `json:"comments"`
	Content     string     `json:"content"`
	Cover       string     `json:"cover"`
	PubTime     int64      `json:"date"`
	UpdateTime  int64      `json:"updateTime"`
	Display     int        `json:"display"`
	Hots        int        `json:"hots"`
	Refer       refer      `json:"refer"`
	Status      int        `json:"status"`
	TagList     []tag      `json:"tagsList"`
	Title       string     `json:"title"`
	Visitors    int        `json:"visitors"`
	CommentList []*Article `json:"commentsList"`
}

func formatContent(contents []models.Segment) string {
	buffer := &bytes.Buffer{}
	images := &bytes.Buffer{}
	j := 1
	for _, seg := range contents {
		switch seg.ContentType {
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

func convertArticle(article *models.Article) *Article {
	u := &models.Account{Id: article.Author}
	u.FindByUserid(article.Author)
	coverText, _ := article.Cover()

	a := &Article{
		Id: "A" + article.Id.Hex(),
		Author: author{
			Id:     u.Id,
			Name:   u.Nickname,
			Avatar: u.Profile,
			Score:  strconv.Itoa(u.Score),
		},
		Comment:    true,
		Comments:   len(article.Reviews),
		Content:    formatContent(article.Contents),
		Cover:      "",
		PubTime:    article.PubTime.Unix() * 1000,
		UpdateTime: article.PubTime.Unix() * 1000,
		Title:      coverText,
	}
	for _, t := range article.Tags {
		a.TagList = append(a.TagList, tag{
			Id:  t,
			Tag: t,
		})
	}
	return a
}

func latestArticlesHandler(w http.ResponseWriter, p pagination) {
	if p.PageIndex == 0 {
		p.PageIndex = 1
	}
	total, articles, _ := models.ArticleList("", p.PageIndex-1, p.PageSize)

	var list []*Article
	for _, article := range articles {
		list = append(list, convertArticle(&article))
	}

	p.Total = total

	writeResponse(w, true, list, &p, nil)
}

func articleInfoHandler(w http.ResponseWriter, params martini.Params) {
	article := &models.Article{}
	article.FindById(params["id"])

	a := convertArticle(article)
	a.CommentList = comments(article)

	p := &pagination{Total: 1, PageSize: 20, PageIndex: 1}
	writeResponse(w, true, a, p, nil)
}

func comments(article *models.Article) (a []*Article) {
	_, list, _ := article.AdminComments(0, 0)
	for _, c := range list {
		art := convertArticle(&c)
		if len(c.Reviews) > 0 {
			art.CommentList = comments(&c)
		}
		a = append(a, art)
	}
	return a
}
