// article
package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"labix.org/v2/mgo/txn"
	"log"
	"time"
)

func init() {
	ensureIndex(articleColl, "author")
	ensureIndex(articleColl, "-pub_time")
}

type Segment struct {
	ContentType string `bson:"seg_type" json:"seg_type"`
	ContentText string `bson:"seg_content" json:"seg_content"`
}

type Article struct {
	Id     bson.ObjectId `bson:"_id,omitempty"`
	Parent string        `bson:",omitempty"`
	Author string
	//Title    string `bson:",omitempty"`
	//Image    string `bson:",omitempty"`
	Contents []Segment
	PubTime  time.Time `bson:"pub_time"`

	Views   []string `bson:",omitempty"`
	Thumbs  []string `bson:",omitempty"`
	Reviews []string `bson:",omitempty"`
}

func FindArticles(ids ...string) (articles []Article, err error) {
	var oid []interface{}
	for _, id := range ids {
		oid = append(oid, bson.ObjectIdHex(id))
	}

	if e := findIds(articleColl, oid, &articles); e != nil {
		err = errors.NewError(errors.DbError, e.Error())
	}
	return
}

func (this *Article) findOne(query interface{}) (bool, error) {
	var articles []Article

	err := search(articleColl, query, nil, 0, 1, nil, nil, &articles)
	if err != nil {
		return false, errors.NewError(errors.DbError, err.Error())
	}
	if len(articles) > 0 {
		*this = articles[0]
	}

	return len(articles) > 0, nil
}

func (this *Article) FindById(id string) (bool, error) {
	if !bson.IsObjectIdHex(id) {
		return false, nil
	}
	return this.findOne(bson.M{"_id": bson.ObjectIdHex(id)})
}

func (this *Article) Save() error {
	this.Id = bson.NewObjectId()
	if len(this.Parent) == 0 {
		if err := save(articleColl, this, true); err != nil {
			return errors.NewError(errors.DbError, err.Error())
		}
		return nil
	}

	if !bson.IsObjectIdHex(this.Parent) {
		return errors.NewError(errors.InvalidMsgError)
	}

	f := func(c *mgo.Collection) error {
		runner := txn.NewRunner(c)
		ops := []txn.Op{
			{
				C:      articleColl,
				Id:     this.Id,
				Assert: txn.DocMissing,
				Insert: this,
			},
			{
				C:      articleColl,
				Id:     bson.ObjectIdHex(this.Parent),
				Assert: txn.DocExists,
				Update: bson.M{
					"$addToSet": bson.M{
						"reviews": this.Id.Hex(),
					},
				},
			},
		}

		return runner.Run(ops, bson.NewObjectId(), nil)
	}

	if err := withCollection("comment_tx", &mgo.Safe{}, f); err != nil {
		log.Println(err)
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Article) Remove() error {
	find, err := this.findOne(bson.M{"author": this.Author, "_id": this.Id})
	if !find {
		return err
	}

	if len(this.Parent) == 0 {
		if err := removeId(articleColl, this.Id.Hex(), true); err != nil {
			if e, ok := err.(*mgo.LastError); ok {
				return errors.NewError(errors.DbError, e.Error())
			}
		}
		return nil
	}

	f := func(c *mgo.Collection) error {
		runner := txn.NewRunner(c)
		ops := []txn.Op{
			{
				C:      articleColl,
				Id:     this.Id,
				Remove: true,
			},
			{
				C:  articleColl,
				Id: bson.ObjectIdHex(this.Parent),
				Update: bson.M{
					"$pull": bson.M{
						"reviews": this.Id.Hex(),
					},
				},
			},
		}

		return runner.Run(ops, bson.NewObjectId(), nil)
	}
	if err := withCollection("comment_tx", &mgo.Safe{}, f); err != nil {
		log.Println(err)
		return errors.NewError(errors.DbError, err.Error())
	}
	return nil
}

func (this *Article) SetThumb(userid string, thumb bool) error {

	var m bson.M

	if thumb {
		m = bson.M{
			"$addToSet": bson.M{
				"thumbs": userid,
			},
		}
	} else {
		m = bson.M{
			"$pull": bson.M{
				"thumbs": userid,
			},
		}
	}

	if err := updateId(articleColl, this.Id, m, true); err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}

	return nil
}

func (this *Article) IsThumbed(userid string) (bool, error) {
	count := 0
	err := search(articleColl, bson.M{"_id": this.Id, "thumbs": userid}, nil, 0, 0, nil, &count, nil)
	if err != nil {
		return false, errors.NewError(errors.DbError, err.Error())
	}
	return count > 0, nil
}

func articlePagingFunc(c *mgo.Collection, first, last string) (query bson.M, err error) {
	article := &Article{}

	if bson.IsObjectIdHex(first) {
		if err := c.FindId(bson.ObjectIdHex(first)).One(article); err != nil {
			return nil, err
		}
		query = bson.M{
			"pub_time": bson.M{
				"$gte": article.PubTime,
			},
			"_id": bson.M{
				"$ne": article.Id,
			},
		}
	} else if bson.IsObjectIdHex(last) {
		if err := c.FindId(bson.ObjectIdHex(last)).One(article); err != nil {
			return nil, err
		}
		query = bson.M{
			"pub_time": bson.M{
				"$lte": article.PubTime,
			},
			"_id": bson.M{
				"$ne": article.Id,
			},
		}
	}

	return
}

func GetArticles(paging *Paging) (int, []Article, error) {
	var articles []Article
	total := 0

	if err := psearch(articleColl, bson.M{"parent": nil}, nil,
		[]string{"-pub_time"}, nil, &articles, articlePagingFunc, paging); err != nil {
		e := errors.NewError(errors.DbError, err.Error())
		if err == mgo.ErrNotFound {
			e = errors.NewError(errors.NotFoundError, err.Error())
		}
		return total, nil, e
	}

	paging.First = ""
	paging.Last = ""
	paging.Count = 0
	if len(articles) > 0 {
		paging.First = articles[0].Id.Hex()
		paging.Last = articles[len(articles)-1].Id.Hex()
		paging.Count = total
	}

	return total, articles, nil
}

func (this *Article) Comments(paging *Paging) (int, []Article, error) {
	var articles []Article
	total := 0

	if err := psearch(articleColl, bson.M{"parent": this.Id.Hex()}, nil,
		[]string{"-pub_time"}, nil, &articles, articlePagingFunc, paging); err != nil {
		e := errors.NewError(errors.DbError, err.Error())
		if err == mgo.ErrNotFound {
			e = errors.NewError(errors.NotFoundError, err.Error())
		}
		return total, nil, e
	}

	paging.First = ""
	paging.Last = ""
	paging.Count = 0
	if len(articles) > 0 {
		paging.First = articles[0].Id.Hex()
		paging.Last = articles[len(articles)-1].Id.Hex()
		paging.Count = total
	}

	return total, articles, nil
}

/*
func (this *Article) ReviewCount() (int, int) {
	total := 0
	err := search(reviewColl, bson.M{"article_id": this.Id.Hex()}, nil, 0, 0, nil, &total, nil)
	if err != nil {
		return 0, errors.DbError
	}
	return total, errors.NoError
}

func (this *Article) LoadBrief() int {
	var articles []Article
	if err := search(articleColl, bson.M{"_id": this.Id}, bson.M{"content": false}, 0, 1, nil, nil, &articles); err != nil {
		return errors.DbError
	}

	if len(articles) > 0 {
		*this = articles[0]
	}
	return errors.NoError
}

func GetArticles(articleIds ...string) (articles []Article, errId int) {
	ids := make([]bson.ObjectId, len(articleIds))
	for i, _ := range articleIds {
		ids[i] = bson.ObjectIdHex(articleIds[i])
	}
	err := search(articleColl,
		bson.M{"_id": bson.M{"$in": ids}},
		bson.M{"content": false},
		0, 0, []string{"-pub_time"}, nil, &articles)
	if err != nil {
		return nil, errors.DbError
	}

	errId = errors.NoError
	return
}


func RandomArticles(excludes []string, max int) (article []Article, errId int) {
	ids := make([]bson.ObjectId, len(excludes))
	for i, _ := range excludes {
		ids[i] = bson.ObjectIdHex(excludes[i])
	}

	selector := bson.M{
		"_id":bson.M{"$nin": ids},
		"random":
	}
}
*/
