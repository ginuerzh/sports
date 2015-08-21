// account
package admin

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/go-martini/martini.v1"
	"io"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	//errs "errors"
)

var defaultCount = 50

type response struct {
	ReqPath  string      `json:"req_path"`
	RespData interface{} `json:"response_data"`
	Error    error       `json:"error"`
}

func BindAccountApi(m *martini.ClassicMartini) {
	m.Post("/admin/login", binding.Json(adminLoginForm{}), adminErrorHandler, adminLoginHandler)
	m.Options("/admin/login", optionsHandler)
	m.Post("/admin/logout", binding.Json(adminLogoutForm{}), adminErrorHandler, adminLogoutHandler)
	m.Options("/admin/logout", optionsHandler)
	m.Get("/admin/user/info", binding.Form(getUserInfoForm{}), adminErrorHandler, singleUserInfoHandler)
	m.Get("/admin/user/list", binding.Form(getUserListForm{}), adminErrorHandler, getUserListHandler)
	m.Get("/admin/user/search", binding.Form(getSearchListForm{}), adminErrorHandler, getSearchListHandler)
	m.Get("/admin/user/friendship", binding.Form(getUserFriendsForm{}), adminErrorHandler, getUserFriendsHandler)
	m.Post("/admin/user/ban", binding.Json(banUserForm{}), adminErrorHandler, banUserHandler)
	m.Options("/admin/user/ban", optionsHandler)
	m.Get("/admin/user/setadmin", binding.Form(setAdminForm{}), setAdminHandler)
	m.Get("/admin/user/auth/list", binding.Form(userAuthListForm{}), userAuthListHandler)
	m.Post("/admin/user/auth", binding.Json(userAuthForm{}), userAuthHandler)
	m.Options("/admin/user/auth", optionsHandler)
	//m.Post("/admin/user/update", updateUserInfoHandler)
	//m.Get("/admin/user/balance", binding.Form(userBalanceForm{}), userBalanceHandler)
	//m.Post("/admin/user/update", binding.Json(userInfoForm{}), adminErrorHandler, updateUserInfoHandler)
}

type setAdminForm struct {
	Id    string `form:"userid"`
	Admin bool   `form:"admin"`
	Token string `form:"access_token"`
}

func setAdminHandler(w http.ResponseWriter, form setAdminForm) {
	user := &models.Account{Id: form.Id}
	err := user.SetAdmin(form.Admin)
	if err != nil {
		writeResponse(w, err)
	}
	writeResponse(w, map[string]interface{}{"userid": user.Id, "admin": form.Admin})
}

// admin login parameter
type adminLoginForm struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func adminLoginHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form adminLoginForm) {
	user := &models.Account{}

	u4, err := uuid.NewV4()
	if err != nil {
		writeResponse(resp, err)
		return
	}
	token := u4.String() + "-" + strconv.FormatInt(time.Now().AddDate(0, 0, 30).Unix(), 10)
	var find bool
	//var err error

	h := md5.New()
	io.WriteString(h, form.Password)
	pwd := fmt.Sprintf("%x", h.Sum(nil))

	if find, err = user.FindByUserPass(strings.ToLower(form.UserName), pwd); !find {
		if err == nil {
			err = errors.NewError(errors.AuthError)
		}
	}

	if err != nil {
		writeResponse(resp, err)
		return
	}

	if !user.IsAdmin() {
		writeResponse(resp, errors.NewError(errors.AuthError, "未授权登录用户"))
		return
	}

	user.SetLastLogin(0, 0, time.Now())
	redis.SetOnlineUser(token, user.Id)
	//redis.LogLogin(user.Id)

	data := map[string]interface{}{
		"userid":       user.Id,
		"access_token": token,
	}
	writeResponse(resp, data)
}

type adminLogoutForm struct {
	Token string `json:"access_token" binding:"required"`
}

func adminLogoutHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form adminLogoutForm) {
	valid, err := checkToken(redis, form.Token)
	if valid {
		data := map[string]interface{}{}
		redis.DelOnlineUser(form.Token)
		writeResponse(resp, data)
	} else {
		writeResponse(resp, err)
	}
}

type getUserInfoForm struct {
	Userid   string `form:"userid"`
	NickName string `form:"nickname"`
	Token    string `form:"access_token" binding:"required"`
}

type equips struct {
	Shoes       []string `json:"shoes"`
	Electronics []string `json:"hardwares"`
	Softwares   []string `json:"softwares"`
}

type userInfoJsonStruct struct {
	Userid   string `json:"userid"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	About    string `json:"about"`
	Profile  string `json:"profile"`
	RegTime  int64  `json:"reg_time"`
	Height   int    `json:"height"`
	Weight   int    `json:"weight"`
	Birth    int64  `json:"birthday"`

	Sign        string `json:"sign"`
	Emotion     string `json:"emotion"`
	Profession  string `json:"profession"`
	Hobby       string `json:"fond"`
	Hometown    string `json:"hometown"`
	OftenAppear string `json:"oftenAppear"`

	Actor string `json:"actor"`
	Admin bool   `json:"admin"`

	Online     bool  `json:"online"`
	OnlineTime int64 `json:"onlinetime"`

	//Props *models.Props `json:"proper_info"`
	Physical int64 `json:"physique_value"`
	Literal  int64 `json:"literature_value"`
	Mental   int64 `json:"magic_value"`
	Wealth   int64 `json:"coin_value"`
	Score    int64 `json:"score"`
	Level    int64 `json:"level"`

	Addr string `json:"address"`
	//models.Location
	Lat float64 `json:"loc_latitude"`
	Lng float64 `json:"loc_longitude"`

	Gender          string `json:"gender"`
	Follows         int    `json:"follows_count"`
	Followers       int    `json:"followers_count"`
	Posts           int    `json:"articles_count"`
	FriendsCount    int    `json:"friends_count"`
	BlacklistsCount int    `json:"blacklist_count"`

	Photos []string `json:"photos"`
	//Equips models.Equip `json:"user_equipInfo"`
	Equip equips `json:"equips"`

	Wallet  string `json:"wallet"`
	LastLog int64  `json:"last_login_time"`
	BanTime int64  `json:"ban_time"`
	//BanTimeStr string `json:"ban_time_str"`
	//RegTimeStr string `json:"reg_time_str"`
	//LastLogStr string `json:"last_login_time_str"`
	BanStatus string `json:"ban_status"`

	Auth *models.UserAuth `json:"auth"`

	Stat *models.UserStat `json:"stat"`

	Action string `json:"action"`
}

func convertUser(user *models.Account, redis *models.RedisLogger) *userInfoJsonStruct {
	info := &userInfoJsonStruct{
		Userid:   user.Id,
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		Actor:    user.Actor,
		Admin:    user.Admin,
		About:    user.About,
		Profile:  user.Profile,
		RegTime:  user.RegTime.Unix(),
		//RegTimeStr: user.RegTime.Format("2006-01-02 15:04:05"),
		Sign:        user.Sign,
		Emotion:     user.Emotion,
		Profession:  user.Profession,
		Hobby:       user.Hobby,
		Hometown:    user.Hometown,
		OftenAppear: user.Oftenappear,

		Action: redis.GetUserAction(user.Id),

		Online:     redis.IsOnline(user.Id),
		OnlineTime: redis.GetOnlineTime(user.Id),

		Height: user.Height,
		Weight: user.Weight,
		Birth:  user.Birth,

		Lat:  user.Loc.Lat,
		Lng:  user.Loc.Lng,
		Addr: user.LocAddr,

		Physical: user.Props.Physical,
		Literal:  user.Props.Literal,
		Mental:   user.Props.Mental,
		Wealth:   redis.GetCoins(user.Id),
		Score:    user.Props.Score,
		Level:    user.Level(),

		Gender: user.Gender,
		Posts:  user.ArticleCount(),

		Photos: user.Photos,

		Wallet:  user.Wallet.Addr,
		LastLog: user.LastLogin.Unix(),

		Auth: user.Auth,

		Stat: user.Stat,
	}

	if len(info.Gender) == 0 {
		info.Gender = "male"
	}

	info.BanTime = user.TimeLimit
	if user.TimeLimit > 0 {
		if user.TimeLimit > time.Now().Unix() {
			info.BanStatus = "normal"
		} else {
			info.BanStatus = "lock"
		}
	} else {
		if user.TimeLimit == 0 {
			info.BanStatus = "normal"
		} else if user.TimeLimit == -1 {
			info.BanStatus = "ban"
		}
	}

	info.Follows, info.Followers, info.FriendsCount, info.BlacklistsCount = redis.FriendCount(user.Id)

	if user.Equips != nil {
		eq := *user.Equips
		info.Equip.Shoes = eq.Shoes
		info.Equip.Electronics = eq.Electronics
		info.Equip.Softwares = eq.Softwares
		//info.Equips = *user.Equips
	}

	/*
		if user.Addr != nil {
			info.Addr = user.Addr.String()
		}
		info.Lat = user.Loc.Lat
		info.Lng = user.Loc.Lng
	*/
	if info.Auth != nil {
		if info.Auth.IdCardTmp != nil {
			info.Auth.IdCard = info.Auth.IdCardTmp
		}
		if info.Auth.CertTmp != nil {
			info.Auth.Cert = info.Auth.CertTmp
		}
		if info.Auth.RecordTmp != nil {
			info.Auth.Record = info.Auth.RecordTmp
		}
	}

	return info
}

func singleUserInfoHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getUserInfoForm) {
	//log.Println("get a single user infomation")

	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(resp, err)
		return
	}

	user := &models.Account{}
	if find, err := user.FindByUserid(form.Userid); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError, "user '"+form.Userid+"' not exists")
		}
		writeResponse(resp, errors.NewError(errors.NotExistsError))
		return
	}

	info := convertUser(user, redis)

	writeResponse(resp, info)
}

type getUserListForm struct {
	//Count      int    `form:"count"`
	Sort string `form:"sort"`
	//NextCursor string `form:"next_cursor"`
	//PrevCursor string `form:"prev_cursor"`
	Token string `form:"access_token" binding:"required"`
	Count int    `form:"page_count"`
	Page  int    `form:"page_index"`
}

type userListJsonStruct struct {
	Users []*userInfoJsonStruct `json:"users"`
	//NextCursor  string               `json:"next_cursor"`
	//PrevCursor  string               `json:"prev_cursor"`
	Page        int `json:"page_index"`
	PageTotal   int `json:"page_total"`
	TotalNumber int `json:"total_number"`
}

func getUserListHandler(r *http.Request, w http.ResponseWriter, redis *models.RedisLogger, form getUserListForm) {
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}

	count := form.Count
	if count == 0 {
		count = defaultCount
	}

	total, users, _ := models.GetUserListBySort(form.Page*count, count, form.Sort)

	list := make([]*userInfoJsonStruct, len(users))
	for i, _ := range users {
		list[i] = convertUser(&users[i], redis)
	}

	totalPage := total / count
	if total%count != 0 {
		totalPage++
	}

	info := &userListJsonStruct{
		Users:       list,
		Page:        form.Page,
		PageTotal:   totalPage,
		TotalNumber: total,
	}

	writeResponse(w, info)
}

type getSearchListForm struct {
	Userid    string `form:"userid"`
	NickName  string `form:"nickname"`
	Role      string `form:"role"`
	Gender    string `form:"gender"`
	Age       string `form:"age"`
	BanStatus string `form:"ban_status"`
	KeyWord   string `form:"keyword"`
	Actor     string `form:"actor"`
	//	Count    int    `form:"count"`
	Sort string `form:"sort"`
	//	NextCursor string `form:"next_cursor"`
	//	PrevCursor string `form:"prev_cursor"`
	Count int    `form:"page_count"`
	Page  int    `form:"page_index"`
	Token string `form:"access_token" binding:"required"`
}

func getSearchListHandler(w http.ResponseWriter, redis *models.RedisLogger, form getSearchListForm) {
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}

	count := form.Count
	if count == 0 {
		count = defaultCount
	}

	total, users, _ := models.GetSearchListBySort(form.Userid, form.NickName, form.KeyWord,
		form.Gender, form.Age, form.BanStatus, form.Role, form.Actor, count*form.Page, count, form.Sort)

	list := make([]*userInfoJsonStruct, len(users))
	for i, _ := range users {
		list[i] = convertUser(&users[i], redis)
	}

	totalPage := total / count
	if total%count != 0 {
		totalPage++
	}

	info := &userListJsonStruct{
		Users:       list,
		Page:        form.Page,
		PageTotal:   totalPage,
		TotalNumber: total,
	}

	writeResponse(w, info)
}

type getUserFriendsForm struct {
	UserId string `form:"userid" binding:"required"`
	Type   string `form:"type"`
	//Count  int    `form:"count"`
	Sort  string `form:"sort"`
	Count int    `form:"page_count"`
	Page  int    `form:"page_index"`
	//	NextCursor string `form:"next_cursor"`
	//	PrevCursor string `form:"prev_cursor"`
	Token string `form:"access_token" binding:"required"`
}

func getUserFriendsHandler(w http.ResponseWriter,
	redis *models.RedisLogger, form getUserFriendsForm) {
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}

	count := form.Count
	if count == 0 {
		count = defaultCount
	}

	var friendType string
	switch form.Type {
	case "follows":
		friendType = models.RelFollowing
	case "followers":
		friendType = models.RelFollower
	case "blacklist":
		friendType = models.RelBlacklist
	case "friends":
		fallthrough
	default:
		friendType = models.RelFriend
	}
	userids := redis.Friends(friendType, form.UserId)
	total := len(userids)
	//count, users, err := models.GetFriendsListBySort(0, getCount, userids, form.Sort, form.PrevCursor, form.NextCursor)
	users, _ := models.GetFriendsListBySort(userids, count*form.Page, count, form.Sort)
	list := make([]*userInfoJsonStruct, len(users))
	for i, _ := range users {
		list[i] = convertUser(&users[i], redis)
	}

	totalPage := total / count
	if total%count != 0 {
		totalPage++
	}

	info := &userListJsonStruct{
		Users:       list,
		Page:        form.Page,
		PageTotal:   totalPage,
		TotalNumber: total,
	}

	writeResponse(w, info)
}

type banUserForm struct {
	Userid   string `json:"userid" binding:"required"`
	Duration int64  `json:"duration"`
	Token    string `json:"access_token" binding:"required"`
}

// This function bans user with a time value or forever by Duration.
func banUserHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form banUserForm) {
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(resp, err)
		return
	}

	user := &models.Account{}
	if find, err := user.FindByUserid(form.Userid); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError, "user '"+form.Userid+"' not exists")
			writeResponse(resp, err)
			return
		}
		writeResponse(resp, errors.NewError(errors.NotExistsError))
		return
	}

	// ws push
	event := &models.Event{
		Type: models.EventStatus,
		Data: models.EventData{
			Type: models.EventBan,
		},
	}

	if form.Duration == 0 {
		user.TimeLimit = 0
		event.Data.Type = models.EventUnban
	} else if form.Duration < 0 {
		user.TimeLimit = -1
		event.Data.Type = models.EventLock
	} else {
		user.TimeLimit = time.Now().Unix() + form.Duration
		event.Data.Type = models.EventBan
	}

	err := user.UpdateBanTime(user.TimeLimit)
	if err != nil {
		writeResponse(resp, err)
		return
	}

	redis.PubMsg(event.Type, user.Id, event.Bytes())

	respData := map[string]interface{}{
		"ban": form.Duration,
	}
	writeResponse(resp, respData)
}

/*
type userInfoForm struct {
	Userid   string `json:"userid" binding:"required"`
	Token    string `json:"access_token" binding:"required"`
	Nickname string `json:"nickname"`

	Shoes       []string `json:"equips_shoes"`
	Electronics []string `json:"equips_hardwares"`
	Softwares   []string `json:"equips_softwares"`

	Phone   string `json:"phone"`
	Type    string `json:"role"`
	About   string `json:"about"`
	Profile string `json:"profile"`
	Hobby   string `json:"hobby"`
	Height  int    `json:"height"`
	Weight  int    `json:"weight"`
	Birth   int64  `json:"birthday"`

	//Props *models.Props `json:"proper_info"`
	Physical int64 `json:"physique_value"`
	Literal  int64 `json:"literature_value"`
	Mental   int64 `json:"magic_value"`
	Wealth   int64 `json:"coin_value"`

	Addr string `json:"address"`
	//models.Location
	Lat float64 `json:"loc_latitude"`
	Lng float64 `json:"loc_longitude"`

	Gender string   `json:"gender"`
	Photos []string `json:"photos"`
}
*/

func updateUserInfoToDB(r *models.RedisLogger, m map[string]interface{}, u *models.Account) error {
	ss := []string{"userid", "access_token", "nickname", "equips_shoes", "equips_hardwares", "equips_softwares",
		"phone", "role", "about", "profile", "hobby", "height", "weight", "birthday",
		"address", "loc_latitude", "loc_longitude", "gender", "photos"}
	changeFields := map[string]interface{}{}

	for _, vv := range ss {

		if value, exists := m[vv]; exists {
			log.Println("value is :", value)
			switch vv {
			case "nickname":
				changeFields["nickname"] = value
			case "equips_shoes":
				changeFields["equips.shoes"] = value
			case "equips_hardwares":
				changeFields["equips.electronics"] = value
			case "equips_softwares":
				changeFields["equips.softwares"] = value
			case "phone":
				changeFields["phone"] = value
			case "role":
				changeFields["role"] = value
			case "about":
				changeFields["about"] = value
			case "profile":
				changeFields["profile"] = value
			case "hobby":
				changeFields["hobby"] = value
			case "height":
				changeFields["height"] = value
			case "weight":
				changeFields["weight"] = value
			case "birthday":
				changeFields["birth"] = value
			case "address":
				v := reflect.ValueOf(value)
				var Addr = new(models.Address)
				Addr.Country = ""
				Addr.Province = ""
				Addr.City = ""
				Addr.Area = ""
				Addr.Desc = v.String()
				changeFields["addr"] = Addr
			case "loc_latitude":
				changeFields["loc.latitude"] = value
			case "loc_longitude":
				changeFields["loc.longitude"] = value
			case "gender":
				changeFields["gender"] = value
			case "photos":
				changeFields["photos"] = value
			}
		}
	}

	change := bson.M{
		"$set": changeFields,
	}
	u.UpdateInfo(change)
	return nil
}

func updateUserInfo(r *models.RedisLogger, req *http.Request, u *models.Account) (err error) {
	if req.Body != nil {
		defer req.Body.Close()

		dec := json.NewDecoder(req.Body)
		for {
			var m map[string]interface{}
			err = dec.Decode(&m)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			token, exist := m["access_token"]
			if !exist {
				err = errors.NewError(errors.AccessError)
				return
			} else {
				v := reflect.ValueOf(token)
				valid, errT := checkToken(r, v.String())
				if !valid {
					err = errT
					return
				}
			}

			if key, exists := m["userid"]; exists {
				var find bool
				v := reflect.ValueOf(key)
				userid := v.String()
				if find, err = u.FindByUserid(userid); !find {
					if err == nil {
						err = errors.NewError(errors.NotExistsError, "user '"+userid+"' not exists")
					}
					return
				}

				log.Println("key is :", key, " uid is :", u.Id)
				if u.Id != key {
					err = errors.NewError(errors.NotExistsError, "user '"+userid+"' not exists")
					return
				}
				err = updateUserInfoToDB(r, m, u)
			}
		}
	}
	return
}

// This function update user info.
func updateUserInfoHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger) {
	user := &models.Account{}
	err := updateUserInfo(redis, request, user)
	if err != nil {
		writeResponse(resp, err)
		return
	}
	data := map[string]interface{}{}
	writeResponse(resp, data)
}

type userBalanceForm struct {
	Id string `form:"id"`
}

func userBalanceHandler(r *http.Request, w http.ResponseWriter, form userBalanceForm) {
	user := &models.Account{}
	user.FindByUserid(form.Id)
	balance, _ := getBalance(user.Wallet.Addrs)
	writeResponse(w, map[string]interface{}{"balance": balance})
}

type balance struct {
	Address     string `json:"address"`
	Confirmed   int64  `json:"confirmed"`
	Unconfirmed int64  `json:"unconfirmed"`
}

type balanceAddrs struct {
	Addrs []balance `json:"addresses"`
}

func getBalance(addrs []string) (b *balanceAddrs, err error) {
	if len(addrs) == 0 {
		return
	}
	resp, err := http.Get("http://localhost:8087" + "/multiaddr?addr=" + strings.Join(addrs, "|"))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	ba := new(balanceAddrs)
	if err = decodeJson(resp.Body, ba); err != nil {
		log.Println(err)
		return
	}

	b = ba

	return
}

type userAuthListForm struct {
	Sort  string `form:"sort"`
	Count int    `form:"page_count"`
	Page  int    `form:"page_index"`
	Token string `form:"access_token" binding:"required"`
}

func userAuthListHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, form userAuthListForm) {
	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}

	count := form.Count
	if count == 0 {
		count = defaultCount
	}

	total, users, _ := models.GetAuthUserList(form.Page, count)

	list := make([]*userInfoJsonStruct, len(users))
	for i, _ := range users {
		list[i] = convertUser(&users[i], redis)
	}

	totalPage := total / count
	if total%count != 0 {
		totalPage++
	}

	info := &userListJsonStruct{
		Users:       list,
		Page:        form.Page,
		PageTotal:   totalPage,
		TotalNumber: total,
	}
	writeResponse(w, info)
}

type userAuthForm struct {
	Userid string `json:"userid" binding:"required"`
	Type   string `json:"auth_type" binding:"required"`
	Status string `json:"auth_status" binding:"required"`
	Review string `json:"auth_review"`
	Token  string `json:"access_token" binding:"required"`
}

func userAuthHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, form userAuthForm) {

	if ok, err := checkToken(redis, form.Token); !ok {
		writeResponse(w, err)
		return
	}

	user := &models.Account{}
	user.FindByUserid(form.Userid)
	err := user.SetAuth(form.Type, form.Status, form.Review)
	writeResponse(w, err)
}
