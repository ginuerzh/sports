// article
package models

import (
	"github.com/ginuerzh/sports/errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"labix.org/v2/mgo/txn"
	"log"
	//"strings"
	"time"
)

const (
	PrivPublic = iota
	_
	PrivPrivate
)

const (
	ArticleCoach   = "coach"
	ArticleRecord  = "record"
	ArticleComment = "comment"
)

const (
	TagRec   = "rec"   // recommend article
	TagTopic = "topic" // interview article
)

type Segment struct {
	ContentType string `bson:"seg_type" json:"seg_type"`
	ContentText string `bson:"seg_content" json:"seg_content"`
}

type rewardUser struct {
	Id   string
	Time time.Time
	Coin int64
}

type Article struct {
	Id           bson.ObjectId `bson:"_id,omitempty"`
	Type         string        // record: running record, coach: coach review, pk: pk result, default is post
	Privilege    int           // 0 - public, 2 - private
	Parent       string        `bson:",omitempty"`
	Author       string
	Refer        string   `bson:",omitempty"`
	ReferArticle string   `bson:"refer_article,omitempty"`
	Title        string   `bson:",omitempty"`
	Image        string   `bson:",omitempty"`
	Images       []string `bson:",omitempty"`
	Contents     []Segment
	Content      string
	Record       string
	PubTime      time.Time `bson:"pub_time"`

	Views            []string `bson:",omitempty"`
	Thumbs           []string `bson:",omitempty"`
	ThumbCount       int      `bson:"thumb_count"`
	Reviews          []string `bson:",omitempty"`
	ReviewCount      int      `bson:"review_count"`
	CoachReviewCount int      `bson:"coach_review_count"`
	Coaches          []string
	//Rewards          []string `bson:",omitempty"`
	RewardUsers []rewardUser
	RewardCount int      `bson:"reward_count"`
	TotalReward int64    `bson:"total_reward"`
	Tags        []string `bson:",omitempty"`
	Loc         Location `bson:",omitempty"`
}

func (this *Article) Exists() (bool, error) {
	b, err := exists(articleColl, bson.M{"_id": this.Id})
	if err != nil {
		return false, errors.NewError(errors.DbError)
	}
	return b, nil
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

func (this *Article) FindRefer() error {
	query := bson.M{
		"refer": bson.M{
			"$ne": nil,
		},
	}
	return findOne(articleColl, query, []string{"-pub_time"}, this)
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

	update := bson.M{
		"$push": bson.M{
			"reviews": this.Id.Hex(),
		},
		"$inc": bson.M{
			"review_count": 1,
		},
	}
	if this.Type == ArticleCoach {
		update = bson.M{
			"$addToSet": bson.M{
				"coaches": this.Author,
			},
			"$inc": bson.M{
				"coach_review_count": 1,
			},
		}
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
				Update: update,
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

func (this *Article) RemoveId() error {
	if err := removeId(articleColl, this.Id, true); err != nil {
		if e, ok := err.(*mgo.LastError); ok {
			return errors.NewError(errors.DbError, e.Error())
		}
	}
	return nil
}

func (this *Article) Remove() error {
	find, err := this.findOne(bson.M{"author": this.Author, "_id": this.Id})
	if !find {
		return err
	}

	if len(this.Parent) == 0 {
		if err := removeId(articleColl, this.Id, true); err != nil {
			if e, ok := err.(*mgo.LastError); ok {
				return errors.NewError(errors.DbError, e.Error())
			}
		}
		return nil
	}

	update := bson.M{
		"$pull": bson.M{
			"reviews": this.Id.Hex(),
		},
		"$inc": bson.M{
			"review_count": -1,
		},
	}
	if this.Type == ArticleCoach {
		update = bson.M{
			"$inc": bson.M{
				"coach_review_count": -1,
			},
		}
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
				C:      articleColl,
				Id:     bson.ObjectIdHex(this.Parent),
				Update: update,
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

func (article *Article) Update() error {
	m := bson.M{}

	if article.Content != "" {
		m["content"] = article.Content
	}
	if len(article.Contents) > 0 {
		m["contents"] = article.Contents
	}
	if article.Title != "" {
		m["title"] = article.Title
	}

	m["image"] = article.Image
	m["images"] = article.Images
	m["tags"] = article.Tags
	m["refer"] = article.Refer
	m["refer_article"] = article.ReferArticle
	m["pub_time"] = time.Now()

	change := bson.M{
		"$set": m,
	}

	if err := updateId(articleColl, article.Id, change, true); err != nil {
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
			"$inc": bson.M{
				"thumb_count": 1,
			},
		}
	} else {
		m = bson.M{
			"$pull": bson.M{
				"thumbs": userid,
			},
			"$inc": bson.M{
				"thumb_count": -1,
			},
		}
	}

	if err := updateId(articleColl, this.Id, m, true); err != nil {
		return errors.NewError(errors.DbError)
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

func articlePagingFunc(c *mgo.Collection, first, last string, args ...interface{}) (query bson.M, err error) {
	article := &Article{}

	if bson.IsObjectIdHex(first) {
		if err := c.FindId(bson.ObjectIdHex(first)).One(article); err != nil {
			return nil, err
		}
		query = bson.M{
			"pub_time": bson.M{
				"$gte": article.PubTime,
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
		}
	}

	return
}

func NewArticles(ids []string, last string) (int, []string, error) {
	total := 0
	var articles []Article

	query := bson.M{
		"parent":    nil,
		"refer":     nil,
		"privilege": bson.M{"$ne": 2},
		"author":    bson.M{"$in": ids},
	}
	if bson.IsObjectIdHex(last) {
		article := &Article{}
		findOne(articleColl, bson.M{"_id": bson.ObjectIdHex(last)}, nil, article)
		if len(article.Id) == 0 {
			return 0, nil, nil
		}

		query["pub_time"] = bson.M{
			"$gte": article.PubTime,
		}
		query["_id"] = bson.M{
			"$ne": article.Id,
		}
	}
	sortFields := []string{"-pub_time", "-_id"}
	search(articleColl, query, nil, 0, 2, sortFields, &total, &articles)

	var uids []string
	for i, _ := range articles {
		uids = append(uids, articles[i].Author)
	}
	//log.Println(len(uids), uids)
	users, _ := FindUsersByIds(0, uids...)
	var profiles []string
	for i, _ := range users {
		profiles = append(profiles, users[i].Profile)
	}
	return total, profiles, nil
}

func GetUserArticles(ids []string, paging *Paging) (int, []Article, error) {
	var articles []Article
	total := 0

	query := bson.M{
		"parent":    nil,
		"refer":     nil,
		"privilege": bson.M{"$ne": 2},
		"author":    bson.M{"$in": ids},
	}

	sortFields := []string{"-pub_time", "-_id"}

	if err := psearch(articleColl, query, nil,
		sortFields, nil, &articles, articlePagingFunc, paging); err != nil {
		e := errors.NewError(errors.DbError, err.Error())
		if err == mgo.ErrNotFound {
			e = errors.NewError(errors.NotFoundError, err.Error())
		}
		return total, nil, e
	}

	for i := 0; i < len(articles); i++ {
		if articles[i].Id.Hex() == paging.First {
			articles = articles[:i]
			break
		} else if articles[i].Id.Hex() == paging.Last {
			articles = articles[i+1:]
			break
		}
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

func GetFollowingsArticles(followings []string, paging *Paging) (int, []Article, error) {
	return GetUserArticles(followings, paging)
}

func GetRecommendArticles(recommends []string, paging *Paging) (int, []Article, error) {
	return GetUserArticles(recommends, paging)
}

func GetArticles(tag string, paging *Paging, withoutContent bool) (int, []Article, error) {
	var articles []Article
	total := 0

	query := bson.M{
		"parent": nil,
		"refer":  nil,
	}
	if len(tag) > 0 {
		query["tags"] = tag
	}

	var selector bson.M

	if withoutContent {
		selector = bson.M{
			"content":  0,
			"contents": 0,
		}
	}

	sortFields := []string{"-pub_time", "-_id"}

	if err := psearch(articleColl, query, selector,
		sortFields, nil, &articles, articlePagingFunc, paging); err != nil {
		e := errors.NewError(errors.DbError, err.Error())
		if err == mgo.ErrNotFound {
			e = errors.NewError(errors.NotFoundError, err.Error())
		}
		return total, nil, e
	}

	for i := 0; i < len(articles); i++ {
		if articles[i].Id.Hex() == paging.First {
			articles = articles[:i]
			break
		} else if articles[i].Id.Hex() == paging.Last {
			articles = articles[i+1:]
			break
		}
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

func (this *Article) CommentCount() (count int) {
	search(articleColl, bson.M{"parent": this.Id.Hex()}, nil, 0, 0, nil, &count, nil)
	return
}

func (this *Article) Comments(typ string, paging *Paging, withoutContent bool) (int, []Article, error) {
	var articles []Article
	total := 0

	sortFields := []string{"-pub_time"}

	var selector bson.M

	if withoutContent {
		selector = bson.M{"content": 0, "contents": 0}
	}

	query := bson.M{
		"parent": this.Id.Hex(),
	}
	switch typ {
	case ArticleCoach:
		query["type"] = typ
	default:
		query["type"] = bson.M{
			"$ne": ArticleCoach,
		}
	}

	if err := psearch(articleColl, query, selector,
		sortFields, &total, &articles, articlePagingFunc, paging); err != nil {
		e := errors.NewError(errors.DbError, err.Error())
		if err == mgo.ErrNotFound {
			e = errors.NewError(errors.NotFoundError, err.Error())
		}
		return total, nil, e
	}

	for i := 0; i < len(articles); i++ {
		if articles[i].Id.Hex() == paging.First {
			articles = articles[:i]
			break
		} else if articles[i].Id.Hex() == paging.Last {
			articles = articles[i+1:]
			break
		}
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

func (this *Article) AdminComments(pageIndex, pageCount int) (total int, articles []Article, err error) {
	err = search(articleColl, bson.M{"parent": this.Id.Hex()}, nil,
		pageIndex*pageCount, pageCount, []string{"-pub_time"}, &total, &articles)
	return
}

func (this *Article) Reward(userid string, amount int64) error {
	change := mgo.Change{
		Update: bson.M{
			"$push": bson.M{
				"rewardusers": &rewardUser{Id: userid, Time: time.Now(), Coin: amount},
			},
			"$inc": bson.M{
				"total_reward": amount,
				"reward_count": 1,
			},
		},
		ReturnNew: true,
	}
	_, err := apply(articleColl, bson.M{"_id": this.Id}, change, this)

	return err
}

func PostCount(start, end time.Time) int {
	c, _ := count(articleColl, bson.M{"pub_time": bson.M{"$gte": start, "$lt": end}})
	return c
}

func SearchArticle(keyword string, paging *Paging) (int, []Article, error) {
	var articles []Article
	total := 0

	query := bson.M{
		"content": bson.M{
			"$regex":   keyword,
			"$options": "i",
		},
	}

	sortFields := []string{"-pub_time", "-_id"}

	if err := psearch(articleColl, query, bson.M{"content": 0, "contents": 0}, sortFields, &total, &articles,
		articlePagingFunc, paging); err != nil {
		if err != mgo.ErrNotFound {
			return total, nil, errors.NewError(errors.DbError, err.Error())
		}
	}

	for i := 0; i < len(articles); i++ {
		if articles[i].Id.Hex() == paging.First {
			articles = articles[:i]
			break
		} else if articles[i].Id.Hex() == paging.Last {
			articles = articles[i+1:]
			break
		}
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

func AdminSearchArticle(keyword string, tag string,
	pageIndex, pageCount int) (total int, articles []Article, err error) {
	query := bson.M{"parent": nil}

	if len(keyword) > 0 {
		query["content"] = bson.M{
			"$regex":   keyword,
			"$options": "i",
		}
	}
	if len(tag) > 0 {
		query["tags"] = tag
	}
	/*
		if len(keyword) == 0 {
			query = bson.M{
				"parent": nil,
				"tags":   tag,
			}
		} else {
			query = bson.M{
				"parent": nil,
				"content": bson.M{
					"$regex":   keyword,
					"$options": "i",
				},
			}
			if len(tag) > 0 {
				query["tags"] = tag
			}
		}
	*/

	err = search(articleColl, query, bson.M{"content": 0, "contents": 0},
		pageIndex*pageCount, pageCount, []string{"-pub_time"}, &total, &articles)
	return
}

func ArticleList(tag string, sort string, pageIndex, pageCount int) (total int, articles []Article, err error) {
	switch sort {
	case "pubtime":
		sort = "pub_time"
	case "-pubtime":
		sort = "-pub_time"
	default:
		sort = "-pub_time"
	}
	query := bson.M{
		"parent": nil,
		"type": bson.M{
			"$ne": "record",
		},
		"refer": nil,
	}
	if tag != "" {
		query["tags"] = tag
	}

	err = search(articleColl, query, bson.M{"content": 0, "contents": 0},
		pageIndex*pageCount, pageCount, []string{sort}, &total, &articles)
	return
}

func (this *Article) FindByRecord(id string) error {
	return findOne(articleColl, bson.M{"record": id}, nil, this)
}

func (this *Article) SetPrivilege(priv int) error {
	change := bson.M{
		"$set": bson.M{
			"privilege": priv,
		},
	}
	return updateId(articleColl, this.Id, change, true)
}

func (this *Article) SetTag(tag string) error {
	change := bson.M{
		"$addToSet": bson.M{
			"tags": tag,
		},
	}

	return updateId(articleColl, this.Id, change, true)
}

func (this *Article) UnsetTag(tag string) error {
	change := bson.M{
		"$pull": bson.M{
			"tags": tag,
		},
	}

	if tag == "" {
		change = bson.M{
			"$unset": bson.M{
				"tags": 1,
			},
		}
	}
	return updateId(articleColl, this.Id, change, true)
}

func TopicList(pageIndex, pageCount int) (total int, articles []Article, err error) {
	query := bson.M{
		"refer": bson.M{
			"$ne": nil,
		},
	}

	err = search(articleColl, query, nil,
		pageIndex*pageCount, pageCount, []string{"-pub_time"}, &total, &articles)
	return
}
