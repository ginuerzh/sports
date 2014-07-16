// account
package controllers

import (
	"encoding/json"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	WeiboUserShowUrl  = "https://api.weibo.com/2/users/show.json"
	WeiboStatusUpdate = "https://api.weibo.com/2/statuses/update.json"
)

var (
	random *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func BindAccountApi(m *martini.ClassicMartini) {
	m.Post("/1/account/register", binding.Json(userRegForm{}), ErrorHandler, registerHandler)
	m.Post("/1/account/login", binding.Json(loginForm{}), ErrorHandler, loginHandler)
	m.Post("/1/user/logout", binding.Json(logoutForm{}), ErrorHandler, logoutHandler)
	m.Get("/1/user/getInfo", binding.Form(getInfoForm{}), ErrorHandler, userInfoHandler)
	m.Post("/1/user/setInfo", binding.Json(setInfoForm{}), ErrorHandler, setInfoHandler)
	m.Post("/1/user/set_profile_image", binding.Json(setProfileForm{}), ErrorHandler, setProfileHandler)

	//m.Get("/1/user/news", binding.Form(userNewsForm{}), ErrorHandler, userNewsHandler)
	m.Get("/1/users", binding.Form(userListForm{}), ErrorHandler, userListHandler)
}

// user register parameter
type userRegForm struct {
	Email    string `json:"email" binding:"required"`
	Nickname string `json:"nikename" binding:"required"`
	Password string `json:"password" binding:"required"`
	//Role     string `json:"role"`
}

func registerHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form userRegForm) {
	user := &models.Account{}

	user.Id = strings.ToLower(form.Email)
	user.Nickname = form.Nickname
	user.Password = Md5(form.Password)
	user.Role = "usrpass"
	user.RegTime = time.Now()
	//user.LastAccess = time.Now()
	//user.Online = true

	if err := user.Save(); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
	} else {
		token := Uuid()
		data := map[string]string{"access_token": token}
		writeResponse(request.RequestURI, resp, data, nil)

		redis.LogRegister(user.Id)
		redis.LogOnlineUser(token, user)
		redis.LogVisitor(user.Id)
	}
}

// user login parameter
type loginForm struct {
	Userid   string `json:"userid"`
	Password string `json:"verfiycode"`
	Type     string `json:"account_type" binding:"required"`
}

type weiboInfo struct {
	ScreenName  string `json:"screen_name"`
	Gender      string `json:"gender"`
	Url         string `json:"url"`
	Avatar      string `json:"avatar_large"`
	Location    string `json:"location"`
	Description string `json:"description"`
	ErrorDesc   string `json:"error"`
	ErrCode     int    `json:"error_code"`
}

func weiboLogin(uid, password string, redis *models.RedisLogger) (bool, *models.Account, error) {
	weibo := weiboInfo{}
	user := &models.Account{}

	v := url.Values{}
	v.Set("uid", uid)
	v.Set("access_token", password)

	url := WeiboUserShowUrl + "?" + v.Encode()
	resp, err := http.Get(url)
	if err != nil {
		return false, nil, errors.NewError(errors.HttpError)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, nil, errors.NewError(errors.HttpError)
	}

	if err := json.Unmarshal(data, &weibo); err != nil {
		return false, nil, errors.NewError(errors.HttpError)
	}

	if weibo.ErrCode != 0 {
		log.Println(weibo.ErrorDesc)
		return false, nil, errors.NewError(errors.AccessError)
	}

	user.Id = strings.ToLower(uid)
	user.Password = Md5(password)
	exist, err := user.Exists()
	if err != nil {
		return false, nil, err
	}

	if exist {
		user.ChangePassword(user.Password)
		return false, user, nil
	}

	user.Nickname = weibo.ScreenName
	user.Gender = weibo.Gender
	user.Url = weibo.Url
	user.Profile = weibo.Avatar
	user.Location = weibo.Location
	user.About = weibo.Description
	user.Role = "weibo"
	user.RegTime = time.Now()

	if err := user.Save(); err != nil {
		return true, nil, err
	}
	redis.LogRegister(user.Id)

	return true, user, nil
}
func guestLogin(redis *models.RedisLogger) (*models.Account, error) {
	user := &models.Account{}
	user.Id = models.GuestUserPrefix + strconv.Itoa(time.Now().Nanosecond()) + ":" + strconv.Itoa(random.Intn(65536))

	return user, nil
}

func loginHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form loginForm) {
	user := &models.Account{}
	var err error
	var reg bool
	token := Uuid()

	switch form.Type {
	case "weibo":
		reg, user, err = weiboLogin(form.Userid, form.Password, redis)
	case "weixin":
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.UnimplementedError))
		return
	case "usrpass":
		var find bool
		if find, err = user.FindByUserPass(strings.ToLower(form.Userid), Md5(form.Password)); !find {
			if err == nil {
				err = errors.NewError(errors.AuthError)
			}
		}
	default: // guest
		//user, err = guestLogin(redis)
		//token = models.GuestUserPrefix + token // start with 'guest:' for redis checking
		err = errors.NewError(errors.UnimplementedError, "unknown type "+form.Type)
	}

	if err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	score := 0
	news, _ := user.UpdateAction(ActLogin, nowDate())
	if news {
		score = actionExps[ActLogin]
		redis.AddScore(user.Id, score)
	}

	data := map[string]interface{}{
		"access_token": token,
		"register":     reg,
		"exp_effect":   score,
	}
	writeResponse(request.RequestURI, resp, data, nil)

	redis.LogOnlineUser(token, user)
	redis.LogVisitor(user.Id)
}

type logoutForm struct {
	Token string `json:"access_token" binding:"required"`
}

func logoutHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form logoutForm) {
	redis.DelOnlineUser(form.Token)
	writeResponse(request.RequestURI, resp, nil, nil)

}

type getInfoForm struct {
	Userid string `form:"userid" binding:"required"`
}

type userJsonStruct struct {
	Userid   string `json:"userid"`
	Nickname string `json:"nikename"`
	Type     string `json:"account_type"`
	Phone    string `json:"phone_number"`
	About    string `json:"about"`
	Location string `json:"location"`
	Profile  string `json:"profile_image"`
	RegTime  int64  `json:"register_time"`
	Hobby    string `json:"hobby"`
	Height   int    `json:"height"`
	Weight   int    `json:"weight"`
	Birth    int64  `json:"birthday"`
	Actor    string `json:"actor"`
	Score    int    `json:"rankscore"`
	Level    int    `json:"rankLevel"`
	Rank     string `json:"rankName"`
	//Online   bool   `json:"online"`
}

func userInfoHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getInfoForm) {
	user := &models.Account{}

	if find, err := user.FindByUserid(form.Userid); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError, "user '"+form.Userid+"' not exists")
		}
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}
	score := redis.UserScore(user.Id)
	level := userLevel(score)
	rank := userRank(level)
	birth := user.Birth.Unix()
	if birth <= 0 {
		birth = 0
	}
	info := &userJsonStruct{
		Userid:   user.Id,
		Nickname: user.Nickname,
		Type:     user.Role,
		Phone:    user.Phone,
		About:    user.About,
		Location: user.Location,
		Profile:  user.Profile,
		RegTime:  user.RegTime.Unix(),
		Hobby:    user.Hobby,
		Height:   user.Height,
		Weight:   user.Weight,
		Birth:    birth,
		Actor:    userActor(user.Actor),
		Score:    score,
		Level:    level,
		Rank:     rank,
	}

	writeResponse(request.RequestURI, resp, info, nil)
}

type setInfoForm struct {
	Token string `json:"access_token" binding:"required"`
	models.UserInfo
}

func setInfoHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form setInfoForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}
	err := user.SetInfo(form.UserInfo)
	score := 0
	if !user.Setinfo && err == nil {
		score = actionExps[ActInfo]
		redis.AddScore(user.Id, score)
	}

	writeResponse(request.RequestURI, resp, map[string]int{"exp_effect": score}, err)

	user.UpdateAction(ActInfo, nowDate())
	redis.LogOnlineUser(form.Token, user)
}

type setProfileForm struct {
	ImageId string `json:"image_id" binding:"required"`
	Token   string `json:"access_token"  binding:"required"`
}

func setProfileHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form setProfileForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	err := user.ChangeProfile(form.ImageId)
	redis.LogOnlineUser(form.Token, user)
	score := 0
	if len(user.Profile) == 0 && err == nil {
		score = actionExps[ActProfile]
		redis.AddScore(user.Id, score)
	}
	writeResponse(request.RequestURI, resp, map[string]int{"exp_effect": score}, err)
}

type userListForm struct {
	PageNumber int `form:"page_number" json:"page_number"`
	//AccessToken string `form:"access_token" json:"access_token"`
}

func userListHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form userListForm) {
	pageSize := models.DefaultPageSize + 2
	total, users, err := models.UserList(pageSize*form.PageNumber, pageSize)
	if err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	jsonStructs := make([]userJsonStruct, len(users))
	for i, _ := range users {
		//view, thumb, review, _ := users[i].ArticleCount()

		jsonStructs[i].Userid = users[i].Id
		jsonStructs[i].Nickname = users[i].Nickname
		jsonStructs[i].Type = users[i].Role
		jsonStructs[i].Profile = users[i].Profile
		jsonStructs[i].Phone = users[i].Phone
		jsonStructs[i].Location = users[i].Location
		jsonStructs[i].About = users[i].About
		jsonStructs[i].RegTime = users[i].RegTime.Unix()
		//jsonStructs[i].Views = view
		//jsonStructs[i].Thumbs = thumb
		//jsonStructs[i].Reviews = review
		//jsonStructs[i].Online = redis.IsOnline(users[i].Id)
	}

	respData := make(map[string]interface{})
	respData["page_number"] = form.PageNumber
	respData["page_more"] = pageSize*(form.PageNumber+1) < total
	//respData["total"] = total
	respData["users"] = jsonStructs
	writeResponse(request.RequestURI, resp, respData, nil)
}
