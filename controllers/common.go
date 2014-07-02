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
