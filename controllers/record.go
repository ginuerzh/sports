// record
package controllers

import (
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func BindRecordApi(m *martini.ClassicMartini) {
	m.Post("/1/record/new",
		binding.Json(newRecordForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		checkLimitHandler,
		newRecordHandler)
	m.Get("/1/record/get",
		binding.Form(getRecordForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		getRecordHandler)
	m.Get("/1/record/timeline",
		binding.Form(recTimelineForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		recTimelineHandler)
	m.Get("/1/record/statistics",
		binding.Form(userRecStatForm{}),
		ErrorHandler,
		userRecStatHandler)
	m.Get("/1/leaderboard/list",
		binding.Form(leaderboardForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		leaderboardHandler)
	m.Get("/1/leaderboard/gameList",
		binding.Form(gamelbForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		gamelbHandler,
	)
}

type record struct {
	Id        string   `json:"record_id"`
	Type      string   `json:"type"`
	Source    string   `json:"source"`
	BeginTime int64    `json:"begin_time"`
	EndTime   int64    `json:"end_time"`
	Duration  int64    `json:"duration"`
	Distance  int      `json:"distance"`
	Weight    int      `json:"weight"`
	Mood      string   `json:"mood"`
	HeartRate int      `json:"heart_rate"`
	Pics      []string `json:"sport_pics"`
	GameType  string   `json:"game_type"`
	GameScore int      `json:"game_score"`
	GameName  string   `json:"game_name"`
	Coins     int64    `json:"coin_value"`
	Magic     int      `json:"magic"`
	Status    string   `json:"status"`
	Result    string   `json:"result"`
}

func convertRecord(r *models.Record) *record {
	rec := &record{}

	rec.Id = r.Id.Hex()
	rec.Type = r.Type
	rec.Status = r.Status
	rec.Coins = r.Coin
	rec.BeginTime = r.StartTime.Unix()
	rec.EndTime = r.EndTime.Unix()

	if r.Sport != nil {
		rec.Source = r.Sport.Source
		rec.Duration = r.Sport.Duration
		rec.Distance = r.Sport.Distance
		rec.Weight = r.Sport.Weight
		rec.Mood = r.Sport.Mood
		rec.HeartRate = r.Sport.HeartRate
		rec.Pics = r.Sport.Pics
		rec.Result = r.Sport.Review
	}
	if r.Game != nil {
		rec.GameName = r.Game.Name
		rec.GameScore = r.Game.Score
		rec.Magic = r.Game.Magic
		if r.Game.Coin > 0 {
			rec.Coins = r.Game.Coin
		}
	}
	return rec
}

type newRecordForm struct {
	Record *record `json:"record_item" binding:"required"`
	Task   int64   `json:"task_id"`
	Public bool    `json:"isPublic"`
	parameter
}

func gameAwards(level int64, gameScore int, isTask bool) Awards {
	awards := Awards{}

	base := 1.0
	scale := 1.0
	factor := 0.5

	if !isTask {
		base = 1.0
		scale = 1.0
		factor = 0.0
	}

	award := int64(scale * (base + factor*float64(level) + factor*(float64(gameScore)/100.0)))
	awards.Wealth = models.Satoshi * award
	awards.Mental = award
	awards.Score = award

	return awards
}

func newRecordHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(newRecordForm)

	rec := &models.Record{
		Uid:       user.Id,
		Task:      form.Task,
		Type:      form.Record.Type,
		StartTime: time.Unix(form.Record.BeginTime, 0),
		EndTime:   time.Unix(form.Record.EndTime, 0),
		PubTime:   time.Now(),
	}

	awards := Awards{}
	level := user.Level()

	switch form.Record.Type {
	case "post":
		rec.Status = models.StatusFinish
		if form.Task > 0 {
			rec.Coin = (2 + level) * models.Satoshi
			awards = Awards{
				Score:   2 + level,
				Literal: 2 + level,
				Wealth:  rec.Coin,
			}
		} else {
			rec.Coin = 2 * models.Satoshi
			awards = Awards{
				Score:   2,
				Literal: 2,
				Wealth:  rec.Coin,
			}
		}

		GiveAwards(user, awards, redis)

		//redis.AddPost(user.Id, "", 1)
		//user.UpdateStat(models.StatArticles, 1)
	case "game":
		if form.Task > 0 {
			awards = gameAwards(level, form.Record.GameScore, true)
			//user.AddTask(models.Tasks[form.Task-1].Type, form.Task, nil)
			rec.Status = models.StatusFinish
		} else {
			if form.Record.GameScore >= 100 &&
				user.Stat != nil && user.Stat.LastGameTime < nowDate().Unix() {
				awards = gameAwards(level, form.Record.GameScore, false)
			}
			user.UpdateStat(models.StatLastGameTime, time.Now().Unix())
		}
		GiveAwards(user, awards, redis)

		rec.Game = &models.GameRecord{
			Type:     form.Record.GameType,
			Name:     form.Record.GameName,
			Duration: form.Record.Duration,
			Score:    form.Record.GameScore,
			Magic:    int(awards.Mental),
		}
		rec.Coin = awards.Wealth

		redis.AddGameTime(user.Id, int(form.Record.Duration))
		user.UpdateStat(models.StatGameTime, form.Record.Duration)
		redis.SetGameMaxScore(gameType(rec.Game.Type), user.Id, rec.Game.Score)
		user.SetGameTime(gameType(rec.Game.Type), time.Now())

	case "run":
		if rec.Task > 0 {
			rec.Delete()
		}
		rec.Sport = &models.SportRecord{
			Source:    form.Record.Source,
			Duration:  form.Record.Duration,
			Distance:  form.Record.Distance,
			Weight:    form.Record.Weight,
			Mood:      form.Record.Mood,
			HeartRate: form.Record.HeartRate,
			Pics:      form.Record.Pics,
		}

		if rec.Sport.Weight == 0 {
			rec.Sport.Weight = user.Weight
		}
		// update weight
		if rec.Sport.Weight != user.Weight {
			user.SetInfo(&models.SetInfo{Weight: rec.Sport.Weight})
		}

		if form.Record.Duration > 0 {
			rec.Sport.Speed = float64(form.Record.Distance) / float64(form.Record.Duration)
		}
		rec.Status = models.StatusAuth

		redis.AddRecord(user.Id, 1)
		user.UpdateStat(models.StatRecords, 1)
		/*
			if len(form.Record.Source) > 0 {
				level := user.Level()
				awards = Awards{
					Physical: 30 + level,
					Wealth:   30 * models.Satoshi,
					Score:    30 + level,
				}

				GiveAwards(user, awards, redis)
				redis.UpdateRecLB(user.Id, rec.Sport.Distance, int(rec.Sport.Duration))

				rec.Coin = awards.Wealth
				rec.Status = models.StatusFinish
			}
		*/

	default:
		log.Println("Unknown record type:", form.Record.Type)
	}
	// assign task id
	if rec.Task == 0 {
		rec.Task = rec.PubTime.Unix()
	}
	if err := rec.Save(); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	if rec.Type == "run" {
		article := &models.Article{
			Author:  user.Id,
			PubTime: time.Now(),
			Loc:     user.Loc,
			Record:  rec.Id.Hex(),
			Type:    models.ArticleRecord,
		}

		if !form.Public {
			article.Privilege = models.PrivPrivate
		}

		if err := article.Save(); err == nil {
			redis.AddPost(user.Id, "", 1)
			user.UpdateStat(models.StatArticles, 1)
		}

		// push event to admins
		admins, _ := models.FindAdmins()
		for i, _ := range admins {
			if admins[i].Id != "1438851594268" {
				continue
			}
			event := &models.Event{
				Type: models.EventArticle,
				Time: time.Now().Unix(),
				Data: models.EventData{
					Type: models.EventRecord,
					Id:   article.Id.Hex(),
					From: user.Id,
					To:   admins[i].Id,
					Body: []models.MsgBody{
						//{Type: "total_count", Content: strconv.Itoa(parent.ReviewCount + 1)},
						{Type: "image", Content: admins[i].Profile},
					},
				},
			}
			event.Save()
			redis.PubMsg(models.EventArticle, admins[i].Id, event.Bytes())
		}
	}

	if form.Task > 0 {
		user.UpdateTask(int(form.Task), rec.Status)
		/*
			if rec.Status == models.StatusFinish {
				// ws push
				event := &models.Event{
					Type: models.EventNotice,
					Time: time.Now().Unix(),
					Data: models.EventData{
						Type: models.EventTaskDone,
						To:   user.Id,
						Body: []models.MsgBody{
							{Type: "task_id", Content: strconv.Itoa(int(form.Task))},
							{Type: "literature_value", Content: strconv.FormatInt(awards.Literal, 10)},
							{Type: "magic_value", Content: strconv.FormatInt(awards.Mental, 10)},
							{Type: "coin_value", Content: strconv.FormatInt(awards.Wealth, 10)},
						},
					},
				}
				event.Save()
				redis.PubMsg(event.Type, event.Data.To, event.Bytes())
			}
		*/
	}

	respData := map[string]interface{}{
		"record_id": rec.Id.Hex(),
		"ExpEffect": awards,
	}
	writeResponse(request.RequestURI, resp, respData, nil)
}

type getRecordForm struct {
	Record string `form:"record_id"`
	parameter
}

func getRecordHandler(r *http.Request, w http.ResponseWriter, p Parameter) {
	rec := &models.Record{}
	form := p.(getRecordForm)
	if bson.IsObjectIdHex(form.Record) {
		rec.Id = bson.ObjectIdHex(form.Record)
		rec.Find()
	}

	writeResponse(r.RequestURI, w, convertRecord(rec), nil)
}

type recTimelineForm struct {
	Userid string `form:"userid" binding:"required"`
	Type   string `form:"type"`
	models.Paging
	parameter
}

func recTimelineHandler(request *http.Request, resp http.ResponseWriter,
	user *models.Account, p Parameter) {
	form := p.(recTimelineForm)
	all := false
	if user.Id == form.Userid {
		all = true
	}

	u := &models.Account{Id: form.Userid}
	_, records, err := u.Records(all, form.Type, &form.Paging)

	recs := make([]*record, len(records))
	for i, _ := range records {
		recs[i] = convertRecord(&records[i])
	}
	respData := map[string]interface{}{
		"record_list":   recs,
		"page_frist_id": form.Paging.First,
		"page_last_id":  form.Paging.Last,
	}
	writeResponse(request.RequestURI, resp, respData, err)
}

type leaderboardResp struct {
	Userid   string `json:"userid"`
	Nickname string `json:"nikename"`
	Profile  string `json:"user_profile_image"`
	Rank     int    `json:"index,omitempty"`
	Score    int64  `json:"score"`
	Level    int64  `json:"rankLevel"`
	Gender   string `json:"sex_type"`
	Birth    int64  `json:"birthday"`
	Actor    string `json:"actor"`
	models.Location
	LastLog  int64  `json:"recent_login_time"`
	Addr     string `json:"locaddr"`
	Distance int    `json:"total_distance"`
	Status   string `json:"status"`
	Phone    string `json:"phone_number"`
	Coins    int64  `json:"coins"`
}

type leaderboardForm struct {
	Type string `form:"query_type"`
	Info string `form:"query_info"`
	models.Paging
	parameter
}

func leaderboardPaging(paging *models.Paging) (start, stop int) {
	start, _ = strconv.Atoi(paging.First)
	stop, _ = strconv.Atoi(paging.Last)
	if start == 0 && stop == 0 {
		stop = paging.Count - 1
		return
	}
	if start > 0 {
		stop = start - 2
		start = stop - paging.Count
		if stop < 0 {
			stop = 0
			start = 1 // start > stop empty set
			return
		}
		if start < 0 {
			start = 0
		}
	}
	if stop > 0 {
		start = stop
		stop = start + paging.Count
	}
	return
}

func leaderboardHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, form leaderboardForm) {
	if form.Paging.Count == 0 {
		form.Paging.Count = models.DefaultPageSize
	}

	start := 0
	stop := 0

	switch form.Type {
	case "FRIEND":
		ids := redis.Friends("friend", user.Id)
		friends, err := models.Users(ids, &form.Paging)
		if err != nil {
			writeResponse(request.RequestURI, resp, nil, err)
			return
		}
		lb := make([]leaderboardResp, len(friends))
		for i, _ := range friends {
			lb[i].Userid = friends[i].Id
			lb[i].Score = friends[i].Props.Score
			lb[i].Level = friends[i].Level()
			lb[i].Profile = friends[i].Profile
			lb[i].Nickname = friends[i].Nickname
			lb[i].Gender = friends[i].Gender
			lb[i].Actor = friends[i].Actor
			lb[i].LastLog = friends[i].LastLogin.Unix()
			lb[i].Birth = friends[i].Birth
			lb[i].Location = friends[i].Loc
			lb[i].Phone = friends[i].Phone
		}

		respData := map[string]interface{}{
			"members_list":  lb,
			"page_frist_id": form.Paging.First,
			"page_last_id":  form.Paging.Last,
		}
		writeResponse(request.RequestURI, resp, respData, nil)

		return

	case "USER_AROUND":
		rank := redis.LBDisRank(form.Info)
		if rank < 0 {
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.NotExistsError, "user not exist"))
			return
		}

		if form.Paging.Count < 0 {
			start = rank
			stop = rank
			break
		}

		start = rank - form.Paging.Count
		if start < 0 {
			start = 0
		}
		stop = rank + form.Paging.Count
	case "TOP":
		fallthrough
	default:
		start, stop = leaderboardPaging(&form.Paging)
	}

	kv := redis.GetDisLB(start, stop)
	ids := make([]string, len(kv))
	for i, _ := range kv {
		ids[i] = kv[i].K
	}

	users, _ := models.FindUsersByIds(0, ids...)

	lb := make([]leaderboardResp, len(kv))
	for i, _ := range kv {
		lb[i].Userid = kv[i].K
		lb[i].Rank = start + i + 1
		lb[i].Score = kv[i].V
		for _, user := range users {
			if user.Id == kv[i].K {
				lb[i].Nickname = user.Nickname
				lb[i].Profile = user.Profile
				break
			}
		}
	}

	page_first := 0
	page_last := 0
	if len(lb) > 0 {
		page_first = lb[0].Rank
		page_last = lb[len(lb)-1].Rank
	}

	respData := map[string]interface{}{
		"members_list":  lb,
		"page_frist_id": strconv.Itoa(page_first),
		"page_last_id":  strconv.Itoa(page_last),
	}
	writeResponse(request.RequestURI, resp, respData, nil)
}

type userRecStatForm struct {
	Userid string `form:"userid" binding:"required"`
	Token  string `form:"access_token"`
}

type statResp struct {
	RecCount      int     `json:"total_records_count"`
	TotalDistance int     `json:"total_distance"`
	TotalDuration int     `json:"total_duration"`
	MaxDistance   *record `json:"max_distance_record"`
	MaxSpeed      *record `json:"max_speed_record"`
	Actor         string  `json:"actor"`
	Score         int64   `json:"rankscore"`
	Level         int64   `json:"rankLevel"`
	Rank          string  `json:"rankName"`
	Index         int     `json:"top_index"`
	LBCount       int     `json:"leaderboard_max_items"`
}

func userRecStatHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form userRecStatForm) {
	user := &models.Account{}
	stats := &statResp{}
	if find, err := user.FindByUserid(form.Userid); !find {
		e := errors.NewError(errors.NotExistsError, "user not found")
		if err != nil {
			e = errors.NewError(errors.DbError, err.Error())
		}
		writeResponse(request.RequestURI, resp, nil, e)
		return
	}

	stats.RecCount, _ = models.TotalRecords(form.Userid)
	stats.TotalDistance, stats.TotalDuration = redis.RecStats(form.Userid)
	maxDisRec, _ := models.MaxDistanceRecord(form.Userid)
	maxSpeedRec, _ := models.MaxSpeedRecord(form.Userid)

	stats.MaxDistance = &record{}
	if maxDisRec.Sport != nil {
		stats.MaxDistance.Type = maxDisRec.Type
		stats.MaxDistance.Source = maxDisRec.Sport.Source
		stats.MaxDistance.BeginTime = maxDisRec.StartTime.Unix()
		stats.MaxDistance.EndTime = maxDisRec.EndTime.Unix()
		stats.MaxDistance.Duration = maxDisRec.Sport.Duration
		stats.MaxDistance.Distance = maxDisRec.Sport.Distance
	}

	stats.MaxSpeed = &record{}
	if maxSpeedRec.Sport != nil {
		stats.MaxDistance.Type = maxDisRec.Type
		stats.MaxDistance.Source = maxDisRec.Sport.Source
		stats.MaxDistance.BeginTime = maxDisRec.StartTime.Unix()
		stats.MaxDistance.EndTime = maxDisRec.EndTime.Unix()
		stats.MaxSpeed.Duration = maxSpeedRec.Sport.Duration
		stats.MaxSpeed.Distance = maxSpeedRec.Sport.Distance
	}

	stats.Score = user.Props.Score
	stats.Actor = user.Actor //userActor(user.Actor)
	stats.Level = user.Level()
	//stats.Rank = userRank(stats.Level)

	stats.Index = redis.LBDisRank(form.Userid) + 1
	stats.LBCount = redis.LBDisCard()

	writeResponse(request.RequestURI, resp, stats, nil)
}

type gamelbForm struct {
	Query string `form:"query_type"`
	Game  string `form:"game_type"`
	Score int    `form:"game_score"`
	Index int    `form:"page_index"`
	Count int    `form:"page_count"`
	parameter
}

func gamelbHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(gamelbForm)
	if form.Index <= 0 {
		form.Index = 0
	}
	if form.Count <= 0 {
		form.Count = 20
	}

	gt := gameType(form.Game)
	var ids []string
	var kvs []models.KV

	switch form.Query {
	case "FRIEND":
		ids = redis.Friends(models.RelFriend, user.Id)
		if len(ids) == 0 {
			break
		}

		ids = append(ids, user.Id)

		if form.Index*form.Count >= len(ids) {
			break
		}

		start := form.Index * form.Count
		end := start + form.Count
		if end > len(ids) {
			end = len(ids)
		}

		scores := redis.UserGameScores(gt, ids...)
		if len(scores) != len(ids) {
			scores = make([]int64, len(ids))
		}

		kvs = make([]models.KV, len(ids))
		for i, _ := range kvs {
			kvs[i].K = ids[i]
			kvs[i].V = scores[i]
			if ids[i] == user.Id {
				kvs[i].V = int64(form.Score)
			}
		}
		sort.Sort(sort.Reverse(models.KVSlice(kvs)))

		kvs = kvs[start:end]
		ids = []string{}
		for _, kv := range kvs {
			ids = append(ids, kv.K)
		}
	case "TOP":
		fallthrough
	default:
		maxScore := 0
		if scores := redis.UserGameScores(gt, user.Id); len(scores) == 1 {
			maxScore = int(scores[0])
		}
		redis.SetGameScore(gt, user.Id, form.Score) // current score
		kvs = redis.GameScores(gt, form.Index*form.Count, form.Count)
		redis.SetGameScore(gt, user.Id, maxScore) // recover max score
		for i, kv := range kvs {
			ids = append(ids, kv.K)
			if kv.K == user.Id {
				kvs[i].V = int64(form.Score)
			}
		}
	}

	var respData struct {
		List []*leaderboardResp `json:"members_list"`
	}

	users, _ := models.FindUsersByIds(1, ids...)
	index := 0
	for _, kv := range kvs {
		for i, _ := range users {
			if users[i].Id == kv.K {
				respData.List = append(respData.List, &leaderboardResp{
					Userid:   users[i].Id,
					Score:    kv.V,
					Rank:     form.Index*form.Count + index + 1,
					Level:    users[i].Level(),
					Profile:  users[i].Profile,
					Nickname: users[i].Nickname,
					Actor:    users[i].Actor,
					Gender:   users[i].Gender,
					LastLog:  users[i].LastGameTime(gt).Unix(),
					Birth:    users[i].Birth,
					Location: users[i].Loc,
					Phone:    users[i].Phone,
				})
				index++

				break
			}
		}
	}

	writeResponse(r.RequestURI, w, respData, nil)
}
