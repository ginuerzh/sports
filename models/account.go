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

type Account struct {
	Id       string    `bson:"_id"`
	Password string    `bson:",omitempty"`
	Nickname string    `bson:",omitempty"`
	Gender   string    `bson:",omitempty"`
	Url      string    `bson:",omitempty"`
	Phone    string    `bson:",omitempty"`
	About    string    `bson:",omitempty"`
	Location string    `bson:",omitempty"`
	Profile  string    `bson:",omitempty"`
	RegTime  time.Time `bson:"reg_time"`
	Role     string    `bson:",omitempty"`
}

func (this *Account) Exists() (bool, error) {
	return this.findOne(bson.M{"_id": this.Id})
}

/*
func FindUsers(ids []string) ([]Account, error) {
	var users []Account
	if err := search(accountColl, bson.M{"_id": bson.M{"$in": ids}}, nil, 0, 0, nil, nil, &users); err != nil {
		return nil, errors.NewError(errors.DbError, err.Error())
	}

	return users, nil
}
*/
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
