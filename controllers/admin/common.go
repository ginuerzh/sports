// common
package admin

import (
	"encoding/json"
	"fmt"
	"github.com/ginuerzh/sports/errors"
	"github.com/martini-contrib/binding"
	"net/http"
	"strings"
)

type AdminPaging struct {
	Next  string `form:"next_cursor" json:"next_cursor"`
	Pre   string `form:"prev_cursor" json:"prev_cursor"`
	Count int    `form:"count" json:"count"`
}

func writeResponse(w http.ResponseWriter, data interface{}) []byte {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	b, _ := json.Marshal(data)
	fmt.Println(string(b))
	w.Write(b)

	return b
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
