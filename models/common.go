// common
package models

import (
	"github.com/nf/geocode"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

const (
	AccessRate = 1 << iota // 001
	ThumbRate              // 010
	ReviewRate             // 100

	AccessRateMask = 6 // 110
	ThumbRateMask  = 5 // 101
	ReviewRateMask = 3 // 011

	DefaultPageSize = 10
	TimeFormat      = "2006-01-02 15:04:05"
)

var (
	mgoSession   *mgo.Session
	databaseName = "sports"
	accountColl  = "accounts"
	userColl     = "users"
	articleColl  = "articles"
	msgColl      = "messages"
	//reviewColl   = "reviews"
	fileColl   = "files"
	recordColl = "records"
	actionColl = "actions"
	groupColl  = "groups"
	//eventColl    = "events"
	//rateColl     = "rates"
)

var (
	GuestUserPrefix = "guest:"
)

type Paging struct {
	First string `form:"page_frist_id" json:"page_frist_id"`
	Last  string `form:"page_last_id" json:"page_last_id"`
	Count int    `form:"page_item_count" json:"-"`
}

type Address struct {
	Country  string `json:"country,omitempty"`
	Province string `json:"province,omitempty"`
	City     string `json:"city,omitempty"`
	Area     string `json:"area,omitempty"`
	Desc     string `json:"location_desc"`
}

func (addr Address) String() string {
	return addr.Country + addr.Province + addr.City + addr.Area + addr.Desc
}

type Location struct {
	Lat float64 `json:"latitude,string"`
	Lng float64 `json:"longitude,string"`
}

func Addr2Loc(addr Address) Location {
	if len(addr.String()) == 0 {
		return Location{}
	}
	req := &geocode.Request{
		Region:   "cn",
		Provider: geocode.GOOGLE,
		Address:  addr.String(),
	}
	resp, err := req.Lookup(nil)
	if err != nil || resp.Status != "OK" || len(resp.Results) == 0 {
		return Location{}
	}

	return Location{
		Lat: resp.Results[0].Geometry.Location.Lat,
		Lng: resp.Results[0].Geometry.Location.Lng,
	}
}

type PagingFunc func(c *mgo.Collection, first, last string) (query bson.M, err error)

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial("localhost")
		if err != nil {
			log.Println(err) // no, not really
		}
	}
	return mgoSession.Clone()
}

func withCollection(collection string, safe *mgo.Safe, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()

	session.SetSafe(safe)
	c := session.DB(databaseName).C(collection)
	return s(c)
}

func exists(collection string, query interface{}) (bool, error) {
	b := false
	q := func(c *mgo.Collection) error {
		n, err := c.Find(query).Count()
		b = n > 0
		return err
	}

	err := withCollection(collection, nil, q)
	return b, err
}

// search with paging
func psearch(collection string, query, selector interface{}, sortFields []string,
	total *int, result interface{}, pagingFunc PagingFunc, paging *Paging) (err error) {
	q := func(c *mgo.Collection) error {
		var pquery bson.M
		if pagingFunc != nil {
			if paging == nil {
				paging = &Paging{}
			}
			pquery, err = pagingFunc(c, paging.First, paging.Last)
			if err != nil {
				return err
			}
		}

		qy := c.Find(query)

		if total != nil {
			if *total, err = qy.Count(); err != nil {
				return err
			}
		}
		if result == nil {
			return nil
		}

		if paging.Count == 0 {
			paging.Count = DefaultPageSize
		}

		if pquery == nil {
			pquery = bson.M{}
		}
		if m, ok := query.(bson.M); ok {
			for k, v := range m {
				pquery[k] = v
			}
		}

		return c.Find(pquery).Select(selector).Sort(sortFields...).Limit(paging.Count).All(result)
	}
	return withCollection(collection, nil, q)
}

func search(collection string, query interface{}, selector interface{},
	skip, limit int, sortFields []string, total *int, result interface{}) error {

	q := func(c *mgo.Collection) error {
		qy := c.Find(query)
		var err error

		if selector != nil {
			qy = qy.Select(selector)
		}

		if total != nil {
			if *total, err = qy.Count(); err != nil {
				return err
			}
		}

		if result == nil {
			return err
		}

		if limit > 0 {
			qy = qy.Limit(limit)
		}
		if skip > 0 {
			qy = qy.Skip(skip)
		}
		if len(sortFields) > 0 {
			qy = qy.Sort(sortFields...)
		}

		return qy.All(result)
	}

	return withCollection(collection, nil, q)
}

func findOne(collection string, query interface{}, sortFields []string, result interface{}) error {
	q := func(c *mgo.Collection) error {
		var err error
		qy := c.Find(query)

		if result == nil {
			return err
		}

		if len(sortFields) > 0 {
			qy = qy.Sort(sortFields...)
		}

		return qy.One(result)
	}

	return withCollection(collection, nil, q)
}

func findIds(c string, ids []interface{}, result interface{}) error {
	return search(c, bson.M{"_id": bson.M{"$in": ids}}, nil, 0, 0, nil, nil, result)
}

func updateId(collection string, id interface{}, change interface{}, safe bool) error {
	update := func(c *mgo.Collection) error {
		return c.UpdateId(id, change)
	}

	if safe {
		return withCollection(collection, &mgo.Safe{}, update)
	}
	return withCollection(collection, nil, update)
}

func update(collection string, selector, change interface{}, safe bool) error {
	update := func(c *mgo.Collection) error {
		return c.Update(selector, change)
	}
	if safe {
		return withCollection(collection, &mgo.Safe{}, update)
	}
	return withCollection(collection, nil, update)
}

func upsert(collection string, selector, change interface{}, safe bool) (*mgo.ChangeInfo, error) {
	var chinfo *mgo.ChangeInfo

	upsert := func(c *mgo.Collection) (err error) {
		chinfo, err = c.Upsert(selector, change)
		//log.Println(chinfo, err)
		return err
	}
	if safe {
		return chinfo, withCollection(collection, &mgo.Safe{}, upsert)
	}
	return chinfo, withCollection(collection, nil, upsert)
}

func save(collection string, o interface{}, safe bool) error {
	insert := func(c *mgo.Collection) error {
		return c.Insert(o)
	}

	if safe {
		return withCollection(collection, &mgo.Safe{}, insert)
	}
	return withCollection(collection, nil, insert)
}

func remove(collection string, selector interface{}, safe bool) error {
	rm := func(c *mgo.Collection) error {
		return c.Remove(selector)
	}
	if safe {
		return withCollection(collection, &mgo.Safe{}, rm)
	}
	return withCollection(collection, nil, rm)
}

func removeId(collection, id string, safe bool) error {
	rm := func(c *mgo.Collection) error {
		return c.RemoveId(bson.ObjectIdHex(id))
	}
	if safe {
		return withCollection(collection, &mgo.Safe{}, rm)
	}
	return withCollection(collection, nil, rm)
}

func ensureIndex(collection string, keys ...string) error {
	ensure := func(c *mgo.Collection) error {
		return c.EnsureIndexKey(keys...)
	}

	return withCollection(collection, nil, ensure)
}

func DateString(t time.Time) string {
	return t.Format("2006-01-02")
}
