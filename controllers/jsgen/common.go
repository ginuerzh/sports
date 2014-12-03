// common
package jsgen

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	AuthError = NewError("错误提示", "登录认证错误")
	DbError   = NewError("错误提示", "数据库错误")
)

type Error struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func NewError(name, msg string) *Error {
	return &Error{
		Name:    name,
		Message: msg,
	}
}

func (e *Error) Error() string {
	return e.Name + ": " + e.Message
}

type pagination struct {
	Total     int `json:"total"`
	PageSize  int `json:"pageSize" form:"s"`
	PageIndex int `json:"pageIndex" form:"p"`
}

type response struct {
	Ack        bool        `json:"ack"`
	Err        error       `json:"error"`
	Time       int64       `json:"timestamp"`
	Data       interface{} `json:"data"`
	Pagination *pagination `json:"pagination"`
}

func writeResponse(w http.ResponseWriter, ack bool, data interface{}, p *pagination, err error) {
	r := &response{
		Ack:        ack,
		Err:        err,
		Time:       time.Now().Unix() * 1000,
		Data:       data,
		Pagination: p,
	}

	b, _ := json.Marshal(r)
	fmt.Println(string(b))
	w.Write(b)

	return
}

type socialItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type social struct {
	Weibo  socialItem `json:"weibo"`
	QQ     socialItem `json:"qq"`
	Google socialItem `json:"google"`
	Baidu  socialItem `json:"baidu"`
}
