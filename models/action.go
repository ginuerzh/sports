// action
package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo/bson"
	"time"
)

func init() {
	ensureIndex(actionColl, "userid")
	ensureIndex(actionColl, "date")
	ensureIndex(actionColl, "userid", "date")
}

type Action struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	Userid  string
	Date    time.Time
	Login   int
	Post    int
	Comment int
	Invite  int
	Profile int
	Info    int
}

func (this *Action) findOne(query interface{}) (bool, error) {
	var actions []Action

	err := search(actionColl, query, nil, 0, 1, nil, nil, &actions)
	if err != nil {
		return false, errors.NewError(errors.DbError, err.Error())
	}
	if len(actions) > 0 {
		*this = actions[0]
	}
	return len(actions) > 0, nil
}
