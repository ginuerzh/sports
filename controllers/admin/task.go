// task
package admin

import (
	"github.com/ginuerzh/sports/controllers"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"labix.org/v2/mgo/bson"
	"net/http"
	//"time"
)

func BindTaskApi(m *martini.ClassicMartini) {
	m.Get("/admin/task/list", binding.Form(tasklistForm{}), adminErrorHandler, tasklistHandler)
	m.Post("/admin/task/auth", binding.Json(taskAuthForm{}), adminErrorHandler, taskAuthHandler)
}

type taskinfo struct {
	Id     int      `json:"task_id"`
	Type   string   `json:"type"`
	Desc   string   `json:"desc"`
	Images []string `json:"images"`
	Status string   `json:"status"`
	Reason string   `json:"reason"`
}

func convertTask(task *models.Task, tl *models.TaskList) *taskinfo {
	t := &taskinfo{
		Id:     task.Id,
		Type:   task.Type,
		Desc:   task.Desc,
		Status: tl.TaskStatus(task.Id),
	}
	proof := tl.GetProof(task.Id)
	t.Images = proof.Pics
	t.Reason = proof.Result

	return t
}

type tasklistForm struct {
	Userid string `form:"userid" binding:"required"`
	Week   int    `form:"week"`
	Token  string `form:"access_token"`
}

func tasklistHandler(w http.ResponseWriter, redis *models.RedisLogger, form tasklistForm) {
	/*
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(w, errors.NewError(errors.AccessError))
			return
		}
	*/

	u := &models.User{Id: form.Userid}
	tl, err := u.GetTasks()
	if err != nil {
		writeResponse(w, err)
		return
	}

	var tasks []*taskinfo
	if form.Week > 0 {
		week := form.Week - 1
		if (week*7 + 7) <= len(models.Tasks) {
			tasks = make([]*taskinfo, 7)
			for i, t := range models.Tasks[week*7 : week*7+7] {
				tasks[i] = convertTask(&t, &tl)
			}
		}
	} else {
		tasks = make([]*taskinfo, len(models.Tasks))
		for i, t := range models.Tasks {
			tasks[i] = convertTask(&t, &tl)
		}
	}

	writeResponse(w, map[string]interface{}{"tasks": tasks})
}

type taskAuthForm struct {
	Userid string `json:"userid" binding:"required"`
	Id     int    `json:"task_id" binding:"required"`
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
	Token  string `json:"access_token"`
}

func taskAuthHandler(w http.ResponseWriter, redis *models.RedisLogger, form taskAuthForm) {
	/*
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(w, errors.NewError(errors.AccessError))
			return
		}
	*/

	u := &models.User{Id: form.Userid}
	if err := u.SetTaskComplete(form.Id, form.Pass, form.Reason); err != nil {
		writeResponse(w, err)
		return
	}

	if form.Pass {
		user := &models.Account{}
		user.FindByUserid(form.Userid)
		awards := controllers.Awards{
			Physical: 30 + user.Props.Level,
			Wealth:   30 * models.Satoshi,
			Score:    30 + user.Props.Level,
		}
		if err := controllers.GiveAwards(user, awards, redis); err != nil {
			writeResponse(w, err)
			return
		}
	}

	writeResponse(w, map[string]bool{"pass": form.Pass})
}
