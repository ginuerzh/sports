package controllers

import (
	//"encoding/json"
	"fmt"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/jinzhu/now"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"log"
	//"math/rand"
	"net/http"
	"strconv"
	"time"
)

func init() {
	now.FirstDayMonday = true
}

func BindTaskApi(m *martini.ClassicMartini) {
	m.Get("/1/tasks/getList",
		binding.Form(getTasksForm{}),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		getTasksHandler)
	m.Get("/1/tasks/getInfo",
		binding.Form(getTaskInfoForm{}),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		getTaskInfoHandler)
	m.Post("/1/tasks/execute",
		binding.Json(completeTaskForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		checkLimitHandler,
		completeTaskHandler)
}

type getTasksForm struct {
	parameter
}

func getTasksHandler(r *http.Request, w http.ResponseWriter, user *models.Account) {

	/*
		tasklist := user.Tasks

		week := len(tasklist.Completed) / 7
		if week > 0 && len(tasklist.Completed)%7 == 0 &&
			tasklist.Last.After(now.BeginningOfWeek()) {
			week -= 1
		}
		list := make([]models.Task, 7)
		if week*7+7 <= len(models.Tasks)+1 {
			copy(list, models.Tasks[week*7:week*7+7])
		}
		for i, _ := range list {
			list[i].Status = tasklist.TaskStatus(list[i].Id)
			if list[i].Type == models.TaskGame && list[i].Status == "FINISH" {
				rec := &models.Record{Uid: user.Id}
				rec.FindByTask(list[i].Id)
				if rec.Game != nil {
					list[i].Desc = fmt.Sprintf("你在%s游戏中得了%d分",
						rec.Game.Name, rec.Game.Score)
				}
			}
		}
	*/
	count, _ := user.TaskRecordCount(models.StatusFinish)
	week := count / 7

	last, _ := user.LastTaskRecord()
	// all weekly tasks are completed
	if week > 0 && count%7 == 0 && last.PubTime.After(now.BeginningOfWeek()) {
		week -= 1
	}
	//log.Println("week", week)
	tasks := make([]models.Task, 7)
	if week*7+7 <= len(models.Tasks) {
		copy(tasks, models.Tasks[week*7:week*7+7])
	}

	for i, task := range tasks {
		tasks[i].Status = models.StatusNormal
		record := &models.Record{Uid: user.Id}

		if find, _ := record.FindByTask(task.Id); find {
			tasks[i].Status = record.Status
		}
		if task.Type == models.TaskGame && task.Status == models.StatusFinish &&
			record.Game != nil {
			tasks[i].Desc = fmt.Sprintf("你在%s游戏中得了%d分",
				record.Game.Name, record.Game.Score)
		}
	}
	//log.Println(tasks)
	//random := rand.New(rand.NewSource(time.Now().Unix()))
	respData := map[string]interface{}{
		"week_id":   week + 1,
		"task_list": tasks,
		//"week_desc": tips[random.Int()%len(tips)],
	}

	writeResponse(r.RequestURI, w, respData, nil)
}

type getTaskInfoForm struct {
	Tid int `form:"task_id" binding:"required"`
	parameter
}

func getTaskInfoHandler(request *http.Request, resp http.ResponseWriter,
	user *models.Account, p Parameter) {

	form := p.(getTaskInfoForm)
	//tasklist := user.Tasks
	var task models.Task
	if form.Tid > 0 && form.Tid <= len(models.Tasks) {
		task = models.Tasks[form.Tid-1]
	}
	task.Status = models.StatusNormal

	record := &models.Record{Uid: user.Id}
	if find, _ := record.FindByTask(task.Id); find {
		task.Status = record.Status
		task.BeginTime = record.StartTime.Unix()
		task.EndTime = record.EndTime.Unix()
		if record.Sport != nil {
			task.Source = record.Sport.Source
			task.Distance = record.Sport.Distance
			task.Duration = record.Sport.Duration
			task.Pics = record.Sport.Pics
			task.Result = record.Sport.Review
		}
		if record.Game != nil && task.Status == models.StatusFinish {
			task.Desc = fmt.Sprintf("你在%s游戏中得了%d分, 得到了%d魔法值和%d贝币",
				record.Game.Name, record.Game.Score, record.Game.Magic, record.Game.Coin/models.Satoshi)
		}
	}
	/*
		task.Status = tasklist.TaskStatus(task.Id)
		proof := tasklist.GetProof(task.Id)
		task.Pics = proof.Pics
		task.Result = proof.Result
		if task.Type == models.TaskGame && task.Status == "FINISH" {
			rec := &models.Record{Uid: user.Id}
			rec.FindByTask(task.Id)
			if rec.Game != nil {
				task.Desc = fmt.Sprintf("你在%s游戏中得了%d分, 得到了%d魔法值和%d贝币",
					rec.Game.Name, rec.Game.Score, rec.Game.Magic, rec.Game.Coin/models.Satoshi)
			}
		}
	*/

	writeResponse(request.RequestURI, resp, map[string]interface{}{"task_info": task}, nil)
}

type completeTaskForm struct {
	Tid    int      `json:"task_id" binding:"required"`
	Proofs []string `json:"task_pics"`
	parameter
}

func completeTaskHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(completeTaskForm)

	if form.Tid < 1 || form.Tid > len(models.Tasks) {
		writeResponse(request.RequestURI, resp, nil, nil)
		return
	}

	record := &models.Record{
		Uid:     user.Id,
		Task:    int64(form.Tid),
		PubTime: time.Now(),
	}

	task := models.Tasks[form.Tid-1]
	switch task.Type {
	case models.TaskPost:
		record.Type = "post"
		record.Status = models.StatusFinish
	default:
	}

	awards := Awards{}
	if form.Tid == 1 {
		awards = Awards{Score: 30, Wealth: 30 * models.Satoshi}
		if err := GiveAwards(user, awards, redis); err != nil {
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.DbError))
			return
		}
		// ws push
		event := &models.Event{
			Type: models.EventStatus,
			Time: time.Now().Unix(),
			Data: models.EventData{
				Type: models.EventTask,
				To:   user.Id,
				Body: []models.MsgBody{
					{Type: "coin_value", Content: strconv.FormatInt(awards.Wealth, 10)},
				},
			},
		}
		redis.PubMsg(event.Type, event.Data.To, event.Bytes())
	}

	err := record.Save()
	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": awards}, err)
}
