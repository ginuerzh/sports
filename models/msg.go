// msg
package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//"labix.org/v2/mgo/txn"
	//"log"
	"time"
)

func init() {
	ensureIndex(msgColl, "from")
	ensureIndex(msgColl, "to")
	ensureIndex(msgColl, "from", "to")
	ensureIndex(msgColl, "-time")
}

type MsgBody struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Message struct {
	Id   bson.ObjectId `bson:"_id,omitempty"`
	From string
	To   string
	Type string
	Body []MsgBody
	Time time.Time
}

func (this *Message) findOne(query interface{}) (bool, error) {
	var msgs []Message

	err := search(msgColl, query, nil, 0, 1, nil, nil, &msgs)
	if err != nil {
		return false, errors.NewError(errors.DbError)
	}
	if len(msgs) > 0 {
		*this = msgs[0]
	}

	return len(msgs) > 0, nil
}

func (this *Message) Last(from string) error {
	var msgs []Message

	err := search(msgColl, bson.M{"from": from}, nil, 0, 1, []string{"-time"}, nil, &msgs)
	if err != nil {
		return errors.NewError(errors.DbError)
	}

	if len(msgs) > 0 {
		*this = msgs[0]
	}
	return nil
}

func (this *Message) Save() error {
	this.Id = bson.NewObjectId()
	if err := save(msgColl, this, true); err != nil {
		return errors.NewError(errors.DbError, err.(*mgo.LastError).Error())
	}

	return nil
}

func msgPagingFunc(c *mgo.Collection, first, last string) (query bson.M, err error) {
	msg := &Message{}

	if bson.IsObjectIdHex(first) {
		if err := c.FindId(bson.ObjectIdHex(first)).One(msg); err != nil {
			return nil, err
		}
		query = bson.M{
			"time": bson.M{
				"$gte": msg.Time,
			},
			"_id": bson.M{
				"$ne": msg.Id,
			},
		}
	} else if bson.IsObjectIdHex(last) {
		if err := c.FindId(bson.ObjectIdHex(last)).One(msg); err != nil {
			return nil, err
		}
		query = bson.M{
			"time": bson.M{
				"$lte": msg.Time,
			},
			"_id": bson.M{
				"$ne": msg.Id,
			},
		}
	}

	return
}
