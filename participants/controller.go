package participants

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const (
	BadJsonErrorMessage   = "The given object could not be interpreted by the server."
	BadAmountErrorMessage = "The given amount could not be interpreted by the server."
	BadScoreErrorMessage  = "The given score could not be interpreted by the server."
	NameParameter         = "name"
	IdParameter           = "id"
	AmountParameter       = "amount"
	ScoreParameter        = "score"
	NewNameParamter       = "newName"
	YearParameter         = "year"
)

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

func (pc *Controller) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest("Create", r)
	jp, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if checkError(w, err) {
		return
	}

	var p Participant
	err = json.Unmarshal(jp, &p)
	if err != nil {
		checkError(w, errors.New(BadJsonErrorMessage))
		return
	}

	year := getYear(ps)
	p, err = pc.service.Create(year, p)
	if checkError(w, err) {
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJson(w, p)
}

func (pc *Controller) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	logRequest("GetAll", r)
	year := getYear(p)
	ps, err := pc.service.GetAll(year)
	if checkError(w, err) {
		return
	}

	writeJson(w, ps)
}

func (pc *Controller) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest("Get", r)

	sId := ps.ByName(IdParameter)
	id, err := strconv.ParseUint(sId, 10, 32)
	if checkError(w, err) {
		return
	}

	year := getYear(ps)
	p, err := pc.service.Get(year, uint32(id))
	if checkError(w, err) {
		return
	}

	writeJson(w, p)
}

func (pc *Controller) IncreaseScore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest("IncreaseScore", r)

	sId := ps.ByName(IdParameter)
	id, err := strconv.ParseUint(sId, 10, 32)
	if checkError(w, err) {
		return
	}

	sAmount := r.URL.Query().Get(AmountParameter)
	amount, err := strconv.ParseUint(sAmount, 10, 16)
	if err != nil {
		checkError(w, errors.New(BadAmountErrorMessage))
		return
	}

	year := getYear(ps)
	p, err := pc.service.IncreaseScore(year, uint32(id), uint16(amount))
	if checkError(w, err) {
		return
	}

	writeJson(w, p)
}

func (pc *Controller) DecreaseScore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest("DecreaseScore", r)

	sId := ps.ByName(IdParameter)
	id, err := strconv.ParseUint(sId, 10, 32)
	if checkError(w, err) {
		return
	}

	sAmount := r.URL.Query().Get(AmountParameter)
	amount, err := strconv.ParseUint(sAmount, 10, 16)
	if err != nil {
		checkError(w, errors.New(BadAmountErrorMessage))
		return
	}

	year := getYear(ps)
	p, err := pc.service.DecreaseScore(year, uint32(id), uint16(amount))
	if checkError(w, err) {
		return
	}

	writeJson(w, p)
}

func (pc *Controller) UpdateScore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest("UpdateScore", r)

	sId := ps.ByName(IdParameter)
	id, err := strconv.ParseUint(sId, 10, 32)
	if checkError(w, err) {
		return
	}

	sScore := r.URL.Query().Get(ScoreParameter)
	score, err := strconv.ParseUint(sScore, 10, 16)
	if err != nil {
		checkError(w, errors.New(BadScoreErrorMessage))
		return
	}

	year := getYear(ps)
	p, err := pc.service.UpdateScore(year, uint32(id), uint16(score))
	if checkError(w, err) {
		return
	}

	writeJson(w, p)
}

func (pc *Controller) UpdateName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest("UpdateName", r)

	sId := ps.ByName(IdParameter)
	id, err := strconv.ParseUint(sId, 10, 32)
	if checkError(w, err) {
		return
	}

	year := getYear(ps)
	newName := r.URL.Query().Get(NewNameParamter)
	p, err := pc.service.UpdateName(year, uint32(id), newName)
	if checkError(w, err) {
		return
	}

	writeJson(w, p)
}

func (pc *Controller) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest("Delete", r)

	sId := ps.ByName(IdParameter)
	id, err := strconv.ParseUint(sId, 10, 32)
	if checkError(w, err) {
		return
	}

	year := getYear(ps)
	err = pc.service.Delete(year, uint32(id))
	if checkError(w, err) {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted"))
}

func logRequest(method string, r *http.Request) {
	log.Println("Request:", method, "from", r.RemoteAddr)
	log.Println("\tUrl:", r.RequestURI)
}
