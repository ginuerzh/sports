// account
package admin

import (
	//"encoding/json"
	"crypto/md5"
	"fmt"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/go-martini/martini.v1"
	"io"
	"log"
	"net/http"
	//"strconv"
	"strings"
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
	m.Get("/admin/user/friendship", binding.Form(getUserFriendsForm{}), adminErrorHandler, getUserFriendsHandler)
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
		writeResponse(resp, nil)
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
		writeResponse(resp, nil)
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

func adminLogoutHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form adminLogoutForm) {
	redis.DelOnlineUser(form.Token)
	writeResponse(resp, nil)
}

type getUserInfoForm struct {
	Userid   string `form:"userid"`
	NickName string `form:"nickname"`
	Token    string `form:"access_token" binding:"required"`
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
	Shoes       []string `json:"shoes"`
	Electronics []string `json:"hardwares"`
	Softwares   []string `json:"softwares"`

	Wallet  string `json:"wallet"`
	LastLog int64  `json:"last_login_time"`
}

func singleUserInfoHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getUserInfoForm) {
	log.Println("get a single user infomation")
	user := &models.Account{}
	if find, err := user.FindByUserid(form.Userid); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError, "user '"+form.Userid+"' not exists")
		}
		writeResponse(resp, nil)
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
		info.Shoes = eq.Shoes
		info.Electronics = eq.Electronics
		info.Softwares = eq.Softwares
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
	getCount := form.Count
	if getCount == 0 {
		getCount = defaultCount
	}
	log.Println("getCount is :", getCount, "sort is :", form.Sort, "pc is :", form.PrevCursor, "nc is :", form.NextCursor)
	count, users, err := models.GetUserListBySort(0, getCount, form.Sort, form.PrevCursor, form.NextCursor)
	if err != nil {
		writeResponse(resp, nil)
		return
	}
	log.Println("count is :", count)

	if count == 0 {
		writeResponse(resp, nil)
		return
	}

	list := make([]userInfoJsonStruct, count)
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
			list[i].Shoes = eq.Shoes
			list[i].Electronics = eq.Electronics
			list[i].Softwares = eq.Softwares
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
		NextCursor:  list[0].Userid,
		PrevCursor:  list[count-1].Userid,
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
		writeResponse(resp, nil)
		return
	}

	count, users, err := models.GetFriendsListBySort(0, getCount, userids, form.Sort, form.PrevCursor, form.NextCursor)
	if err != nil {
		writeResponse(resp, nil)
		return
	}
	log.Println("count is :", count)
	if count == 0 {
		writeResponse(resp, nil)
		return
	}

	list := make([]userInfoJsonStruct, count)
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
			list[i].Shoes = eq.Shoes
			list[i].Electronics = eq.Electronics
			list[i].Softwares = eq.Softwares
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
		NextCursor:  list[count-1].Userid,
		PrevCursor:  list[0].Userid,
		TotalNumber: count,
	}

	writeResponse(resp, info)
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
