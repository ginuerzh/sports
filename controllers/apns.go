// apns
package controllers

import (
	"github.com/zhengying/apns"
	"log"
)

type ApnClient struct {
	Dev     *apns.Client
	Release *apns.Client
}

func (c *ApnClient) Send(token, alert string, badge int, sound string) error {
	payload := apns.NewPayload()
	payload.Alert = alert
	payload.Badge = badge
	payload.Sound = sound

	pn := apns.NewPushNotification()
	pn.DeviceToken = token
	pn.AddPayload(payload)

	resp := c.Dev.Send(pn)
	log.Println("apns push", resp.AppleResponse)
	if !resp.Success {
		log.Println("apns dev:", resp.Error)
	}
	resp = c.Release.Send(pn)
	if !resp.Success {
		log.Println("apns release:", resp.Error)
	}
	return resp.Error
}

func sendApn(client *ApnClient, msg string, badge int, devs ...string) error {
	for _, dev := range devs {
		if err := client.Send(dev, msg, badge, ""); err != nil {
			return err
		}
	}

	return nil
}
