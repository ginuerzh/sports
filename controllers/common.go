// common
package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/ginuerzh/sports/errors"
	"github.com/martini-contrib/binding"
	"github.com/nu7hatch/gouuid"
	"github.com/zhengying/apns"
	"io"
	//"log"
	"net/http"
	//"strconv"
	"strings"
	"time"
)

const (
	ActLogin   = "login"
	ActPost    = "post"
	ActComment = "comment"
	ActInvite  = "invite"
	ActProfile = "profile"
	ActInfo    = "info"
)

type response struct {
	ReqPath  string      `json:"req_path"`
	RespData interface{} `json:"response_data"`
	Error    error       `json:"error"`
}

func writeResponse(uri string, resp http.ResponseWriter, data interface{}, err error) []byte {
	if err == nil {
		err = errors.NewError(errors.NoError)
	}
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, _ := json.Marshal(response{ReqPath: uri, RespData: data, Error: err})
	fmt.Println(string(b))
	resp.Write(b)

	return b
}

func writeRawResponse(resp http.ResponseWriter, raw []byte) {
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp.Write(raw)
}

func Md5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func FileMd5(file io.Reader) string {
	h := md5.New()
	io.Copy(h, file)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Uuid() string {
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}

	return u4.String()
}

func sendApns(client *apns.Client, token, alert string, badge int, sound string) error {
	payload := apns.NewPayload()
	payload.Alert = alert
	payload.Badge = badge
	payload.Sound = sound

	pn := apns.NewPushNotification()
	pn.DeviceToken = token
	pn.AddPayload(payload)

	resp := client.Send(pn)
	return resp.Error
}

func ErrorHandler(err binding.Errors, request *http.Request, resp http.ResponseWriter) {
	if err.Len() > 0 {
		e := err[0]
		s := e.Classification + ": "
		if len(e.FieldNames) > 0 {
			s += strings.Join(e.FieldNames, ",")
		}
		s += " " + e.Message
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.JsonError, s))
	}
}

func userActor(actor string) string {
	/*
		switch actor {
		case "PRO":
			return "专业运动员"
		case "MID":
			return "业余运动员"
		case "AMATEUR":
			fallthrough
		default:
			return "爱好者"
		}
	*/
	return actor
}

var levelScores = []int{
	0, 20, 30, 45, 67, 101, 151, 227, 341, 512, /* 1 - 10 */
	768, 1153, 1729, 2594, 3892, 5838, 8757, 13136, 19705, 29557, /* 11 - 20 */
	44336, 66505, 99757, 149636, 224454, 336682, 505023, 757535, 1136302, 1704453, /* 21 - 30 */
	2556680, 3835021, 5752531, 8628797, 12943196, /* 31 - 35 */
	19414794, 29122192, 43683288, 65524932, 98287398, /* 36 - 40 */
}

func userLevel(score int) int {
	for i, s := range levelScores {
		if s > score {
			return i
		}
		if s == score {
			return i + 1
		}
	}
	return len(levelScores)
}

func userRank(level int) string {
	if level <= 10 {
		return "初级"
	} else if level <= 20 {
		return "中级"
	} else if level <= 30 {
		return "高级"
	} else {
		return "至尊"
	}
}

func nowDate() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location())
}

var actionExps = map[string]int{
	ActLogin:   1,
	ActPost:    10,
	ActComment: 1,
	ActInvite:  30,
	ActProfile: 20,
	ActInfo:    20,
}
