package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/wim07101993/flies.api/participants"
)

func setupLogging() {
	f, err := os.OpenFile("flies.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Println("Failed to open log file")
	}
	//defer f.Close()

	log.SetOutput(f)
}

func getSettings(path string) Settings {
	var s Settings

	if path == "" {
		s = CreateDefaultSettings()
	} else {
		s = readSettingsFromFile(path)
	}

	log.Println("File set to:", s.ParticipantsDirectory)

	return s
}

func createRouter(c participants.Controller) *httprouter.Router {
	r := httprouter.New()

	r.POST("/api/participants/:"+participants.YearParameter, c.Create)

	r.GET("/api/participants/:"+participants.YearParameter, c.GetAll)
	r.GET("/api/participants/:"+participants.YearParameter+"/:"+participants.IdParameter, c.Get)

	r.PUT("/api/participants/:"+participants.YearParameter+"/:"+participants.IdParameter+"/score", c.UpdateScore)
	r.PUT("/api/participants/:"+participants.YearParameter+"/:"+participants.IdParameter+"/increaseScore", c.IncreaseScore)
	r.PUT("/api/participants/:"+participants.YearParameter+"/:"+participants.IdParameter+"/name", c.UpdateName)
	r.PUT("/api/participants/:"+participants.YearParameter+"/:"+participants.IdParameter+"/decreaseScore", c.DecreaseScore)

	r.DELETE("/api/participants/:"+participants.YearParameter+"/:"+participants.IdParameter, c.Delete)

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
	// set up te loggin
	setupLogging()

	// read program args
	p := getProgramArgs()
	// create settings
	set := getSettings(p.settingFilePath)

	// create service
	s := participants.NewService(set.ParticipantsDirectory)
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
