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

var (
	WalletAddr string
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
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Headers", "Origin,X-Requested-With,X_Requested_With,Content-Type,Accept")
	b, _ := json.Marshal(response{ReqPath: uri, RespData: data, Error: err})

	s := strings.Replace(string(b), "172.24.222.54:8082", "172.24.222.42:8082", -1)

	fmt.Println("<<<", string(s))
	resp.Write([]byte(s))

	return []byte(s)
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

	if r.URL.Path == "/1/file/upload" ||
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
	fmt.Println("===", p.TokenId(), uid)
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
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.AccessError, "禁止操作，该帐号已被管理员锁定"))
	}
}

func userActor(actor string) string {
	return actor
}

func nowDate() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location())
}

func decodeJson(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

type Awards struct {
	Physical int64 `json:"exp_physique,omitempty"`
	Literal  int64 `json:"exp_literature,omitempty"`
	Mental   int64 `json:"exp_magic,omitempty"`
	Wealth   int64 `json:"exp_coin,omitempty"`
	Score    int64 `json:"exp_rankscore,omitempty"`
	Level    int64 `json:"exp_rankLevel,omitempty"`
}

func GiveAwards(user *models.Account, awards Awards, redis *models.RedisLogger) error {
	if awards.Level < 0 || awards.Score < 0 {
		panic("invalid level or score")
	}
	if _, err := sendCoin(user.Wallet.Addr, awards.Wealth); err != nil {
		return err
	}
	redis.SendCoins(user.Id, awards.Wealth)

	err := user.UpdateProps(models.Props{
		Physical: awards.Physical,
		Literal:  awards.Literal,
		Mental:   awards.Mental,
		//Wealth:   awards.Wealth,
		Score: awards.Score,
		//Level: awards.Level,
	})
	if err != nil {
		return err
	}

	if lvl := models.Score2Level(user.Props.Score + awards.Score); lvl > user.Level() {
		// ws push
		event := &models.Event{
			Type: models.EventNotice,
			Time: time.Now().Unix(),
			Data: models.EventData{
				Type: models.EventLevelUP,
				To:   user.Id,
			},
		}
		event.Save()
		redis.PubMsg(event.Type, event.Data.To, event.Bytes())
	}

	return nil
}

func sendCoin(toAddr string, amount int64) (string, error) {
	if len(toAddr) == 0 || amount <= 0 {
		return "", nil
	}
	resp, err := http.PostForm(CoinAddr+"/send", url.Values{"to": {toAddr}, "amount": {strconv.FormatInt(amount, 10)}})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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

func consumeCoin(fromAddr string, value int64) (string, error) {
	if fromAddr == "" || value <= 0 {
		return "", nil
	}

	user := &models.Account{}
	if find, _ := user.FindByWalletAddr(fromAddr); !find {
		return "", errors.NewError(errors.NotFoundError)
	}
	wal, err := getWallet(user.Wallet.Id, user.Wallet.Key)
	if err != nil {
		return "", err
	}
	outputs, amount, err := getUnspent(fromAddr, wal.Keys, value)
	if value > amount {
		return "", errors.NewError(errors.AccessError, "余额不足")
	}
	changeAddr := fromAddr
	if len(changeAddr) == 0 {
		changeAddr = wal.Keys[0].PubKey
	}
	rawtx, err := CreateRawTx2(outputs, amount, value, WalletAddr, changeAddr)
	if err != nil {
		return "", err
	}
	return sendRawTx(rawtx)
}

func gameType(t string) int {
	switch t {
	case "七夕跳跳跳", "QIXI":
		return 0x01
	case "密室逃脱", "MISHI":
		return 0x02
	case "熊出没", "XIONGCHUMO":
		return 0x03
	case "蜘蛛侠", "SPIDERMAN":
		return 0x04
	case "转你妹", "ZNM":
		return 0x05
	}
	return 0
}
