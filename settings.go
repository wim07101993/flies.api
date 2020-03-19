package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Settings struct {
	ParticipantsDirectory string `json:"participantsDirectory"`
	PortNumber            string `json:"portNumber"`
	IpAddress             string `json:"ipAddress"`
}

func CreateDefaultSettings() Settings {
	return Settings{
		ParticipantsDirectory: "participants",
		PortNumber:            "5000",
		IpAddress:             "",
	}
}

func readSettingsFromFile(filePath string) Settings {
	log.Println("Reading settings at:", filePath)
	bsettings, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var settings Settings
	err = json.Unmarshal(bsettings, &settings)

	if err != nil {
		panic(err)
	}

	return settings
}
