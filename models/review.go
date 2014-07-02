// review
package models

/*
import (
	"github.com/ginuerzh/sports/errors"
	//"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

func init() {
	ensureIndex(reviewColl, "userid", "-ctime")
}

type Review struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	ArticleId string        `bson:"article_id"`
	Userid    string
	Content   string
	Mentions  []string `bson:",omitempty"`
	Thumbs    []string `bson:",omitempty"`
	Ctime     time.Time
	Mtime     time.Time `bson:",omitempty"`
}

func (this *Review) findOne(query interface{}) (bool, int) {
	var reviews []Review

	err := search(reviewColl, query, nil, 0, 1, nil, nil, &reviews)
	if err != nil {
		return false, errors.DbError
	}
	if len(reviews) > 0 {
		*this = reviews[0]
	}

	return len(reviews) > 0, errors.NoError
}

func (this *Review) FindById(id string) (bool, int) {
	return this.findOne(bson.M{"_id": bson.ObjectIdHex(id)})
}

func (this *Review) Save() (errId int) {
	errId = errors.NoError

	this.Id = bson.NewObjectId()
	if err := save(reviewColl, this, true); err != nil {
		errId = errors.DbError
	}
	return
}

func GetReviewList(articleId string, skip, limit int) (reviews []Review, errId int) {
	err := search(reviewColl, bson.M{"article_id": articleId}, nil, skip, limit, []string{"-ctime"}, nil, &reviews)
	if err != nil {
		return nil, errors.DbError
	}

	errId = errors.NoError
	return
}

func (this *Review) SetThumb(userid string, thumb bool) int {
	var change bson.M

	if thumb {
		change = bson.M{
			"$addToSet": bson.M{
				"thumbs": userid,
			},
		}
	} else {
		change = bson.M{
			"$pull": bson.M{
				"thumbs": userid,
			},
		}
	}
	err := updateId(reviewColl, this.Id, change, true)
	if err != nil {
		return errors.DbError
	}
	return errors.NoError
}

func (this *Review) IsThumbed(userid string) (bool, int) {
	count := 0
	err := search(reviewColl, bson.M{"_id": this.Id, "thumbs": userid}, nil, 0, 0, nil, &count, nil)
	if err != nil {
		return false, errors.DbError
	}
	return count > 0, errors.NoError
}
*/
