// group
package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo/bson"
	"time"
)

func init() {
}

type Group struct {
	Id      bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Gid     string        `json:"-"`
	Name    string        `json:"name,omitempty"`
	Profile string        `bson:",omitempty" json:"profile,omitempty"`
	Desc    string        `bson:",omitempty" json:"desc,omitempty"`
	Creator string        `json:"-"`
	Level   int           `json:"-"`
	Addr    *Address      `bson:",omitempty" json:"addr,omitempty"`
	Loc     *Location     `bson:",omitempty" json:"loc,omitempty"`
	Time    time.Time     `json:"-"`
	Members []string      `bson:",omitempty" json:"-"`
}

func (group *Group) Exists() (bool, error) {
	return exists(groupColl, bson.M{"gid": group.Gid})
}

func (group *Group) FindById(gid string) error {
	if err := findOne(groupColl, bson.M{"gid": gid}, nil, group); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (group *Group) Save() error {
	group.Id = bson.NewObjectId()
	group.Gid = group.Id.Hex()
	if err := save(groupColl, group, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (group *Group) Update() error {
	change := bson.M{
		"$set": Struct2Map(group),
	}

	if err := update(groupColl, bson.M{"gid": group.Gid}, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (group *Group) Remove(userid string) error {
	if err := remove(groupColl, bson.M{"gid": group.Gid, "creator": userid}, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
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
