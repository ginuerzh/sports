// account
package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"labix.org/v2/mgo/txn"
	"log"
	"time"
)

var (
//dur time.Duration
)

func init() {
	//dur, _ = time.ParseDuration("-30h") // auto logout after 15 minutes since last access

	ensureIndex(accountColl, "_id", "password")
	ensureIndex(accountColl, "-score")
	ensureIndex(accountColl, "nickname")
	ensureIndex(accountColl, "-reg_time")
	ensureIndex(accountColl, "-lastlogin")
	ensureIndex2D(accountColl, "loc")
}

type UserInfo struct {
	Hobby    string `json:"hobby"`
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

	Lng float64
	Lat float64
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

type Account struct {
	Id          string    `bson:"_id,omitempty" json:"-"`
	Nickname    string    `bson:",omitempty" json:"nickname,omitempty"`
	Password    string    `bson:",omitempty" json:"password,omitempty"`
	Profile     string    `bson:",omitempty" json:"profile,omitempty"`
	RegTime     time.Time `bson:"reg_time,omitempty" json:"-"`
	Role        string    `bson:",omitempty" json:"-"`
	Hobby       string    `bson:",omitempty" json:"hobby,omitempty"`
	Height      int       `bson:",omitempty" json:"height,omitempty"`
	Weight      int       `bson:",omitempty" json:"weight,omitempty"`
	Birth       int64     `bson:",omitempty" json:"birth,omitempty"`
	Actor       string    `bson:",omitempty" json:"actor,omitempty"`
	Gender      string    `bson:",omitempty" json:"gender,omitempty"`
	Url         string    `bson:",omitempty" json:"url,omitempty"`
	Phone       string    `bson:",omitempty" json:"phone,omitempty"`
	About       string    `bson:",omitempty" json:"about,omitempty"`
	Addr        *Address  `bson:",omitempty" json:"addr,omitempty"`
	Loc         *Location `bson:",omitempty" json:"-"`
	Photos      []string  `json:"-"`
	Setinfo     bool      `json:"setinfo,omitempty"`
	Wallet      DbWallet  `json:"-"`
	LastLogin   time.Time `bson:"lastlogin" json:"-"`
	LoginCount  int       `bson:"login_count" json:"-"`
	LoginAwards []int     `bson:"login_awards" json:"-"`

	Score int `json:"-"`
	Level int `json:"-"`

	Equips *Equip `bson:",omitempty" json:"-"`

	TimeLimit int64 `bson:"timelimit" json:"timelimit"`
}

func (this *Account) Exists() (bool, error) {
	return this.findOne(bson.M{"_id": this.Id})
}

func FindUsers(ids []string) ([]Account, error) {
	var users []Account
	if err := search(accountColl, bson.M{"_id": bson.M{"$in": ids}}, nil, 0, 0, nil, nil, &users); err != nil {
		return nil, errors.NewError(errors.DbError, err.Error())
	}

	return users, nil
}

func (this *Account) findOne(query interface{}) (bool, error) {
	var users []Account

	err := search(accountColl, query, nil, 0, 1, nil, nil, &users)
	if err != nil {
		return false, errors.NewError(errors.DbError, err.Error())
	}
	if len(users) > 0 {
		*this = users[0]
	}
	return len(users) > 0, nil
}

func (this *Account) FindByUserid(userid string) (bool, error) {
	return this.findOne(bson.M{"_id": userid})
}

func (this *Account) FindByNickname(nickname string) (bool, error) {
	return this.findOne(bson.M{"nickname": nickname})
}

func (this *Account) FindByUserPass(userid, password string) (bool, error) {
	return this.findOne(bson.M{"_id": userid, "password": password})
}

func (this *Account) FindByWalletAddr(addr string) (bool, error) {
	return this.findOne(bson.M{"wallet.addrs": addr})
}

func (this *Account) CheckExists() (bool, error) {
	return this.findOne(bson.M{"$or": []bson.M{{"_id": this.Id}, {"nickname": this.Nickname}}})
}

func (this *Account) Save() error {
	f := func(c *mgo.Collection) error {
		runner := txn.NewRunner(c)
		ops := []txn.Op{
			{
				C:      accountColl,
				Id:     this.Id,
				Assert: txn.DocMissing,
				Insert: this,
			},
			{
				C:      userColl,
				Id:     this.Id,
				Assert: txn.DocMissing,
				Insert: &User{Id: this.Id, Push: true},
			},
		}

		return runner.Run(ops, bson.NewObjectId(), nil)
	}

	if err := withCollection("reg_tx", &mgo.Safe{}, f); err != nil {
		log.Println(err)
		e := errors.NewError(errors.DbError, err.Error())
		if err == txn.ErrAborted {
			e = errors.NewError(errors.UserExistError)
		}
		return e
	}
	return nil
}

func (this *Account) Update() error {
	change := bson.M{
		"$set": Struct2Map(this),
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Account) UpdateBanTime(banTime int64) error {
	change := bson.M{
		"$set": bson.M{
			"timelimit": banTime,
		},
	}

	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Account) UpdateLevel(score int, level int) error {
	change := bson.M{
		"$set": bson.M{
			"score": score,
			"level": level,
		},
	}
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Account) UpdateLocation(loc Location) error {
	change := bson.M{
		"$set": bson.M{
			"loc": loc,
		},
	}
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Account) SetLogin(count int, lastlog time.Time) (int64, error) {
	this.LoginCount = count
	this.LoginAwards = []int{1, 2, 3, 4, 5, 6, 7}
	change := bson.M{
		"$set": bson.M{
			"lastlogin":    lastlog,
			"login_count":  count,
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

func (this *Account) SetWallet(wallet DbWallet) error {
	change := bson.M{
		"$set": bson.M{
			"wallet": wallet,
		},
	}
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Account) ChangePassword(newPass string) error {
	change := bson.M{
		"$set": bson.M{
			"password": newPass,
		},
	}

	if err := update(accountColl, bson.M{"_id": this.Id}, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
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
		return errors.NewError(errors.DbError, err.Error())
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
	}
	if err := update(accountColl, bson.M{"_id": this.Id}, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
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
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func UserList(skip, limit int) (total int, users []Account, err error) {
	if err := search(accountColl, nil, nil, skip, limit, []string{"-reg_time"}, &total, &users); err != nil {
		return 0, nil, errors.NewError(errors.DbError, err.Error())
	}

	return
}

// This function returns users after preCursor or nextCursor sorted by sortOrder. The return count total should not be more than limit.
func GetUserListBySort(skip, limit int, sortOrder, preCursor, nextCursor string) (total int, users []Account, err error) {
	user := &Account{}
	var query bson.M
	var sortby string

	if len(nextCursor) > 0 {
		user.findOne(bson.M{"_id": nextCursor})
	} else if len(preCursor) > 0 {
		user.findOne(bson.M{"_id": preCursor})
	}

	switch sortOrder {
	case "logintime":
		if len(nextCursor) > 0 {
			query = bson.M{
				"lastlogin": bson.M{
					"$lte": user.LastLogin,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "-lastlogin"
		} else if len(preCursor) > 0 {
			query = bson.M{
				"lastlogin": bson.M{
					"$gte": user.LastLogin,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "lastlogin"
		} else {
			query = bson.M{}
			sortby = "-lastlogin"
		}
		query["reg_time"] = bson.M{
			"$gt": time.Unix(0, 0),
		}

	case "userid":
		if len(nextCursor) > 0 {
			query = bson.M{
				"_id": bson.M{
					"$gt": user.Id,
				},
			}
			sortby = "_id"
		} else if len(preCursor) > 0 {
			query = bson.M{
				"_id": bson.M{
					"$lt": user.Id,
				},
			}
			sortby = "-_id"
		} else {
			query = bson.M{}
			sortby = "_id"
		}
		query["reg_time"] = bson.M{
			"$gt": time.Unix(0, 0),
		}

	case "nickname":
		if len(nextCursor) > 0 {
			query = bson.M{
				"nickname": bson.M{
					"$gte": user.Nickname,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "nickname"
		} else if len(preCursor) > 0 {
			query = bson.M{
				"nickname": bson.M{
					"$lte": user.Nickname,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "-nickname"
		} else {
			query = bson.M{}
			sortby = "nickname"
		}
		query["reg_time"] = bson.M{
			"$gt": time.Unix(0, 0),
		}

	case "score":
		if len(nextCursor) > 0 {
			query = bson.M{
				"score": bson.M{
					"$lte": user.Score,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "-score"
		} else if len(preCursor) > 0 {
			query = bson.M{
				"score": bson.M{
					"$gte": user.Score,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "score"
		} else {
			query = bson.M{}
			sortby = "-score"
		}
		query["reg_time"] = bson.M{
			"$gt": time.Unix(0, 0),
		}

	case "regtime":
		log.Println("regtime")
		fallthrough
	default:
		log.Println("default")
		if len(nextCursor) > 0 {
			query = bson.M{
				"reg_time": bson.M{
					"$lte": user.RegTime,
					"$gt":  time.Unix(0, 0),
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "-reg_time"
		} else if len(preCursor) > 0 {
			query = bson.M{
				"reg_time": bson.M{
					"$gte": user.RegTime,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "reg_time"
		} else {
			query = bson.M{
				"reg_time": bson.M{
					"$gt": time.Unix(0, 0),
				},
			}
			sortby = "-reg_time"
		}
	}

	q := func(c *mgo.Collection) error {
		pq := bson.M{
			"reg_time": bson.M{
				"$gt": time.Unix(0, 0),
			},
		}
		qy := c.Find(pq)

		if total, err = qy.Count(); err != nil {
			return err
		}
		return err
	}

	if err = withCollection(accountColl, nil, q); err != nil {
		return 0, nil, errors.NewError(errors.DbError, err.Error())
	}

	if err = search(accountColl, query, nil, skip, limit, []string{sortby}, nil, &users); err != nil {
		return 0, nil, errors.NewError(errors.DbError, err.Error())
	}

	if len(preCursor) > 0 {
		totalCount := len(users)
		for i := 0; i < totalCount/2; i++ {
			users[i], users[totalCount-1-i] = users[totalCount-1-i], users[i]
		}
	}
	return
}

// This function search users with userid or nickname after preCursor or nextCursor sorted by sortOrder. The return count total should not be more than limit.
func GetSearchListBySort(id, nickname, keywords string, skip, limit int, sortOrder, preCursor, nextCursor string) (total int, users []Account, err error) {
	user := &Account{}
	var query bson.M
	var sortby string
	keyType := -1

	if len(nextCursor) > 0 {
		user.findOne(bson.M{"_id": nextCursor})
	} else if len(preCursor) > 0 {
		user.findOne(bson.M{"_id": preCursor})
	}

	switch sortOrder {
	case "logintime":
		if len(nextCursor) > 0 {
			query = bson.M{
				"lastlogin": bson.M{
					"$lte": user.LastLogin,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "-lastlogin"
		} else if len(preCursor) > 0 {
			query = bson.M{
				"lastlogin": bson.M{
					"$gte": user.LastLogin,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "lastlogin"
		} else {
			query = bson.M{}
			sortby = "-lastlogin"
		}
		query["reg_time"] = bson.M{
			"$gt": time.Unix(0, 0),
		}

	case "userid":
		if len(nextCursor) > 0 {
			query = bson.M{
				"_id": bson.M{
					"$gt": user.Id,
				},
			}
			sortby = "_id"
		} else if len(preCursor) > 0 {
			query = bson.M{
				"_id": bson.M{
					"$lt": user.Id,
				},
			}
			sortby = "-_id"
		} else {
			query = bson.M{}
			sortby = "_id"
		}
		query["reg_time"] = bson.M{
			"$gt": time.Unix(0, 0),
		}

	case "nickname":
		if len(nextCursor) > 0 {
			query = bson.M{
				"nickname": bson.M{
					"$gte": user.Nickname,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "nickname"
		} else if len(preCursor) > 0 {
			query = bson.M{
				"nickname": bson.M{
					"$lte": user.Nickname,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "-nickname"
		} else {
			query = bson.M{}
			sortby = "nickname"
		}
		query["reg_time"] = bson.M{
			"$gt": time.Unix(0, 0),
		}

	case "score":
		if len(nextCursor) > 0 {
			query = bson.M{
				"score": bson.M{
					"$lte": user.Score,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "-score"
		} else if len(preCursor) > 0 {
			query = bson.M{
				"score": bson.M{
					"$gte": user.Score,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "score"
		} else {
			query = bson.M{}
			sortby = "-score"
		}
		query["reg_time"] = bson.M{
			"$gt": time.Unix(0, 0),
		}

	case "regtime":
		log.Println("regtime")
		fallthrough
	default:
		log.Println("default")
		if len(nextCursor) > 0 {
			query = bson.M{
				"reg_time": bson.M{
					"$lte": user.RegTime,
					"$gt":  time.Unix(0, 0),
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "-reg_time"
		} else if len(preCursor) > 0 {
			query = bson.M{
				"reg_time": bson.M{
					"$gte": user.RegTime,
				},
				"_id": bson.M{
					"$ne": user.Id,
				},
			}
			sortby = "reg_time"
		} else {
			query = bson.M{
				"reg_time": bson.M{
					"$gt": time.Unix(0, 0),
				},
			}
			sortby = "-reg_time"
		}
	}

	if len(keywords) > 0 {
		keyType = 0
		query["$or"] = []bson.M{
			bson.M{
				"_id": bson.M{
					"$ne":      user.Id,
					"$regex":   keywords,
					"$options": "i",
				},
			},

			bson.M{
				"nickname": bson.M{
					"$ne":      user.Id,
					"$regex":   keywords,
					"$options": "i",
				},
			},
			bson.M{
				"phone": bson.M{
					"$ne":      user.Id,
					"$regex":   keywords,
					"$options": "i",
				},
			},
			bson.M{
				"about": bson.M{
					"$ne":      user.Id,
					"$regex":   keywords,
					"$options": "i",
				},
			},
			bson.M{
				"hobby": bson.M{
					"$ne":      user.Id,
					"$regex":   keywords,
					"$options": "i",
				},
			},
		}
	} else if len(nickname) > 0 && len(id) > 0 {
		keyType = 0
		query["$or"] = []bson.M{
			bson.M{
				"_id": bson.M{
					"$ne":      user.Id,
					"$regex":   id,
					"$options": "i",
				},
			},

			bson.M{
				"nickname": bson.M{
					"$ne":      user.Id,
					"$regex":   nickname,
					"$options": "i",
				},
			},
		}
	} else if len(nickname) > 0 {
		keyType = 1
		query["nickname"] = bson.M{
			"$ne":      user.Id,
			"$regex":   nickname,
			"$options": "i",
		}
	} else if len(id) > 0 {
		keyType = 2
		query["_id"] = bson.M{
			"$ne":      user.Id,
			"$regex":   id,
			"$options": "i",
		}
	}

	q := func(c *mgo.Collection) error {
		pq := bson.M{
			"reg_time": bson.M{
				"$gt": time.Unix(0, 0),
			},
		}
		if keyType == 0 {
			pq["$or"] = query["$or"]
		} else if keyType == 1 {
			pq["$nickname"] = query["$nickname"]
		} else if keyType == 2 {
			pq["$_id"] = query["$_id"]
		}
		qy := c.Find(pq)

		if total, err = qy.Count(); err != nil {
			return err
		}
		return err
	}

	if err = withCollection(accountColl, nil, q); err != nil {
		return 0, nil, errors.NewError(errors.DbError, err.Error())
	}

	if err = search(accountColl, query, nil, skip, limit, []string{sortby}, nil, &users); err != nil {
		return 0, nil, errors.NewError(errors.DbError, err.Error())
	}

	if len(preCursor) > 0 {
		totalCount := len(users)
		for i := 0; i < totalCount/2; i++ {
			users[i], users[totalCount-1-i] = users[totalCount-1-i], users[i]
		}
	}

	return
}

// This function returns the friends list of the user. Return users after preCursor or nextCursor and sorted by sortOrder.
// The return count total should not be more than limit
func GetFriendsListBySort(skip, limit int, ids []string, sortOrder, preCursor, nextCursor string) (total int, users []Account, err error) {
	user := &Account{}
	var query bson.M
	var sortby string
	var uids []string

	if len(nextCursor) > 0 {
		user.findOne(bson.M{"_id": nextCursor})
	} else if len(preCursor) > 0 {
		user.findOne(bson.M{"_id": preCursor})
	} else {
		user.Id = ""
	}

	for i := 0; i < len(ids); i++ {
		if ids[i] != user.Id {
			uids = append(uids, ids[i])
		}
	}

	switch sortOrder {
	case "logintime":
		if len(nextCursor) > 0 {
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
			query = bson.M{
				"_id": bson.M{
					"$in": ids,
				},
			}
			sortby = "-lastlogin"
		}
		query["reg_time"] = bson.M{
			"$gt": time.Unix(0, 0),
		}

	case "userid":
		if len(nextCursor) > 0 {
			query = bson.M{
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "_id"
		} else if len(preCursor) > 0 {
			query = bson.M{
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "-_id"
		} else {
			query = bson.M{
				"_id": bson.M{
					"$in": ids,
				},
			}
			sortby = "_id"
		}
		query["reg_time"] = bson.M{
			"$gt": time.Unix(0, 0),
		}

	case "nickname":
		if len(nextCursor) > 0 {
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
			query = bson.M{
				"_id": bson.M{
					"$in": ids,
				},
			}
			sortby = "nickname"
		}
		query["reg_time"] = bson.M{
			"$gt": time.Unix(0, 0),
		}

	case "score":
		if len(nextCursor) > 0 {
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
			query = bson.M{
				"_id": bson.M{
					"$in": ids,
				},
			}
			sortby = "-score"
		}
		query["reg_time"] = bson.M{
			"$gt": time.Unix(0, 0),
		}

	case "regtime":
		log.Println("regtime")
		fallthrough
	default:
		log.Println("default")
		if len(nextCursor) > 0 {
			query = bson.M{
				"reg_time": bson.M{
					"$lte": user.RegTime,
					"$gt":  time.Unix(0, 0),
				},
				"_id": bson.M{
					"$in": uids,
				},
			}
			sortby = "-reg_time"
		} else if len(preCursor) > 0 {
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
			query = bson.M{
				"_id": bson.M{
					"$in": ids,
					"reg_time": bson.M{
						"$gt": time.Unix(0, 0),
					},
				},
			}
			sortby = "-reg_time"
		}
	}

	q := func(c *mgo.Collection) error {
		pq := bson.M{
			"reg_time": bson.M{
				"$gt": time.Unix(0, 0),
			},
		}
		qy := c.Find(pq)

		if total, err = qy.Count(); err != nil {
			return err
		}
		return err
	}

	if err = withCollection(accountColl, nil, q); err != nil {
		return 0, nil, errors.NewError(errors.DbError, err.Error())
	}

	if err := search(accountColl, query, nil, skip, limit, []string{sortby}, nil, &users); err != nil {
		return 0, nil, errors.NewError(errors.DbError, err.Error())
	}

	if len(preCursor) > 0 {
		totalCount := len(users)
		for i := 0; i < totalCount/2; i++ {
			users[i], users[totalCount-1-i] = users[totalCount-1-i], users[i]
		}
	}

	return
}

func recordPagingFunc(c *mgo.Collection, first, last string, args ...interface{}) (query bson.M, err error) {
	record := &Record{}

	if bson.IsObjectIdHex(first) {
		if err := c.FindId(bson.ObjectIdHex(first)).One(record); err != nil {
			return nil, err
		}
		query = bson.M{
			"time": bson.M{
				"$gte": record.PubTime,
			},
			"_id": bson.M{
				"$ne": record.Id,
			},
		}
	} else if bson.IsObjectIdHex(last) {
		if err := c.FindId(bson.ObjectIdHex(last)).One(record); err != nil {
			return nil, err
		}
		query = bson.M{
			"time": bson.M{
				"$lte": record.PubTime,
			},
			"_id": bson.M{
				"$ne": record.Id,
			},
		}
	}

	return
}

func (this *Account) Records(paging *Paging) (int, []Record, error) {
	var records []Record
	total := 0

	pageUp := false
	sortFields := []string{"-time"}
	if len(paging.First) > 0 {
		pageUp = true
		sortFields = []string{"time"}
	}

	if err := psearch(recordColl, bson.M{"uid": this.Id}, nil,
		sortFields, nil, &records, recordPagingFunc, paging); err != nil {
		e := errors.NewError(errors.DbError, err.Error())
		if err == mgo.ErrNotFound {
			e = errors.NewError(errors.NotFoundError, err.Error())
		}
		return total, nil, e
	}

	paging.First = ""
	paging.Last = ""
	paging.Count = 0
	if len(records) > 0 {
		if pageUp {
			for i := 0; i < len(records)/2; i++ {
				t := records[i]
				records[i] = records[len(records)-i-1]
				records[len(records)-i-1] = t
			}
		}
		paging.First = records[0].Id.Hex()
		paging.Last = records[len(records)-1].Id.Hex()
		paging.Count = total
	}

	return total, records, nil
}

func (this *Account) UpdateAction(action string, date time.Time) (bool, error) {
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
			"score": bson.M{
				"$gte": user.Score,
			},
			"_id": bson.M{
				"$in": ids,
				"$ne": user.Id,
			},
		}
	} else if len(last) > 0 {
		if err := c.FindId(last).One(user); err != nil {
			return nil, err
		}
		query = bson.M{
			"score": bson.M{
				"$lte": user.Score,
			},
			"_id": bson.M{
				"$in": ids,
				"$ne": user.Id,
			},
		}
	}
	return
}

func UserCount() (count int) {
	search(accountColl, bson.M{"reg_time": bson.M{"$gt": time.Unix(0, 0)}}, nil, 0, 0, nil, &count, nil)
	return
}

func Users(ids []string, paging *Paging) ([]Account, error) {
	var users []Account
	total := 0

	pageUp := false
	sortFields := []string{"-score"}
	if len(paging.First) > 0 {
		pageUp = true
		sortFields = []string{"score"}
	}

	if err := psearch(accountColl, nil, nil, sortFields, nil, &users, friendsPagingFunc, paging, ids); err != nil {
		e := errors.NewError(errors.DbError, err.Error())
		if err == mgo.ErrNotFound {
			e = errors.NewError(errors.NotFoundError, err.Error())
		}
		return nil, e
	}

	paging.First = ""
	paging.Last = ""
	paging.Count = 0
	if len(users) > 0 {
		if pageUp {
			for i := 0; i < len(users)/2; i++ {
				t := users[i]
				users[i] = users[len(users)-i-1]
				users[len(users)-i-1] = t
			}
		}
		paging.First = users[0].Id
		paging.Last = users[len(users)-1].Id
		paging.Count = total
	}

	return users, nil
}

func (this *Account) ArticleCount() (count int) {
	query := bson.M{"author": this.Id, "parent": nil}
	search(accountColl, query, nil, 0, 0, nil, &count, nil)
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
			"_id": bson.M{
				"$ne": user.Id,
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
			"_id": bson.M{
				"$ne": user.Id,
			},
		}
	}

	return
}

func Search(nickname string, paging *Paging) ([]Account, error) {
	var users []Account
	total := 0

	query := bson.M{
		"reg_time": bson.M{
			"$gt": time.Unix(0, 0),
		},
	}

	if len(nickname) > 0 {
		query["nickname"] = bson.M{
			"$regex":   nickname,
			"$options": "i",
		}
	}

	pageUp := false
	sortFields := []string{"-lastlogin"}
	if len(paging.First) > 0 {
		pageUp = true
		sortFields = []string{"lastlogin"}
	}

	if err := psearch(accountColl, query, nil, sortFields, nil, &users, searchPagingFunc, paging); err != nil {
		if err != mgo.ErrNotFound {
			return nil, errors.NewError(errors.DbError, err.Error())
		}
	}

	paging.First = ""
	paging.Last = ""
	paging.Count = 0
	if len(users) > 0 {
		if pageUp {
			for i := 0; i < len(users)/2; i++ {
				t := users[i]
				users[i] = users[len(users)-i-1]
				users[len(users)-i-1] = t
			}
		}

		paging.First = users[0].Id
		paging.Last = users[len(users)-1].Id
		paging.Count = total
	}
	return users, nil
}

func (this *Account) SearchNear(paging *Paging) ([]Account, error) {
	var users []Account
	total := 0
	if this.Loc == nil || this.Loc.Lat == 0 || this.Loc.Lng == 0 {
		return nil, nil
	}
	query := bson.M{
		"loc": bson.M{
			"$near": []float64{this.Loc.Lat, this.Loc.Lng},
		},
	}

	if err := psearch(accountColl, query, nil, nil, nil, &users, nil, nil); err != nil {
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

func (this *Account) ClearEvent(eventType string, eventId string) int {
	info, err := removeAll(eventColl, bson.M{"data.to": this.Id, "data.type": eventType, "data.id": eventId}, true)
	if err != nil {
		return 0
	}
	return info.Removed
}

func (this *Account) UpdateInfo(change bson.M) error {
	if err := updateId(accountColl, this.Id, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}
