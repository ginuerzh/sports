// article_rate
package models

/*
import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo/bson"
)

func init() {
	ensureIndex(rateColl, "userid")
}

type ArticleRate struct {
	Article string
	Rate    int
}

type UserRate struct {
	Userid string
	Rates  []ArticleRate
}

func (this *UserRate) findOne(query interface{}) (bool, int) {
	var rates []UserRate

	err := search(rateColl, query, nil, 0, 1, nil, nil, &rates)
	if err != nil {
		return false, errors.DbError
	}
	if len(rates) > 0 {
		*this = rates[0]
	}

	return len(rates) > 0, errors.NoError
}

func (this *UserRate) FindByUserid(userid string) (bool, int) {
	return this.findOne(bson.M{"userid": userid})
}
*/
