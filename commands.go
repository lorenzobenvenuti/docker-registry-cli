package main

import (
	"fmt"
	"strings"

	"github.com/lorenzobenvenuti/docker-registry-client/registry"
)

func getRepositories(hub *registry.Registry) []string {
	repositories, err := hub.Repositories()
	if err != nil {
		panic(err)
	}
	return repositories
}

func printRepositories(hub *registry.Registry) {
	for _, repository := range getRepositories(hub) {
		fmt.Printf("%s\n", repository)
	}
}

func getTags(hub *registry.Registry, repository string) []string {
	tags, err := hub.Tags(repository)
	if err != nil {
		panic(err)
	}
	return tags
}

func getImages(hub *registry.Registry) []string {
	images := []string{}
	for _, repository := range getRepositories(hub) {
		for _, tag := range getTags(hub, repository) {
			images = append(images, fmt.Sprintf("%s:%s", repository, tag))
		}
	}
	return images
}

func printImages(hub *registry.Registry) {
	for _, image := range getImages(hub) {
		fmt.Printf("%s\n", image)
	}
}

func searchExpression(hub *registry.Registry, expression string) {
	for _, image := range getImages(hub) {
		if strings.Index(strings.ToLower(image), strings.ToLower(expression)) != -1 {
			fmt.Printf("%s\n", image)
		}
	}
}

func deleteManifest(hub *registry.Registry, image string) {
	tokens := strings.Split(image, ":")
	if len(tokens) != 2 {
		panic("Invalid image, must be in the form image:tag")
	}
	digest, err := hub.ManifestDigest(tokens[0], tokens[1])
	if err != nil {
		panic(err)
	}
	err = hub.DeleteManifest(tokens[0], digest)
	if err != nil {
		panic(err)
	}
}
