// account
package controllers

import (
	//"encoding/json"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	//"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	random *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func BindAccountApi(m *martini.ClassicMartini) {
	m.Post("/1/account/register", binding.Json(userRegForm{}), ErrorHandler, registerHandler)
	m.Post("/1/account/login", binding.Json(loginForm{}, (*GetToken)(nil)), ErrorHandler, CheckUserIDHandler, loginHandler)
	m.Get("/1/user/getDailyLoginRewardInfo", binding.Form(loginAwardsForm{}), ErrorHandler, loginAwardsHandler)
	m.Post("/1/user/logout", binding.Json(logoutForm{}), ErrorHandler, logoutHandler)
	m.Get("/1/user/getInfo", binding.Form(getInfoForm{}), ErrorHandler, userInfoHandler)
	m.Get("/1/user/getRelatedMembersCount", binding.Form(friendCountForm{}), ErrorHandler, friendCountHandler)
	m.Post("/1/user/setInfo", binding.Json(setInfoForm{}, (*GetToken)(nil)), ErrorHandler, CheckHandler, setInfoHandler)
	m.Post("/1/user/setProfileImage", binding.Json(setProfileForm{}, (*GetToken)(nil)), ErrorHandler, CheckHandler, setProfileHandler)
	m.Post("/1/account/importFriends", binding.Json(importFriendsForm{}, (*GetToken)(nil)), ErrorHandler, CheckHandler, importFriendsHandler)

	m.Post("/1/user/setLifePhotos", binding.Json(setPhotosForm{}, (*GetToken)(nil)), ErrorHandler, CheckHandler, setPhotosHandler)
	m.Post("/1/user/deleteLifePhoto", binding.Json(delPhotoForm{}, (*GetToken)(nil)), ErrorHandler, CheckHandler, delPhotoHandler)
	//m.Get("/1/user/news", binding.Form(userNewsForm{}), ErrorHandler, userNewsHandler)
	m.Get("/1/users", binding.Form(userListForm{}), ErrorHandler, userListHandler)

	m.Get("/1/user/getPKPropertiesInfo", binding.Form(scoreDiffForm{}), ErrorHandler, scoreDiffHandler)
	m.Get("/1/user/getPropertiesValue", binding.Form(getPropsForm{}), getPropsHandler)
	m.Post("/1/user/updateEquipment", binding.Json(setEquipForm{}, (*GetToken)(nil)), ErrorHandler, CheckHandler, setEquipHandler)
	m.Get("/1/user/search", binding.Form(searchForm{}), ErrorHandler, searchHandler)
}

// user register parameter
type userRegForm struct {
	Email    string `json:"email" binding:"required"`
	Nickname string `json:"nikename"`
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
	dbw, err := getNewWallet()
	if err != nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.DbError, "wallet: "+err.Error()))
		return
	}
	user.Wallet = *dbw

	if err := user.Save(); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
	} else {
		token := Uuid()
		data := map[string]string{"access_token": token}
		writeResponse(request.RequestURI, resp, data, nil)

		redis.LogRegister(user.Id)
		//redis.SetOnlineUser(token, user, true)

		// ws push
		notice := &models.Event{
			Type: models.EventWallet,
			Time: time.Now().Unix(),
			Data: models.EventData{
				Type: models.EventTx,
				Id:   user.Id,
				From: user.Id,
				Body: []models.MsgBody{
					{Type: "rule", Content: "1"},
					{Type: "nikename", Content: user.Id},
					{Type: "total_count", Content: "1"},
					{Type: "image", Content: user.Profile},
				},
			},
		}
		redis.Notice(notice.Bytes())
	}
}

func getNewWallet() (*models.DbWallet, error) {
	w := NewWallet()
	id, err := saveWallet("", w)
	var addrs []string
	for _, key := range w.Keys {
		addrs = append(addrs, key.PubKey)
	}
	return &models.DbWallet{Id: id, Key: w.SharedKey, Addr: w.Keys[0].PubKey, Addrs: addrs}, err
}

// user login parameter
type loginForm struct {
	Userid   string `json:"userid"`
	Password string `json:"verfiycode"`
	Type     string `json:"account_type"`
}

func (this loginForm) getUserId() string {
	return this.Userid
}

func weiboLogin(uid, password string, redis *models.RedisLogger) (bool, *models.Account, error) {
	user := &models.Account{Id: strings.ToLower(uid)}
	exists, err := user.Exists()
	if err != nil {
		return false, nil, err
	}

	p := Md5(password)
	registered := user.RegTime.Unix() > 0

	if registered {
		if user.Password != p {
			user.ChangePassword(p)
		}
		return false, user, nil
	}
	weiboUser, err := GetWeiboUserInfo(uid, password)
	if err != nil {
		return false, nil, err
	}

	user.Nickname = weiboUser.ScreenName
	user.Password = p
	user.Gender = weiboUser.Gender
	user.Url = weiboUser.Url
	user.Profile = weiboUser.Avatar
	user.Addr = &models.Address{Desc: weiboUser.Location}
	user.About = weiboUser.Description
	user.Role = "weibo"
	user.RegTime = time.Now()

	dbw, err := getNewWallet()
	if err != nil {
		return true, nil, err
	}
	user.Wallet = *dbw

	if !exists {
		if err := user.Save(); err != nil {
			return true, nil, err
		}
	} else {
		if err := user.Update(); err != nil {
			return true, nil, err
		}
	}
	redis.LogRegister(user.Id)

	return true, user, nil
}
func guestLogin(redis *models.RedisLogger) (*models.Account, error) {
	user := &models.Account{}
	user.Id = models.GuestUserPrefix + strconv.Itoa(time.Now().Nanosecond()) + ":" + strconv.Itoa(random.Intn(65536))

	return user, nil
}

func loginHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, getU getUser) {
	form := getU.(loginForm)
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
		fallthrough
	default:
		var find bool
		if find, err = user.FindByUserPass(strings.ToLower(form.Userid), Md5(form.Password)); !find {
			if err == nil {
				err = errors.NewError(errors.AuthError)
			}
		}
	}

	if err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	//user.UpdateAction(ActLogin, d)
	redis.SetOnlineUser(token, user, true)
	redis.LogLogin(user.Id)

	lastlog := time.Now()
	d := nowDate()
	count := user.LoginCount

	if user.LastLogin.Unix() < d.Unix()-24*3600 {
		count = 1
	} else if user.LastLogin.Unix() < d.Unix() {
		count++
	}

	award, _ := user.SetLogin(count, lastlog)
	awards := Awards{}
	if user.LastLogin.Unix() < d.Unix() {
		log.Println("award")
		awards.Wealth = award * models.Satoshi
		if err := giveAwards(user, &awards, redis); err != nil {
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.DbError, err.Error()))
			return
		}
	}
	data := map[string]interface{}{
		"access_token":    token,
		"userid":          user.Id,
		"register":        reg,
		"last_login_time": user.LastLogin.Unix(),
		"ExpEffect":       awards,
	}
	writeResponse(request.RequestURI, resp, data, nil)
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
	Token  string `form:"access_token"`
}

type userJsonStruct struct {
	Userid   string `json:"userid"`
	Nickname string `json:"nikename"`

	Phone   string `json:"phone_number"`
	Type    string `json:"account_type"`
	About   string `json:"about"`
	Profile string `json:"profile_image"`
	RegTime int64  `json:"register_time"`
	Hobby   string `json:"hobby"`
	Height  int    `json:"height"`
	Weight  int    `json:"weight"`
	Birth   int64  `json:"birthday"`

	Actor string `json:"actor"`
	Rank  string `json:"rankName"`
	//Followed bool   `json:"beFriend"`
	Online bool `json:"beOnline"`

	Props *models.Props `json:"proper_info"`

	Addr string `json:"location_desc"`
	models.Location

	Gender    string `json:"sex_type"`
	Follows   int    `json:"attention_count"`
	Followers int    `json:"fans_count"`
	Posts     int    `json:"post_count"`

	Photos []string     `json:"user_images"`
	Equips models.Equip `json:"user_equipInfo"`

	Wallet   string `json:"wallet"`
	Relation string `json:"relation"`
	LastLog  int64  `json:"last_login_time"`
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

	info := &userJsonStruct{
		Userid:   user.Id,
		Nickname: user.Nickname,
		Phone:    user.Phone,
		Type:     user.Role,
		About:    user.About,
		Profile:  user.Profile,
		RegTime:  user.RegTime.Unix(),
		Hobby:    user.Hobby,
		Height:   user.Height,
		Weight:   user.Weight,
		Birth:    user.Birth,
		Actor:    userActor(user.Actor),

		Rank:   userRank(user.Level),
		Online: redis.IsOnline(user.Id),
		Gender: user.Gender,
		//Follows:   len(redis.Follows(user.Id)),
		//Followers: len(redis.Followers(user.Id)),
		Posts: user.ArticleCount(),

		Props: redis.UserProps(user.Id),

		Photos: user.Photos,

		Wallet:  user.Wallet.Addr,
		LastLog: user.LastLogin.Unix(),
	}

	info.Follows, info.Followers, _, _ = redis.FriendCount(user.Id)

	if user.Equips != nil {
		info.Equips = *user.Equips
	}

	if user.Addr != nil {
		info.Addr = user.Addr.String()
	}
	if user.Loc != nil {
		info.Location = *user.Loc
	}

	if user.Equips != nil {
		info.Equips = *user.Equips
	}

	if u := redis.OnlineUser(form.Token); u != nil {
		relation := redis.Relationship(u.Id, user.Id)
		switch relation {
		case models.RelFriend:
			info.Relation = "FRIENDS"
		case models.RelFollowing:
			info.Relation = "ATTENTION"
		case models.RelFollower:
			info.Relation = "FANS"
		case models.RelBlacklist:
			info.Relation = "DEFRIEND"
		}
	}

	writeResponse(request.RequestURI, resp, info, nil)
}

type friendCountForm struct {
	Token string `form:"access_token" binding:"required"`
}

func friendCountHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form friendCountForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	follows, followers, friends, blacklist := redis.FriendCount(user.Id)
	respData := map[string]int{
		"friend_count":    friends,
		"attention_count": follows,
		"fans_count":      followers,
		"defriend_count":  blacklist,
	}
	writeResponse(request.RequestURI, resp, respData, nil)
}

type setInfoForm struct {
	Token string `json:"access_token" binding:"required"`
	models.UserInfo
}

func (this setInfoForm) getTokenId() string {
	return this.Token
}

func setInfoHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, getT GetToken) {
	form := getT.(setInfoForm)
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	user.Nickname = form.UserInfo.Nickname
	user.Hobby = form.UserInfo.Hobby
	user.Height = form.UserInfo.Height
	user.Weight = form.UserInfo.Weight
	user.Birth = form.UserInfo.Birth
	user.Actor = form.UserInfo.Actor
	user.Gender = form.UserInfo.Gender
	user.Phone = form.UserInfo.Phone
	user.About = form.UserInfo.About

	addr := &models.Address{
		Country:  form.UserInfo.Country,
		Province: form.UserInfo.Province,
		City:     form.UserInfo.City,
		Area:     form.UserInfo.Area,
		Desc:     form.UserInfo.LocDesc,
	}
	if addr.String() != "" {
		user.Addr = addr
	}
	setinfo := user.Setinfo
	user.Setinfo = true
	err := user.Update()

	score := 0
	if !setinfo && err == nil {
		score = actionExps[ActInfo]
		//redis.AddScore(user.Id, score)
	}

	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": Awards{Wealth: int64(score)}}, err)

	user.UpdateAction(ActInfo, nowDate())
	redis.SetOnlineUser(form.Token, user, false)
}

type setProfileForm struct {
	ImageId string `json:"image_id" binding:"required"`
	Token   string `json:"access_token" binding:"required"`
}

func (this setProfileForm) getTokenId() string {
	return this.Token
}

func setProfileHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, getT GetToken) {
	form := getT.(setProfileForm)
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	err := user.ChangeProfile(form.ImageId)
	redis.SetOnlineUser(form.Token, user, false)
	/*
		score := 0
		if len(user.Profile) == 0 && err == nil {
			score = actionExps[ActProfile]
			//redis.AddScore(user.Id, score)
		}
	*/
	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": Awards{}}, err)
}

type setPhotosForm struct {
	Token string   `json:"access_token" binding:"required"`
	Pics  []string `json:"pic_ids"`
}

func (this setPhotosForm) getTokenId() string {
	return this.Token
}

func setPhotosHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, getT GetToken) {
	form := getT.(setPhotosForm)
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}
	err := user.AddPhotos(form.Pics)
	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": Awards{}}, err)
}

type delPhotoForm struct {
	Token string `json:"access_token" binding:"required"`
	Photo string `json:"pic_id"`
}

func (this delPhotoForm) getTokenId() string {
	return this.Token
}

func delPhotoHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, getT GetToken) {
	form := getT.(delPhotoForm)
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}
	err := user.DelPhoto(form.Photo)
	writeResponse(request.RequestURI, resp, nil, err)
}

type loginAwardsForm struct {
	Token string `form:"access_token" binding:"required"`
}

func loginAwardsHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form loginAwardsForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}
	user.FindByUserid(user.Id)
	respData := map[string]interface{}{
		"continuous_logined_days": user.LoginCount,
		"login_reward_list":       user.LoginAwards,
	}
	writeResponse(request.RequestURI, resp, respData, nil)
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
		//jsonStructs[i].Location = users[i].Location
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

type scoreDiffForm struct {
	Token string `form:"access_token" binding:"required"`
	Uid   string `form:"userid" binding:"required"`
}

func scoreDiffHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form scoreDiffForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	me := redis.UserProps(user.Id)
	you := redis.UserProps(form.Uid)

	respData := map[string]int64{
		"physique_times":   you.Physical - me.Physical,
		"literature_times": you.Literal - me.Literal,
		"magic_times":      you.Mental - me.Mental,
	}
	writeResponse(request.RequestURI, resp, respData, nil)
}

type getPropsForm struct {
	Uid string `form:"userid" binding:"required"`
}

func getPropsHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getPropsForm) {
	writeResponse(request.RequestURI, resp, redis.UserProps(form.Uid), nil)
}

type setEquipForm struct {
	Token  string       `json:"access_token" binding:"required"`
	Equips models.Equip `json:"user_equipInfo"`
}

func (this setEquipForm) getTokenId() string {
	return this.Token
}

func setEquipHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, getT GetToken) {
	form := getT.(setEquipForm)
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	err := user.SetEquip(form.Equips)
	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": Awards{}}, err)
}

type searchForm struct {
	Token    string `form:"access_token"`
	Nearby   bool   `form:"search_nearby"`
	Nickname string `form:"search_nickname"`
	models.Paging
}

func searchHandler(r *http.Request, w http.ResponseWriter, redis *models.RedisLogger, form searchForm) {

	users := []models.Account{}
	var err error

	if form.Nearby {
		user := redis.OnlineUser(form.Token)
		if user == nil {
			writeResponse(r.RequestURI, w, nil, errors.NewError(errors.AccessError))
			return
		}
		form.Paging.Count = 50
		users, err = user.SearchNear(&form.Paging)
	} else {
		users, err = models.Search(form.Nickname, &form.Paging)
	}

	list := make([]leaderboardResp, len(users))
	for i, _ := range users {
		list[i].Userid = users[i].Id
		list[i].Score = users[i].Score
		list[i].Level = users[i].Level
		list[i].Profile = users[i].Profile
		list[i].Nickname = users[i].Nickname
		list[i].Gender = users[i].Gender
		list[i].LastLog = users[i].LastLogin.Unix()
		list[i].Birth = users[i].Birth
		if users[i].Loc != nil {
			list[i].Location = *users[i].Loc
		}
	}

	respData := map[string]interface{}{
		"members_list":  list,
		"page_frist_id": form.Paging.First,
		"page_last_id":  form.Paging.Last,
	}
	writeResponse(r.RequestURI, w, respData, err)
}

type importFriendsForm struct {
	Type     string `json:"account_type"`
	Uid      string `json:"userid" binding:"required"`
	AppKey   string `json:"appkey"`
	AppToken string `json:"verfiycode" binding:"required"`
	Token    string `json:"access_token" binding:"required"`
}

func (this importFriendsForm) getTokenId() string {
	return this.Token
}

func importFriendsHandler(r *http.Request, w http.ResponseWriter, redis *models.RedisLogger, getT GetToken) {
	form := getT.(importFriendsForm)
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.AccessError))
		return
	}

	switch form.Type {
	case "weibo":
		log.Println("get weibo friends")
		friends, err := GetWeiboFriends(form.AppKey, form.Uid, form.AppToken)
		if err != nil {
			writeResponse(r.RequestURI, w, nil, errors.NewError(errors.DbError, err.Error()))
			return
		}
		log.Println("import weibo friends", len(friends))
		for _, friend := range friends {
			//log.Println(friend.Id, friend.ScreenName)
			u := &models.Account{
				Id:       strconv.FormatInt(int64(friend.Id), 10),
				Nickname: friend.ScreenName,
				Profile:  friend.Avatar,
				Role:     "weibo",
				Gender:   friend.Gender,
				Addr:     &models.Address{Desc: friend.Location},
			}
			if find, _ := u.Exists(); find {
				if u.RegTime.Unix() > 0 { // registered users only
					redis.ImportFriend(user.Id, u.Id)
				}
				continue
			}
			if err := u.Save(); err == nil {
				redis.SetWBImport(user.Id, u.Id)
			}
		}
	default:
	}
	writeResponse(r.RequestURI, w, map[string]interface{}{"ExpEffect": Awards{}}, nil)
}
