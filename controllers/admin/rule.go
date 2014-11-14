// rule
package admin

import (
	//"encoding/json"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"net/http"
)

func BindRuleApi(m *martini.ClassicMartini) {
	m.Post("/admin/rule/set", binding.Json(setRuleForm{}), adminErrorHandler, setRuleHandler)
}

// admin login parameter
type setRuleForm struct {
	Id      int    `json:"rule_id" binding:"required"`
	Message string `json:"message"`
	Token   string `form:"access_token"`
}

func setRuleHandler(w http.ResponseWriter, redis *models.RedisLogger, form setRuleForm) {
	rule := &models.Rule{
		RuleId:  form.Id,
		Message: form.Message,
	}

	if err := rule.Save(); err != nil {
		writeResponse(w, err)
		return
	}

	writeResponse(w, map[string]interface{}{})
}
