// file
package models

import (
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/weedo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func init() {

}

type File struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Fid         string
	Name        string `bson:"filename"`
	Length      int64  `bson:"length"`
	Md5         string
	Owner       string
	Count       int
	ContentType string    `bson:"contentType"`
	UploadDate  time.Time `bson:"uploadDate"`
}

func (this *File) Exists() (bool, error) {
	count := 0
	err := search(fileColl, bson.M{"fid": this.Fid}, nil, 0, 0, nil, &count, nil)
	if err != nil {
		return false, errors.NewError(errors.DbError)
	}
	return count > 0, nil
}

func (this *File) findOne(query interface{}) (bool, error) {
	var files []File

	err := search(fileColl, query, nil, 0, 1, nil, nil, &files)
	if err != nil {
		return false, errors.NewError(errors.DbError)
	}
	if len(files) > 0 {
		*this = files[0]
	}

	return len(files) > 0, nil
}

func (this *File) FindByFid(fid string) (bool, error) {
	return this.findOne(bson.M{"fid": fid})
}

func (this *File) Save() error {
	this.Id = bson.NewObjectId()
	if err := save(fileColl, this, true); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

func (this *File) Delete() error {
	remove := func(c *mgo.Collection) error {
		err := c.Remove(bson.M{"fid": this.Fid})
		if err == nil {
			weedo.Delete(this.Fid, this.Count) //TODO: fail process
		}
		return err
	}

	if err := withCollection(fileColl, nil, remove); err != nil {
		if err != mgo.ErrNotFound {
			return errors.NewError(errors.DbError)
		}
	}
	return nil
}

func (this *File) OwnedBy(userid string) (bool, error) {
	return this.findOne(bson.M{"fid": this.Fid, "owner": userid})
}
