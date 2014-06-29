package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type HeroDB interface{}

type heroDB struct {
	heros    map[int]Hero
	numHeros int
}

type Hero struct {
	Id    int
	Name  string
	Power int
}

type BatmanRequest struct {
	Hero Hero
}

func allHeros(db HeroDB, parms martini.Params, r render.Render) {
	u := (db).(*heroDB)
	var heros []Hero

	for _, v := range u.heros {
		heros = append(heros, v)
	}
	r.JSON(http.StatusOK, heros)
}

func findHero(db HeroDB, parms martini.Params, r render.Render) {
	u := (db).(*heroDB)

	id, _ := strconv.Atoi(parms["id"])

	hero, ok := u.heros[id]
	if !ok {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, hero)
}

func updateHero(db HeroDB, req *http.Request, r render.Render) {
	u := (db).(*heroDB)

	breq := new(BatmanRequest)
	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()

	err := json.Unmarshal(body, &breq)
	hero := breq.Hero

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	u.heros[hero.Id] = hero
	r.JSON(http.StatusOK, hero)
}

func createHero(db HeroDB, req *http.Request, r render.Render) {
	u := (db).(*heroDB)

	breq := new(BatmanRequest)
	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	err := json.Unmarshal(body, &breq)
	hero := breq.Hero

	u.numHeros++

	hero.Id = u.numHeros

	if err != nil {
		r.Error(http.StatusInternalServerError)
	}

	u.heros[hero.Id] = hero
	r.JSON(http.StatusCreated, hero)
}

func removeHero(db HeroDB, parms martini.Params, r render.Render) {
	u := (db).(*heroDB)
	id, _ := strconv.Atoi(parms["id"])

	delete(u.heros, id)

	r.JSON(http.StatusOK, nil)
}

func main() {

	m := martini.Classic()

	m.Use(render.Renderer())

	// Setup routes
	r := martini.NewRouter()
	r.Get(`/api/heros`, allHeros)
	r.Get(`/api/heros/:id`, findHero)
	r.Post(`/api/heros`, createHero)
	r.Put(`/api/heros/:id`, updateHero)
	r.Delete(`/api/heros/:id`, removeHero)

	// Inject database
	u := &heroDB{make(map[int]Hero), 0}
	m.MapTo(u, (*HeroDB)(nil))

	// Add the router action
	m.Action(r.Handle)

	m.Run()
}
