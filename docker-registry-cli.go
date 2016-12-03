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
	search       = app.Command("search", "Search the registry")
	expression   = search.Arg("expression", "Search expression").String()
	images       = app.Command("images", "Lists the registry images (repository:tag)")
	delete       = app.Command("delete", "Deletes an image")
	image        = delete.Arg("image", "Image to delete (name:tag)").String()
)

func getRegistry() *registry.Registry {
	hub, err := registry.New((*registryUrl).String(), *user, *password)
	if err != nil {
		panic(err)
	}
	if !(*debug) {
		hub.Logf = func(format string, args ...interface{}) {}
	}
	return hub
}

func print(items []string) {
	for _, item := range items {
		fmt.Printf("%s\n", item)
	}
}

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	api := NewRegistryApi(getRegistry())
	switch cmd {
	case repositories.FullCommand():
		print(api.GetAllRepositories())
	case images.FullCommand():
		print(api.GetAllImages())
	case search.FullCommand():
		print(api.SearchImages(*expression))
	case delete.FullCommand():
		api.DeleteImage(*image)
	}
}
