// convert user contacts struct
package main

import (
	"labix.org/v2/mgo/bson"
	"log"
)

type Contact struct {
	Id       string
	Profile  string
	Nickname string
}

type Account struct {
	Id       string    `bson:"_id,omitempty"`
	Contacts []Contact `bson:",omitempty"`
}

func main() {
	var users []Account

	total := 0
	search("accounts", nil, bson.M{"contacts": 1}, 0, 0, nil, &total, &users)

	for _, user := range users {
		var ids []string
		for _, contact := range user.Contacts {
			ids = append(ids, contact.Id)
		}
		change := bson.M{"$set": bson.M{"contacts": ids}}
		if err := updateId("accounts", user.Id, change, true); err != nil {
			log.Println(err)
		}
	}
}
