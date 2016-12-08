package main

import (
	"fmt"
	"os"

	"github.com/heroku/docker-registry-client/registry"
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

func getRegistry() (*registry.Registry, error) {
	hub, err := registry.New((*registryUrl).String(), *user, *password)
	if err != nil {
		return nil, err
	}
	if !(*debug) {
		hub.Logf = func(format string, args ...interface{}) {}
	}
	return hub, nil
}

func handleResult(items []string, err error) {
	failIfErrorIsNotNil(err)
	for _, item := range items {
		fmt.Printf("%s\n", item)
	}
}

func failIfErrorIsNotNil(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	registry, err := getRegistry()
	failIfErrorIsNotNil(err)
	api := NewRegistryApi(registry)
	switch cmd {
	case repositories.FullCommand():
		handleResult(api.GetAllRepositories())
	case images.FullCommand():
		handleResult(api.GetAllImages())
	case search.FullCommand():
		handleResult(api.SearchImages(*expression))
	case delete.FullCommand():
		failIfErrorIsNotNil(api.DeleteImage(*image))
	}
}
