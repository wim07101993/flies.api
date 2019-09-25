package participants

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sort"
)

const (
	ParticipantNotFoundErrorMessage = "Participant not found"
	NameAlreadyTakenErrorMessage    = "Name already taken"
	NameCannotBeEmptyErrorMessage   = "The name of a participant cannot be empty"
)

type Service struct {
	filePath string
}

func NewService(filePath string) Service {
	return Service{
		filePath: filePath,
	}
}

func (pc *Service) Create(p Participant) (Participant, error) {
	ps, err := pc.readFile()
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

	if err = pc.writeFile(ps); err != nil {
		return Participant{}, err
	} else {
		return p, nil
	}
}

func (pc *Service) GetAll() ([]Participant, error) {
	ps, err := pc.readFile()
	if err != nil {
		return nil, err
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].Score > ps[j].Score })
	return ps, nil
}

func (pc *Service) Get(id uint32) (Participant, error) {
	ps, err := pc.readFile()
	if err != nil {
		return Participant{}, err
	}

	i := findParticipant(id, ps)
	if i < 0 {
		return Participant{}, errors.New(ParticipantNotFoundErrorMessage)
	}

	return ps[i], nil
}

func (pc *Service) IncreaseScore(id uint32, amount uint16) (Participant, error) {
	return pc.updateParticipant(id, func(p *Participant) error {
		p.Score += amount
		return nil
	})
}

func (pc *Service) DecreaseScore(id uint32, amount uint16) (Participant, error) {
	return pc.updateParticipant(id, func(p *Participant) error {
		p.Score -= amount
		return nil
	})
}

func (pc *Service) UpdateScore(id uint32, newScore uint16) (Participant, error) {
	return pc.updateParticipant(id, func(p *Participant) error {
		p.Score = newScore
		return nil
	})
}

func (pc *Service) UpdateName(id uint32, newName string) (Participant, error) {
	if isEmptyOrWhiteSpace(newName) {
		return Participant{}, errors.New(NameCannotBeEmptyErrorMessage)
	}

	return pc.updateParticipant(id, func(p *Participant) error {
		p.Name = newName
		return nil
	})
}

func (pc *Service) updateParticipant(id uint32, updater func(*Participant) error) (Participant, error) {
	ps, err := pc.readFile()
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

	if err = pc.writeFile(ps); err != nil {
		return Participant{}, err
	}
	return ps[i], nil
}

func (pc *Service) Delete(id uint32) error {
	ps, err := pc.readFile()
	if err != nil {
		return err
	}

	i := findParticipant(id, ps)
	if i < 0 {
		return errors.New(ParticipantNotFoundErrorMessage)
	}

	ps = removeAt(ps, i)
	return pc.writeFile(ps)
}

func (pc *Service) readFile() ([]Participant, error) {
	ps := []Participant{}
	if _, err := os.Stat(pc.filePath); os.IsNotExist(err) {
		return ps, pc.writeFile(ps)
	}

	jps, err := ioutil.ReadFile(pc.filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jps, &ps)
	return ps, err
}

func (pc *Service) writeFile(ps []Participant) error {
	jps, err := json.Marshal(ps)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(pc.filePath, jps, 755)
	return err
}
