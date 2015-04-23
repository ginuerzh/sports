// common
package admin

import (
	//"bytes"
	"encoding/json"
	//"fmt"
	"github.com/ginuerzh/sports/errors"
	"github.com/martini-contrib/binding"
	"io"
	"net/http"
	"strings"
)

type AdminPaging struct {
	Next      string `form:"next_cursor" json:"next_cursor"`
	Pre       string `form:"prev_cursor" json:"prev_cursor"`
	Count     int    `form:"count" json:"count"`
	PageIndex int    `form:"page_index" json:"page_index"`
	PageCount int    `form:"page_count" json:"page_count"`
}

func writeRawResponse(w http.ResponseWriter, contentType string, data []byte) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,X_Requested_With")
	w.Write(data)
}

func writeResponse(w http.ResponseWriter, data interface{}) []byte {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin,X-Requested-With,X_Requested_With,Content-Type,Accept")
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET")

	b, _ := json.Marshal(data)

	s := strings.Replace(string(b), "172.24.222.54:8082", "172.24.222.42:8082", -1)

	w.Write([]byte(s))

	return b
}

func optionsHandler(w http.ResponseWriter) {
	writeResponse(w, nil)
}

func adminErrorHandler(err binding.Errors, w http.ResponseWriter) {
	if err.Len() > 0 {
		e := err[0]
		s := e.Classification + ": "
		if len(e.FieldNames) > 0 {
			s += strings.Join(e.FieldNames, ",")
		}
		s += " " + e.Message
		writeResponse(w, errors.NewError(errors.JsonError, s))
	}
}

func decodeJson(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
