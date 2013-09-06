package main

import (
	"encoding/json"
	"flag"
	"github.com/cznic/kv"
	"github.com/dustin/go-nma"
	"github.com/jdiez17/go-pushover"
	"github.com/xconstruct/go-pushbullet"
	"log"
	"time"
)

type DistributionList struct {
	Nma        NMAConfig
	Pushbullet PBConfig
	Pushover   POConfig
}

type PBUser struct {
	APIKey  string `json:"apikey"`
	Devices []int  `json:"devices"`
}

type PBConfig struct {
	Users []PBUser `json:"users"`
}

type NMAConfig struct {
	APIKeys []string `json:"apikeys"`
}

type POConfig struct {
	APIKey string   `json:"apikey"`
	Users  []string `json:"users"`
}

func loadDistributionList(datastore *kv.DB, list []byte) (DistributionList, error) {

	v, err := datastore.Get(nil, list)
	if err != nil {
		return DistributionList{}, err
	}

	var l DistributionList
	if err = json.Unmarshal(v, &l); err != nil {
		return DistributionList{}, err
	}

	return l, nil
}

// don't actually send the notifications
// TODO(dgryski): replace with mocks
const debug = true

func SendNMA(title, body string, nmaconf NMAConfig, done chan struct{}) {
	for _, u := range nmaconf.APIKeys {
		c := nma.New(u)
		n := &nma.Notification{Application: "mpush", Event: title, Description: body}
		if debug {
			log.Println("pretending to notify nma user", u)
		} else {
			err := c.Notify(n)
			if err != nil {
				log.Println("Unable to notify nma:", u, ":", err)
			}
		}
	}

	done <- struct{}{}
}

func SendPushover(title, body string, poconf POConfig, done chan struct{}) {
	for _, u := range poconf.Users {
		c := pushover.New(u, poconf.APIKey)
		n := pushover.Notification{
			Title:     title,
			Message:   body,
			Timestamp: time.Now(),
		}
		if debug {
			log.Println("pretending to notify pushover user", u)
		} else {
			_, err := c.Notify(n)
			if err != nil {
				log.Println("Unable to notify pushover:", u, ":", err)
			}
		}
	}
	done <- struct{}{}
}

func SendPushBullet(title, body string, pbconf PBConfig, done chan struct{}) {
	for _, u := range pbconf.Users {
		c := pushbullet.New(u.APIKey)
		for _, d := range u.Devices {
			if debug {
				log.Println("pretending to notify pushbullet user", u)
			} else {
				err := c.PushNote(d, title, body)
				if err != nil {
					log.Println("Unable to notify pushbullet:", u.APIKey, "/", d, ":", err)
				}
			}
		}
	}
	done <- struct{}{}
}

func main() {

	list := flag.String("list", "", "the target list")
	title := flag.String("t", "", "title of the message")
	body := flag.String("b", "", "the body of the message")

	flag.Parse()

	datastore, err := kv.Open("mpush.kvdb", &kv.Options{})

	defer datastore.Close()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("looking for distribution list: ", *list)

	conf, err := loadDistributionList(datastore, []byte(*list))
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go SendNMA(*title, *body, conf.Nma, done)
	go SendPushover(*title, *body, conf.Pushover, done)
	go SendPushBullet(*title, *body, conf.Pushbullet, done)

	for i := 0; i < 3; i++ {
		<-done
	}
}
