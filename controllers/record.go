// record
package controllers

import (
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"net/http"
	"strconv"
	"time"
)

func BindRecordApi(m *martini.ClassicMartini) {
	m.Post("/1/record/new", binding.Json(newRecordForm{}), ErrorHandler, newRecordHandler)
	m.Get("/1/record/timeline", binding.Form(recTimelineForm{}), ErrorHandler, recTimelineHandler)
	m.Get("/1/leaderboard/list", binding.Form(leaderboardForm{}), ErrorHandler, leaderboardHandler)
	m.Get("/1/record/statistics", binding.Form(userRecStatForm{}), ErrorHandler, userRecStatHandler)
}

type record struct {
	Type     string `json:"type"`
	Time     int64  `json:"action_time"`
	Duration int64  `json:"duration"`
	Distance int    `json:"distance"`
}

type newRecordForm struct {
	Record *record `json:"record_item" binding:"required"`
	Token  string  `json:"access_token" binding:"required"`
}

func newRecordHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form newRecordForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	rec := &models.Record{
		Uid:      user.Id,
		Type:     form.Record.Type,
		Time:     time.Unix(form.Record.Time, 0),
		Duration: form.Record.Duration,
		Distance: form.Record.Distance,
		PubTime:  time.Now(),
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

	score := actionExps[ActPost]
	redis.AddScore(user.Id, score)

	respData := map[string]int{
		"leaderboard_effect": rankDiff,
		"self_record_effect": recDiff,
		"exp_effect":         score,
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
		recs[i].Duration = records[i].Duration
		recs[i].Distance = records[i].Distance
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
	Rank     int    `json:"index"`
	Score    int    `json:"score"`
}

type leaderboardForm struct {
	Type  string `form:"query_type"`
	Info  string `form:"query_info"`
	Token string `form:"access_token"`
	models.Paging
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

func leaderboardHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form leaderboardForm) {
	/*
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
			return
		}
	*/

	if form.Paging.Count == 0 {
		form.Paging.Count = models.DefaultPageSize
	}

	start := 0
	stop := 0

	switch form.Type {
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

	users, _ := models.FindUsers(ids)

	lb := make([]leaderboardResp, len(kv))
	for i, _ := range kv {
		lb[i].Userid = kv[i].K
		lb[i].Rank = start + i + 1
		lb[i].Score = int(kv[i].V)
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
		"leaderboard_list": lb,
		"page_frist_id":    strconv.Itoa(page_first),
		"page_last_id":     strconv.Itoa(page_last),
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
	Score         int     `json:"rankscore"`
	Level         int     `json:"rankLevel"`
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
	if maxDisRec != nil {
		stats.MaxDistance = &record{
			Type:     maxDisRec.Type,
			Time:     maxDisRec.Time.Unix(),
			Duration: maxDisRec.Duration,
			Distance: maxDisRec.Distance,
		}
	}
	stats.Score = redis.UserScore(form.Userid)
	stats.Actor = userActor(user.Actor)
	stats.Level = userLevel(stats.Score)
	stats.Rank = userRank(stats.Level)

	stats.Index = redis.LBDisRank(form.Userid) + 1
	stats.LBCount = redis.LBDisCard()

	writeResponse(request.RequestURI, resp, stats, nil)
}
