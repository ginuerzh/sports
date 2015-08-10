// convert user contacts struct
package main

import (
	"flag"
	"labix.org/v2/mgo/bson"
	"log"
)

func init() {
	flag.StringVar(&MongoAddr, "mongo", "localhost:27017", "mongodb server")
	flag.Parse()
}

type Account struct {
	Id    string `bson:"_id,omitempty"`
	Actor string `bson:",omitempty"`
}

func main() {
	var users []Account

	total := 0
	search("accounts", nil, bson.M{"actor": 1}, 0, 0, nil, &total, &users)

	for _, user := range users {
		if user.Actor == "" {
			continue
		}
		var actor []string
		actor = append(actor, user.Actor)

		change := bson.M{"$set": bson.M{"actor": actor}}
		if err := updateId("accounts", user.Id, change, true); err != nil {
			log.Println(err)
		}
	}
}
