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
	"gopkg.in/go-martini/martini.v1"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	//"os"
	"bytes"
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
	fmt.Println("<<<", string(b))
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

type requestBody struct {
	bytes.Buffer
}

func (rb *requestBody) Close() error {
	rb.Reset()
	return nil
}

func DumpReqBodyHandler(r *http.Request) {
	fmt.Println("###", r.URL)

	if r.MultipartForm != nil ||
		(r.URL.Path == "/ueditor/controller" && r.URL.Query().Get("action") == "uploadimage") {
		return
	}
	if r.Method == "GET" || r.Body == nil {
		return
	}
	rb := &requestBody{}
	if _, err := io.Copy(rb, r.Body); err != nil {
		log.Println(err)
	}
	fmt.Println(">>>", rb.String())
	r.Body = rb
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

type Parameter interface {
	TokenId() string
}

type parameter struct {
	Token string `form:"access_token" json:"access_token"`
}

func (p parameter) TokenId() string {
	return p.Token
}

func checkTokenHandler(c martini.Context, p Parameter, redis *models.RedisLogger, r *http.Request, w http.ResponseWriter) {
	uid := redis.OnlineUser(p.TokenId())
	if len(uid) == 0 {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.AccessError))
	}
	c.Map(&models.Account{Id: uid})
}

func loadUserHandler(c martini.Context, user *models.Account, redis *models.RedisLogger, r *http.Request, w http.ResponseWriter) {
	if find, err := user.FindByUserid(user.Id); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError)
		}
		writeResponse(r.RequestURI, w, nil, err)
		return
	}
}

func checkLimitHandler(user *models.Account, r *http.Request, w http.ResponseWriter) {
	if user.TimeLimit < 0 || user.TimeLimit > time.Now().Unix() {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.AccessError))
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

/*
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
*/
func nowDate() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location())
}

/*
var actionExps = map[string]int{
	ActLogin:   1,
	ActPost:    10,
	ActComment: 1,
	ActThumb:   1,
	ActInvite:  30,
	ActProfile: 20,
	ActInfo:    20,
}
*/

func decodeJson(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

/*
func updateProps(userid string, props *models.Props, redis *models.RedisLogger) (int, int, error) {
	props, err := redis.AddProps(userid, props)
	if err != nil {
		return 0, 0, err
	}
	score := models.UserScore(props)
	level := models.UserLevel(score)

	return score, level, nil
}
*/
type Awards struct {
	Physical int64 `json:"exp_physique,omitempty"`
	Literal  int64 `json:"exp_literature,omitempty"`
	Mental   int64 `json:"exp_magic,omitempty"`
	Wealth   int64 `json:"exp_coin,omitempty"`
	Score    int64 `json:"exp_rankscore,omitempty"`
	Level    int64 `json:"exp_rankLevel,omitempty"`
}

func GiveAwards(user *models.Account, awards Awards, redis *models.RedisLogger) error {
	if _, err := sendCoin(user.Wallet.Addr, awards.Wealth); err != nil {
		return err
	}
	redis.AddCoins(user.Id, awards.Wealth)

	return user.UpdateProps(models.Props{
		Physical: awards.Physical,
		Literal:  awards.Literal,
		Mental:   awards.Mental,
		//Wealth:   awards.Wealth,
		Score: awards.Score,
		Level: awards.Level,
	})
}

func sendCoin(toAddr string, amount int64) (string, error) {
	if len(toAddr) == 0 || amount <= 0 {
		return "", nil
	}
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
