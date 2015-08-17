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

type Segment struct {
	ContentType string `bson:"seg_type" json:"seg_type"`
	ContentText string `bson:"seg_content" json:"seg_content"`
}
type Article struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Image    string        `bson:",omitempty"`
	Images   []string      `bson:",omitempty"`
	Contents []Segment
	Content  string
}

func main() {
	var articles []Article
	total := 0
	search("articles", bson.M{"type": bson.M{"$in": []interface{}{nil, ""}}}, nil, 0, 0, nil, &total, &articles)
	for _, article := range articles {
		set := bson.M{}
		if article.Image != "" {
			set["image"] = strings.Replace(article.Image, replacing, replaced, -1)
		}
		if len(article.Images) > 0 {
			for i, _ := range article.Images {
				article.Images[i] = strings.Replace(article.Images[i], replacing, replaced, -1)
			}
			set["images"] = article.Images
		}
		if article.Content != "" {
			set["content"] = strings.Replace(article.Content, replacing, replaced, -1)
		}
		for i, _ := range article.Contents {
			article.Contents[i].ContentText = strings.Replace(article.Contents[i].ContentText, replacing, replaced, -1)
		}
		set["contents"] = article.Contents
		change := bson.M{
			"$set": set,
		}
		b, _ := json.Marshal(change)
		log.Println(article.Id.Hex(), string(b))
		if err := updateId("articles", article.Id, change, true); err != nil {
			log.Println(err)
		}
	}
}
