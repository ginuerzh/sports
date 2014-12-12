// task
package admin

import (
	"github.com/ginuerzh/sports/controllers"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"labix.org/v2/mgo/bson"
	"github.com/jinzhu/now"
	"net/http"
	//"time"
	"log"
)

func BindTaskApi(m *martini.ClassicMartini) {
	m.Get("/admin/task/list", binding.Form(tasklistForm{}), adminErrorHandler, tasklistHandler)
	m.Get("/admin/task/timeline", binding.Form(taskTimelineForm{}), adminErrorHandler, taskTimelineHandler)
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
	AdminPaging
	Token string `form:"access_token"`
}

type userTask struct {
	Id       string      `json:"userid"`
	Nickname string      `json:"nickname"`
	Profile  string      `json:"profile"`
	Tasks    []*taskinfo `json:"tasks"`
}

func tasklistHandler(w http.ResponseWriter, redis *models.RedisLogger, form tasklistForm) {
	if form.PageCount == 0 {
		form.PageCount = 50
	}
	total, users, _ := models.UserList("", form.PageIndex, form.PageCount)
	log.Println(total, len(users))
	usertasks := make([]*userTask, len(users))
	for i, user := range users {
		usertasks[i] = &userTask{}
		usertasks[i].Id = user.Id
		usertasks[i].Nickname = user.Nickname
		usertasks[i].Profile = user.Profile

		tasklist := user.Tasks
		week := len(tasklist.Completed) / 7
		if week > 0 && len(tasklist.Completed)%7 == 0 &&
			tasklist.Last.After(now.BeginningOfWeek()) {
			week -= 1
		}
		usertasks[i].Tasks = make([]*taskinfo, 7)
		for j, t := range models.Tasks[week*7 : week*7+7] {
			usertasks[i].Tasks[j] = convertTask(&t, &tasklist)
		}
	}

	pages := total / form.PageCount
	if total%form.PageCount > 0 {
		pages++
	}
	resp := map[string]interface{}{
		"users":        usertasks,
		"page_index":   form.PageIndex,
		"page_total":   pages,
		"total_number": total,
	}

	writeResponse(w, resp)
}

type taskTimelineForm struct {
	Userid string `form:"userid" binding:"required"`
	Week   int    `form:"week"`
	Token  string `form:"access_token"`
}

func taskTimelineHandler(w http.ResponseWriter, redis *models.RedisLogger, form taskTimelineForm) {
	/*
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(w, errors.NewError(errors.AccessError))
			return
		}
	*/
	/*
		u := &models.User{Id: form.Userid}
		tl, err := u.GetTasks()
	*/
	u := &models.Account{}
	u.FindByUserid(form.Userid)

	tl := u.Tasks

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

	//u := &models.User{Id: form.Userid}
	user := &models.Account{}
	user.FindByUserid(form.Userid)
	if err := user.SetTaskComplete(form.Id, form.Pass, form.Reason); err != nil {
		writeResponse(w, err)
		return
	}

	if form.Pass {

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
