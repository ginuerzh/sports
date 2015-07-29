package controllers

import (
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"log"
	"net/http"
)

func BindSysApi(m *martini.ClassicMartini) {
	m.Get("/1/sys/config",
		binding.Form(getSysConfForm{}),
		ErrorHandler,
		checkTokenHandler,
		getSysConfHandler)
	m.Get("/m",
		mobileRedirectHandler)
}

type getSysConfForm struct {
	parameter
}

func getSysConfHandler(r *http.Request, w http.ResponseWriter) {
	var sysConf struct {
		LevelScore []int64  `json:"level_score"`
		PKEffects  []string `json:"pk_effects"`
	}

	sysConf.LevelScore = models.GetLevelScore()
	sysConf.PKEffects = []string{
		"http://7xj9pe.com1.z0.glb.clouddn.com/pk/1-14.zip",
		"http://7xj9pe.com1.z0.glb.clouddn.com/pk/2-14.zip",
		"http://7xj9pe.com1.z0.glb.clouddn.com/pk/3-12.zip",
		"http://7xj9pe.com1.z0.glb.clouddn.com/pk/4-12.zip",
		"http://7xj9pe.com1.z0.glb.clouddn.com/pk/5-16.zip",
		"http://7xj9pe.com1.z0.glb.clouddn.com/pk/6-13.zip",
		"http://7xj9pe.com1.z0.glb.clouddn.com/pk/7-12.zip",
		"http://7xj9pe.com1.z0.glb.clouddn.com/pk/8-11.zip",
		"http://7xj9pe.com1.z0.glb.clouddn.com/pk/9-12.zip",
		"http://7xj9pe.com1.z0.glb.clouddn.com/pk/10-11.zip",
	}

	writeResponse(r.RequestURI, w, sysConf, nil)
}

func mobileRedirectHandler(r *http.Request, w http.ResponseWriter) {
	http.Redirect(w, r, "/sport_phone/index.html", http.StatusMovedPermanently)
}
