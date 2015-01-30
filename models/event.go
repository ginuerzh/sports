// Event
package models

import (
	"encoding/json"
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	//"time"
)

const (
	EventMsg     = "message"
	EventArticle = "article"
	EventWallet  = "wallet"
	EventStatus  = "status"

	EventChat    = "chat"
	EventGChat   = "groupchat"
	EventSub     = "subscribe"
	EventUnsub   = "unsubscribe"
	EventThumb   = "thumb"
	EventComment = "comment"
	EventTx      = "tx"
	EventReward  = "reward"
	EventBan     = "ban"
	EventUnban   = "unban"
	EventLock    = "lock"
	EventTask    = "task"
)

func init() {
	ensureIndex(eventColl, "-time")
}

type EventData struct {
	Type string    `json:"type"`
	Id   string    `bson:"pid" json:"pid"`
	From string    `json:"from"`
	To   string    `json:"to"`
	Body []MsgBody `json:"body"`
}

type Event struct {
	Id   bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Type string        `json:"type"`
	Data EventData     `bson:"push" json:"push"`
	Time int64         `json:"time"`
}

func (e *Event) Bytes() []byte {
	b, _ := json.Marshal(e)
	return b
}

func (e *Event) Save() error {
	e.Id = bson.NewObjectId()
	if err := save(eventColl, e, true); err != nil {
		log.Println(err)
		return errors.NewError(errors.DbError, err.(*mgo.LastError).Error())
	}

	return nil
}

func (this *Event) Upsert() error {
	query := bson.M{
		"push.type": this.Data.Type,
		"push.pid":  this.Data.Id,
		"push.from": this.Data.From,
		"push.to":   this.Data.To,
	}
	//log.Println("event upsert", query, Struct2Map(this))
	_, err := upsert(eventColl, query, Struct2Map(this), true)
	return err
}

func (this *Event) Delete() int {
	info, err := removeAll(eventColl,
		bson.M{
			"push.type": this.Data.Type,
			"push.pid":  this.Data.Id,
			"push.from": this.Data.From,
			"push.to":   this.Data.To,
		},
		true)

	if err != nil {
		return 0
	}
	return info.Removed
}

func (this *Event) Clear() int {
	info, err := removeAll(eventColl,
		bson.M{
			"push.type": this.Data.Type,
			"push.pid":  this.Data.Id,
			"push.to":   this.Data.To,
		},
		true)

	if err != nil {
		return 0
	}
	return info.Removed
}

func Events(userid string) (events []Event, err error) {
	err = search(eventColl, bson.M{"push.to": userid}, nil, 0, 0, []string{"-time"}, nil, &events)
	return
}

func EventCount(typ string, id string, to string) int {
	n, _ := count(eventColl, bson.M{"push.type": typ, "push.pid": id, "push.to": to})
	return n
}
