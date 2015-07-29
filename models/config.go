package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo/bson"
	//"log"
)

type Video struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Desc  string `json:"desc"`
}

type Config struct {
	//Id     bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Videos []Video  `json:"videos"`
	Pets   []string `json:"pets"`
	Addr   string   `json:"addr"`
}

func (config *Config) Find() error {
	return findOne(configColl, nil, nil, config)
}

func (config *Config) Update() error {
	change := bson.M{
		"$set": config,
	}

	if _, err := upsert(configColl, nil, change, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}
