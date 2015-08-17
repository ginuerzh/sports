// task
package admin

import (
	"github.com/ginuerzh/sports/controllers"
	"github.com/ginuerzh/sports/errors"
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
	m.Options("/admin/task/auth", optionsHandler)
	m.Post("/admin/task/auth_list", binding.Json(taskAuthListForm{}), adminErrorHandler, taskAuthListHandler)
	m.Options("/admin/task/auth_list", optionsHandler)
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
	Mood      string   `json:"mood"`
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
	if info.Id > 0 && info.Id <= int64(len(models.NewTasks)) {
		info.Desc = models.NewTasks[info.Id-1].Desc
	}
	if record.Sport != nil {
		info.Source = record.Sport.Source
		info.Distance = record.Sport.Distance
		info.Duration = record.Sport.Duration
		info.Images = record.Sport.Pics
		info.Reason = record.Sport.Review
		info.Mood = record.Sport.Mood
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
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}

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
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}

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
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}
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

type taskAuth struct {
	Userid string `json:"userid" binding:"required"`
	Id     int64  `json:"task_id" binding:"required"`
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
}

type taskAuthForm struct {
	taskAuth
	Token string `json:"access_token"`
}

func taskAuthOptionsHandler(w http.ResponseWriter) {
	writeResponse(w, nil)
}

func taskAuthFunc(userid string, auth *taskAuth, redis *models.RedisLogger) error {
	user := &models.Account{}
	user.FindByUserid(auth.Userid)

	record := &models.Record{Uid: user.Id, Task: auth.Id}
	record.FindByTask(auth.Id)
	awards := controllers.Awards{}

	parent := &models.Article{}
	parent.FindByRecord(record.Id.Hex())
	if len(parent.Id) > 0 && len(auth.Reason) > 0 && parent.Author != userid {
		review := &models.Article{
			Parent:   parent.Id.Hex(),
			Author:   userid,
			Title:    auth.Reason,
			Type:     models.ArticleCoach,
			Contents: []models.Segment{{ContentType: "TEXT", ContentText: auth.Reason}},
			PubTime:  time.Now(),
		}
		review.Save()
	}

	if auth.Pass {
		level := user.Level()
		awards = controllers.Awards{
			Physical: 3 + level,
			Wealth:   3 * models.Satoshi,
			Score:    3 + level,
		}
		awards.Level = models.Score2Level(user.Props.Score+awards.Score) - level

		controllers.GiveAwards(user, awards, redis)

		if record.Sport != nil {
			redis.UpdateRecLB(user.Id, record.Sport.Distance, int(record.Sport.Duration))
		}

		record.SetStatus(models.StatusFinish, auth.Reason, awards.Wealth)
		if auth.Id < 1000 {
			user.UpdateTask(int(auth.Id), models.StatusFinish)
		}
	} else {
		record.SetStatus(models.StatusUnFinish, auth.Reason, 0)
		if auth.Id < 1000 {
			user.UpdateTask(int(auth.Id), models.StatusUnFinish)
		}
		parent.SetPrivilege(models.PrivPrivate)
	}

	// ws push
	event := &models.Event{
		Type: models.EventNotice,
		Time: time.Now().Unix(),
		Data: models.EventData{
			Type: models.EventTaskDone,
			To:   user.Id,
			Body: []models.MsgBody{
				{Type: "physique_value", Content: strconv.FormatInt(awards.Physical, 10)},
				{Type: "coin_value", Content: strconv.FormatInt(awards.Wealth, 10)},
			},
		},
	}
	if auth.Id < 1000 {
		event.Data.Body = append(event.Data.Body, models.MsgBody{Type: "task_id", Content: strconv.Itoa(int(auth.Id))})
	}

	if !auth.Pass {
		event.Data.Type = models.EventTaskFailure
	}
	redis.PubMsg(event.Type, event.Data.To, event.Bytes())
	event.Save()

	event = &models.Event{
		Type: models.EventArticle,
		Time: time.Now().Unix(),
		Data: models.EventData{
			Type: models.EventCoachPass,
			Id:   parent.Id.Hex(),
			From: userid,
			To:   parent.Author,
			Body: []models.MsgBody{
				{Type: "total_count", Content: strconv.Itoa(parent.CoachReviewCount + 1)},
				{Type: "image", Content: ""},
			},
		},
	}
	if !auth.Pass {
		event.Data.Type = models.EventCoachNPass
	}
	event.Save()
	redis.PubMsg(event.Type, event.Data.To, event.Bytes())

	return nil
}

func taskAuthHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, form taskAuthForm) {
	userid := redis.OnlineUser(form.Token)
	if len(userid) == 0 {
		writeResponse(w, errors.NewError(errors.AccessError))
		return
	}

	a := &taskAuth{
		Userid: form.Userid,
		Id:     form.Id,
		Pass:   form.Pass,
		Reason: form.Reason,
	}
	if err := taskAuthFunc(userid, a, redis); err != nil {
		writeResponse(w, err)
		return
	}
	writeResponse(w, map[string]bool{"pass": form.Pass})
}

type taskAuthListForm struct {
	Auths []taskAuth `json:"auths"`
	Token string     `json:"access_token"`
}

func taskAuthListHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, form taskAuthListForm) {

	userid := redis.OnlineUser(form.Token)
	if len(userid) == 0 {
		writeResponse(w, errors.NewError(errors.AccessError))
		return
	}

	pass := make([]bool, len(form.Auths))
	for i, _ := range form.Auths {
		taskAuthFunc(userid, &form.Auths[i], redis)
		pass[i] = form.Auths[i].Pass
	}

	writeResponse(w, map[string][]bool{"pass": pass})
}
