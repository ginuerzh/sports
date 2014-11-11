// common
package controllers

import (
	"crypto/md5"
	"encoding/json"
	errs "errors"
	"fmt"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"github.com/nu7hatch/gouuid"
	"github.com/zhengying/apns"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	//"os"
	"strings"
	"time"
)

const (
	ActLogin   = "login"
	ActPost    = "post"
	ActComment = "comment"
	ActThumb   = "thumb"
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

type GetToken interface {
	getTokenId() string
}

func CheckHandler(getT GetToken, redis *models.RedisLogger, request *http.Request, resp http.ResponseWriter) {
	token := getT.getTokenId()
	user := redis.OnlineUser(token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
	}
	log.Println("user.TimeLimit is :", user.TimeLimit, "cur time is:", time.Now().Unix())
	if user.TimeLimit == -1 || user.TimeLimit > time.Now().Unix() {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
	}
}

type getUser interface {
	getUserId() string
}

func CheckUserIDHandler(getU getUser, redis *models.RedisLogger, request *http.Request, resp http.ResponseWriter) {
	uid := getU.getUserId()
	user := &models.Account{}
	if find, err := user.FindByUserid(uid); !find {
		if err == nil {
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.NotExistsError, "user '"+uid+"' not exists"))
			return
		}
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.NotExistsError))
		return
	}

	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}
	log.Println("user.TimeLimit is :", user.TimeLimit, "cur time is:", time.Now().Unix())
	//if user.TimeLimit == -1 || user.TimeLimit > time.Now().Unix() {
	if user.TimeLimit == -1 {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
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
	ActThumb:   1,
	ActInvite:  30,
	ActProfile: 20,
	ActInfo:    20,
}

type Awards struct {
	Physical int64 `json:"exp_physique,omitempty"`
	Literal  int64 `json:"exp_literature,omitempty"`
	Mental   int64 `json:"exp_magic,omitempty"`
	Wealth   int64 `json:"exp_coin,omitempty"`
	Score    int64 `json:"exp_rankscore,omitempty"`
	Level    int64 `json:"exp_rankLevel,omitempty"`
}

func decodeJson(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func sendCoin(toAddr string, amount int64) (string, error) {
	resp, err := http.PostForm(CoinAddr+"/send", url.Values{"to": {toAddr}, "amount": {strconv.FormatInt(amount, 10)}})
	if err != nil {
		return "", err
	}

	var r struct {
		Txid   string `json:"txid"`
		Result string `json:"result"`
	}

	if err = decodeJson(resp.Body, &r); err != nil {
		return r.Txid, err
	}
	if r.Result != "ok" {
		return r.Txid, errs.New(r.Result)
	}

	return r.Txid, nil
}

func updateProps(userid string, props *models.Props, redis *models.RedisLogger) (int, int, error) {
	props, err := redis.AddProps(userid, props)
	if err != nil {
		return 0, 0, err
	}
	score := models.UserScore(props)
	level := models.UserLevel(score)

	return score, level, nil
}

func giveAwards(user *models.Account, awards *Awards, redis *models.RedisLogger) error {
	if awards.Wealth > 0 {
		_, err := sendCoin(user.Wallet.Addr, awards.Wealth)
		if err != nil {
			return err
		}
	}

	props := &models.Props{
		Physical: awards.Physical,
		Literal:  awards.Literal,
		Mental:   awards.Mental,
		Wealth:   awards.Wealth,
	}

	score, level, err := updateProps(user.Id, props, redis)
	if err != nil {
		return err
	}
	awards.Score = int64(score - user.Score)
	awards.Level = int64(level - user.Level)

	return user.UpdateLevel(score, level)
}
