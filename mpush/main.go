package main

import (
	"encoding/json"
	"github.com/dustin/go-nma"
	"github.com/jdiez17/go-pushover"
	"github.com/xconstruct/go-pushbullet"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Pushbullet PBConfig
	Nma        NMAConfig
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

func readConfig() (Config, error) {
	cfgfile := filepath.Join(os.Getenv("HOME"), ".mpush.config.json")
	f, err := os.Open(cfgfile)
	if err != nil {
		return Config{}, err
	}
	defer func() {
		f.Close()
	}()

	var cfg Config
	dec := json.NewDecoder(f)
	if err = dec.Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func main() {

	conf, err := readConfig()

	if err != nil {
		log.Fatal(err)
	}

	title := "Test Message"
	body := "Test Body Hello World"

	for _, u := range conf.Nma.APIKeys {
		c := nma.New(u)
		n := &nma.Notification{Event: title, Description: body}
		err := c.Notify(n)
		if err != nil {
			log.Println("Unable to notify nma:", u, ":", err)
		}
	}

	for _, u := range conf.Pushbullet.Users {
		c := pushbullet.New(u.APIKey)
		for _, d := range u.Devices {
			err := c.PushNote(d, title, body)
			if err != nil {
				log.Println("Unable to notify pushbullet:", u.APIKey, "/", d, ":", err)
			}
		}
	}

	for _, u := range conf.Pushover.Users {
		c := pushover.New(u, conf.Pushover.APIKey)
		n := pushover.Notification{
			Title:     title,
			Message:   body,
			Timestamp: time.Now(),
		}
		_, err := c.Notify(n)
		if err != nil {
			log.Println("Unable to notify pushover:", u, ":", err)
		}
	}
}
