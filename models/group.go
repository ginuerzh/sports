// group
package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo/bson"
	"time"
)

func init() {
	ensureIndex(groupColl, "gid")
	ensureIndex(groupColl, "creator")
	ensureIndex(groupColl, "-time")
}

type Group struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	Gid     string
	Name    string
	Profile string
	Desc    string
	Creator string
	Level   int
	Addr    Address
	Loc     Location
	Time    time.Time
	Members []string
}

func (group *Group) Exists() (bool, error) {
	return exists(groupColl, bson.M{"gid": group.Gid})
}

func (group *Group) FindById(gid string) error {
	return findOne(groupColl, bson.M{"gid": gid}, nil, group)
}

func (group *Group) Save() error {
	group.Id = bson.NewObjectId()
	group.Gid = group.Id.Hex()
	return save(groupColl, group, true)
}

func (group *Group) Remove(userid string) error {
	return remove(groupColl, bson.M{"gid": group.Gid, "creator": userid}, true)
}

func (group *Group) SetMember(userid string, remove bool) error {
	var m bson.M

	if remove {
		m = bson.M{
			"$pull": bson.M{
				"members": userid,
			},
		}
	} else {
		m = bson.M{
			"$addToSet": bson.M{
				"members": userid,
			},
		}
	}

	if err := update(groupColl, bson.M{"gid": group.Gid}, m, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}

	return nil
}
