// record
package controllers

import (
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"log"
	"net/http"
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
	m.Get("/1/record/timeline",
		binding.Form(recTimelineForm{}),
		ErrorHandler,
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
}

type record struct {
	Type      string   `json:"type"`
	Time      int64    `json:"action_time"`
	Duration  int64    `json:"duration"`
	Distance  int      `json:"distance"`
	Pics      []string `json:"sport_pics"`
	GameScore int      `json:"game_score"`
	GameName  string   `json:"game_name"`
	Magic     int      `json:"magic"`
	Status    string   `json:"status"`
}

type newRecordForm struct {
	Record *record `json:"record_item" binding:"required"`
	Task   int64   `json:"task_id"`
	parameter
}

func gameAwards(level int64, gameScore int, isTask bool) Awards {
	awards := Awards{}

	base := 5.0
	scale := 1.2
	factor := 1.5

	if !isTask {
		base = 1.0
		scale = 0.5
		factor = 0.5
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
		Uid:     user.Id,
		Task:    form.Task,
		Type:    form.Record.Type,
		Time:    time.Unix(form.Record.Time, 0),
		PubTime: time.Now(),
	}

	awards := Awards{}

	switch form.Record.Type {
	case "game":
		level := user.Level()
		if form.Task > 0 {
			awards = gameAwards(level, form.Record.GameScore, true)
			//user.AddTask(models.Tasks[form.Task-1].Type, form.Task, nil)
			rec.Status = models.StatusFinish
		} else {
			awards = gameAwards(level, form.Record.GameScore, false)
		}
		awards.Level = models.Score2Level(user.Props.Score+awards.Score) - level
		GiveAwards(user, awards, redis)

		rec.Game = &models.GameRecord{
			Name:  form.Record.GameName,
			Score: form.Record.GameScore,
			Magic: int(awards.Mental),
			Coin:  awards.Wealth,
		}

	case "run":
		rec.Sport = &models.SportRecord{
			Duration: form.Record.Duration,
			Distance: form.Record.Distance,
			Pics:     form.Record.Pics,
		}
		if form.Record.Duration > 0 {
			rec.Sport.Speed = float64(form.Record.Distance) / float64(form.Record.Duration)
		}
		rec.Status = models.StatusAuth
	default:
		log.Println("Unknown record type:", form.Record.Type)
	}
	if rec.Task == 0 {
		rec.Task = rec.PubTime.Unix()
	}
	if err := rec.Save(); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	rank := redis.LBDisRank(user.Id)
	maxDis := redis.MaxDisRecord(user.Id)
	redis.UpdateRecLB(user.Id, form.Record.Distance, int(form.Record.Duration))
	rankDiff := 0
	if rank >= 0 {
		rankDiff = redis.LBDisRank(user.Id) - rank
	}

	recDiff := 0
	if maxDis > 0 {
		recDiff = redis.MaxDisRecord(user.Id) - maxDis
	}

	respData := map[string]interface{}{
		"leaderboard_effect": rankDiff,
		"self_record_effect": recDiff,
		"ExpEffect":          awards,
	}
	writeResponse(request.RequestURI, resp, respData, nil)
}

type recTimelineForm struct {
	Userid string `form:"userid" binding:"required"`
	Token  string `form:"access_token"`
	models.Paging
}

func recTimelineHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form recTimelineForm) {
	user := &models.Account{Id: form.Userid}
	_, records, err := user.Records(&form.Paging)

	recs := make([]record, len(records))
	for i, _ := range records {
		recs[i].Type = records[i].Type
		recs[i].Time = records[i].Time.Unix()
		recs[i].Status = records[i].Status
		if records[i].Sport != nil {
			recs[i].Duration = records[i].Sport.Duration
			recs[i].Distance = records[i].Sport.Distance
			recs[i].Pics = records[i].Sport.Pics
		}
		if records[i].Game != nil {
			recs[i].GameName = records[i].Game.Name
			recs[i].GameScore = records[i].Game.Score
			recs[i].Magic = records[i].Game.Magic
		}
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
	models.Location
	LastLog  int64  `json:"recent_login_time"`
	Addr     string `json:"locaddr"`
	Distance int    `json:"total_distance"`
	Status   string `json:"status"`
	Phone    string `json:"phone_number"`
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

	users, _ := models.FindUsersByIds(ids)

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

	stats.MaxDistance = &record{
		Type: maxDisRec.Type,
		Time: maxDisRec.Time.Unix(),
	}
	if maxDisRec.Sport != nil {
		stats.MaxDistance.Duration = maxDisRec.Sport.Duration
		stats.MaxDistance.Distance = maxDisRec.Sport.Distance
	}

	stats.MaxSpeed = &record{
		Type: maxSpeedRec.Type,
		Time: maxSpeedRec.Time.Unix(),
	}
	if maxSpeedRec.Sport != nil {
		stats.MaxSpeed.Duration = maxSpeedRec.Sport.Duration
		stats.MaxSpeed.Distance = maxSpeedRec.Sport.Distance
	}

	stats.Score = user.Props.Score
	stats.Actor = userActor(user.Actor)
	stats.Level = user.Level()
	//stats.Rank = userRank(stats.Level)

	stats.Index = redis.LBDisRank(form.Userid) + 1
	stats.LBCount = redis.LBDisCard()

	writeResponse(request.RequestURI, resp, stats, nil)
}
