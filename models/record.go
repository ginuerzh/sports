// record
package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

func init() {

}

const (
	StatusNormal   = "NORMAL"
	StatusFinish   = "FINISH"
	StatusUnFinish = "UNFINISH"
	StatusAuth     = "AUTHENTICATION"
)

type SportRecord struct {
	Source   string
	Duration int64
	Distance int
	Speed    float64
	Pics     []string
	Review   string
}

type GameRecord struct {
	Type  string
	Name  string
	Score int
	Magic int
	Coin  int64
}

type Record struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Uid       string
	Task      int64
	Status    string
	Type      string
	StartTime time.Time
	EndTime   time.Time
	Sport     *SportRecord `bson:",omitempty"`
	Game      *GameRecord  `bson:",omitempty"`
	Coin      int64
	PubTime   time.Time `bson:"pub_time"`
}

func (this *Record) findOne(query interface{}) (bool, error) {
	var records []Record

	err := search(recordColl, query, nil, 0, 1, nil, nil, &records)
	if err != nil {
		return false, errors.NewError(errors.DbError)
	}
	if len(records) > 0 {
		*this = records[0]
	}

	return len(records) > 0, nil
}
func (this *Record) FindByTask(tid int64) (bool, error) {
	return this.findOne(bson.M{"uid": this.Uid, "task": tid})
}

func (this *Record) SetStatus(status string, review string, coin int64) error {
	query := bson.M{"uid": this.Uid, "task": this.Task}
	change := bson.M{
		"$set": bson.M{
			"status":       status,
			"sport.review": review,
			"coin":         coin,
		},
	}
	if err := update(recordColl, query, change, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

func TotalRecords(userid string) (int, error) {
	total := 0
	err := search(recordColl, bson.M{"uid": userid}, nil, 0, 0, nil, &total, nil)
	return total, err
}

func TaskRecords(pageIndex, pageCount int) (int, []Record, error) {
	var records []Record
	total := 0

	if err := search(recordColl, bson.M{"type": "run"}, nil,
		pageIndex*pageCount, pageCount, []string{"-pub_time"}, &total, &records); err != nil && err != mgo.ErrNotFound {
		return total, nil, errors.NewError(errors.DbError)
	}

	return total, records, nil
}

func SearchTaskByUserid(userid string, finish bool, pageIndex, pageCount int) (int, []Record, error) {
	var records []Record
	total := 0
	if len(userid) == 0 {
		return total, records, nil
	}

	query := bson.M{
		"type":   "run",
		"uid":    userid,
		"status": bson.M{"$in": []interface{}{StatusFinish, StatusUnFinish}},
	}
	if !finish {
		query["status"] = StatusAuth
	}
	if err := search(recordColl, query, nil,
		pageIndex*pageCount, pageCount, []string{"-pub_time"}, &total, &records); err != nil && err != mgo.ErrNotFound {
		return total, nil, errors.NewError(errors.DbError)
	}

	return total, records, nil
}

func MaxDistanceRecord(userid string) (*Record, error) {
	record := &Record{}
	err := findOne(recordColl, bson.M{"uid": userid, "status": StatusFinish}, []string{"-sport.distance"}, record)

	return record, err
}

func MaxSpeedRecord(userid string) (*Record, error) {
	record := &Record{}
	err := findOne(recordColl, bson.M{"uid": userid, "status": StatusFinish}, []string{"-sport.speed"}, record)

	return record, err
}

func (this *Record) Save() error {
	this.Id = bson.NewObjectId()
	if err := save(recordColl, this, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

func (this *Record) Delete() error {
	if err := remove(recordColl, bson.M{"uid": this.Uid, "task": this.Task}, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

// This function returns records of type recType between fromTime and toTime, at same time if nextCursor or  preCursor is not nil, the records should after
// the cursor. The count is the max count it returns this time.
func GetRecords(id, recType string, nextCursor, preCursor string, count int, fromTime, toTime int64, skip, limit int) (int, []Record, error) {
	var records []Record
	total := 0

	ft := time.Unix(0, 0)
	if fromTime > 0 {
		ft = time.Unix(fromTime, 0)
	}

	tt := time.Now()
	if toTime > 0 {
		tt = time.Unix(toTime, 0)
	}

	sortby := "-pub_time"

	var pc, nc bson.ObjectId
	var pcValid, ncValid bool
	if len(nextCursor) > 0 {
		if bson.IsObjectIdHex(nextCursor) {
			nc = bson.ObjectIdHex(nextCursor)
			ncValid = true
		}
	}

	if len(preCursor) > 0 {
		if bson.IsObjectIdHex(preCursor) {
			pc = bson.ObjectIdHex(preCursor)
			pcValid = true
			sortby = "pub_time"
		}
	}

	var query bson.M
	if len(recType) > 0 {
		if ncValid {
			query = bson.M{
				"_id": bson.M{
					"$ne": nc,
				},
				"uid":  id,
				"type": recType,
				"pub_time": bson.M{
					"$gt": ft,
					"$lt": tt,
				},
			}
		} else if pcValid {
			query = bson.M{
				"_id": bson.M{
					"$ne": pc,
				},
				"uid":  id,
				"type": recType,
				"pub_time": bson.M{
					"$gt": ft,
					"$lt": tt,
				},
			}

		} else {
			query = bson.M{
				"uid":  id,
				"type": recType,
				"pub_time": bson.M{
					"$gt": ft,
					"$lt": tt,
				},
			}

		}
	} else {
		if ncValid {
			query = bson.M{
				"_id": bson.M{
					"$ne": nc,
				},
				"uid": id,
				"pub_time": bson.M{
					"$gt": ft,
					"$lt": tt,
				},
			}
		} else if pcValid {
			query = bson.M{
				"_id": bson.M{
					"$ne": pc,
				},
				"uid": id,
				"pub_time": bson.M{
					"$gt": ft,
					"$lt": tt,
				},
			}

		} else {
			query = bson.M{
				"uid": id,
				"pub_time": bson.M{
					"$gt": ft,
					"$lt": tt,
				},
			}

		}
	}

	var err error
	q := func(c *mgo.Collection) error {
		pq := bson.M{
			"uid": id}
		qy := c.Find(pq)

		if total, err = qy.Count(); err != nil {
			return err
		}
		return err
	}

	if err = withCollection(recordColl, nil, q); err != nil {
		return 0, nil, errors.NewError(errors.DbError)
	}

	if err = search(recordColl, query, nil, skip, limit, []string{sortby}, nil, &records); err != nil {
		return 0, nil, errors.NewError(errors.DbError)
	}

	if pcValid {
		totalCount := len(records)
		for i := 0; i < totalCount/2; i++ {
			records[i], records[totalCount-1-i] = records[totalCount-1-i], records[i]
		}
	}
	return total, records, nil
}

// This function removes the recType records between fromTime and toTime of  user "id".
func RemoveRecordsByID(id, recType string, fromTime, toTime int64) (int, error) {
	total := 0
	var records []Record

	ft := time.Unix(0, 0)
	if fromTime > 0 {
		ft = time.Unix(fromTime, 0)
	}

	tt := time.Now()
	if toTime > 0 {
		tt = time.Unix(toTime, 0)
	}
	var rm bson.M
	if len(recType) > 0 {
		rm = bson.M{
			"uid":  id,
			"type": recType,
			"pub_time": bson.M{
				"$gt": ft,
				"$lt": tt,
			},
		}
	} else {
		rm = bson.M{
			"uid": id,
			"pub_time": bson.M{
				"$gt": ft,
				"$lt": tt,
			},
		}
	}
	err := search(recordColl, rm, nil, 0, 0, nil, &total, &records)
	if err != nil {
		return 0, err
	}
	if total == 0 {
		return 0, nil
	}
	err = remove(recordColl, rm, true)
	if err != nil {
		return 0, err
	}
	return total, err
}
