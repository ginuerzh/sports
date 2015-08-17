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

type UserAuth struct {
	IdCard    *AuthInfo `bson:",omitempty" json:"idcard"`
	IdCardTmp *AuthInfo `bson:",omitempty" json:"-"`
	Cert      *AuthInfo `bson:",omitempty" json:"cert"`
	CertTmp   *AuthInfo `bson:",omitempty" json:"-"`
	Record    *AuthInfo `bson:",omitempty" json:"record"`
	RecordTmp *AuthInfo `bson:",omitempty" json:"-"`
}

type AuthInfo struct {
	Images []string `bson:",omitempty" json:"auth_images"`
	Desc   string   `bson:",omitempty" json:"auth_desc"`
	Status string   `bson:",omitempty" json:"auth_status"`
	Review string   `bson:",omitempty" json:"auth_review"`
}

type Account struct {
	Id      string `bson:"_id,omitempty"`
	Profile string
	Photos  []string
	Auth    *UserAuth
}

func main() {
	/*
		it, err := iter("accounts", nil, nil, 0, 0, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer it.Close()
		user := &Account{}
	*/
	var users []Account
	total := 0
	search("accounts", nil, nil, 0, 0, nil, &total, &users)

	for _, user := range users {
		set := bson.M{}
		if strings.Contains(user.Profile, replacing) {
			set["profile"] = strings.Replace(user.Profile, replacing, replaced, -1)
		}
		find := false
		for i, photo := range user.Photos {
			if strings.Contains(photo, replacing) {
				user.Photos[i] = strings.Replace(user.Photos[i], replacing, replaced, -1)
				find = true
			}
		}
		if find {
			set["photos"] = user.Photos
		}

		find = false

		if user.Auth != nil {
			if b := fixAuthImage(user.Auth.Cert); b {
				find = b
			}
			if b := fixAuthImage(user.Auth.CertTmp); b {
				find = b
			}
			if b := fixAuthImage(user.Auth.IdCard); b {
				find = b
			}
			if b := fixAuthImage(user.Auth.IdCardTmp); b {
				find = b
			}
			if b := fixAuthImage(user.Auth.Record); b {
				find = b
			}
			if b := fixAuthImage(user.Auth.RecordTmp); b {
				find = b
			}
		}

		if find {
			set["auth"] = user.Auth
		}

		if len(set) > 0 {
			change := bson.M{"$set": set}
			b, _ := json.Marshal(change)
			log.Println(user.Id, string(b))
			if err := updateId("accounts", user.Id, change, true); err != nil {
				log.Println(err)
			}
		}
	}
}

func fixAuthImage(auth *AuthInfo) (find bool) {
	if auth == nil {
		return
	}
	for i, image := range auth.Images {
		if strings.Contains(image, replacing) {
			auth.Images[i] = strings.Replace(auth.Images[i], replacing, replaced, -1)
			find = true
		}
	}

	return
}
