package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Settings struct {
	ParticiPantsFilePath string `json:"particiPantsFilePath"`
	PortNumber           string `json:"portNumber"`
	IpAddress            string `json:"ipAddress"`
}

func CreateDefaultSettings() Settings {
	return Settings{
		ParticiPantsFilePath: "participants.json",
		PortNumber:           "5000",
		IpAddress:            "",
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
