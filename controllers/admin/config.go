package admin

import (
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"net/http"
)

func BindConfigApi(m *martini.ClassicMartini) {
	m.Get("/admin/config/get", binding.Form(getConfigForm{}), getConfigHandler)
	m.Post("/admin/config/set", binding.Json(setConfigForm{}), setConfigHandler)
	m.Options("/admin/config/set", optionsHandler)
}

type getConfigForm struct {
	Token string `form:"access_token" binding:"required"`
}

func getConfigHandler(w http.ResponseWriter, redis *models.RedisLogger, form getConfigForm) {
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}

	config := &models.Config{}
	config.Find()
	if config.Videos == nil {
		config.Videos = []models.Video{}
	}

	writeResponse(w, config)
}

type setConfigForm struct {
	Config *models.Config `json:"config"`
	Token  string         `json:"access_token" binding:"required"`
}

func setConfigHandler(w http.ResponseWriter, redis *models.RedisLogger, form setConfigForm) {
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}

	var err error
	if form.Config != nil {
		err = form.Config.Update()
	}
	writeResponse(w, err)
}
