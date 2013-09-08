package main

// mostly https://github.com/emicklei/go-restful/blob/master/examples/restful-user-resource.go

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"log"
	"net/http"
	"path"
	"strconv"
)

type HeroResource struct {
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

// Global Filter
func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("[global-filter (logger)] %s,%s\n", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}

// WebService Filter
func webserviceLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("[webservice-filter (logger)] %s,%s\n", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}

func (u HeroResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/api/heros").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Filter(webserviceLogging)

	ws.Route(ws.GET("").To(u.allHeros).
		// docs
		Doc("get all heros").
		Writes([]Hero{})) // on the response

	ws.Route(ws.GET("/{hero-id}").To(u.findHero).
		// docs
		Doc("get a hero").
		Param(ws.PathParameter("list-id", "identifier of the hero").DataType("string")).
		Writes(Hero{})) // on the response

	ws.Route(ws.POST("").To(u.createHero).
		// docs
		Doc("create a hero").
		Param(ws.BodyParameter("Hero", "representation of a hero").DataType("main.Hero")).
		Reads(Hero{})) // from the request

	ws.Route(ws.PUT("/{hero-id}").To(u.updateHero).
		// docs
		Doc("update a hero").
		Param(ws.PathParameter("hero-id", "identifier of the hero").DataType("string")).
		Param(ws.BodyParameter("Hero", "representation of a hero").DataType("main.Hero")).
		Reads(Hero{})) // from the request

	ws.Route(ws.DELETE("/{hero-id}").To(u.removeHero).
		// docs
		Doc("delete a hero").
		Param(ws.PathParameter("hero-id", "identifier of the hero").DataType("string")))

	container.Add(ws)
}

func (u HeroResource) allHeros(request *restful.Request, response *restful.Response) {
	var heros []Hero

	for _, v := range u.heros {
		heros = append(heros, v)
	}
	response.WriteEntity(heros)
}

func (u HeroResource) findHero(request *restful.Request, response *restful.Response) {
	ids := request.PathParameter("hero-id")

	id, _ := strconv.Atoi(ids)

	hero, ok := u.heros[id]

	if !ok {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Hero could not be found.")
		return
	}

	response.WriteEntity(hero)
}

func (u *HeroResource) updateHero(request *restful.Request, response *restful.Response) {

	breq := new(BatmanRequest)
	err := request.ReadEntity(&breq)
	hero := breq.Hero

	if err == nil {
		u.heros[hero.Id] = hero
		response.WriteEntity(hero)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (u *HeroResource) createHero(request *restful.Request, response *restful.Response) {

	breq := new(BatmanRequest)
	err := request.ReadEntity(&breq)
	hero := breq.Hero

	u.numHeros++

	hero.Id = u.numHeros

	if err == nil {
		u.heros[hero.Id] = hero
		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(hero)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (u *HeroResource) removeHero(request *restful.Request, response *restful.Response) {
	ids := request.PathParameter("hero-id")

	id, _ := strconv.Atoi(ids)

	delete(u.heros, id)
}

var rootdir = "."

func staticFromPathParam(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(rootdir, req.PathParameter("resource")))
}

func staticFromQueryParam(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(rootdir, req.QueryParameter("resource")))
}

func main() {

	wsContainer := restful.NewContainer()

	restful.Filter(globalLogging)

	herostore := make(map[int]Hero)

	herows := HeroResource{heros: herostore, numHeros: 0}
	herows.Register(wsContainer)

	// static files
	wsStatic := new(restful.WebService)
	wsStatic.Route(wsStatic.GET("/static/{resource}").To(staticFromPathParam))
	wsStatic.Route(wsStatic.GET("/static").To(staticFromQueryParam))
	wsContainer.Add(wsStatic)

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	config := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://localhost:8080",
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/home/dgryski/work/src/cvs/swagger-ui/dist"}
	swagger.RegisterSwaggerService(config, wsContainer)

	log.Printf("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
