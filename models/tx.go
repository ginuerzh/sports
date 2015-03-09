package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo/bson"
	"time"
)

type Tx struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Uid   string
	Coins int64
	Value int
	Time  time.Time
}

func (this *Tx) Save() error {
	this.Id = bson.NewObjectId()
	if err := save(txColl, this, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}
