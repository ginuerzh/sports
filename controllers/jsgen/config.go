// config
package jsgen

import (
	"gopkg.in/go-martini/martini.v1"
	"net/http"
)

type infoAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Url   string `json:"url"`
}

type dependencies struct {
	Marked     string `json:"marked"`
	Mongoskin  string `json:"mongoskin"`
	Nodemailer string `json:"nodemailer"`
	Redis      string `json:"redis"`
	Rrestjs    string `json:"rrestjs"`
	Thenjs     string `json:"thenjs"`
	Xss        string `json:"xss"`
}

type info struct {
	Author       infoAuthor   `json:"author"`
	CdnHost      string       `json:"cdnHost"`
	Dependencies dependencies `json:"dependencies"`
	Desc         string       `json:"description"`
	HomePage     string       `json:"homepage"`
	Keywords     []string     `json:"keywords"`
	Main         string       `json:main"`
	Nodejs       string       `json:"nodejs"`
	Rrestjs      string       `json:"rrestjs"`
	Version      string       `json:"version"`
}

type config struct {
	ArticleTagsMax int
	ContentMaxLen  int
	ContentMinLen  int
	SummaryMaxLen  int
	TitleMaxLen    int
	TitleMinLen    int
	UserNameMaxLen int
	UserNameMinLen int
	UserTagsMax    int
	Articles       int         `json:"articles"`
	Beian          string      `json:"beian"`
	CloudDomian    string      `json:"cloudDomian"`
	Comments       int         `json:"comments"`
	Date           int64       `json:"date"`
	Desc           string      `json:"description"`
	Domain         string      `json:"domain"`
	Keywords       string      `json:"keywords"`
	Logo           string      `json:"logo"`
	MaxOnlineNum   int         `json:"maxOnlineNum"`
	maxOnlineTime  int64       `json:"maxOnlineTime"`
	MetaDesc       string      `json:"metadesc"`
	MetaTitle      string      `json:"metatitle"`
	OnlineNum      int         `json:"onlineNum"`
	OnlineUsers    int         `json:"onlineUsers"`
	Register       bool        `json:"register"`
	Taglist        []tag       `json:"tagsList"`
	Title          string      `json:"title"`
	Upload         bool        `json:"upload"`
	Url            string      `json:"url"`
	User           interface{} `json:"user"`
	Users          int         `json:"users"`
	Visitors       int         `json:"visitors"`
	Info           info        `json:"info"`
}

func BindConfigApi(m *martini.ClassicMartini) {
	m.Get("/api/index", configHandler)
}

func configHandler(w http.ResponseWriter) {
	writeResponse(w, true, &config{}, nil, nil)
}
