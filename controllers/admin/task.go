// task
package admin

import (
	"github.com/ginuerzh/sports/controllers"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"labix.org/v2/mgo/bson"
	//"github.com/jinzhu/now"
	"net/http"
	"time"
	//"log"
	"strconv"
)

func BindTaskApi(m *martini.ClassicMartini) {
	m.Get("/admin/task/list", binding.Form(tasklistForm{}), adminErrorHandler, tasklistHandler)
	m.Get("/admin/task/timeline", binding.Form(userTaskForm{}), adminErrorHandler, userTaskHandler)
	//m.Get("/admin/task/timeline", binding.Form(taskTimelineForm{}), adminErrorHandler, taskTimelineHandler)
	m.Post("/admin/task/auth", binding.Json(taskAuthForm{}), adminErrorHandler, taskAuthHandler)
	m.Options("/admin/task/auth", taskAuthOptionsHandler)
}

type taskinfo struct {
	Id        int64    `json:"task_id"`
	Type      string   `json:"type"`
	Desc      string   `json:"desc"`
	BeginTime int64    `json:"begin_time"`
	EndTime   int64    `json:"end_time"`
	Duration  int64    `json:"duration"`
	Distance  int      `json:"distance"`
	Source    string   `json:"source"`
	Images    []string `json:"images"`
	Status    string   `json:"status"`
	Reason    string   `json:"reason"`
}

/*
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
*/

func convertTask(record *models.Record) *taskinfo {
	info := &taskinfo{
		Id:        record.Task,
		Status:    record.Status,
		BeginTime: record.StartTime.Unix(),
		EndTime:   record.EndTime.Unix(),
	}
	if record.Type == "run" {
		if record.Task < 1000 {
			info.Type = models.TaskRunning
		} else {
			info.Type = models.TaskNormal
		}
	} else if record.Type == "game" {
		info.Type = models.TaskGame
	}
	if len(info.Status) == 0 {
		info.Status = models.StatusFinish
	}
	if info.Id > 0 && info.Id <= int64(len(models.Tasks)) {
		info.Desc = models.Tasks[info.Id-1].Desc
	}
	if record.Sport != nil {
		info.Source = record.Sport.Source
		info.Distance = record.Sport.Distance
		info.Duration = record.Sport.Duration
		info.Images = record.Sport.Pics
		info.Reason = record.Sport.Review
	}

	return info
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

	users := make(map[string]*models.Account)
	var tasks []*userTask
	total, records, _ := models.TaskRecords(form.PageIndex, form.PageCount)
	for _, record := range records {

		info := convertTask(&record)

		user := users[record.Uid]
		if user == nil {
			user = &models.Account{}
			user.FindByUserid(record.Uid)
			users[record.Uid] = user
		}
		task := &userTask{
			Id:       record.Uid,
			Nickname: user.Nickname,
			Profile:  user.Profile,
			Tasks:    []*taskinfo{info},
		}
		tasks = append(tasks, task)
	}
	/*
		total, users, _ := models.UserList("-task", form.PageIndex, form.PageCount)
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

			for _, t := range models.Tasks[week*7 : week*7+7] {
				if t.Type == models.TaskRunning {
					usertasks[i].Tasks = append(usertasks[i].Tasks, convertTask(&t, &tasklist))
				}
			}
		}
	*/

	pages := total / form.PageCount
	if total%form.PageCount > 0 {
		pages++
	}
	resp := map[string]interface{}{
		"users":        tasks,
		"page_index":   form.PageIndex,
		"page_total":   pages,
		"total_number": total,
	}

	writeResponse(w, resp)
}

type userTaskForm struct {
	Nickname string `form:"nickname"`
	Finish   bool   `form:"finish"`
	AdminPaging
	Token string `form:"access_token"`
}

func userTaskHandler(w http.ResponseWriter, redis *models.RedisLogger, form userTaskForm) {
	if form.PageCount == 0 {
		form.PageCount = 50
	}

	user := &models.Account{}
	user.FindByNickname(form.Nickname)

	var tasks []*userTask
	total, records, _ := models.SearchTaskByUserid(user.Id, form.Finish, form.PageIndex, form.PageCount)
	for _, record := range records {

		info := convertTask(&record)

		task := &userTask{
			Id:       record.Uid,
			Nickname: user.Nickname,
			Profile:  user.Profile,
			Tasks:    []*taskinfo{info},
		}
		tasks = append(tasks, task)
	}

	pages := total / form.PageCount
	if total%form.PageCount > 0 {
		pages++
	}
	resp := map[string]interface{}{
		"users":        tasks,
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
	/*
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
	*/
	writeResponse(w, map[string]interface{}{"tasks": nil})
}

type taskAuthForm struct {
	Userid string `json:"userid" binding:"required"`
	Id     int64  `json:"task_id" binding:"required"`
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
	Token  string `json:"access_token"`
}

func taskAuthOptionsHandler(w http.ResponseWriter) {
	writeResponse(w, nil)
}

func taskAuthHandler(r *http.Request, w http.ResponseWriter, redis *models.RedisLogger, form taskAuthForm) {
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
	/*
		if err := user.SetTaskComplete(form.Id, form.Pass, form.Reason); err != nil {
			writeResponse(w, err)
			return
		}
	*/
	record := &models.Record{Uid: user.Id, Task: form.Id}
	record.FindByTask(form.Id)
	awards := controllers.Awards{}

	if form.Pass {

		level := user.Level()
		awards = controllers.Awards{
			Physical: 30 + level,
			Wealth:   30 * models.Satoshi,
			Score:    30 + level,
		}
		awards.Level = models.Score2Level(user.Props.Score+awards.Score) - level

		if err := controllers.GiveAwards(user, awards, redis); err != nil {
			writeResponse(w, err)
			return
		}
		if record.Sport != nil {
			redis.UpdateRecLB(user.Id, record.Sport.Distance, int(record.Sport.Duration))
		}

		record.SetStatus(models.StatusFinish, form.Reason, awards.Wealth)
	} else {
		record.SetStatus(models.StatusUnFinish, form.Reason, 0)
	}

	// ws push
	event := &models.Event{
		Type: models.EventStatus,
		Time: time.Now().Unix(),
		Data: models.EventData{
			Type: models.EventTask,
			To:   user.Id,
			Body: []models.MsgBody{
				{Type: "physique_value", Content: strconv.FormatInt(awards.Physical, 10)},
				{Type: "coin_value", Content: strconv.FormatInt(awards.Wealth, 10)},
			},
		},
	}
	redis.PubMsg(event.Type, event.Data.To, event.Bytes())

	writeResponse(w, map[string]bool{"pass": form.Pass})
}
