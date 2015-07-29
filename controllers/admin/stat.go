// stat
package admin

import (
	"github.com/ginuerzh/sports/models"
	//"github.com/jinzhu/now"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"net/http"
	"time"
)

func BindStatApi(m *martini.ClassicMartini) {
	m.Get("/admin/stat/summary", binding.Form(summaryForm{}), adminErrorHandler, summaryHandler)
	m.Get("/admin/stat/retention", binding.Form(retentionForm{}), adminErrorHandler, retentionHandler)
}

type summaryForm struct {
	Days  int    `form:"days"`
	Token string `form:"access_token"`
}

func summaryHandler(w http.ResponseWriter, redis *models.RedisLogger, form summaryForm) {
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}

	var stats struct {
		RegPhone      []int64 `json:"reg_phone"`
		RegEmail      []int64 `json:"reg_email"`
		RegWeibo      []int64 `json:"reg_weibo"`
		Registers     []int64 `json:"registers"`
		Logins        []int64 `json:"logins"`
		CoachLogins   []int64 `json:"coach_logins"`
		Actives       []int64 `json:"actives"`
		PostUsers     []int64 `json:"post_users"`
		Posts         []int64 `json:"posts"`
		Gamers        []int64 `json:"gamers"`
		GameTime      []int64 `json:"game_time"`
		RecordUsers   []int64 `json:"record_users"`
		AuthCoaches   []int64 `json:"auth_coaches"`
		Coins         []int64 `json:"coins"`
		Users         int     `json:"users"`
		Onlines       int     `json:"onlines"`
		OnlineCoaches int     `json:"online_coaches"`
	}
	days := form.Days
	if days <= 0 {
		days = 3
	}
	//var start, end time.Time
	//start, end = now.BeginningOfDay(), now.EndOfDay()

	stats.RegPhone = redis.RegisterCount(days, models.AccountPhone)
	stats.RegEmail = redis.RegisterCount(days, models.AccountEmail)
	stats.RegWeibo = redis.RegisterCount(days, models.AccountWeibo)
	stats.Registers = redis.RegisterCount(days, "")
	stats.Logins = redis.LoginCount(days)
	stats.CoachLogins = redis.CoachLoginCount(days)
	actives := make([]int64, days)
	for i := 0; i < days; i++ {
		actives[i] = stats.Logins[i] - stats.Registers[i]
	}
	stats.Actives = actives
	stats.PostUsers = redis.PostUserCount(days)
	stats.Posts = redis.PostsCount(days)
	stats.Gamers = redis.GamersCount(days)
	stats.GameTime = redis.GameTime(days)
	stats.RecordUsers = redis.RecordUsersCount(days)
	stats.AuthCoaches = redis.AuthCoachesCount(days)
	stats.Coins = redis.CoinsCount(days)
	stats.Users = models.UserCount()
	stats.Onlines = redis.Onlines()

	writeResponse(w, stats)
}

type retentionForm struct {
	Date  int64  `form:"date"`
	Token string `form:"access_token"`
}

func retentionHandler(w http.ResponseWriter, redis *models.RedisLogger, form retentionForm) {
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}

	date := time.Now()
	if form.Date > 0 {
		date = time.Unix(form.Date, 0)
	}
	r := redis.Retention(date)
	writeResponse(w, map[string]interface{}{"retention": r})
}
