package participants

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"unicode"
)

const (
	RanOutOfIds = "The maximum number of players is reached..."
)

func hasParticipantWithName(name string, participants []Participant) bool {
	for _, participant := range participants {
		if participant.Name == name {
			return true
		}
	}
	return false
}

func isEmptyOrWhiteSpace(s string) bool {
	if s == "" {
		return true
	}

	for _, c := range s {
		if !unicode.IsSpace(c) {
			return false
		}
	}

	return true
}

func getNewId(participants []Participant) (uint32, error) {
	var maxId uint32 = 0
	var ids = map[uint32]struct{}{}
	for _, p := range participants {
		ids[p.Id] = struct{}{}
		if p.Id > maxId {
			maxId = p.Id
		}
	}

	if maxId != uint32(math.MaxUint32) {
		return maxId + 1, nil
	}

	var i uint32
	for i = 0; i < math.MaxUint32; i++ {
		if _, ok := ids[i]; !ok {
			return uint32(i), nil
		}
	}

	return 0, errors.New(RanOutOfIds)
}

func findParticipant(id uint32, participants []Participant) int {
	for i, participant := range participants {
		if participant.Id == id {
			return i
		}
	}
	return -1
}

func removeAt(ps []Participant, i int) []Participant {
	ps[i] = ps[len(ps)-1]
	return ps[:len(ps)-1]
}

func checkError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	errMes := err.Error()
	log.Println("\tError:", errMes)

	if errMes == ParticipantNotFoundErrorMessage {
		http.Error(w, errMes, http.StatusNotFound)
	} else if errMes == NameAlreadyTakenErrorMessage ||
		errMes == BadJsonErrorMessage ||
		errMes == NameCannotBeEmptyErrorMessage {
		http.Error(w, errMes, http.StatusBadRequest)
	} else {
		http.Error(w, errMes, http.StatusInternalServerError)
	}

	return true
}

func writeJson(w http.ResponseWriter, v interface{}) {
	jv, err := json.Marshal(v)
	if checkError(w, err) {
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(jv)
}
