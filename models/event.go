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
	Id   string    `json:"pid"`
	From string    `json:"from"`
	To   string    `json:"to"`
	Body []MsgBody `json:"body"`
}

type Event struct {
	Id   bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Type string        `json:"type"`
	Data EventData     `json:"push"`
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
		"data.type": this.Data.Type,
		"data.id":   this.Data.Id,
		"data.from": this.Data.From,
		"data.to":   this.Data.To,
	}
	return upsert(eventColl, query, Struct2Map(this), true)
}

func Events(userid string) (events []Event, err error) {
	err = search(eventColl, bson.M{"data.to": userid}, nil, 0, 0, []string{"-time"}, nil, &events)
	return
}

func EventCount(typ string, id string) (count int) {
	search(eventColl, bson.M{"data.type": typ, "data.id": id}, nil, 0, 0, nil, &count, nil)
	return
}
