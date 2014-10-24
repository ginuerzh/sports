// stat
package admin

import (
	"github.com/ginuerzh/sports/models"
	"github.com/jinzhu/now"
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
	Token string `form:"access_token"`
}

func summaryHandler(w http.ResponseWriter, redis *models.RedisLogger, form summaryForm) {
	/*
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(w, errors.NewError(errors.AccessError))
			return
		}
	*/

	var start, end time.Time
	start, end = now.BeginningOfDay(), now.EndOfDay()

	registers := redis.RegisterCount(3)
	logins := redis.LoginCount(3)
	summary := map[string]interface{}{
		"registers": registers,
		"logins":    logins,
		"actives":   []int64{logins[0] - registers[0], logins[1] - registers[1], logins[2] - registers[2]},
		"onlines":   redis.Onlines(),
		"users":     models.UserCount(),
		"posts": []int{
			models.PostCount(start, end),
			models.PostCount(start.AddDate(0, 0, 1), end.AddDate(0, 0, 1)),
			models.PostCount(start.AddDate(0, 0, 2), end.AddDate(0, 0, 2)),
		},
	}

	writeResponse(w, summary)
}

type retentionForm struct {
	Date  int64  `form:"date"`
	Token string `form:"access_token"`
}

func retentionHandler(w http.ResponseWriter, redis *models.RedisLogger, form retentionForm) {
	/*
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(w, errors.NewError(errors.AccessError))
			return
		}
	*/
	date := time.Now()
	if form.Date > 0 {
		date = time.Unix(form.Date, 0)
	}
	r := redis.Retention(date)
	writeResponse(w, map[string]interface{}{"retention": r})
}
