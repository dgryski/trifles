package main

import (
	"encoding/json"
	"github.com/cznic/kv"
	"github.com/dustin/go-nma"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/gorilla/schema"
	"github.com/jdiez17/go-pushover"
	"github.com/xconstruct/go-pushbullet"
	"log"
	"net/http"
	"path"
	"time"
)

type DistributionListResource struct {
	datastore *kv.DB
}

type DistributionList struct {
	Id         string
	Nma        NMAConfig
	Pushbullet PBConfig
	Pushover   POConfig
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

func (u DistributionListResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/api/lists").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Filter(webserviceLogging)

	ws.Route(ws.GET("").To(u.allLists).
		// docs
		Doc("get all lists").
		Writes([]DistributionList{})) // on the response

	ws.Route(ws.GET("/{list-id}").To(u.findList).
		// docs
		Doc("get a list").
		Param(ws.PathParameter("list-id", "identifier of the list").DataType("string")).
		Writes(DistributionList{})) // on the response

	ws.Route(ws.POST("").To(u.updateList).
		// docs
		Doc("update a list").
		Param(ws.BodyParameter("List", "representation of a list").DataType("main.List")).
		Reads(DistributionList{})) // from the request

	ws.Route(ws.PUT("/{list-id}").To(u.createList).
		// docs
		Doc("create a list").
		Param(ws.PathParameter("list-id", "identifier of the list").DataType("string")).
		Param(ws.BodyParameter("List", "representation of a list").DataType("main.List")).
		Reads(DistributionList{})) // from the request

	ws.Route(ws.DELETE("/{list-id}").To(u.removeList).
		// docs
		Doc("delete a list").
		Param(ws.PathParameter("list-id", "identifier of the list").DataType("string")))

	container.Add(ws)
}

func (u DistributionListResource) allLists(request *restful.Request, response *restful.Response) {
	var lists []DistributionList
	e, _ := u.datastore.SeekFirst()

	for {
		_, v, err := e.Next()
		if err != nil {
			break
		}

		var d DistributionList
		json.Unmarshal(v, &d)
		lists = append(lists, d)
	}

	response.WriteEntity(lists)
}

func (u DistributionListResource) findList(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("list-id")

	dst, err := u.datastore.Get(nil, []byte(id))

	if dst == nil || err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "List could not be found.")
	} else {
		var l DistributionList
		json.Unmarshal(dst, &l)
		response.WriteEntity(l)
	}
}

func (u *DistributionListResource) updateList(request *restful.Request, response *restful.Response) {
	l := new(DistributionList)
	err := request.ReadEntity(&l)
	if err == nil {
		j, _ := json.Marshal(l)
		u.datastore.Set([]byte(l.Id), j)
		response.WriteEntity(l)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (u *DistributionListResource) createList(request *restful.Request, response *restful.Response) {
	l := DistributionList{Id: request.PathParameter("list-id")}
	err := request.ReadEntity(&l)
	if err == nil {
		j, _ := json.Marshal(l)
		u.datastore.Set([]byte(l.Id), j)

		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(l)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (u *DistributionListResource) removeList(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("list-id")
	u.datastore.Delete([]byte(id))
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

type Getter interface {
	Get(dst, src []byte) ([]byte, error)
}

func loadDistributionList(datastore Getter, list []byte) (DistributionList, error) {

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

	n := pushover.Notification{
		Title:     pn.Title,
		Message:   pn.Body,
		Timestamp: time.Now(),
	}

	for _, u := range poconf.Users {
		c := pushover.New(u, poconf.APIKey)
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

var decoder = schema.NewDecoder()

func pushHandler(request *restful.Request, response *restful.Response, datastore *kv.DB) {

	err := request.Request.ParseForm()
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	pushNotification := new(PushNotification)

	err = decoder.Decode(pushNotification, request.Request.PostForm)

	log.Println(pushNotification)

	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	if pushNotification.List == "" {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
		return
	}

	targets, err := loadDistributionList(datastore, []byte(pushNotification.List))
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	done := make(chan struct{})

	go sendNMA(pushNotification, targets.Nma, done)
	go sendPushover(pushNotification, targets.Pushover, done)
	go sendPushBullet(pushNotification, targets.Pushbullet, done)

	for i := 0; i < 3; i++ {
		<-done
	}
}

var rootdir = "."

func staticFromPathParam(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(rootdir+"/static/", req.PathParameter("resource")))
}

func staticFromQueryParam(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(rootdir+"/static/", req.QueryParameter("resource")))
}

func main() {

	datastore, err := kv.Open("mpush.kvdb", &kv.Options{})

	if err != nil {
		log.Fatal(err)
	}
	defer datastore.Close()

	wsContainer := restful.NewContainer()

	restful.Filter(globalLogging)

	dsl := DistributionListResource{datastore}
	dsl.Register(wsContainer)

	wsPush := new(restful.WebService)

	wsPush.Path("/api/push")
	wsPush.Route(wsPush.POST("").
		Consumes("application/x-www-form-urlencoded").
		To(func(request *restful.Request, response *restful.Response) { pushHandler(request, response, datastore) }).
		// docs
		Doc("push to a distribution list").
		Param(wsPush.BodyParameter("List", "list to send to").DataType("string")).
		Param(wsPush.BodyParameter("Title", "title of notification").DataType("string")).
		Param(wsPush.BodyParameter("Body", "body of notification").DataType("string")).
		Reads(PushNotification{})) // from the request
	wsContainer.Add(wsPush)

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
