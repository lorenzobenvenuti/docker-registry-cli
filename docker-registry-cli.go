package main

import (
	"fmt"
	"os"

	"github.com/lorenzobenvenuti/docker-registry-client/registry"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app          = kingpin.New("docker-registry-cli", "Docker registry V2 command line interface")
	registryUrl  = app.Flag("registry", "Registy url").Short('r').URL()
	user         = app.Flag("user", "Username").Short('u').Default("").String()
	password     = app.Flag("password", "Password").Short('p').Default("").String()
	debug        = app.Flag("debug", "Debug mode").Bool()
	repositories = app.Command("repositories", "Lists the registry repositories")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case repositories.FullCommand():
		hub, err := registry.New((*registryUrl).String(), *user, *password)
		if !(*debug) {
			hub.Logf = func(format string, args ...interface{}) {}
		}
		if err != nil {
			panic(err)
		}
		repositories, err := hub.Repositories()
		if err != nil {
			panic(err)
		}
		for _, repository := range repositories {
			fmt.Printf("%s\n", repository)
		}
	}
}
