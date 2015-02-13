package admin

import (
	//"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"log"
	"net/http"
)

var defaultRecordsCount = 50

func BindRecordsApi(m *martini.ClassicMartini) {
	m.Get("/admin/record/timeline", binding.Form(getRecordsForm{}), adminErrorHandler, getRecordsListHandler)
	m.Post("/admin/record/delete", binding.Json(deleteRecordsForm{}), adminErrorHandler, deleteRecordsHandler)
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

func getRecordsListHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getRecordsForm) {
	valid, errT := checkToken(redis, form.Token)
	if !valid {
		writeResponse(resp, errT)
		return
	}

	getCount := form.Count
	if getCount == 0 {
		getCount = defaultRecordsCount
	}

	//tn, records, err := models.GetRecords(form.Userid, form.Type, form.NextCursor, form.PrevCursor, form.Count, form.FromTime, form.ToTime, 0, getCount)
	tn, records, err := models.GetRecords(form.Userid, form.Type, "", "", form.Count, form.FromTime, form.ToTime, getCount*form.Page, getCount)
	if err != nil {
		writeResponse(resp, err)
		return
	}
	/*
		if tn == 0 {
			writeResponse(resp, errors.NewError(errors.NotExistsError))
			return

		}
	*/
	tnvalid := len(records)
	recs := make([]record, tnvalid)
	for i, _ := range records {
		recs[i].ID = records[i].Uid
		recs[i].Type = records[i].Type
		recs[i].RecTime = records[i].PubTime.Unix()
		recs[i].RecTimeStr = records[i].PubTime.Format("2006-01-02 15:04:05")
		if records[i].Sport != nil {
			recs[i].Duration = int(records[i].Sport.Duration)
			recs[i].Distance = records[i].Sport.Distance
			recs[i].Images = records[i].Sport.Pics
		}
		if records[i].Game != nil {
			recs[i].GameName = records[i].Game.Name
			recs[i].GameScore = records[i].Game.Score
		}
		recs[i].PubTime = records[i].PubTime.Unix()
		recs[i].PubTimeStr = records[i].PubTime.Format("2006-01-02 15:04:05")
	}

	totalPage := tn / getCount
	if tn%getCount != 0 {
		totalPage++
	}

	if tnvalid == 0 {
		respData := &recordsListJsonStruct{
			Records:   recs,
			Page:      form.Page,
			PageTotal: totalPage,
			//NextCursor:  "",
			//PrevCursor:  "",
			TotalNumber: tn,
		}
		writeResponse(resp, respData)
	} else {
		respData := &recordsListJsonStruct{
			Records:   recs,
			Page:      form.Page,
			PageTotal: totalPage,
			//NextCursor:  records[tnvalid-1].Id.String(),
			//PrevCursor:  records[0].Id.String(),
			TotalNumber: tn,
		}
		writeResponse(resp, respData)
	}
}

type deleteRecordsForm struct {
	Userid   string `json:"userid" binding:"required"`
	Type     string `json:"type"`
	FromTime int64  `json:"from_time"`
	ToTime   int64  `json:"to_time"`
	Token    string `json:"access_token" binding:"required"`
}

func deleteRecordsHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form deleteRecordsForm) {
	valid, errT := checkToken(redis, form.Token)
	if !valid {
		writeResponse(resp, errT)
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
