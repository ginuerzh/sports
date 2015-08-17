// convert user image url
package main

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"log"
	"strings"
)

const (
	replacing = "106.187.48.51"
	replaced  = "ice139.com"
)

type SportRecord struct {
	Source    string
	Duration  int64
	Distance  int
	Weight    int
	Mood      string
	HeartRate int
	Speed     float64
	Pics      []string
	Review    string
}

type Record struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Sport *SportRecord  `bson:",omitempty"`
}

func main() {
	var records []Record
	total := 0
	search("records", bson.M{"sport": bson.M{"$ne": nil}}, nil, 0, 0, nil, &total, &records)
	for _, record := range records {
		if len(record.Sport.Pics) == 0 {
			continue
		}
		for i, _ := range record.Sport.Pics {
			record.Sport.Pics[i] = strings.Replace(record.Sport.Pics[i], replacing, replaced, -1)
		}
		change := bson.M{
			"$set": bson.M{
				"sport.pics": record.Sport.Pics,
			},
		}
		b, _ := json.Marshal(change)
		log.Println(record.Id.Hex(), string(b))
		if err := updateId("records", record.Id, change, true); err != nil {
			log.Println(err)
		}
	}
}
