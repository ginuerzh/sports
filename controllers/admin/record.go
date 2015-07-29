package admin

import (
	//"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"log"
	"net/http"
	"time"
)

var defaultRecordsCount = 50

func BindRecordsApi(m *martini.ClassicMartini) {
	m.Get("/admin/record/timeline", binding.Form(getRecordsForm{}), adminErrorHandler, getRecordsListHandler)
	m.Post("/admin/record/delete", binding.Json(deleteRecordsForm{}), adminErrorHandler, deleteRecordsHandler)
	m.Options("/admin/record/delete", optionsHandler)
}

type getRecordsForm struct {
	Userid string `form:"userid" binding:"required"`
	Type   string `form:"type"`
	Count  int    `form:"page_count"`
	Page   int    `form:"page_index"`
	//Count      int    `form:"count"`
	//NextCursor string `form:"next_cursor"`
	//PrevCursor string `form:"prev_cursor"`
	FromTime int64  `form:"from_time"`
	ToTime   int64  `form:"to_time"`
	Token    string `form:"access_token" binding:"required"`
}

type record struct {
	ID         string   `json:"record_id"`
	Type       string   `json:"type"`
	Duration   int      `json:"duration"`
	Distance   int      `json:"distance"`
	Images     []string `json:"images"`
	GameName   string   `json:"game_name"`
	GameScore  int      `json:"game_score"`
	RecTime    int64    `json:"time"`
	PubTime    int64    `json:"pub_time"`
	RecTimeStr string   `json:"time_str"`
	PubTimeStr string   `json:"pub_time_str"`
}

type recordsListJsonStruct struct {
	Records []record `json:"records"`
	//NextCursor  string   `json:"next_cursor"`
	//PrevCursor  string   `json:"prev_cursor"`
	Page        int `json:"page_index"`
	PageTotal   int `json:"page_total"`
	TotalNumber int `json:"total_number"`
}

func getRecordsListHandler(w http.ResponseWriter, redis *models.RedisLogger, form getRecordsForm) {
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}

	count := form.Count
	if count == 0 {
		count = defaultCount
	}

	fromTime := time.Unix(0, 0)
	toTime := time.Now()
	//tn, records, err := models.GetRecords(form.Userid, form.Type, form.NextCursor, form.PrevCursor, form.Count, form.FromTime, form.ToTime, 0, getCount)
	total, records, _ := models.GetRecords(form.Userid, form.Type, form.Count, fromTime, toTime, count*form.Page, count)

	list := make([]record, len(records))
	for i, _ := range records {
		list[i].ID = records[i].Uid
		list[i].Type = records[i].Type
		list[i].RecTime = records[i].PubTime.Unix()
		if records[i].Sport != nil {
			list[i].Duration = int(records[i].Sport.Duration)
			list[i].Distance = records[i].Sport.Distance
			list[i].Images = records[i].Sport.Pics
		}
		if records[i].Game != nil {
			list[i].GameName = records[i].Game.Name
			list[i].GameScore = records[i].Game.Score
		}
		list[i].PubTime = records[i].PubTime.Unix()
	}

	totalPage := total / count
	if total%count != 0 {
		totalPage++
	}

	info := &recordsListJsonStruct{
		Records:     list,
		Page:        form.Page,
		PageTotal:   totalPage,
		TotalNumber: total,
	}

	writeResponse(w, info)
}

type deleteRecordsForm struct {
	Userid   string `json:"userid" binding:"required"`
	Type     string `json:"type"`
	FromTime int64  `json:"from_time"`
	ToTime   int64  `json:"to_time"`
	Token    string `json:"access_token" binding:"required"`
}

func deleteRecordsHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form deleteRecordsForm) {
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(resp, err)
		return
	}

	count, err := models.RemoveRecordsByID(form.Userid, form.Type, form.FromTime, form.ToTime)
	if err != nil {
		writeResponse(resp, err)
		return
	}

	respData := map[string]interface{}{
		"count": count,
	}
	writeResponse(resp, respData)
}
