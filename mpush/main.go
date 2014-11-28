package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/AlekSi/pushover"
	"github.com/boltdb/bolt"
	"github.com/dustin/go-nma"
	"github.com/xconstruct/go-pushbullet"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type DistributionList struct {
	Id         string
	Nma        NMAConfig
	Pushbullet PBConfig
	Pushover   POConfig
}

const (
	bucketDistributionLists = "DistributionList"
)

var errNotFound = errors.New("list not found")

func allLists(db *bolt.DB, params martini.Params, req *http.Request, r render.Render) {
	var lists []DistributionList

	err := db.View(
		func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketDistributionLists))

			return b.ForEach(func(k, v []byte) error {
				var d DistributionList
				err := json.Unmarshal(v, &d)
				// TODO(dgryski): skip invalid keys instead of aborting on corrupt db?[
				if err != nil {
					return err
				}
				lists = append(lists, d)
				return nil
			})
		})

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, lists)
}

func findList(db *bolt.DB, params martini.Params, req *http.Request, r render.Render) {
	id := params["id"]

	var dst []byte
	err := db.View(
		func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketDistributionLists))
			dst = b.Get([]byte(id))
			return nil
		})

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	if dst == nil {
		r.Error(http.StatusNotFound)
		return
	}

	var l DistributionList
	err = json.Unmarshal(dst, &l)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, l)
}

func updateList(db *bolt.DB, params martini.Params, req *http.Request, r render.Render) {

	var l DistributionList

	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	err := json.Unmarshal(body, &l)

	if err != nil {
		r.Error(http.StatusBadRequest)
		return
	}

	if params["id"] != l.Id {
		r.Error(http.StatusBadRequest)
		return
	}

	// marshal back out to json to normalize our data
	j, err := json.Marshal(l)

	db.Update(
		func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketDistributionLists))
			return b.Put([]byte(l.Id), j)
		})

	r.Status(http.StatusOK)
}

func createList(db *bolt.DB, params martini.Params, req *http.Request, r render.Render) {

	var l DistributionList

	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	err := json.Unmarshal(body, &l)

	if err != nil {
		r.Error(http.StatusBadRequest)
		return
	}

	if params["id"] != l.Id {
		r.Error(http.StatusBadRequest)
		return
	}

	j, _ := json.Marshal(l)
	db.Update(
		func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketDistributionLists))
			return b.Put([]byte(l.Id), j)
		})

	r.JSON(http.StatusCreated, l)
}

func removeList(db *bolt.DB, params martini.Params, req *http.Request, r render.Render) {
	id := params["id"]

	err := db.Update(
		func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketDistributionLists))
			return b.Delete([]byte(id))
		})
	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.Status(http.StatusOK)
}

type PBUser struct {
	APIKey  string   `json:"apikey"`
	Devices []string `json:"devices"`
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

type Getter interface {
	Get(key []byte) []byte
}

func loadDistributionList(bucket Getter, list []byte) (DistributionList, error) {

	v := bucket.Get(list)

	if v == nil {
		return DistributionList{}, errNotFound
	}

	var l DistributionList
	if err := json.Unmarshal(v, &l); err != nil {
		return DistributionList{}, err
	}

	return l, nil
}

// don't actually send the notifications
// TODO(dgryski): replace with mocks
const debug = true

func sendNMA(pn *PushNotification, nmaconf NMAConfig, done chan struct{}) {

	n := &nma.Notification{Application: "mpush", Event: pn.Title, Description: pn.Body}

	for _, u := range nmaconf.APIKeys {
		c := nma.New(u)
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

func sendPushover(pn *PushNotification, poconf POConfig, done chan struct{}) {

	pushover.DefaultClient.ApplicationToken = poconf.APIKey

	for _, u := range poconf.Users {
		if debug {
			log.Println("pretending to notify pushover user", u)
		} else {
			err := pushover.Send(&pushover.Message{User: u, Title: pn.Title, Message: pn.Body, Timestamp: time.Now()})
			if err != nil {
				log.Println("Unable to notify pushover:", u, ":", err)
			}
		}
	}
	done <- struct{}{}
}

func sendPushBullet(pn *PushNotification, pbconf PBConfig, done chan struct{}) {
	for _, u := range pbconf.Users {
		c := pushbullet.New(u.APIKey)
		for _, d := range u.Devices {
			if debug {
				log.Println("pretending to notify pushbullet user", u)
			} else {
				err := c.PushNote(d, pn.Title, pn.Body)
				if err != nil {
					log.Println("Unable to notify pushbullet:", u.APIKey, "/", d, ":", err)
				}
			}
		}
	}
	done <- struct{}{}
}

type PushNotification struct {
	List  string
	Title string
	Body  string
}

func pushHandler(db *bolt.DB, params martini.Params, req *http.Request, r render.Render) {

	var pushNotification PushNotification

	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	err := json.Unmarshal(body, &pushNotification)

	if err != nil {
		r.Error(http.StatusBadRequest)
		return
	}

	if pushNotification.List == "" {
		r.Error(404)
		return
	}

	var targets DistributionList
	err = db.View(
		func(tx *bolt.Tx) error {
			targets, err = loadDistributionList(tx.Bucket([]byte(bucketDistributionLists)), []byte(pushNotification.List))
			return err
		})

	if err != nil {
		log.Println(err)
		if err == errNotFound {
			r.Error(http.StatusNotFound)
		} else {
			r.Error(http.StatusInternalServerError)
		}
		return
	}

	done := make(chan struct{})

	go sendNMA(&pushNotification, targets.Nma, done)
	go sendPushover(&pushNotification, targets.Pushover, done)
	go sendPushBullet(&pushNotification, targets.Pushbullet, done)

	for i := 0; i < 3; i++ {
		<-done
	}
}

func main() {

	db, err := bolt.Open("mpush.db", 0666, nil)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	m := martini.Classic()

	m.Use(render.Renderer())

	// Setup routes
	r := martini.NewRouter()
	r.Get(`/api/v0/lists`, allLists)
	r.Get(`/api/v0/lists/:id`, findList)
	r.Post(`/api/v0/lists`, createList)
	r.Put(`/api/v0/lists/:id`, updateList)
	r.Delete(`/api/v0/lists/:id`, removeList)

	r.Post("/api/v0/push", pushHandler)

	// Inject database
	m.Map(db)

	// Add the router action
	m.Action(r.Handle)

	m.Run()
}
