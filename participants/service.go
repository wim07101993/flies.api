package participants

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

const (
	ParticipantNotFoundErrorMessage = "Participant not found"
	NameAlreadyTakenErrorMessage    = "Name already taken"
	NameCannotBeEmptyErrorMessage   = "The name of a participant cannot be empty"
	BaseFileName                    = "Participants"
	Extension                       = ".json"
)

type Service struct {
	directory string
}

func NewService(directory string) Service {
	if directory[len(directory)-1:] != "/" {
		directory = directory + "/"
	}
	return Service{
		directory: directory,
	}
}

func (pc *Service) Create(year int, p Participant) (Participant, error) {
	ps, err := pc.readFile(year)
	if err != nil {
		return Participant{}, err
	}
	if hasParticipantWithName(p.Name, ps) {
		return Participant{}, errors.New(NameAlreadyTakenErrorMessage)
	}
	if isEmptyOrWhiteSpace(p.Name) {
		return Participant{}, errors.New(NameCannotBeEmptyErrorMessage)
	}

	p.Id, err = getNewId(ps)
	if err != nil {
		return Participant{}, err
	}

	ps = append(ps, p)

	if err = pc.writeFile(year, ps); err != nil {
		return Participant{}, err
	} else {
		return p, nil
	}
}

func (pc *Service) GetAll(year int) ([]Participant, error) {
	ps, err := pc.readFile(year)
	if err != nil {
		return nil, err
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].Score > ps[j].Score })
	return ps, nil
}

func (pc *Service) Get(year int, id uint32) (Participant, error) {
	ps, err := pc.readFile(year)
	if err != nil {
		return Participant{}, err
	}

	i := findParticipant(id, ps)
	if i < 0 {
		return Participant{}, errors.New(ParticipantNotFoundErrorMessage)
	}

	return ps[i], nil
}

func (pc *Service) IncreaseScore(year int, id uint32, amount uint16) (Participant, error) {
	return pc.updateParticipant(year, id, func(p *Participant) error {
		p.Score += amount
		return nil
	})
}

func (pc *Service) DecreaseScore(year int, id uint32, amount uint16) (Participant, error) {
	return pc.updateParticipant(year, id, func(p *Participant) error {
		p.Score -= amount
		return nil
	})
}

func (pc *Service) UpdateScore(year int, id uint32, newScore uint16) (Participant, error) {
	return pc.updateParticipant(year, id, func(p *Participant) error {
		p.Score = newScore
		return nil
	})
}

func (pc *Service) UpdateName(year int, id uint32, newName string) (Participant, error) {
	if isEmptyOrWhiteSpace(newName) {
		return Participant{}, errors.New(NameCannotBeEmptyErrorMessage)
	}

	return pc.updateParticipant(year, id, func(p *Participant) error {
		p.Name = newName
		return nil
	})
}

func (pc *Service) updateParticipant(year int, id uint32, updater func(*Participant) error) (Participant, error) {
	ps, err := pc.readFile(year)
	if err != nil {
		return Participant{}, err
	}

	i := findParticipant(id, ps)
	if i < 0 {
		return Participant{}, errors.New(ParticipantNotFoundErrorMessage)
	}

	if err = updater(&ps[i]); err != nil {
		return Participant{}, err
	}

	if err = pc.writeFile(year, ps); err != nil {
		return Participant{}, err
	}
	return ps[i], nil
}

func (pc *Service) Delete(year int, id uint32) error {
	ps, err := pc.readFile(year)
	if err != nil {
		return err
	}

	i := findParticipant(id, ps)
	if i < 0 {
		return errors.New(ParticipantNotFoundErrorMessage)
	}

	ps = removeAt(ps, i)
	return pc.writeFile(year, ps)
}

func (pc *Service) readFile(year int) ([]Participant, error) {
	ps := []Participant{}
	path := pc.directory + BaseFileName + strconv.Itoa(year) + Extension
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ps, pc.writeFile(year, ps)
	}

	jps, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jps, &ps)
	return ps, err
}

func (pc *Service) writeFile(year int, ps []Participant) error {
	jps, err := json.Marshal(ps)
	if err != nil {
		return err
	}

	path := pc.directory + BaseFileName + strconv.Itoa(year) + Extension
	err = ioutil.WriteFile(path, jps, 0755)
	return err
}
