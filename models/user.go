// user
package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//"labix.org/v2/mgo/txn"
	"log"
	"time"
)

const (
	TaskRunning = "PHYSIQUE"
	TaskPost    = "LITERATURE"
	TaskGame    = "MAGIC"
)

const (
	TaskCompleted = iota
	TaskUncompleted
)

type Contact struct {
	Id       string
	Profile  string
	Nickname string
	Count    int
	Last     *Message `bson:",omitempty"`
}

/*
type Event struct {
	Id      string
	Thumbs  []string `bson:",omitempty"`
	Reviews []string `bson:",omitempty"`
}
*/

type Proof struct {
	Tid    int
	Pics   []string
	Result string `bson:",omitempty"`
}

type TaskList struct {
	Completed   []int
	Uncompleted []int
	Waited      []int
	Proofs      []Proof
	Last        time.Time
}

type User struct {
	Id string `bson:"_id"`

	Contacts []Contact `bson:",omitempty"`
	//Events   []Event   `bson:",omitempty"`
	Devs  []string `bson:",omitempty"`
	Push  bool
	Tasks TaskList
}

func (this *User) findOne(query interface{}) (bool, error) {
	var users []User

	err := search(userColl, query, nil, 0, 1, nil, nil, &users)
	if err != nil {
		return false, errors.NewError(errors.DbError, err.Error())
	}
	if len(users) > 0 {
		*this = users[0]
	}

	return len(users) > 0, nil
}

func (this *User) FindByUserid(userid string) (bool, error) {
	return this.findOne(bson.M{"_id": userid})
}

func (this *User) Articles(typ string, paging *Paging) (int, []Article, error) {
	var articles []Article
	total := 0
	var query bson.M
	switch typ {
	case "COMMENTS":
		query = bson.M{"author": this.Id, "parent": bson.M{"$ne": nil}}
	case "ARTICLES":
		query = bson.M{"author": this.Id, "parent": nil}
	default:
		query = bson.M{"author": this.Id}
	}

	if err := psearch(articleColl, query, nil, []string{"-pub_time"}, nil, &articles,
		articlePagingFunc, paging); err != nil {
		e := errors.NewError(errors.DbError, err.Error())
		if err == mgo.ErrNotFound {
			e = errors.NewError(errors.NotFoundError, err.Error())
		}
		return total, nil, e
	}

	if len(articles) > 0 {
		paging.First = articles[0].Id.Hex()
		paging.Last = articles[len(articles)-1].Id.Hex()
		paging.Count = total
	}

	return total, articles, nil
}

/*
func (this *User) LastMessage(userid string) *Message {
	_, msgs, _ := this.Messages(userid, &Paging{Count: 1})
	if len(msgs) == 0 {
		return nil
	}
	return &msgs[0]
}
*/

func (this *User) Messages(userid string, paging *Paging) (int, []Message, error) {
	var msgs []Message
	total := 0
	query := bson.M{
		"$or": []bson.M{
			bson.M{"from": userid, "to": this.Id},
			bson.M{"from": this.Id, "to": userid},
		},
	}

	if err := psearch(msgColl, query, nil, []string{"-time"}, nil, &msgs,
		msgPagingFunc, paging); err != nil {
		e := errors.NewError(errors.DbError, err.Error())
		if err == mgo.ErrNotFound {
			e = errors.NewError(errors.NotFoundError, err.Error())
		}
		return total, nil, e
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

/*
func (this *User) AddEvent(event *Event) error {
	var change bson.M
	selector := bson.M{
		"_id":       this.Id,
		"events.id": event.Id,
	}
	if len(event.Thumbs) > 0 {
		change = bson.M{
			"$addToSet": bson.M{
				"events.$.thumbs": bson.M{
					"$each": event.Thumbs,
				},
			},
		}
	}
	if len(event.Reviews) > 0 {
		change = bson.M{
			"$addToSet": bson.M{
				"events.$.reviews": bson.M{
					"$each": event.Reviews,
				},
			},
		}
	}

	err := update(userColl, selector, change, true)
	//log.Println(err)
	if err == nil {
		return nil
	}

	if err != mgo.ErrNotFound {
		return errors.NewError(errors.DbError, err.Error())
	}

	// not found
	selector = bson.M{
		"_id": this.Id,
	}
	change = bson.M{
		"$push": bson.M{
			"events": event,
		},
	}
	err = update(userColl, selector, change, true)
	if err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}

	return nil
}
*/
func (this *User) AddContact(contact *Contact) error {
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
	err := update(userColl, selector, change, true)
	if err == nil {
		return nil
	}
	log.Println(err)
	if err != mgo.ErrNotFound {
		return errors.NewError(errors.DbError, err.Error())
	}

	// not found
	selector = bson.M{
		"_id": this.Id,
	}
	change = bson.M{
		"$push": bson.M{
			"contacts": contact,
		},
	}
	err = update(userColl, selector, change, true)
	if err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *User) MarkRead(typ, id string) error {
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

	if err := update(userColl, selector, change, true); err != nil {
		if err != mgo.ErrNotFound {
			return errors.NewError(errors.DbError, err.Error())
		}
	}

	return nil
}

func (this *User) SetPush(push bool) error {
	selector := bson.M{
		"_id": this.Id,
	}
	change := bson.M{
		"$set": bson.M{
			"push": push,
		},
	}
	if err := update(userColl, selector, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *User) PushEnabled() (bool, error) {
	var users []User
	var enabled bool

	err := search(userColl, bson.M{"_id": this.Id}, bson.M{"push": true},
		0, 1, nil, nil, &users)
	if err != nil {
		return false, errors.NewError(errors.DbError, err.Error())
	}

	if len(users) > 0 {
		enabled = users[0].Push
	}

	return enabled, nil
}

func (this *User) Devices() ([]string, bool, error) {
	var users []User
	var devs []string
	var enabled bool

	err := search(userColl, bson.M{"_id": this.Id}, bson.M{"devs": true, "push": true},
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

func (this *User) AddDevice(dev string) error {
	selector := bson.M{
		"_id": this.Id,
	}
	change := bson.M{
		"$addToSet": bson.M{
			"devs": dev,
		},
	}
	if err := update(userColl, selector, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *User) RmDevice(dev string) error {
	selector := bson.M{
		"_id": this.Id,
	}
	change := bson.M{
		"$pull": bson.M{
			"devs": dev,
		},
	}
	if err := update(userColl, selector, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *User) GetTasks() (TaskList, error) {
	_, err := this.FindByUserid(this.Id)
	return this.Tasks, err
}

func (this *User) AddTask(typ string, tid int, proofs []string) error {
	selector := bson.M{
		"_id": this.Id,
	}

	update(userColl, selector, bson.M{"$pull": bson.M{"tasks.proofs": bson.M{"tid": tid}}}, true)

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

	if err := update(userColl, selector, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *User) SetTaskComplete(tid int) error {
	selector := bson.M{
		"_id": this.Id,
	}

	change := bson.M{
		"$pull": bson.M{
			"tasks.waited": tid,
		},
		"$addToSet": bson.M{
			"tasks.completed": tid,
		},
	}

	if err := update(userColl, selector, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}
