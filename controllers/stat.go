// stat
package controllers

import (
//"github.com/ginuerzh/sports/errors"
//"github.com/ginuerzh/sports/models"
//"gopkg.in/go-martini/martini.v1"
//"net/http"
//"time"
)

const (
//ServerStatV1Uri = "/1/stat"
)

/*
func BindStatApi(m *martini.ClassicMartini) {
	m.Get(ServerStatV1Uri, serverStatHandler)
}

func serverStatHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger) {
	respData := make(map[string]interface{})
	respData["visitors"] = redis.VisitorsCount(3)
	respData["pv"] = redis.PV(models.DateString(time.Now()))
	respData["registers"] = redis.RegisterCount(3)

	respData["top_views"] = redis.ArticleTopView(3, 3)
	respData["top_reviews"] = redis.ArticleTopReview(3)
	respData["top_thumbs"] = redis.ArticleTopThumb(3)
	respData["onlines"] = redis.Onlines()
	//respData["users"] = redis.Users()

	writeResponse(request.RequestURI, resp, respData, nil)
}
*/
