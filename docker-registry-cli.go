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

func print(items []string) {
	for _, item := range items {
		fmt.Printf("%s\n", item)
	}
}

func processOutput(items []string, err error) {
	if err != nil {
		handleError(err)
		return
	}
	print(items)
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, err.Error())
	os.Exit(1)
}

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	registry, err := getRegistry()
	if err != nil {
		handleError(err)
		return
	}
	api := NewRegistryApi(registry)
	switch cmd {
	case repositories.FullCommand():
		processOutput(api.GetAllRepositories())
	case images.FullCommand():
		processOutput(api.GetAllImages())
	case search.FullCommand():
		processOutput(api.SearchImages(*expression))
	case delete.FullCommand():
		err := api.DeleteImage(*image)
		if err != nil {
			panic(err)
		}
	}
}
