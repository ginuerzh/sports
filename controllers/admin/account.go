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
	//"strconv"
	"reflect"
	"strings"
	"time"
	//errs "errors"
)

var defaultCount = 20

type response struct {
	ReqPath  string      `json:"req_path"`
	RespData interface{} `json:"response_data"`
	Error    error       `json:"error"`
}

func BindAccountApi(m *martini.ClassicMartini) {
	m.Post("/admin/login", binding.Json(adminLoginForm{}), adminErrorHandler, adminLoginHandler)
	m.Post("/admin/logout", binding.Json(adminLogoutForm{}), adminErrorHandler, adminLogoutHandler)
	m.Get("/admin/user/info", binding.Form(getUserInfoForm{}), adminErrorHandler, singleUserInfoHandler)
	m.Get("/admin/user/list", binding.Form(getUserListForm{}), adminErrorHandler, getUserListHandler)
	m.Get("/admin/user/search", binding.Form(getSearchListForm{}), adminErrorHandler, getSearchListHandler)
	m.Get("/admin/user/friendship", binding.Form(getUserFriendsForm{}), adminErrorHandler, getUserFriendsHandler)
	m.Post("/admin/user/ban", binding.Json(banUserForm{}), adminErrorHandler, banUserHandler)
	m.Post("/admin/user/update", updateUserInfoHandler)
	//m.Post("/admin/user/update", binding.Json(userInfoForm{}), adminErrorHandler, updateUserInfoHandler)
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
	token := u4.String()
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

	redis.SetOnlineUser(token, user, true)
	redis.LogLogin(user.Id)

	data := map[string]interface{}{
		"userid":       user.Id,
		"access_token": token,
	}
	writeResponse(resp, data)
}

type adminLogoutForm struct {
	Token string `json:"access_token" binding:"required"`
}

func checkToken(r *models.RedisLogger, t string) (valid bool, err error) {
	user := r.OnlineUser(t)
	if user == nil {
		err = errors.NewError(errors.AccessError)
		valid = false
		return
	}
	valid = true
	return
}

func adminLogoutHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form adminLogoutForm) {
	valid, err := checkToken(redis, form.Token)
	if valid {
		redis.DelOnlineUser(form.Token)
		writeResponse(resp, nil)
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

	Phone   string `json:"phone"`
	Type    string `json:"role"`
	About   string `json:"about"`
	Profile string `json:"profile"`
	RegTime int64  `json:"reg_time"`
	Hobby   string `json:"hobby"`
	Height  int    `json:"height"`
	Weight  int    `json:"weight"`
	Birth   int64  `json:"birthday"`

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
}

func singleUserInfoHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getUserInfoForm) {
	log.Println("get a single user infomation")

	valid, errT := checkToken(redis, form.Token)
	if !valid {
		writeResponse(resp, errT)
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

	info := &userInfoJsonStruct{
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

		Physical: redis.UserProps(user.Id).Physical,
		Literal:  redis.UserProps(user.Id).Literal,
		Mental:   redis.UserProps(user.Id).Mental,
		Wealth:   redis.UserProps(user.Id).Wealth,
		Score:    redis.UserProps(user.Id).Score,
		Level:    redis.UserProps(user.Id).Level,

		Gender: user.Gender,
		Posts:  user.ArticleCount(),

		Photos: user.Photos,

		Wallet:  user.Wallet.Addr,
		LastLog: user.LastLogin.Unix(),
	}

	info.Follows, info.Followers, info.FriendsCount, info.BlacklistsCount = redis.FriendCount(user.Id)

	if user.Equips != nil {
		eq := *user.Equips
		info.Equip.Shoes = eq.Shoes
		info.Equip.Electronics = eq.Electronics
		info.Equip.Softwares = eq.Softwares
		//info.Equips = *user.Equips
	}

	if user.Addr != nil {
		info.Addr = user.Addr.String()
	}
	if user.Loc != nil {
		loction := *user.Loc
		info.Lat = loction.Lat
		info.Lng = loction.Lng
	}

	writeResponse(resp, info)
}

type getUserListForm struct {
	Count      int    `form:"count"`
	Sort       string `form:"sort"`
	NextCursor string `form:"next_cursor"`
	PrevCursor string `form:"prev_cursor"`
	Token      string `form:"access_token" binding:"required"`
}

type userListJsonStruct struct {
	Users       []userInfoJsonStruct `json:"users"`
	NextCursor  string               `json:"next_cursor"`
	PrevCursor  string               `json:"prev_cursor"`
	TotalNumber int                  `json:"total_number"`
}

func getUserListHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getUserListForm) {
	valid, errT := checkToken(redis, form.Token)
	if !valid {
		writeResponse(resp, errT)
		return
	}

	getCount := form.Count
	if getCount == 0 {
		getCount = defaultCount
	}
	log.Println("getCount is :", getCount, "sort is :", form.Sort, "pc is :", form.PrevCursor, "nc is :", form.NextCursor)
	count, users, err := models.GetUserListBySort(0, getCount, form.Sort, form.PrevCursor, form.NextCursor)
	if err != nil {
		writeResponse(resp, err)
		return
	}
	countvalid := len(users)
	log.Println("countvalid is :", countvalid)

	if countvalid == 0 {
		writeResponse(resp, err)
		return
	}

	list := make([]userInfoJsonStruct, countvalid)
	for i, user := range users {
		list[i].Userid = user.Id
		list[i].Nickname = user.Nickname
		list[i].Phone = user.Phone
		list[i].Type = user.Role
		list[i].About = user.About
		list[i].Profile = user.Profile
		list[i].RegTime = user.RegTime.Unix()
		list[i].Hobby = user.Hobby
		list[i].Height = user.Height
		list[i].Weight = user.Weight
		list[i].Birth = user.Birth
		list[i].Gender = user.Gender
		list[i].Posts = user.ArticleCount()
		list[i].Photos = user.Photos
		list[i].Wallet = user.Wallet.Addr
		list[i].LastLog = user.LastLogin.Unix()
		list[i].Follows, list[i].Followers, list[i].FriendsCount, list[i].BlacklistsCount = redis.FriendCount(user.Id)
		pups := redis.UserProps(user.Id)
		if pups != nil {
			ups := *pups
			list[i].Physical = ups.Physical
			list[i].Literal = ups.Literal
			list[i].Mental = ups.Mental
			list[i].Wealth = ups.Wealth
			list[i].Score = ups.Score
			list[i].Level = ups.Level
		}

		if user.Equips != nil {
			eq := *user.Equips
			list[i].Equip.Shoes = eq.Shoes
			list[i].Equip.Electronics = eq.Electronics
			list[i].Equip.Softwares = eq.Softwares
			//info.Equips = *user.Equips
		}

		if user.Addr != nil {
			list[i].Addr = user.Addr.String()
		}
		if user.Loc != nil {
			loction := *user.Loc
			list[i].Lat = loction.Lat
			list[i].Lng = loction.Lng
		}
	}
	/*
		var pc, nc string
		switch form.Sort {
		case "logintime":
			pc = strconv.FormatInt(list[0].LastLog, 10)
			nc = strconv.FormatInt(list[count-1].LastLog, 10)
		case "userid":
			pc = list[0].Userid
			nc = list[count-1].Userid
		case "nickname":
			pc = list[0].Nickname
			nc = list[count-1].Nickname
		case "score":
			pc = strconv.FormatInt(list[0].Score, 10)
			nc = strconv.FormatInt(list[count-1].Score, 10)
		case "regtime":
			fallthrough
		default:
			pc = strconv.FormatInt(list[0].RegTime, 10)
			nc = strconv.FormatInt(list[count-1].RegTime, 10)
		}
	*/
	info := &userListJsonStruct{
		Users:       list,
		NextCursor:  list[countvalid-1].Userid,
		PrevCursor:  list[0].Userid,
		TotalNumber: count,
	}

	writeResponse(resp, info)
}

type getSearchListForm struct {
	Userid     string `form:"userid"`
	NickName   string `form:"nickname"`
	Count      int    `form:"count"`
	Sort       string `form:"sort"`
	NextCursor string `form:"next_cursor"`
	PrevCursor string `form:"prev_cursor"`
	Token      string `form:"access_token" binding:"required"`
}

func getSearchListHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getSearchListForm) {
	valid, errT := checkToken(redis, form.Token)
	if !valid {
		writeResponse(resp, errT)
		return
	}

	getCount := form.Count
	if getCount == 0 {
		getCount = defaultCount
	}
	log.Println("getCount is :", getCount, "sort is :", form.Sort, "pc is :", form.PrevCursor, "nc is :", form.NextCursor)
	count, users, err := models.GetSearchListBySort(form.Userid, form.NickName, 0, getCount, form.Sort, form.PrevCursor, form.NextCursor)
	if err != nil {
		writeResponse(resp, err)
		return
	}
	countvalid := len(users)
	log.Println("countvalid is :", countvalid)

	if countvalid == 0 {
		writeResponse(resp, err)
		return
	}

	list := make([]userInfoJsonStruct, countvalid)
	for i, user := range users {
		list[i].Userid = user.Id
		list[i].Nickname = user.Nickname
		list[i].Phone = user.Phone
		list[i].Type = user.Role
		list[i].About = user.About
		list[i].Profile = user.Profile
		list[i].RegTime = user.RegTime.Unix()
		list[i].Hobby = user.Hobby
		list[i].Height = user.Height
		list[i].Weight = user.Weight
		list[i].Birth = user.Birth
		list[i].Gender = user.Gender
		list[i].Posts = user.ArticleCount()
		list[i].Photos = user.Photos
		list[i].Wallet = user.Wallet.Addr
		list[i].LastLog = user.LastLogin.Unix()
		list[i].Follows, list[i].Followers, list[i].FriendsCount, list[i].BlacklistsCount = redis.FriendCount(user.Id)
		pups := redis.UserProps(user.Id)
		if pups != nil {
			ups := *pups
			list[i].Physical = ups.Physical
			list[i].Literal = ups.Literal
			list[i].Mental = ups.Mental
			list[i].Wealth = ups.Wealth
			list[i].Score = ups.Score
			list[i].Level = ups.Level
		}

		if user.Equips != nil {
			eq := *user.Equips
			list[i].Equip.Shoes = eq.Shoes
			list[i].Equip.Electronics = eq.Electronics
			list[i].Equip.Softwares = eq.Softwares
			//info.Equips = *user.Equips
		}

		if user.Addr != nil {
			list[i].Addr = user.Addr.String()
		}
		if user.Loc != nil {
			loction := *user.Loc
			list[i].Lat = loction.Lat
			list[i].Lng = loction.Lng
		}
	}

	info := &userListJsonStruct{
		Users:       list,
		NextCursor:  list[countvalid-1].Userid,
		PrevCursor:  list[0].Userid,
		TotalNumber: count,
	}

	writeResponse(resp, info)
}

type getUserFriendsForm struct {
	UserId     string `form:"userid" binding:"required"`
	Type       string `form:"type"`
	Count      int    `form:"count"`
	Sort       string `form:"sort"`
	NextCursor string `form:"next_cursor"`
	PrevCursor string `form:"prev_cursor"`
	Token      string `form:"access_token" binding:"required"`
}

func getUserFriendsHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getUserFriendsForm) {
	valid, errT := checkToken(redis, form.Token)
	if !valid {
		writeResponse(resp, errT)
		return
	}

	getCount := form.Count
	if getCount == 0 {
		getCount = defaultCount
	}

	log.Println("getCount is ", getCount)
	var friendType string
	switch form.Type {
	case "follows":
		friendType = models.RelFollower
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
	if getCount > len(userids) {
		log.Println("userid length is littler than getCount")
		getCount = len(userids)
	}

	if getCount == 0 {
		writeResponse(resp, errors.NewError(errors.NotExistsError))
		return
	}

	count, users, err := models.GetFriendsListBySort(0, getCount, userids, form.Sort, form.PrevCursor, form.NextCursor)
	if err != nil {
		writeResponse(resp, err)
		return
	}
	countvalid := len(users)
	log.Println("countvalid is :", countvalid)
	if countvalid == 0 {
		writeResponse(resp, errors.NewError(errors.DbError))
		return
	}

	list := make([]userInfoJsonStruct, countvalid)
	for i, user := range users {
		list[i].Userid = user.Id
		list[i].Nickname = user.Nickname
		list[i].Phone = user.Phone
		list[i].Type = user.Role
		list[i].About = user.About
		list[i].Profile = user.Profile
		list[i].RegTime = user.RegTime.Unix()
		list[i].Hobby = user.Hobby
		list[i].Height = user.Height
		list[i].Weight = user.Weight
		list[i].Birth = user.Birth
		list[i].Gender = user.Gender
		list[i].Posts = user.ArticleCount()
		list[i].Photos = user.Photos
		list[i].Wallet = user.Wallet.Addr
		list[i].LastLog = user.LastLogin.Unix()
		list[i].Follows, list[i].Followers, list[i].FriendsCount, list[i].BlacklistsCount = redis.FriendCount(user.Id)
		pps := *redis.UserProps(user.Id)
		list[i].Physical = pps.Physical
		list[i].Literal = pps.Literal
		list[i].Mental = pps.Mental
		list[i].Wealth = pps.Wealth
		list[i].Score = pps.Score
		list[i].Level = pps.Level

		if user.Equips != nil {
			eq := *user.Equips
			list[i].Equip.Shoes = eq.Shoes
			list[i].Equip.Electronics = eq.Electronics
			list[i].Equip.Softwares = eq.Softwares
			//info.Equips = *user.Equips
		}

		if user.Addr != nil {
			list[i].Addr = user.Addr.String()
		}
		if user.Loc != nil {
			loction := *user.Loc
			list[i].Lat = loction.Lat
			list[i].Lng = loction.Lng
		}
	}
	/*
		var pc, nc string
		switch form.Sort {
		case "logintime":
			pc = strconv.FormatInt(list[0].LastLog, 10)
			nc = strconv.FormatInt(list[count-1].LastLog, 10)
		case "userid":
			pc = list[0].Userid
			nc = list[count-1].Userid
		case "nickname":
			pc = list[0].Nickname
			nc = list[count-1].Nickname
		case "score":
			pc = strconv.FormatInt(list[0].Score, 10)
			nc = strconv.FormatInt(list[count-1].Score, 10)
		case "regtime":
			fallthrough
		default:
			pc = strconv.FormatInt(list[0].RegTime, 10)
			nc = strconv.FormatInt(list[count-1].RegTime, 10)
		}
	*/
	info := &userListJsonStruct{
		Users:       list,
		NextCursor:  list[countvalid-1].Userid,
		PrevCursor:  list[0].Userid,
		TotalNumber: count,
	}

	writeResponse(resp, info)
}

type banUserForm struct {
	Userid   string `json:"userid" binding:"required"`
	Duration int64  `json:"duration"`
	Token    string `json:"access_token" binding:"required"`
}

// This function bans user with a time value or forever by Duration.
func banUserHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form banUserForm) {
	valid, errT := checkToken(redis, form.Token)
	if !valid {
		writeResponse(resp, errT)
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

	if form.Duration == 0 {
		user.TimeLimit = 0
	} else if form.Duration == -1 {
		user.TimeLimit = -1
	} else if form.Duration > 0 {
		user.TimeLimit = time.Now().Unix() + form.Duration
	} else {
		writeResponse(resp, errors.NewError(errors.NotFoundError))
		return
	}

	err := user.Update()
	if err != nil {
		writeResponse(resp, err)
		return
	}
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
		"phone", "role", "about", "profile", "hobby", "height", "weight", "birthday", "physique_value", "literature_value",
		"magic_value", "coin_value", "address", "loc_latitude", "loc_longitude", "gender", "photos"}
	changeFields := map[string]interface{}{}
	propChanged := false
	var prop1, prop2 *models.Props
	for _, vv := range ss {
		log.Println("vv is :", vv)
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
			case "physique_value":
				v, _ := value.(float64)
				//v := reflect.ValueOf(value).Int()
				if !propChanged {
					propChanged = true
					prop1 = r.UserProps(u.Id)
				}
				if int64(v)-prop1.Physical != 0 {
					if prop2 == nil {
						prop2 = new(models.Props)
					}
					prop2.Physical = int64(v) - prop1.Physical
				}
			case "literature_value":
				v, _ := value.(float64)
				//v := reflect.ValueOf(value).Int()
				if !propChanged {
					propChanged = true
					prop1 = r.UserProps(u.Id)
				}
				if int64(v)-prop1.Literal != 0 {
					if prop2 == nil {
						prop2 = new(models.Props)
					}
					prop2.Literal = int64(v) - prop1.Literal
				}
			case "magic_value":
				v, _ := value.(float64)
				//v := reflect.ValueOf(value).Int()
				if !propChanged {
					propChanged = true
					prop1 = r.UserProps(u.Id)
				}
				if int64(v)-prop1.Mental != 0 {
					if prop2 == nil {
						prop2 = new(models.Props)
					}
					prop2.Mental = int64(v) - prop1.Mental
				}
			case "coin_value":
				/*
					v, _ := value.(float64)
							//v := reflect.ValueOf(value).Int()
						if !propChanged {
							propChanged = true
							prop1 = r.UserProps(u.Id)
						}
							if int64(v)-prop1.Wealth != 0 {
								if prop2 == nil {
									prop2 = new(models.Props)
								}
								prop2.Wealth = int64(v) - prop1.Wealth
							}
				*/
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
	if prop2 != nil {
		_, err := r.AddProps(u.Id, prop2)
		return err
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
	writeResponse(resp, nil)
}

/*
func GetFriendsListBySort(skip, limit int, ids []string, sortOrder, preCursor, nextCursor string) (total int, users []Account, err error) {
	user := &Account{}
	var query bson.M
	var sortby string
	var uids []string

	pc, _ := strconv.Atoi(preCursor)
	nc, _ := strconv.Atoi(nextCursor)

	switch sortOrder {
	case "logintime":
		if len(nextCursor) > 0 {
			user.findOne(bson.M{"lastlogin": nc})
			for i := 0; i < len(ids); i++ {
				if ids[i] != user.Id {
					uids = append(uids, ids[i])
				}
			}
			query = bson.M{
				"lastlogin": bson.M{
					"$lte": user.LastLogin,
				},
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "-lastlogin"
		} else if len(preCursor) > 0 {
			user.findOne(bson.M{"lastlogin": pc})
			for i := 0; i < len(ids); i++ {
				if ids[i] != user.Id {
					uids = append(uids, ids[i])
				}
			}
			query = bson.M{
				"lastlogin": bson.M{
					"$gte": user.LastLogin,
				},
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "lastlogin"
		} else {
			user.LastLogin = time.Now()
			query = bson.M{
				"lastlogin": bson.M{
					"$lte": user.LastLogin,
				},
				"_id": bson.M{
					"$in": ids,
				},
			}
			sortby = "-lastlogin"
		}

	case "userid":
		if len(nextCursor) > 0 {
			user.findOne(bson.M{"_id": nextCursor})
			for i := 0; i < len(ids); i++ {
				if ids[i] != user.Id {
					uids = append(uids, ids[i])
				}
			}
			query = bson.M{
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "_id"
		} else if len(preCursor) > 0 {
			user.findOne(bson.M{"_id": preCursor})
			for i := 0; i < len(ids); i++ {
				if ids[i] != user.Id {
					uids = append(uids, ids[i])
				}
			}
			query = bson.M{
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "-_id"
		} else {
			user.Id = ""
			query = bson.M{
				"_id": bson.M{
					"$in": ids,
				},
			}
			sortby = "_id"
		}

	case "nickname":
		if len(nextCursor) > 0 {
			user.findOne(bson.M{"nickname": nextCursor})
			for i := 0; i < len(ids); i++ {
				if ids[i] != user.Id {
					uids = append(uids, ids[i])
				}
			}
			query = bson.M{
				"nickname": bson.M{
					"$gte": user.Nickname,
				},
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "nickname"
		} else if len(preCursor) > 0 {
			user.findOne(bson.M{"nickname": preCursor})
			for i := 0; i < len(ids); i++ {
				if ids[i] != user.Id {
					uids = append(uids, ids[i])
				}
			}
			query = bson.M{
				"nickname": bson.M{
					"$lte": user.Nickname,
				},
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "-nickname"
		} else {
			user.Nickname = ""
			query = bson.M{
				"nickname": bson.M{
					"$gt": user.Nickname,
				},
				"_id": bson.M{
					"$in": ids,
				},
			}
			sortby = "nickname"
		}

	case "score":
		if len(nextCursor) > 0 {
			user.findOne(bson.M{"score": nc})
			for i := 0; i < len(ids); i++ {
				if ids[i] != user.Id {
					uids = append(uids, ids[i])
				}
			}
			query = bson.M{
				"score": bson.M{
					"$lte": user.Score,
				},
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "-score"
		} else if len(preCursor) > 0 {
			user.findOne(bson.M{"score": pc})
			for i := 0; i < len(ids); i++ {
				if ids[i] != user.Id {
					uids = append(uids, ids[i])
				}
			}
			query = bson.M{
				"score": bson.M{
					"$gte": user.Score,
				},
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "score"
		} else {
			user.Score = 0
			query = bson.M{
				"score": bson.M{
					"$lte": user.Score,
				},
				"_id": bson.M{
					"$in": ids,
				},
			}
			sortby = "-score"
		}

	case "regtime":
		log.Println("regtime")
		fallthrough
	default:
		log.Println("default")
		if len(nextCursor) > 0 {
			user.findOne(bson.M{"reg_time": nc})
			for i := 0; i < len(ids); i++ {
				if ids[i] != user.Id {
					uids = append(uids, ids[i])
				}
			}
			query = bson.M{
				"reg_time": bson.M{
					"$lte": user.RegTime,
				},
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "-reg_time"
		} else if len(preCursor) > 0 {
			user.findOne(bson.M{"reg_time": pc})
			for i := 0; i < len(ids); i++ {
				if ids[i] != user.Id {
					uids = append(uids, ids[i])
				}
			}
			query = bson.M{
				"reg_time": bson.M{
					"$gte": user.RegTime,
				},
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "reg_time"
		} else {
			user.RegTime = time.Now()
			query = bson.M{
				"reg_time": bson.M{
					"$lte": user.RegTime,
				},
				"_id": bson.M{
					"$in": ids,
				},
			}
			sortby = "-reg_time"
		}
	}

	if err := search(accountColl, query, nil, skip, limit, []string{sortby}, &total, &users); err != nil {
		return 0, nil, errors.NewError(errors.DbError, err.Error())
	}

	return
}
*/
