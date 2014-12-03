// account
package jsgen

import (
	"crypto/md5"
	"fmt"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"io"
	"log"
	"net/http"
	"strings"
)

func BindAccountApi(m *martini.ClassicMartini) {
	m.Post("/api/user/login", binding.Json(jsgenLoginForm{}), adminLoginHandler)
}

type jsgenLoginForm struct {
	Auto     bool   `json:"logauto"`
	Username string `json:"logname"`
	Password string `json:"logpwd"`
	Time     int64  `json:"logtime"`
}

type User struct {
	Id          string   `json:"_id"`
	Name        string   `json:"name"`
	Nickname    string   `json:"nickname"`
	Email       string   `json:"email"`
	Locked      bool     `json:"locked"`
	Social      social   `json:"social"`
	Sex         string   `json:"sex"`
	Role        int      `json:"role"`
	Avatar      string   `json:"avatar"`
	Desc        string   `json:"desc"`
	Date        int64    `json:"date"`
	Score       int      `json:"score"`
	ReadTime    int64    `json:"readtimestamp"`
	LastLogin   int64    `json:"lastLoginDate"`
	Fans        int      `json:"fans"`
	Follow      int      `json:"follow"`
	FollowList  []string `json:"followlist"`
	TagsList    []string `json:"tagsList"`
	Articles    int      `json:"articles"`
	Collections int      `json:"collections"`
	MarkList    []int    `json:"markList"`
	Unread      []string `json:"unread"`
	ReceiveList []string `json:"receiveList"`
	SendList    []string `json:"sendList"`
}

func adminLoginHandler(w http.ResponseWriter, redis *models.RedisLogger, form jsgenLoginForm) {
	user := &models.Account{}

	log.Println(form.Username, form.Password)
	h := md5.New()
	io.WriteString(h, form.Password)
	pwd := fmt.Sprintf("%x", h.Sum(nil))

	find := false
	var err error
	if find, err = user.FindByUserPass(strings.ToLower(form.Username), pwd); !find {
		if err == nil {
			err = AuthError
		} else {
			err = DbError
		}
	}

	if err != nil {
		writeResponse(w, false, nil, nil, err)
		return
	}

	redis.LogLogin(user.Id)

	info := &User{
		Id:          user.Id,
		Name:        user.Nickname,
		Nickname:    user.Nickname,
		Email:       "",
		Locked:      user.TimeLimit != 0,
		Social:      social{},
		Sex:         user.Gender,
		Role:        1,
		Avatar:      user.Profile,
		Desc:        user.About,
		Date:        user.RegTime.Unix() * 1000,
		Score:       user.Score,
		ReadTime:    0,
		LastLogin:   user.LastLogin.Unix() * 1000,
		FollowList:  redis.Friends(models.RelFollowing, user.Id),
		TagsList:    []string{},
		Articles:    user.ArticleCount(),
		Collections: 0,
		MarkList:    []int{},
		Unread:      []string{},
		ReceiveList: []string{},
		SendList:    []string{},
	}
	info.Follow, info.Fans, _, _ = redis.FriendCount(user.Id)

	writeResponse(w, true, info, nil, nil)
}
