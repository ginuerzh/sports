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
	m.Get("/1/tasks/get",
		binding.Form(getTaskForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		getTaskHandler)
	m.Get("/1/tasks/getList",
		binding.Form(getTasksForm{}),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		getTasksHandler)
	m.Get("/1/tasks/getInfo",
		binding.Form(getTaskInfoForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		//loadUserHandler,
		getTaskInfoHandler)
	m.Get("/1/tasks/result",
		binding.Form(getTaskResultForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		//loadUserHandler,
		getTaskResultHandler)
	m.Post("/1/tasks/execute",
		binding.Json(completeTaskForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		checkLimitHandler,
		completeTaskHandler)
	m.Get("/1/tasks/referrals",
		binding.Form(taskReferralForm{}),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		taskReferralsHandler)
	m.Post("/1/tasks/referral/pass",
		binding.Json(referralPassForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		taskReferralPassHandler)
	m.Post("/1/tasks/share",
		binding.Json(taskShareForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		taskShareHandler)
	m.Post("/1/tasks/shared",
		binding.Json(taskSharedForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		taskSharedHandler)
}

type getTaskForm struct {
	Next bool `form:"next"`
	parameter
}

func getTaskHandler(r *http.Request, w http.ResponseWriter,
	user *models.Account, p Parameter) {
	form := p.(getTaskForm)

	tid := user.Taskid
	status := user.TaskStatus

	if tid == 0 {
		rec, _ := user.LastTaskRecord2()
		tid = int(rec.Task + 1) // next task
		status = rec.Status
		if status == "" {
			status = models.StatusNormal
		}
		user.UpdateTask(tid, status)
	}

	if status == "" {
		status = models.StatusNormal
	}

	if form.Next {
		if status == models.StatusFinish {
			tid++
		}
		if status == models.StatusFinish || status == models.StatusUnFinish {
			status = models.StatusNormal
		}

		user.UpdateTask(tid, status)
	}

	if tid > len(models.NewTasks) {
		writeResponse(r.RequestURI, w, nil, nil)
	}

	task := models.NewTasks[tid-1]
	task.Status = status

	config := &models.Config{}
	config.Find()
	if task.Index < len(config.Videos) {
		video := config.Videos[task.Index]
		task.Video = video.Url
		/*
			if len(video.Desc) > 0 {
				task.Desc = video.Desc
			}
		*/
	}

	var stat struct {
		Distance int `json:"distance"`
		Run      int `json:"run"`
		Article  int `json:"article"`
		Game     int `json:"game"`
	}

	stat.Article, _ = user.TaskRecordCount("post", models.StatusFinish)
	stat.Game, _ = user.TaskRecordCount("game", models.StatusFinish)
	records, _ := user.TaskRecords("run")
	stat.Run = len(records)
	for i, _ := range records {
		stat.Distance += records[i].Sport.Distance
	}

	respData := map[string]interface{}{
		"task": task,
		"stat": stat,
	}
	writeResponse(r.RequestURI, w, respData, nil)
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
	count, _ := user.TaskRecordCount("", models.StatusFinish)
	week := count / 7

	var target, actual int

	last, _ := user.LastTaskRecord()
	// all weekly tasks are completed
	if week > 0 && count%7 == 0 && last.AuthTime.After(now.BeginningOfWeek()) {
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
			tasks[i].Result = fmt.Sprintf("你在%s游戏中得了%d分",
				record.Game.Name, record.Game.Score)
		}
		if task.Type == models.TaskRunning {
			target += task.Distance
			if tasks[i].Status == models.StatusFinish && record.Sport != nil {
				actual += record.Sport.Distance
			}
		}
	}
	//log.Println(tasks)
	//random := rand.New(rand.NewSource(time.Now().Unix()))
	respData := map[string]interface{}{
		"week_id":              week + 1,
		"task_list":            tasks,
		"task_target_distance": target,
		"task_actual_distance": actual,
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
			task.Result = fmt.Sprintf("你在%s游戏中得了%d分, 得到了%d魔法值和%d贝币",
				record.Game.Name, record.Game.Score, record.Game.Magic, record.Coin/models.Satoshi)
		}
		if task.Type == models.TaskPost {
			task.Result = fmt.Sprintf("你得到了%d文学值和%d贝币", record.Coin/models.Satoshi, record.Coin/models.Satoshi)
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

type getTaskResultForm struct {
	Tid int `form:"task_id" binding:"required"`
	parameter
}

func getTaskResultHandler(request *http.Request, resp http.ResponseWriter,
	user *models.Account, p Parameter) {

	form := p.(getTaskResultForm)
	var task models.Task
	if form.Tid > 0 && form.Tid <= len(models.NewTasks) {
		task = models.NewTasks[form.Tid-1]
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
			task.Result = fmt.Sprintf("你在%s游戏中得了%d分, 得到了%d魔法值和%d贝币",
				record.Game.Name, record.Game.Score, record.Game.Magic, record.Coin/models.Satoshi)
		}
		if task.Type == models.TaskPost {
			task.Result = fmt.Sprintf("你得到了%d文学值和%d贝币", record.Coin/models.Satoshi, record.Coin/models.Satoshi)
		}
	}

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
		record.Coin = 10 * models.Satoshi
	default:
	}

	awards := Awards{}
	if form.Tid == 1 {
		awards = Awards{Wealth: 3 * models.Satoshi}
		GiveAwards(user, awards, redis)
		record.Coin = awards.Wealth

		/*
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
		*/
	}

	err := record.Save()
	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": awards}, err)
}

type taskReferralForm struct {
	parameter
}

type referral struct {
	Userid   string   `json:"userid"`
	Nickname string   `json:"nikename"`
	Profile  string   `json:"profile_image"`
	Gender   string   `json:"sex_type"`
	Images   []string `json:"user_images"`
	Birthday int64    `json:"birthday"`
	Lastlog  int64    `json:"last_login_time"`
	models.Location

	RunRatio     float32 `json:"run_ratio"`
	PostRatio    float32 `json:"post_ratio"`
	PkRatio      float32 `json:"pk_ratio"`
	LastTime     int64   `json:"last_time"`
	LastDistance int     `json:"last_distance"`
	LastId       string  `json:"last_id"`
}

func taskReferralsHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	if user.Taskid > len(models.NewTasks) {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.AccessError))
		return
	}
	task := models.NewTasks[user.Taskid]
	excludes := redis.TaskSharers()
	friends := redis.Friends(models.RelFriend, user.Id)
	excludes = append(excludes, user.Id)
	excludes = append(excludes, friends...)
	users, _ := user.TaskReferrals(task.Type, excludes)
	var referrals []*referral

	for i, _ := range users {
		ref := &referral{
			Userid:   users[i].Id,
			Nickname: users[i].Nickname,
			Profile:  users[i].Profile,
			Gender:   users[i].Gender,
			Images:   users[i].Photos,
			Birthday: users[i].Birth,
			Lastlog:  users[i].LastLogin.Unix(),
			Location: users[i].Loc,
		}
		if users[i].Ratios.RunRecv > 0 {
			ref.RunRatio = float32(users[i].Ratios.RunAccept) / float32(users[i].Ratios.RunRecv)
		}
		if users[i].Ratios.PostRecv > 0 {
			ref.PostRatio = float32(users[i].Ratios.PostAccept) / float32(users[i].Ratios.PostRecv)
		}
		if users[i].Ratios.PKRecv > 0 {
			ref.PkRatio = float32(users[i].Ratios.PKAccept) / float32(users[i].Ratios.PKRecv)
		}
		if task.Type == models.TaskRunning || task.Type == models.TaskGame {
			rec, _ := users[i].LastRecord("run")
			ref.LastId = rec.Id.Hex()
			ref.LastTime = rec.PubTime.Unix()
		} else if task.Type == models.TaskPost {
			article := users[i].LatestArticle()
			ref.LastId = article.Id.Hex()
			ref.LastTime = article.PubTime.Unix()
		}
		referrals = append(referrals, ref)
	}

	respData := map[string]interface{}{"referrals": referrals}
	writeResponse(r.RequestURI, w, respData, nil)
}

type referralPassForm struct {
	Userid string `json:"userid"`
	parameter
}

func taskReferralPassHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {
	form := p.(referralPassForm)
	u := &models.Account{Id: form.Userid}
	err := u.SetBlock(user.Id, true)
	writeResponse(r.RequestURI, w, nil, err)
}

type taskShareForm struct {
	Type   string `json:"type"`
	Userid string `json:"userid" binding:"required"`
	TaskId int    `json:"task_id" binding:"required"`
	models.Location
	Addr  string `json:"addr"`
	Time  int64  `json:"time"`
	Image string `json:"image"`
	Coin  int64  `json:"coin"`
	parameter
}

func taskShareHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {
	form := p.(taskShareForm)

	if form.TaskId > len(models.NewTasks) {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.AccessError))
		return
	}

	//task := models.NewTasks[form.TaskId]
	// ws push
	event := &models.Event{
		Type: models.EventSystem,
		Time: time.Now().Unix(),
		Data: models.EventData{
			//Type: models.EventRunShare,
			Id:   user.Id,
			From: user.Id,
			To:   form.Userid,
		},
	}

	switch form.Type {
	case models.TaskRunning:
		record, _ := user.LastRecord("run")
		latlng := strconv.FormatFloat(form.Lat, 'f', 7, 64) + "," +
			strconv.FormatFloat(form.Lng, 'f', 7, 64)
		event.Data.Type = models.EventRunShare
		event.Data.Body = []models.MsgBody{
			{Type: "record_id", Content: record.Id.Hex()},
			{Type: "latlng", Content: latlng},
			{Type: "locaddr", Content: form.Addr},
			{Type: "time", Content: strconv.FormatInt(form.Time, 10)},
			{Type: "addr_image", Content: form.Image},
		}
	case models.TaskPost:
		event.Data.Type = models.EventPostShare
		article := user.LatestArticle()
		event.Data.Body = []models.MsgBody{
			{Type: "article_id", Content: article.Id.Hex()},
		}
	case models.TaskGame:
		record, _ := user.LastRecord("run")
		event.Data.Type = models.EventPKShare
		event.Data.Body = []models.MsgBody{
			{Type: "record_id", Content: record.Id.Hex()},
		}
	default:
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.AccessError))
		return
	}

	event.Save()
	redis.PubMsg(event.Type, event.Data.To, event.Bytes())

	u := &models.Account{Id: form.Userid}
	u.UpdateRatio(form.Type, false)
	redis.SetTaskShare(form.Userid, true)

	if _, err := consumeCoin(user.Wallet.Addr, form.Coin); err == nil {
		redis.ConsumeCoins(user.Id, form.Coin)
	}

	writeResponse(r.RequestURI, w, nil, nil)
}

type taskSharedForm struct {
	Sender    string `json:"sender" binding:"required"`
	Accept    bool   `json:"accept"`
	Type      string `json:"type" binding:"required"`
	ArticleId string `json:"article_id"`
	Addr      string `json:"addr"`
	Time      int64  `json:"time"`
	Image     string `json:"image"`
	parameter
}

func taskSharedHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {
	form := p.(taskSharedForm)

	u := &models.Account{}
	u.FindByUserid(form.Sender)

	// ws push
	event := &models.Event{
		Type: models.EventSystem,
		Time: time.Now().Unix(),
		Data: models.EventData{
			Id:   user.Id,
			From: user.Id,
			To:   form.Sender,
		},
	}

	var article *models.Article

	switch form.Type {
	case models.TaskRunning:
		event.Data.Type = models.EventRunShared
		article = &models.Article{
			Author:  user.Id,
			PubTime: time.Now(),
			Contents: []models.Segment{
				{ContentType: "TEXT", ContentText: "我和" + u.Nickname + "约好一起跑步，有想一起参加的吗？" +
					"\n跑步地点： " + form.Addr +
					"\n跑步时间： " + time.Unix(form.Time, 0).Format("2006-01-02 3:04 PM")},
				{ContentType: "IMAGE", ContentText: form.Image},
			},
		}
	case models.TaskPost:
		article := &models.Article{}
		if find, _ := article.FindById(form.ArticleId); find {
			article.SetThumb(user.Id, true)
		}
		event.Data.Type = models.EventPostShared
	case models.TaskGame:
		event.Data.Type = models.EventPKShared
		result := u.Nickname + " 主动PK " + user.Nickname + "大获全胜。"
		if u.Props.Score < user.Props.Score {
			result = u.Nickname + " 主动PK " + user.Nickname + "大败亏输。"
		}
		article = &models.Article{
			Author:  user.Id,
			Type:    "pk",
			PubTime: time.Now(),
			Contents: []models.Segment{
				{ContentType: "TEXT", ContentText: result},
				{ContentType: "IMAGE", ContentText: form.Image},
			},
		}
	default:
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.AccessError))
		return
	}

	awards := Awards{}
	if form.Accept {
		redis.SetRelationship(user.Id, []string{form.Sender}, models.RelFriend, true)
		event.Save()
		redis.PubMsg(models.EventSystem, form.Sender, event.Bytes())
		user.UpdateRatio(form.Type, true)
		if article != nil {
			article.Save()
		}
		awards.Wealth = 1 * models.Satoshi
		GiveAwards(user, awards, redis)
	}

	redis.SetTaskShare(user.Id, false)

	writeResponse(r.RequestURI, w, map[string]interface{}{"ExpEffect": awards}, nil)
}
