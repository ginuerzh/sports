// user
package models

import ()

/*
type Event struct {
	Id      string
	Thumbs  []string `bson:",omitempty"`
	Reviews []string `bson:",omitempty"`
}
*/
/*
type User struct {
	Id string `bson:"_id"`

	//Events   []Event   `bson:",omitempty"`

}


func (this *User) findOne(query interface{}) (bool, error) {
	var users []User

	err := search(userColl, query, nil, 0, 1, nil, nil, &users)
	if err != nil {
		return false, errors.NewError(errors.DbError, err.Error())
	}
	if len(users) > 0 {
		*this = users[0]
	}

	return len(users) > 0, nil
}

func (this *User) FindByUserid(userid string) (bool, error) {
	return this.findOne(bson.M{"_id": userid})
}


func (this *User) LastMessage(userid string) *Message {
	_, msgs, _ := this.Messages(userid, &Paging{Count: 1})
	if len(msgs) == 0 {
		return nil
	}
	return &msgs[0]
}
*/

/*
func (this *User) AddEvent(event *Event) error {
	var change bson.M
	selector := bson.M{
		"_id":       this.Id,
		"events.id": event.Id,
	}
	if len(event.Thumbs) > 0 {
		change = bson.M{
			"$addToSet": bson.M{
				"events.$.thumbs": bson.M{
					"$each": event.Thumbs,
				},
			},
		}
	}
	if len(event.Reviews) > 0 {
		change = bson.M{
			"$addToSet": bson.M{
				"events.$.reviews": bson.M{
					"$each": event.Reviews,
				},
			},
		}
	}

	err := update(userColl, selector, change, true)
	//log.Println(err)
	if err == nil {
		return nil
	}

	if err != mgo.ErrNotFound {
		return errors.NewError(errors.DbError, err.Error())
	}

	// not found
	selector = bson.M{
		"_id": this.Id,
	}
	change = bson.M{
		"$push": bson.M{
			"events": event,
		},
	}
	err = update(userColl, selector, change, true)
	if err != nil {
		return errors.NewError(errors.DbError, err.Error())
	}

	return nil
}
*/
