// rule
package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo/bson"
	"time"
)

func init() {
}

type Rule struct {
	Id      string   `bson:"_id,omitempty" json:"-"`
	RuleId  int      `bson:"rule_id" json:"rule_id"`
	Users   []string `json:"users"`
	Message string   `json:"message,omitempty"`
}

func (r *Rule) Save() error {
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	query := bson.M{
		"reg_time": bson.M{"$gt": time.Unix(0, 0), "$lt": t},
	}
	var users []Account
	search(accountColl, query, bson.M{"_id": true}, 0, 0, nil, nil, &users)
	for _, user := range users {
		r.Users = append(r.Users, user.Id)
	}

	m := Struct2Map(r)
	if _, err := upsert(ruleColl, bson.M{"rule_id": r.RuleId}, m, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}
