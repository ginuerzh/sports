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

func (this *Message) RemoveId() error {
	if err := removeId(msgColl, this.Id, true); err != nil {
		if e, ok := err.(*mgo.LastError); ok {
			return errors.NewError(errors.DbError, e.Error())
		}
	}
	return nil
}

func (this *Message) Delete(from, to string, start, end time.Time) (count int, err error) {
	selector := bson.M{
		"$or": []bson.M{
			bson.M{"from": from, "to": to},
			bson.M{"from": to, "to": from},
		},
		"time": bson.M{
			"$lte": end,
			"$gte": start,
		},
	}
	info, err := removeAll(msgColl, selector, true)
	if info != nil {
		count = info.Removed
	}

	return
}

func msgPagingFunc(c *mgo.Collection, first, last string, args ...interface{}) (query bson.M, err error) {
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

func AdminMessages(from, to string, pageIndex, pageCount int) (total int, msgs []Message, err error) {
	var or []bson.M

	if len(from) > 0 {
		or = append(or, bson.M{"from": from})
	}
	if len(to) > 0 {
		or = append(or, bson.M{"to": to})
	}
	err = search(msgColl, bson.M{"type": "chat", "$or": or}, nil,
		pageIndex*pageCount, pageCount, []string{"-time"}, &total, &msgs)

	return
}
