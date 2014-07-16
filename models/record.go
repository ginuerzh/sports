// record
package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo/bson"
	"time"
)

type Record struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Uid      string
	Type     string
	Time     time.Time
	Duration int64
	Distance int
	PubTime  time.Time `bson:"pub_time"`
}

func init() {
	ensureIndex(recordColl, "uid")
	ensureIndex(recordColl, "-time")
	ensureIndex(recordColl, "-pub_time")
}

func (this *Record) findOne(query interface{}) (bool, error) {
	var records []Record

	err := search(recordColl, query, nil, 0, 1, nil, nil, &records)
	if err != nil {
		return false, errors.NewError(errors.DbError, err.Error())
	}
	if len(records) > 0 {
		*this = records[0]
	}

	return len(records) > 0, nil
}

func TotalRecords(userid string) (int, error) {
	total := 0
	err := search(recordColl, bson.M{"uid": userid}, nil, 0, 0, nil, &total, nil)
	return total, err
}

func MaxDistanceRecord(userid string) (rec *Record, err error) {
	var records []Record
	err = search(recordColl, bson.M{"uid": userid}, nil, 0, 1, []string{"-distance"}, nil, &records)
	if len(records) > 0 {
		rec = &records[0]
	}
	return
}

func (this *Record) Save() error {
	this.Id = bson.NewObjectId()
	if err := save(recordColl, this, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}
