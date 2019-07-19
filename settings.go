package main

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
