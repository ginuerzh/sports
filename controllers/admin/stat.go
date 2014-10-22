// stat
package admin

import (
	"github.com/ginuerzh/sports/models"
	"github.com/jinzhu/now"
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

func BindStatApi(m *martini.ClassicMartini) {
	m.Get("/admin/stat/summary", binding.Form(summaryForm{}), adminErrorHandler, summaryHandler)
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
	Date int64 `form:"date"`
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

	
}
