package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/wim07101993/fly-swatting-contest.api/participants"
)

func createRouter(c participants.Controller) *httprouter.Router {
	r := httprouter.New()

	r.POST("/api/participants", c.Create)

	r.GET("/api/participants/", c.GetAll)
	r.GET("/api/participants/:"+participants.IdParameter, c.Get)

	r.PUT("/api/participants/:"+participants.IdParameter+"/score", c.UpdateScore)
	r.PUT("/api/participants/:"+participants.IdParameter+"/name", c.UpdateName)
	r.PUT("/api/participants/:"+participants.IdParameter+"/increaseScore", c.IncreaseScore)
	r.PUT("/api/participants/:"+participants.IdParameter+"/decreaseScore", c.DecreaseScore)

	r.DELETE("/api/participants/:"+participants.IdParameter, c.Delete)

	return r
}

func addCorsToHandler(h http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler(h)
}

func main() {
	// create settings
	set := CreateDefaultSettings()

	// create service
	s := participants.NewService(set.ParticiPantsFilePath)
	// create controller
	c := participants.NewController(s)
	// create router
	r := createRouter(c)

	// add cors to router
	handler := addCorsToHandler(r)

	// start serving
	log.Println("Start listening at", set.IpAddress+":"+set.PortNumber)
	log.Fatal(http.ListenAndServe(set.IpAddress+":"+set.PortNumber, handler))
}
