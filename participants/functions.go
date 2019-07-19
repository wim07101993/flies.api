package participants

import (
	"errors"
	"math"
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
