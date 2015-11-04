// account
package models

import (
	//"encoding/json"
	"github.com/ginuerzh/sports/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"labix.org/v2/mgo/txn"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	ActorAdmin = "admin"
	ActorCoach = "coach"
)

const (
	AccountEmail = "email"
	AccountPhone = "phone"
	AccountWeibo = "weibo"
)

const (
	AuthIdCard = "idcard"
	AuthCert   = "cert"
	AuthRecord = "record"

	AuthUnverified = "unverified"
	AuthVerifying  = "verifying"
	AuthVerified   = "verified"
	AuthRefused    = "refused"
)

const (
	StatOnlineTime      = "onlinetime"
	StatRecords         = "records"
	StatArticles        = "articles"
	StatComments        = "comments"
	StatPosts           = "posts"
	StatGameTime        = "gametime"
	StatLastArticleTime = "lastarticletime"
	StatLastCommentTime = "lastcommenttime"
	StatLastThumbTime   = "lastthumbtime"
	StatLastGameTime    = "lastgametime"
)

func init() {
	//dur, _ = time.ParseDuration("-30h") // auto logout after 15 minutes since last access

}

type UserInfo struct {
	Height   int    `json:"height"`
	Weight   int    `json:"weight"`
	Birth    int64  `json:"birthday"`
	Actor    string `json:"actor"`
	Gender   string `json:"sex_type"`
	Phone    string `json:"phone_number"`
	About    string `json:"about"`
	Nickname string `json:"nikename"`

	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Area     string `json:"area"`
	LocDesc  string `json:"location_desc"`

	Sign        string `json:"sign"`
	Emotion     string `json:"emotion"`
	Profession  string `json:"profession"`
	Hobby       string `json:"fond"`
	Hometown    string `json:"hometown"`
	OftenAppear string `json:"oftenAppear"`
	CoverImage  string `json:"cover_image"`
}

// for mongodb json -> bson
type SetInfo struct {
	Phone    string   `json:"phone,omitempty"`
	Nickname string   `json:"nickname,omitempty"`
	Height   int      `json:"height,omitempty"`
	Weight   int      `json:"weight,omitempty"`
	Birth    int64    `json:"birth,omitempty"`
	Gender   string   `json:"gender,omitempty"`
	Address  *Address `json:"addr,omitempty"`

	Sign        string `json:"sign,omitempty"`
	Emotion     string `json:"emotion,omitempty"`
	Profession  string `json:"profession,omitempty"`
	Hobby       string `json:"hobby,omitempty"`
	Hometown    string `json:"hometown,omitempty"`
	OftenAppear string `json:"oftenappear,omitempty"`
	CoverImage  string `json:"coverimage,omitempty"`

	Setinfo    bool `json:"setinfo,omitempty"`
	SetinfoAll bool `json:"setinfoall,omitempty"`
}

type DbWallet struct {
	Id    string   `bson:"wallet_id"`
	Key   string   `bson:"shared_key"`
	Addr  string   `bson:"addr"`
	Addrs []string `bson:"addrs"`
}

type Equip struct {
	Shoes       []string `json:"run_shoe"`
	Electronics []string `json:"ele_product"`
	Softwares   []string `json:"step_tool"`
}

type Props struct {
	Physical int64 `json:"physique_value"`
	Literal  int64 `json:"literature_value"`
	Mental   int64 `json:"magic_value"`
	Wealth   int64 `json:"coin_value"`
	Score    int64 `json:"rankscore"`
	Level    int64 `json:"rankLevel"`
}

/*
type Contact struct {
	Id       string
	Profile  string
	Nickname string
	Count    int
	Last     *Message `bson:",omitempty"`
}
*/
type GameTime struct {
	Game01 time.Time
	Game02 time.Time
	Game03 time.Time
	Game04 time.Time
	Game05 time.Time
}

type UserAuth struct {
	IdCard    *AuthInfo `bson:",omitempty" json:"idcard"`
	IdCardTmp *AuthInfo `bson:",omitempty" json:"-"`
	Cert      *AuthInfo `bson:",omitempty" json:"cert"`
	CertTmp   *AuthInfo `bson:",omitempty" json:"-"`
	Record    *AuthInfo `bson:",omitempty" json:"record"`
	RecordTmp *AuthInfo `bson:",omitempty" json:"-"`
}

type AuthInfo struct {
	Images []string `bson:",omitempty" json:"auth_images"`
	Desc   string   `bson:",omitempty" json:"auth_desc"`
	Status string   `bson:",omitempty" json:"auth_status"`
	Review string   `bson:",omitempty" json:"auth_review"`
}

type UserStat struct {
	OnlineTime      int `json:"onlinetime"`
	Records         int `json:"records"`
	Articles        int `json:"articles"`
	Comments        int `json:"comments"`
	Posts           int `json:"posts"`
	GameTime        int `json:"gametime"`
	LastArticleTime int64
	LastCommentTime int64
	LastThumbTime   int64
	LastGameTime    int64
}

type Ratios struct {
	RunRecv    int
	RunAccept  int
	PostRecv   int
	PostAccept int
	PKRecv     int
	PKAccept   int
}

type Account struct {
	Id       string `bson:"_id,omitempty"`
	Email    string
	Phone    string `bson:",omitempty"`
	Weibo    string
	Nickname string    `bson:",omitempty"`
	Profile  string    `bson:",omitempty"`
	Gender   string    `bson:",omitempty"`
	Birth    int64     `bson:",omitempty"`
	RegTime  time.Time `bson:"reg_time,omitempty"`
	Loc      Location  `bson:",omitempty"`

	Password   string   `bson:",omitempty"`
	Role       string   `bson:",omitempty"`
	Height     int      `bson:",omitempty"`
	Weight     int      `bson:",omitempty"`
	Actor      string   `bson:",omitempty"` // coach
	Admin      bool     `bson:",omitempty"`
	Url        string   `bson:",omitempty"`
	About      string   `bson:",omitempty"`
	Addr       *Address `bson:",omitempty"`
	LocAddr    string   `bson:"locaddr"`
	Photos     []string
	CoverImage string
	Wallet     DbWallet
	Chips      int64

	LastLogin   time.Time `bson:"lastlogin"`
	LoginCount  int       `bson:"login_count"`
	LoginDays   int       `bson:"login_days"`
	LoginAwards []int64   `bson:"login_awards"`

	Props    Props
	Equips   *Equip   `bson:",omitempty"`
	GameTime GameTime `bson:"game_time"`

	Sign        string `bson:",omitempty"`
	Emotion     string `bson:",omitempty"`
	Profession  string `bson:",omitempty"`
	Hobby       string `bson:",omitempty"`
	Hometown    string `bson:",omitempty"`
	Oftenappear string `bson:",omitempty"`

	Contacts []string `bson:",omitempty"`
	Devs     []string `bson:",omitempty"` // apple device id

	Auth *UserAuth `bson:",omitempty"`

	TimeLimit int64 `bson:"timelimit"`
	Privilege int

	Setinfo    bool
	SetinfoAll bool
	PhotoSet   bool
	Push       bool

	Stat *UserStat `bson:",omitempty"`

	Taskid     int      `bson:"task_id,omitempty"`
	TaskStatus string   `bson:"task_status,omitempty"`
	Blocks     []string `bson:",omitempty"`

	Ratios Ratios
}

func (this *Account) IsActor(actor string) bool {
	return this.Actor == actor
}

func (this *Account) IsAdmin() bool {
	return this.Admin
}

func (this *Account) Level() int64 {
	return Score2Level(this.Props.Score)
}

func (this *Account) LastGameTime(typ int) time.Time {
	switch typ {
	case 0x01:
		return this.GameTime.Game01
	case 0x02:
		return this.GameTime.Game02
	case 0x03:
		return this.GameTime.Game03
	case 0x04:
		return this.GameTime.Game04
	case 0x05:
		return this.GameTime.Game05
	}

	return time.Unix(0, 0)
}

func (this *Account) Exists(t string) (bool, error) {
	var query bson.M

	switch t {
	case AccountWeibo:
		query = bson.M{"weibo": this.Weibo}
	case AccountEmail:
		query = bson.M{"email": this.Email}
	case AccountPhone:
		query = bson.M{"phone": this.Phone}
	case "nickname":
		query = bson.M{"nickname": this.Nickname}
	default:
		query = bson.M{"_id": this.Id}
	}

	c, err := count(accountColl, query)
	return c > 0, err
}

func CheckUserExists(id, types string) (bool, error) {
	var query bson.M

	switch types {
	case AccountWeibo:
		query = bson.M{"weibo": id}
	case AccountEmail:
		query = bson.M{"email": id}
	case AccountPhone:
		query = bson.M{"phone": id}
	case "nickname":
		query = bson.M{"nickname": id}
	default:
		query = bson.M{"_id": id}
	}

	c, err := count(accountColl, query)
	return c > 0, err
}

func (this *Account) Find(id, types string) (bool, error) {
	var query bson.M

	switch types {
	case AccountEmail:
		query = bson.M{"email": id}
	case AccountPhone:
		query = bson.M{"phone": id}
	case "nickname":
		query = bson.M{"nickname": id}
	case AccountWeibo:
		query = bson.M{"weibo": id}
	default:
		query = bson.M{"_id": id}
	}

	if err := findOne(accountColl, query, nil, this); err != nil {
		return false, err
	}
	return true, nil
}

func (this *Account) FindPass(id, types, password string) (bool, error) {
	query := bson.M{"password": password}

	switch types {
	case AccountEmail:
		query["email"] = id
	case AccountPhone:
		query["phone"] = id
	case "nickname":
		query["nickname"] = id
	case AccountWeibo:
		query["weibo"] = id
		delete(query, "password")
	default:
		query["_id"] = id
	}

	if err := findOne(accountColl, query, nil, this); err != nil {
		return false, err
	}
	return true, nil
}

func FindByActor(actor string, verbose bool) ([]Account, error) {
	var users []Account
	var selector bson.M

	if actor == "" {
		return nil, nil
	}

	if !verbose {
		selector = bson.M{"_id": 1, "nickname": 1, "profile": 1}
	}

	err := search(accountColl, bson.M{"actor": actor}, selector, 0, 0, nil, nil, &users)
	return users, err
}

func FindUsersByIds(verbose int, ids ...string) ([]Account, error) {
	var users []Account
	var selector bson.M

	if len(ids) == 0 {
		return nil, nil
	}

	switch verbose {
	case 0:
		selector = bson.M{"_id": 1, "nickname": 1, "profile": 1}
	case 1:
		selector = bson.M{"photos": 0, "contacts": 0, "wallet": 0, "equips": 0, "devs": 0}
	case 2:
		selector = bson.M{"contacts": 0, "wallet": 0}
	}

	err := search(accountColl, bson.M{"_id": bson.M{"$in": ids}}, selector, 0, 0, nil, nil, &users)
	return users, err
}

func FindUsersByPhones(phones []string) ([]Account, error) {
	var users []Account

	if err := search(accountColl, bson.M{"phone": bson.M{"$in": phones}}, nil, 0, 0, nil, nil, &users); err != nil {
		return nil, errors.NewError(errors.DbError)
	}

	return users, nil
}

func FindAdmins() ([]Account, error) {
	var users []Account
	selector := bson.M{"_id": 1, "nickname": 1, "profile": 1}
	err := search(accountColl, bson.M{"admin": true}, selector, 0, 0, nil, nil, &users)
	return users, err
}

func (this *Account) FindByWeibo(weibo string) (bool, error) {
	if len(weibo) == 0 {
		return false, nil
	}
	return this.findOne(bson.M{"weibo": weibo})
}

func (this *Account) findOne(query interface{}) (bool, error) {
	var users []Account

	err := search(accountColl, query, nil, 0, 1, nil, nil, &users)
	if err != nil {
		return false, errors.NewError(errors.DbError)
	}
	if len(users) > 0 {
		*this = users[0]
	}
	return len(users) > 0, nil
}

func (this *Account) FindByUserid(userid string) (bool, error) {
	if len(userid) == 0 {
		return false, nil
	}
	return this.findOne(bson.M{"_id": userid})
}

func (this *Account) FindByNickname(nickname string) (bool, error) {
	if len(nickname) == 0 {
		return false, nil
	}
	return this.findOne(bson.M{"nickname": nickname})
}

func (this *Account) FindByPhone(phone string) (bool, error) {
	if len(phone) == 0 {
		return false, nil
	}
	return this.findOne(bson.M{"phone": phone})
}

func (this *Account) FindByUserPass(userid, password string) (bool, error) {
	if len(userid) == 0 || len(password) == 0 {
		return false, nil
	}
	query := bson.M{
		"$or": []bson.M{
			{"email": userid},
			{"phone": userid},
		},
		"password": password,
	}
	return this.findOne(query)
}

func (this *Account) FindByWalletAddr(addr string) (bool, error) {
	if len(addr) == 0 {
		return false, nil
	}
	return this.findOne(bson.M{"wallet.addrs": addr})
}

func (this *Account) CheckExists() (bool, error) {
	if len(this.Id) == 0 || len(this.Nickname) == 0 {
		return false, nil
	}
	return this.findOne(bson.M{"$or": []bson.M{{"_id": this.Id}, {"nickname": this.Nickname}}})
}

var random = rand.New(rand.NewSource(time.Now().Unix()))

func (this *Account) Save() error {
	now := time.Now()

	this.Id = fmt.Sprintf("%d%03d", now.Unix(), now.Nanosecond()%1000)
	this.Push = true
	this.LastLogin = time.Now()
	this.LoginDays = 1
	if len(this.Gender) == 0 {
		this.Gender = "male"
	}
	return save(accountColl, this, true)
}

func (this *Account) SetInfo(setinfo *SetInfo) error {
	change := bson.M{
		"$set": Struct2Map(setinfo),
	}
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

/*
func (this *Account) Update() error {
	change := bson.M{
		"$set": Struct2Map(this),
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}
*/
func (this *Account) UpdateBanTime(banTime int64) error {
	change := bson.M{
		"$set": bson.M{
			"timelimit": banTime,
		},
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

/*
func (this *Account) UpdateLevel(score int, level int) error {
	change := bson.M{
		"$set": bson.M{
			"props.score": score,
			"props.level": level,
		},
	}
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}
*/
func (this *Account) UpdateLocation(loc Location, locaddr string) error {
	change := bson.M{
		"$set": bson.M{
			"loc":     loc,
			"locaddr": locaddr,
		},
	}
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

func (this *Account) SetLoginAwards(awards []int64) error {
	change := bson.M{
		"$set": bson.M{
			"login_wards": awards,
		},
	}
	this.LoginAwards = awards
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

/*
func (this *Account) SetLogin(count int, lastlog time.Time) (int64, error) {
	this.LoginCount = count
	this.LoginAwards = []int{1, 2, 3, 4, 5, 6, 7}
	change := bson.M{
		"$set": bson.M{
			"lastlogin":    lastlog,
			"login_days":  count,
			"login_awards": this.LoginAwards,
		},
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return 0, errors.NewError(errors.DbError, err.Error())
	}
	award := 0
	if count > 7 {
		award = this.LoginAwards[6]
	} else {
		award = this.LoginAwards[count-1]
	}
	return int64(award), nil
}
*/

func (this *Account) SetLastLogin(days int, loginCount int, lastlog time.Time) error {
	change := bson.M{
		"$set": bson.M{
			"lastlogin":   lastlog,
			"login_days":  days,
			"login_count": loginCount,
		},
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}

	return nil
}

func (this *Account) UpdateProps(awards Props) error {
	change := bson.M{
		"$inc": bson.M{
			"props.physical": awards.Physical,
			"props.literal":  awards.Literal,
			"props.mental":   awards.Mental,
			//"props.wealth":   awards.Wealth,
			"props.score": awards.Score,
			"props.level": awards.Level,
		},
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}

	return nil
}

func (this *Account) SetWallet(wallet DbWallet) error {
	change := bson.M{
		"$set": bson.M{
			"wallet": wallet,
		},
	}
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

/*
func (this *Account) ChangePassword(newPass string) error {
	change := bson.M{
		"$set": bson.M{
			"password": newPass,
		},
	}

	if err := update(accountColl, bson.M{"_id": this.Id}, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}
*/

func (this *Account) SetPassword(newPass string) error {
	change := bson.M{
		"$set": bson.M{
			"password": newPass,
		},
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

func (this *Account) ChangeProfile(profile string) error {
	this.Profile = profile

	change := bson.M{
		"$set": bson.M{
			"profile": profile,
		},
	}

	if err := update(accountColl, bson.M{"_id": this.Id}, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

func (this *Account) AddPhotos(photos []string) error {
	change := bson.M{
		"$addToSet": bson.M{
			"photos": bson.M{
				"$each": photos,
			},
		},
		"$set": bson.M{
			"photoset": true,
		},
	}
	if err := update(accountColl, bson.M{"_id": this.Id}, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

func (this *Account) DelPhoto(id string) error {
	change := bson.M{
		"$pull": bson.M{
			"photos": id,
		},
	}
	if err := update(accountColl, bson.M{"_id": this.Id}, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

func (this *Account) Recommend(excludes []string) (users []Account, err error) {
	var list []Account

	query := bson.M{
		"_id": bson.M{
			"$nin": excludes,
		},
		"actor": "coach",
	}
	err = search(accountColl, query, bson.M{"contacts": 0}, 0, 10, nil, nil, &list)
	users = append(users, list...)

	for _, user := range list {
		excludes = append(excludes, user.Id)
	}

	if this.Loc.Lat != 0 && this.Loc.Lng != 0 {
		query := bson.M{
			"loc": bson.D{
				{"$near", []float64{this.Loc.Lat, this.Loc.Lng}},
				{"$maxDistance", float64(50000) / float64(111319)},
			},
			"_id":       bson.M{"$nin": excludes},
			"timelimit": 0,
			"nickname":  bson.M{"$ne": nil},
			"birth": bson.M{
				"$gte": time.Unix(this.Birth, 0).AddDate(-5, 0, 0).Unix(),
				"$lte": time.Unix(this.Birth, 0).AddDate(5, 0, 0).Unix(),
			},
		}
		list = nil
		err = search(accountColl, query, bson.M{"contacts": 0}, 0, 20, []string{"-phone"}, nil, &list)
		users = append(users, list...)
		for _, user := range list {
			excludes = append(excludes, user.Id)
		}
	}

	if len(users) < 12 {
		query := bson.M{
			"_id": bson.M{"$nin": excludes},
			"birth": bson.M{
				"$gte": time.Unix(this.Birth, 0).AddDate(-5, 0, 0).Unix(),
				"$lte": time.Unix(this.Birth, 0).AddDate(5, 0, 0).Unix(),
			},
			"nickname":  bson.M{"$ne": nil},
			"timelimit": 0,
		}
		list = nil
		err = search(accountColl, query, bson.M{"contacts": 0}, 0, 50, []string{"-phone", "-score"}, nil, &list)
		users = append(users, list...)
		for _, user := range list {
			excludes = append(excludes, user.Id)
		}
	}

	if len(users) < 12 {
		query := bson.M{
			"_id":       bson.M{"$nin": excludes},
			"nickname":  bson.M{"$ne": nil},
			"timelimit": 0,
		}
		list = nil
		err = search(accountColl, query, bson.M{"contacts": 0}, 0, 50, []string{"-phone", "-score"}, nil, &list)
		users = append(users, list...)
	}

	return
}

func UserList(sort string, pageIndex, pageCount int) (total int, users []Account, err error) {
	var query bson.M

	switch sort {
	case "regtime":
		sort = "reg_time"
	case "-regtime":
		sort = "-reg_time"
	case "score":
		sort = "props.score"
	case "-score":
		sort = "-props.score"
	case "task":
		sort = "tasks.last"
		/*
			query = bson.M{
				"tasks.count": bson.M{
					"$gt": 0,
				},
			}
		*/
	case "-task":
		sort = "-tasks.last"
		/*
			query = bson.M{
				"tasks.count": bson.M{
					"$gt": 0,
				},
			}
		*/
	default:
		sort = "-reg_time"
	}
	err = search(accountColl, query, nil, pageIndex*pageCount, pageCount, []string{sort}, &total, &users)
	return
}

func UserLeaderBoard(types string, paging *Paging) (users []Account, err error) {
	var sort string
	index := 0
	limit := DefaultPageSize
	if paging.Count > 0 {
		limit = paging.Count
	}

	if len(paging.First) > 0 {
		index, _ = strconv.Atoi(paging.First)
		index -= DefaultPageSize
		if index < 0 {
			index = 0
		}
	} else if len(paging.Last) > 0 {
		index, _ = strconv.Atoi(paging.Last)
	}

	switch types {
	case "physique":
		sort = "-props.physical"
	case "literature":
		sort = "-props.literal"
	case "magic":
		sort = "-props.mental"
	}

	err = search(accountColl, nil, nil, index, limit, []string{sort}, nil, &users)

	paging.First = strconv.Itoa(index)
	paging.Last = strconv.Itoa(index + len(users))

	return
}

// This function returns users after preCursor or nextCursor sorted by sortOrder. The return count total should not be more than limit.
func GetUserListBySort(skip, limit int, sort string) (total int, users []Account, err error) {

	var sortby string

	switch sort {
	case "logintime":
		sortby = "lastlogin"
	case "-logintime":
		sortby = "-lastlogin"
	case "userid":
		sortby = "_id"
	case "-userid":
		sortby = "-_id"
	case "nickname":
		sortby = "nickname"
	case "-nickname":
		sortby = "-nickname"
	case "score":
		sortby = "props.score"
	case "-score":
		sortby = "-props.score"
	case "regtime":
		sortby = "reg_time"
	case "-regtime":
		sortby = "-reg_time"
	case "age":
		sortby = "-birth"
	case "-age":
		sortby = "birth"
	case "gender":
		sortby = "gender"
	case "-gender":
		sortby = "-gender"
	case "ban":
		sortby = "timelimit"
	case "-ban":
		sortby = "-timelimit"
	case "-onlinetime":
		sortby = "-stat.onlinetime"
	case "onlinetime":
		sortby = "stat.onlinetime"
	case "-record":
		sortby = "-stat.records"
	case "record":
		sortby = "stat.records"
	case "-post":
		sortby = "-stat.posts"
	case "post":
		sortby = "stat.posts"
	case "-gametime":
		sortby = "-stat.gametime"
	case "gametime":
		sortby = "stat.gametime"
	default:
		sortby = "-reg_time"
	}

	query := bson.M{"reg_time": bson.M{"$gt": time.Unix(0, 0)}}

	if err = search(accountColl, query, bson.M{"contacts": 0}, skip, limit, []string{sortby}, &total, &users); err != nil {
		return 0, nil, errors.NewError(errors.DbError)
	}

	return
}

// This function search users with userid or nickname after preCursor or nextCursor sorted by sortOrder. The return count total should not be more than limit.
func GetSearchListBySort(id, nickname, keywords string,
	gender, age, banStatus, role string, actor string, skip, limit int, sortOrder string) (total int, users []Account, err error) {

	var sortby string

	switch sortOrder {
	case "logintime":
		sortby = "lastlogin"
	case "-logintime":
		sortby = "-lastlogin"
	case "userid":
		sortby = "_id"
	case "-userid":
		sortby = "-_id"
	case "nickname":
		sortby = "nickname"
	case "-nickname":
		sortby = "-nickname"
	case "score":
		sortby = "props.score"
	case "-score":
		sortby = "-props.score"
	case "regtime":
		sortby = "reg_time"
	case "-regtime":
		sortby = "-reg_time"
	case "age":
		sortby = "-birth"
	case "-age":
		sortby = "birth"
	case "gender":
		sortby = "gender"
	case "-gender":
		sortby = "-gender"
	case "ban":
		sortby = "timelimit"
	case "-ban":
		sortby = "-timelimit"
	default:
		sortby = "-reg_time"
	}

	and := []bson.M{
		{"reg_time": bson.M{"$gt": time.Unix(0, 0)}},
	}

	if len(keywords) > 0 {
		q := bson.M{"$or": []bson.M{
			{"_id": bson.M{"$regex": keywords, "$options": "i"}},
			{"nickname": bson.M{"$regex": keywords, "$options": "i"}},
			{"phone": bson.M{"$regex": keywords, "$options": "i"}},
			{"about": bson.M{"$regex": keywords, "$options": "i"}},
			{"hobby": bson.M{"$regex": keywords, "$options": "i"}},
		}}
		and = append(and, q)
	}

	if len(gender) > 0 {
		if strings.HasPrefix(gender, "f") {
			and = append(and, bson.M{"gender": bson.M{"$in": []interface{}{"f", "female"}}})
		} else {
			and = append(and, bson.M{"gender": bson.M{"$in": []interface{}{"m", "male", nil}}})
		}
	}
	if len(age) > 0 {
		s := strings.Split(age, "-")
		if len(s) == 1 {
			if a, err := strconv.Atoi(s[0]); err == nil {
				if a == 0 {
					and = append(and, bson.M{"birth": bson.M{"$exists": false}})
				} else {
					start, end := AgeToTimeRange(a)
					and = append(and, bson.M{"birth": bson.M{"$gte": start.Unix(), "$lte": end.Unix()}})
				}

			}
		}
		if len(s) == 2 {
			low, _ := strconv.Atoi(s[0])
			high, _ := strconv.Atoi(s[1])
			if low == high {
				start, end := AgeToTimeRange(low)
				and = append(and, bson.M{"birth": bson.M{
					"$gte": start.Unix(),
					"$lte": end.Unix(),
				}})
			} else {
				if low > high {
					low, high = high, low
				}
				start, _ := AgeToTimeRange(high)
				_, end := AgeToTimeRange(low)

				if low == 0 {
					and = append(and, bson.M{"$or": []bson.M{
						{"birth": bson.M{"$gte": start.Unix(), "$lte": end.Unix()}},
						{"birth": bson.M{"$exists": false}},
					}})
				} else {
					and = append(and, bson.M{"birth": bson.M{"$gte": start.Unix(), "$lte": end.Unix()}})
				}
			}
		}
	}
	if len(banStatus) > 0 {
		switch banStatus {
		case "normal":
			and = append(and, bson.M{"timelimit": bson.M{"$in": []interface{}{0, nil}}})
		case "lock":
			and = append(and, bson.M{"timelimit": bson.M{"$gt": 0}})
		case "ban":
			and = append(and, bson.M{"timelimit": bson.M{"$lt": 0}})
		}
	}
	if len(role) > 0 {
		and = append(and, bson.M{"role": role})
	}

	if len(actor) > 0 {
		switch actor {
		case "user":
			and = append(and, bson.M{"actor": bson.M{"$in": []interface{}{nil, ""}}, "admin": nil})
		case "coach":
			and = append(and, bson.M{"actor": "coach"})
		case "admin":
			and = append(and, bson.M{"admin": true})
		}
	}

	query := bson.M{"$and": and}

	//b, _ := json.Marshal(query)
	//log.Println("query:", string(b))
	if err = search(accountColl, query, bson.M{"contacts": 0}, skip, limit, []string{sortby}, &total, &users); err != nil {
		return 0, nil, errors.NewError(errors.DbError)
	}

	return
}

// This function returns the friends list of the user. Return users after preCursor or nextCursor and sorted by sortOrder.
// The return count total should not be more than limit
func GetFriendsListBySort(ids []string, skip, limit int, sortOrder string) (users []Account, err error) {
	var sortby string
	query := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	switch sortOrder {
	case "logintime":
		sortby = "-lastlogin"

	case "userid":
		sortby = "_id"

	case "nickname":
		sortby = "nickname"

	case "score":
		sortby = "-props.score"

	case "regtime":
		log.Println("regtime")
		fallthrough
	default:
		sortby = "-reg_time"
	}

	err = search(accountColl, query, nil, skip, limit, []string{sortby}, nil, &users)

	return
}

func recordPagingFunc(c *mgo.Collection, first, last string, args ...interface{}) (query bson.M, err error) {
	record := &Record{}

	if bson.IsObjectIdHex(first) {
		if err := c.FindId(bson.ObjectIdHex(first)).One(record); err != nil {
			return nil, err
		}
		query = bson.M{
			"starttime": bson.M{
				"$gte": record.PubTime,
			},
		}
	} else if bson.IsObjectIdHex(last) {
		if err := c.FindId(bson.ObjectIdHex(last)).One(record); err != nil {
			return nil, err
		}
		query = bson.M{
			"starttime": bson.M{
				"$lte": record.PubTime,
			},
		}
	}

	return
}

func (this *Account) Records(all bool, types string, paging *Paging) (int, []Record, error) {
	var records []Record
	total := 0
	var query bson.M

	switch types {
	case "game":
		query = bson.M{"uid": this.Id, "type": "game"}
	case "run":
		query = bson.M{"uid": this.Id, "type": "run"}
	default:
		query = bson.M{"uid": this.Id, "type": bson.M{"$ne": "post"}}
	}
	if !all {
		query["status"] = bson.M{"$in": []string{StatusFinish, ""}}
	}
	sortFields := []string{"-starttime", "-_id"}

	if err := psearch(recordColl, query, nil,
		sortFields, nil, &records, recordPagingFunc, paging); err != nil && err != mgo.ErrNotFound {

		return total, nil, errors.NewError(errors.DbError)
	}

	for i := 0; i < len(records); i++ {
		if records[i].Id.Hex() == paging.First {
			records = records[:i]
			break
		} else if records[i].Id.Hex() == paging.Last {
			records = records[i+1:]
			break
		}
	}

	paging.First = ""
	paging.Last = ""
	paging.Count = 0
	if len(records) > 0 {
		paging.First = records[0].Id.Hex()
		paging.Last = records[len(records)-1].Id.Hex()
		paging.Count = total
	}

	return total, records, nil
}

func (this *Account) TaskRecords(types string) (records []Record, err error) {
	query := bson.M{
		"uid": this.Id,
		"task": bson.M{
			"$lte": 1000,
		},
		"type":   types,
		"status": StatusFinish,
	}

	err = search(recordColl, query, nil, 0, 0, nil, nil, &records)
	return
}

func (this *Account) TaskRecordCount(types, status string) (int, error) {
	query := bson.M{
		"uid": this.Id,
		"task": bson.M{
			"$lte": 1000,
		},
		"status": status,
	}
	if len(types) > 0 {
		query["type"] = types
	}

	count, err := count(recordColl, query)

	return count, err
}

func (this *Account) LastRecord(types string) (*Record, error) {
	query := bson.M{
		"uid":  this.Id,
		"type": types,
	}
	record := &Record{}
	err := findOne(recordColl, query, []string{"pub_time"}, record)
	return record, err
}

func (this *Account) LastTaskRecord() (*Record, error) {
	query := bson.M{
		"uid": this.Id,
		"task": bson.M{
			"$gt": 0,
			"$lt": 1000,
		},
	}
	record := &Record{}
	err := findOne(recordColl, query, []string{"-auth_time"}, record)

	return record, err
}

func (this *Account) LastTaskRecord2() (*Record, error) {
	query := bson.M{
		"uid": this.Id,
		"task": bson.M{
			"$gt": 0,
			"$lt": 1000,
		},
	}
	record := &Record{}
	err := findOne(recordColl, query, []string{"-task", "-pub_time"}, record)

	return record, err
}

func (this *Account) UpdateAction(action string, date time.Time) (bool, error) {
	/*
		selector := bson.M{
			"userid": this.Id,
			"date":   date,
		}
		update := bson.M{
			"$inc": bson.M{
				action: 1,
			},
		}
		chinfo, err := upsert(actionColl, selector, update, true)
		//log.Println(chinfo, err)
		if err != nil {
			return false, errors.NewError(errors.DbError, err.Error())
		}

		return chinfo.UpsertedId != nil, nil
	*/
	return false, nil
}

func friendsPagingFunc(c *mgo.Collection, first, last string, args ...interface{}) (query bson.M, err error) {
	user := &Account{}

	var ids interface{}
	if len(args) > 0 {
		ids = args[0]
	}

	query = bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	if len(first) > 0 {
		if err := c.FindId(first).One(user); err != nil {
			return nil, err
		}
		query = bson.M{
			"props.score": bson.M{
				"$gte": user.Props.Score,
			},
			"_id": bson.M{
				"$in": ids,
			},
		}
	} else if len(last) > 0 {
		if err := c.FindId(last).One(user); err != nil {
			return nil, err
		}
		query = bson.M{
			"props.score": bson.M{
				"$lte": user.Props.Score,
			},
			"_id": bson.M{
				"$in": ids,
			},
		}
	}

	return
}

func UserCount() (n int) {
	n, _ = count(accountColl, bson.M{"reg_time": bson.M{"$gt": time.Unix(0, 0)}})
	return
}

func Users(ids []string, paging *Paging) (users []Account, err error) {
	//total := 0
	index := 0
	limit := DefaultPageSize
	if paging.Count > 0 {
		limit = paging.Count
	}

	if len(paging.First) > 0 {
		index, _ = strconv.Atoi(paging.First)
		index -= DefaultPageSize
		if index < 0 {
			index = 0
		}
	} else if len(paging.Last) > 0 {
		index, _ = strconv.Atoi(paging.Last)
	}

	query := bson.M{"_id": bson.M{"$in": ids}}
	sortFields := []string{"-props.score", "-_id"}
	err = search(accountColl, query, nil, index, limit, sortFields, nil, &users)

	paging.First = strconv.Itoa(index)
	paging.Last = strconv.Itoa(index + len(users))

	return
}

func (this *Account) ArticleCount() (count int) {
	query := bson.M{"author": this.Id, "parent": nil}
	search(articleColl, query, nil, 0, 0, nil, &count, nil)
	return
}

func (this *Account) SetEquip(equip Equip) error {
	change := bson.M{
		"$set": bson.M{
			"equips": equip,
		},
	}

	if err := update(accountColl, bson.M{"_id": this.Id}, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func searchPagingFunc(c *mgo.Collection, first, last string, args ...interface{}) (query bson.M, err error) {
	user := &Account{}

	if len(first) > 0 {
		if err := c.FindId(first).One(user); err != nil {
			return nil, err
		}
		query = bson.M{
			"lastlogin": bson.M{
				"$gte": user.LastLogin,
			},
		}
	} else if len(last) > 0 {
		if err := c.FindId(last).One(user); err != nil {
			return nil, err
		}
		query = bson.M{
			"lastlogin": bson.M{
				"$lte": user.LastLogin,
			},
		}
	}

	return
}

func SearchUsers(nickname string, paging *Paging) ([]Account, error) {
	var users []Account
	total := 0

	query := bson.M{}

	if len(nickname) > 0 {
		query["nickname"] = bson.M{
			"$regex":   charFilter(nickname),
			"$options": "i",
		}
	}

	sortFields := []string{"-lastlogin", "-_id"}

	if err := psearch(accountColl, query, bson.M{"contacts": 0}, sortFields, nil, &users, searchPagingFunc, paging); err != nil {
		if err != mgo.ErrNotFound {
			return nil, errors.NewError(errors.DbError, err.Error())
		}
	}

	for i := 0; i < len(users); i++ {
		if users[i].Id == paging.First {
			users = users[:i]
			break
		} else if users[i].Id == paging.Last {
			users = users[i+1:]
			break
		}
	}

	paging.First = ""
	paging.Last = ""
	paging.Count = 0
	if len(users) > 0 {
		paging.First = users[0].Id
		paging.Last = users[len(users)-1].Id
		paging.Count = total
	}
	return users, nil
}

func (this *Account) SearchNear(paging *Paging, distance int) ([]Account, error) {
	var users []Account
	total := 0
	fmt.Println("search nearby:", this.Loc.Lat, this.Loc.Lng, distance)
	if this.Loc.Lat == 0 && this.Loc.Lng == 0 {
		return nil, nil
	}
	query := bson.M{
		"loc": bson.M{
			"$near": []float64{this.Loc.Lat, this.Loc.Lng},
		},
		"nickname": bson.M{"$ne": nil},
	}
	if distance > 0 {
		query = bson.M{
			"loc": bson.D{
				{"$near", []float64{this.Loc.Lat, this.Loc.Lng}},
				{"$maxDistance", float64(distance) / float64(111319)},
			},
			"nickname": bson.M{"$ne": nil},
		}
	}

	if err := psearch(accountColl, query, bson.M{"contacts": 0}, nil, nil, &users, nil, paging); err != nil {
		if err != mgo.ErrNotFound {
			return nil, errors.NewError(errors.DbError, err.Error())
		}
	}

	paging.First = ""
	paging.Last = ""
	paging.Count = 0
	if len(users) > 0 {
		paging.First = users[0].Id
		paging.Last = users[len(users)-1].Id
		paging.Count = total
	}
	return users, nil
}

func (this *Account) AddWalletAddr(addr string) error {
	change := bson.M{
		"$addToSet": bson.M{
			"wallet.addrs": addr,
		},
	}
	if err := update(accountColl, bson.M{"_id": this.Id}, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Account) EventCount(eventType, pushType string) int {
	query := bson.M{
		"push.to": this.Id,
		"type":    bson.M{"$ne": EventNotice},
	}
	if len(eventType) > 0 {
		query["type"] = eventType
	}
	if len(pushType) > 0 {
		query["push.type"] = pushType
	}

	n, _ := count(eventColl, query)
	return n
}

func (this *Account) UpdateInfo(change bson.M) error {
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Account) ArticleTimeline(pageIndex, pageCount int) (total int, articles []Article, err error) {
	err = search(articleColl, bson.M{"author": this.Id, "parent": nil}, nil,
		pageIndex*pageCount, pageCount, []string{"-pub_time"}, &total, &articles)
	return
}

/*
func (this *Account) GetTasks() (TaskList, error) {
	_, err := this.FindByUserid(this.Id)
	return this.Tasks, err
}

func (this *Account) AddTask(typ string, tid int, proofs []string) error {
	selector := bson.M{
		"_id": this.Id,
	}

	update(accountColl, selector, bson.M{"$pull": bson.M{"tasks.proofs": bson.M{"tid": tid}}}, true)

	var change bson.M
	if typ == TaskRunning {
		change = bson.M{
			"$pull": bson.M{
				"tasks.uncompleted": tid,
			},
			"$addToSet": bson.M{
				"tasks.waited": tid,
				"tasks.proofs": Proof{Tid: tid, Pics: proofs},
			},
			"$set": bson.M{
				"tasks.last": time.Now(),
			},
			"$inc": bson.M{
				"tasks.count": 1,
			},
		}
	} else {
		change = bson.M{
			"$addToSet": bson.M{
				"tasks.completed": tid,
			},
			"$set": bson.M{
				"tasks.last": time.Now(),
			},
		}
	}

	if err := update(accountColl, selector, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Account) SetTaskComplete(tid int, completed bool, reason string) error {
	if len(reason) > 0 {
		selector := bson.M{
			"_id":              this.Id,
			"tasks.proofs.tid": tid,
		}
		update(accountColl, selector, bson.M{"$set": bson.M{"tasks.proofs.$.result": reason}}, true)
	}

	selector := bson.M{
		"_id": this.Id,
	}
	var change bson.M

	if completed {
		change = bson.M{
			"$pull": bson.M{
				"tasks.waited": tid,
			},
			"$addToSet": bson.M{
				"tasks.completed": tid,
			},
			"$inc": bson.M{
				"tasks.count": -1,
			},
		}
	} else {
		change = bson.M{
			"$pull": bson.M{
				"tasks.waited": tid,
			},
			"$addToSet": bson.M{
				"tasks.uncompleted": tid,
			},
			"$inc": bson.M{
				"tasks.count": -1,
			},
		}
	}

	if err := update(accountColl, selector, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}
*/
func (this *Account) Articles(typ string, paging *Paging) (int, []Article, error) {
	var articles []Article
	var query, selector bson.M
	total := 0

	switch typ {
	case "COMMENTS":
		query = bson.M{"author": this.Id, "parent": bson.M{"$ne": nil}}
		//selector = bson.M{"content": 0, "contents": 0}
	case "ARTICLES":
		query = bson.M{"author": this.Id, "parent": nil, "refer": nil}
	default:
		query = bson.M{"author": this.Id}
	}

	sortFields := []string{"-pub_time", "-_id"}

	if err := psearch(articleColl, query, selector, sortFields, &total, &articles,
		articlePagingFunc, paging); err != nil {
		e := errors.NewError(errors.DbError, err.Error())
		if err == mgo.ErrNotFound {
			e = errors.NewError(errors.NotFoundError, err.Error())
		}
		return total, nil, e
	}

	for i := 0; i < len(articles); i++ {
		if articles[i].Id.Hex() == paging.First {
			articles = articles[:i]
			break
		} else if articles[i].Id.Hex() == paging.Last {
			articles = articles[i+1:]
			break
		}
	}

	paging.First = ""
	paging.Last = ""
	paging.Count = 0
	if len(articles) > 0 {
		paging.First = articles[0].Id.Hex()
		paging.Last = articles[len(articles)-1].Id.Hex()
		paging.Count = total
	}

	return total, articles, nil
}

func (this *Account) Messages(userid string, paging *Paging) (int, []Message, error) {
	var msgs []Message
	total := 0
	query := bson.M{
		"$or": []bson.M{
			bson.M{"from": userid, "to": this.Id},
			bson.M{"from": this.Id, "to": userid},
		},
		"owners": this.Id,
	}

	sortFields := []string{"-time", "-_id"}

	if err := psearch(msgColl, query, nil, sortFields, &total, &msgs,
		msgPagingFunc, paging); err != nil {
		e := errors.NewError(errors.DbError, err.Error())
		if err == mgo.ErrNotFound {
			e = errors.NewError(errors.NotFoundError, err.Error())
		}
		return total, nil, e
	}

	for i := 0; i < len(msgs); i++ {
		if msgs[i].Id.Hex() == paging.First {
			msgs = msgs[:i]
			break
		} else if msgs[i].Id.Hex() == paging.Last {
			msgs = msgs[i+1:]
			break
		}
	}

	paging.First = ""
	paging.Last = ""
	paging.Count = 0
	if len(msgs) > 0 {
		paging.First = msgs[0].Id.Hex()
		paging.Last = msgs[len(msgs)-1].Id.Hex()
		paging.Count = total
	}

	return total, msgs, nil
}

func (this *Account) AddContact(contact string) error {
	/*
		selector := bson.M{
			"_id":         this.Id,
			"contacts.id": contact.Id,
		}
		change := bson.M{

			"$inc": bson.M{
				"contacts.$.count": contact.Count,
			},

			"$set": bson.M{
				"contacts.$.profile":  contact.Profile,
				"contacts.$.nickname": contact.Nickname,
				"contacts.$.last":     contact.Last,
			},
		}
		err := update(accountColl, selector, change, true)
		if err == nil {
			return nil
		}
		//log.Println(err)
		if err != mgo.ErrNotFound {
			return errors.NewError(errors.DbError, err.Error())
		}

		// not found
	*/

	change := bson.M{
		"$addToSet": bson.M{
			"contacts": contact,
		},
	}
	err := updateId(accountColl, this.Id, change, true)
	if err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Account) DelContact(contact string) error {
	change := bson.M{
		"$pull": bson.M{
			"contacts": contact,
		},
	}
	err := updateId(accountColl, this.Id, change, true)
	if err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

/*
func (this *Account) MarkRead(typ, id string) error {
	var selector, change bson.M

	switch typ {
	case "chat":
		selector = bson.M{
			"_id":         this.Id,
			"contacts.id": id,
		}
		change = bson.M{
			"$set": bson.M{
				"contacts.$.count": 0,
			},
		}
	case "article":
		selector = bson.M{
			"_id":       this.Id,
			"events.id": id,
		}
		change = bson.M{
			"$unset": bson.M{
				"events.$.reviews": 1,
				"events.$.thumbs":  1,
			},
		}
	default:
		return nil
	}

	if err := update(accountColl, selector, change, true); err != nil {
		if err != mgo.ErrNotFound {
			return errors.NewError(errors.DbError)
		}
	}

	return nil
}
*/
func (this *Account) SetPush(push bool) error {
	selector := bson.M{
		"_id": this.Id,
	}
	change := bson.M{
		"$set": bson.M{
			"push": push,
		},
	}
	if err := update(accountColl, selector, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Account) PushEnabled() (bool, error) {
	var users []Account
	var enabled bool

	err := search(accountColl, bson.M{"_id": this.Id}, bson.M{"push": 1},
		0, 1, nil, nil, &users)
	if err != nil {
		return false, errors.NewError(errors.DbError, err.Error())
	}

	if len(users) > 0 {
		enabled = users[0].Push
	}

	return enabled, nil
}

/*
func (this *Account) Devices() ([]string, bool, error) {
	var users []Account
	var devs []string
	var enabled bool

	err := search(accountColl, bson.M{"_id": this.Id}, bson.M{"devs": true, "push": true},
		0, 1, nil, nil, &users)
	if err != nil {
		return nil, false, errors.NewError(errors.DbError, err.Error())
	}

	if len(users) > 0 {
		devs = users[0].Devs
		enabled = users[0].Push
	}

	return devs, enabled, nil
}
*/
func (this *Account) AddDevice(dev string) error {
	selector := bson.M{
		"_id": this.Id,
	}
	change := bson.M{
		"$addToSet": bson.M{
			"devs": dev,
		},
	}
	if err := update(accountColl, selector, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Account) RmDevice(dev string) error {
	selector := bson.M{
		"_id": this.Id,
	}
	change := bson.M{
		"$pull": bson.M{
			"devs": dev,
		},
	}
	if err := update(accountColl, selector, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Account) LatestArticle() *Article {
	article := &Article{}
	findOne(articleColl, bson.M{"author": this.Id, "parent": nil},
		[]string{"-pub_time"}, article)

	return article
}

func (this *Account) ContactList() ([]string, error) {
	err := findOne(accountColl, bson.M{"_id": this.Id}, nil, this)
	return this.Contacts, err
}

func (this *Account) SetGameTime(typ int, t time.Time) error {
	var change bson.M

	switch typ {
	case 0x01:
		change = bson.M{
			"$set": bson.M{"game_time.game01": t},
		}
	case 0x02:
		change = bson.M{
			"$set": bson.M{"game_time.game02": t},
		}
	case 0x03:
		change = bson.M{
			"$set": bson.M{"game_time.game03": t},
		}
	case 0x04:
		change = bson.M{
			"$set": bson.M{"game_time.game04": t},
		}
	case 0x05:
		change = bson.M{
			"$set": bson.M{"game_time.game05": t},
		}
	default:
		change = bson.M{
			"$set": bson.M{"game_time.game00": t},
		}
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}

	return nil
}

func (this *Account) SetAdmin(isAdmin bool) error {
	change := bson.M{
		"$set": bson.M{
			"admin": isAdmin,
		},
	}
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}

	return nil
}

func txPagingFunc(c *mgo.Collection, first, last string, args ...interface{}) (query bson.M, err error) {
	tx := &Tx{}

	if bson.IsObjectIdHex(first) {
		if err := c.FindId(bson.ObjectIdHex(first)).One(tx); err != nil {
			return nil, err
		}
		query = bson.M{
			"time": bson.M{
				"$gte": tx.Time,
			},
		}
	} else if bson.IsObjectIdHex(last) {
		if err := c.FindId(bson.ObjectIdHex(last)).One(tx); err != nil {
			return nil, err
		}
		query = bson.M{
			"time": bson.M{
				"$lte": tx.Time,
			},
		}
	}

	return
}

func (this *Account) Txs(paging *Paging) (int, []Tx, error) {
	var txs []Tx
	total := 0

	sortFields := []string{"-time", "-_id"}
	query := bson.M{"uid": this.Id}
	psearch(txColl, query, nil, sortFields, nil, &txs, txPagingFunc, paging)

	for i := 0; i < len(txs); i++ {
		if txs[i].Id.Hex() == paging.First {
			txs = txs[:i]
			break
		} else if txs[i].Id.Hex() == paging.Last {
			txs = txs[i+1:]
			break
		}
	}

	paging.First = ""
	paging.Last = ""
	paging.Count = 0
	if len(txs) > 0 {
		paging.First = txs[0].Id.Hex()
		paging.Last = txs[len(txs)-1].Id.Hex()
		paging.Count = total
	}

	return total, txs, nil
}

func GetAuthUserList(index, count int) (total int, users []Account, err error) {
	query := bson.M{
		"$or": []bson.M{
			{"auth.idcardtmp": bson.M{"$ne": nil}},
			{"auth.certtmp": bson.M{"$ne": nil}},
			{"auth.recordtmp": bson.M{"$ne": nil}},
		},
	}
	err = search(accountColl, query, nil, index*count, count, nil, &total, &users)
	return
}

func (this *Account) SetAuthInfo(types string, images []string, desc string) error {
	var change bson.M

	authInfo := &AuthInfo{Images: images, Desc: desc, Status: AuthVerifying}

	switch types {
	case AuthIdCard:
		change = bson.M{
			"$set": bson.M{"auth.idcardtmp": authInfo},
		}
	case AuthCert:
		change = bson.M{
			"$set": bson.M{"auth.certtmp": authInfo},
		}
	case AuthRecord:
		change = bson.M{
			"$set": bson.M{"auth.recordtmp": authInfo},
		}
	default:
		return nil
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}

	return nil
}

func (this *Account) SetAuth(types, status, review string) error {
	var change bson.M

	switch types {
	case AuthIdCard:
		if this.Auth.IdCardTmp == nil {
			return nil
		}
		this.Auth.IdCardTmp.Review = review
		this.Auth.IdCardTmp.Status = status
		change = bson.M{
			"$set": bson.M{
				"auth.idcardtmp.status": status,
				"auth.idcardtmp.review": review,
			},
		}
		if status == AuthVerified && this.Auth.IdCardTmp != nil {
			//this.Auth.IdCardTmp.Status = AuthVerified
			change = bson.M{
				"$set":   bson.M{"auth.idcard": this.Auth.IdCardTmp},
				"$unset": bson.M{"auth.idcardtmp": 1},
			}
			if (this.Auth.Cert != nil && this.Auth.Cert.Status == AuthVerified) ||
				(this.Auth.Record != nil && this.Auth.Record.Status == AuthVerified) {
				change = bson.M{
					"$set": bson.M{
						"auth.idcard": this.Auth.IdCardTmp,
						"actor":       ActorCoach,
					},
					"$unset": bson.M{"auth.idcardtmp": 1},
				}
			}
		}
	case AuthCert:
		if this.Auth.CertTmp == nil {
			return nil
		}
		this.Auth.CertTmp.Review = review
		this.Auth.CertTmp.Status = status
		change = bson.M{
			"$set": bson.M{
				"auth.certtmp.status": status,
				"auth.certtmp.review": review,
			},
		}
		if status == AuthVerified && this.Auth.CertTmp != nil {
			//this.Auth.CertTmp.Status = AuthVerified
			change = bson.M{
				"$set":   bson.M{"auth.cert": this.Auth.CertTmp},
				"$unset": bson.M{"auth.certtmp": 1},
			}
			if this.Auth.IdCard != nil && this.Auth.IdCard.Status == AuthVerified {
				change = bson.M{
					"$set": bson.M{
						"auth.cert": this.Auth.CertTmp,
						"actor":     ActorCoach,
					},
					"$unset": bson.M{"auth.certtmp": 1},
				}
			}
		}
	case AuthRecord:
		if this.Auth.RecordTmp == nil {
			return nil
		}
		this.Auth.RecordTmp.Review = review
		this.Auth.RecordTmp.Status = status
		change = bson.M{
			"$set": bson.M{
				"auth.recordtmp.status": status,
				"auth.recordtmp.review": review,
			},
		}
		if status == AuthVerified && this.Auth.RecordTmp != nil {
			//this.Auth.RecordTmp.Status = AuthVerified
			change = bson.M{
				"$set":   bson.M{"auth.record": this.Auth.RecordTmp},
				"$unset": bson.M{"auth.recordtmp": 1},
			}
			if this.Auth.IdCard != nil && this.Auth.IdCard.Status == AuthVerified {
				change = bson.M{
					"$set": bson.M{
						"auth.record": this.Auth.RecordTmp,
						"actor":       ActorCoach,
					},
					"$unset": bson.M{"auth.recordtmp": 1},
				}
			}
		}
	default:
		return nil
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}

	return nil
}

func (this *Account) UpdateStat(types string, count int64) error {
	var change bson.M

	switch types {
	case StatOnlineTime:
		change = bson.M{"$inc": bson.M{"stat.onlinetime": count}}
	case StatArticles:
		change = bson.M{
			"$inc": bson.M{"stat.articles": count, "stat.posts": count},
		}
	case StatComments:
		change = bson.M{
			"$inc": bson.M{"stat.comments": count, "stat.posts": count},
		}
	case StatRecords:
		change = bson.M{
			"$inc": bson.M{"stat.records": count},
		}
	case StatGameTime:
		change = bson.M{
			"$inc": bson.M{"stat.gametime": count},
		}
	case StatLastArticleTime:
		change = bson.M{
			"$set": bson.M{"stat.lastarticletime": count},
		}
	case StatLastCommentTime:
		change = bson.M{
			"$set": bson.M{"stat.lastcommenttime": count},
		}
	case StatLastThumbTime:
		change = bson.M{
			"$set": bson.M{"stat.lastthumbtime": count},
		}
	case StatLastGameTime:
		change = bson.M{
			"$set": bson.M{"stat.lastgametime": count},
		}
	default:
		return nil
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}

	return nil
}

func (this *Account) UpdateTask(tid int, status string) error {
	change := bson.M{
		"$set": bson.M{
			"task_id":     tid,
			"task_status": status,
		},
	}
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

func (this *Account) TaskReferrals(taskType string, excludes []string) (accounts []Account, err error) {
	query := bson.M{
		"task_id":     this.Taskid,
		"task_status": this.TaskStatus,
		"_id": bson.M{
			"$nin": excludes,
		},
		"blocks": bson.M{
			"$ne": this.Id,
		},
	}

	if this.Loc.Lat > 0 {
		query["loc"] = bson.D{
			{"$near", []float64{this.Loc.Lat, this.Loc.Lng}},
			// {"$maxDistance", float64(100000) / float64(111319)},
		}
	}

	err = search(accountColl, query, nil, 0, 10, nil, nil, &accounts)

	return
}

func (this *Account) PropIndex(prop string, value int64) (int, error) {
	var query bson.M

	switch prop {
	case "physique":
		query = bson.M{
			"props.literal": bson.M{
				"$gt": value,
			},
		}
	case "literature":
		query = bson.M{
			"props.literal": bson.M{
				"$gt": value,
			},
		}
	case "magic":
		query = bson.M{
			"props.mental": bson.M{
				"$gt": value,
			},
		}
	case "score":
		query = bson.M{
			"props.score": bson.M{
				"$gt": value,
			},
		}
	default:
		return 0, nil
	}

	n, err := count(accountColl, query)
	if err != nil {
		return 0, err
	}
	return n + 1, nil
}

func (this *Account) UpdateRatio(types string, accept bool) error {
	var change bson.M

	switch types {
	case TaskRunning:
		if accept {
			change = bson.M{
				"$inc": bson.M{
					"ratios.runaccept": 1,
				},
			}
		} else {
			change = bson.M{
				"$inc": bson.M{
					"ratios.runrecv": 1,
				},
			}
		}
	case TaskPost:
		if accept {
			change = bson.M{
				"$inc": bson.M{
					"ratios.postaccept": 1,
				},
			}
		} else {
			change = bson.M{
				"$inc": bson.M{
					"ratios.postrecv": 1,
				},
			}
		}
	case TaskGame:
		if accept {
			change = bson.M{
				"$inc": bson.M{
					"ratios.pkaccept": 1,
				},
			}
		} else {
			change = bson.M{
				"$inc": bson.M{
					"ratios.pkrecv": 1,
				},
			}
		}
	default:
		return nil
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

func (this *Account) SetBlock(userid string, block bool) error {
	change := bson.M{
		"$addToSet": bson.M{
			"blocks": userid,
		},
	}
	if !block {
		change = bson.M{
			"$pull": bson.M{
				"blocks": userid,
			},
		}
	}
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}
