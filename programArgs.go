package main

import "os"

type programArgs struct {
	settingFilePath string
}

func getProgramArgs() programArgs {
	p := programArgs{}

	if len(os.Args) > 1 {
		p.settingFilePath = os.Args[1]
	}

	return p
}
