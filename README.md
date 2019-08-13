# flies.api
The api that is used by the flies.browser and flies.webapp repos


## Installation

### Windows
Go to the buid/windows directory and copy the flies.api.exe file to where you want the program to run be installed.

### Linux
A linux build is not yet available, for now you will have to build it yourself (see below)


## Build
It is also possible to build the program from source. 

To build the code, Golang will need to be installed. The method to install Golang can be found at [the installation page of the 
Golang language](https://golang.org/doc/install).

After the installation of Golang, the code and all of it's dependencies need to be downloaded in the ```GOPATH```. The easiest 
way to do this is by running the command ```go get github.com/wim07101993/flies.api```.

Once the packages are downloaded, it is possible to run the application by navigating to the foder where the source code is located
(```$GOPATH/src/github.com/wim07101993/flies.api```) and running the ```go run .``` command. Then the server will automatically start
on localhost:5000.
To build the source code to have an executable, run the ```go build``` command.

## Usage

The executable can be run without any paramters, in that case the default settings are used. These settings are:
```
{
  "participantsFilePath": "participants.json",
  "ipAddress": "",
  "portNumber": "5000"
}
```

A settings file can also be specified by running the executable in a terminal and appending the command with the path 
where the settings file can be found: ```./flies.api settings/file/path/settings.json```.
This settings file must be formatted in json like the example above.
