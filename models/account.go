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
	ensureIndex(accountColl, "nickname")
	ensureIndex(accountColl, "-reg_time")
}

type UserInfo struct {
	Hobby    string `json:"hobby"`
	Height   int    `json:"height"`
	Weight   int    `json:"weight"`
	Birth    int64  `json:"birthday"`
	Actor    string `json:"actor"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone_number"`
	About    string `json:"about"`
	Nickname string `json:"nikename"`

	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Area     string `json:"area"`
	LocDesc  string `json:"location_desc"`
	Lng      float64
	Lat      float64
}

type Account struct {
	Id       string    `bson:"_id,omitempty"`
	Nickname string    `bson:",omitempty"`
	Password string    `bson:",omitempty"`
	Profile  string    `bson:",omitempty"`
	RegTime  time.Time `bson:"reg_time,omitempty"`
	Role     string    `bson:",omitempty"`
	Hobby    string    `bson:",omitempty"`
	Height   int       `bson:",omitempty"`
	Weight   int       `bson:",omitempty"`
	Birth    time.Time `bson:",omitempty"`
	Actor    string    `bson:",omitempty"`
	Gender   string    `bson:",omitempty"`
	Url      string    `bson:",omitempty"`
	Phone    string    `bson:",omitempty"`
	About    string    `bson:",omitempty"`
	Addr     Address
	Loc      Location
	Setinfo  bool
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

func (this *Account) SetInfo(info UserInfo) error {
	mset := bson.M{
		"setinfo": true,
	}

	if len(info.About) > 0 {
		mset["about"] = info.About
	}
	if len(info.Actor) > 0 {
		mset["actor"] = info.Actor
	}
	if info.Birth > 0 {
		mset["birth"] = time.Unix(info.Birth, 0)
	}
	if len(info.Gender) > 0 {
		mset["gender"] = info.Gender
	}
	if info.Height > 0 {
		mset["height"] = info.Height
	}
	if len(info.Hobby) > 0 {
		mset["hobby"] = info.Hobby
	}
	if len(info.Nickname) > 0 {
		mset["nickname"] = info.Nickname
	}
	if len(info.Phone) > 0 {
		mset["phone"] = info.Phone
	}
	if info.Weight > 0 {
		mset["weight"] = info.Weight
	}

	change := bson.M{
		"$set": mset,
	}

	if err := update(accountColl, bson.M{"_id": this.Id}, change, true); err != nil {
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

func UserList(skip, limit int) (total int, users []Account, err error) {
	if err := search(accountColl, nil, nil, skip, limit, []string{"-reg_time"}, &total, &users); err != nil {
		return 0, nil, errors.NewError(errors.DbError, err.Error())
	}

	return
}

func recordPagingFunc(c *mgo.Collection, first, last string) (query bson.M, err error) {
	record := &Record{}

	if bson.IsObjectIdHex(first) {
		if err := c.FindId(bson.ObjectIdHex(first)).One(record); err != nil {
			return nil, err
		}
		query = bson.M{
			"time": bson.M{
				"$gte": record.Time,
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
				"$lte": record.Time,
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

	if err := psearch(recordColl, bson.M{"uid": this.Id}, nil,
		[]string{"-time"}, nil, &records, recordPagingFunc, paging); err != nil {
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
