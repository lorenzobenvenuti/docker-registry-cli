package main

import (
	"fmt"
	"strings"

	"github.com/lorenzobenvenuti/docker-registry-client/registry"
)

type RegistryApi interface {
	GetAllRepositories() []string
	GetAllImages() []string
	SearchImages(key string) []string
	DeleteImage(image string)
}

type registryApiImpl struct {
	hub *registry.Registry
}

func NewRegistryApi(hub *registry.Registry) RegistryApi {
	return &registryApiImpl{hub}
}

func (r *registryApiImpl) GetAllRepositories() []string {
	repositories, err := r.hub.Repositories()
	if err != nil {
		panic(err)
	}
	return repositories
}

func (r *registryApiImpl) GetAllImages() []string {
	images := []string{}
	for _, repository := range r.GetAllRepositories() {
		for _, tag := range r.getTags(repository) {
			images = append(images, fmt.Sprintf("%s:%s", repository, tag))
		}
	}
	return images
}

func (r *registryApiImpl) getTags(repository string) []string {
	tags, err := r.hub.Tags(repository)
	if err != nil {
		panic(err)
	}
	return tags
}

func (r *registryApiImpl) SearchImages(key string) []string {
	images := []string{}
	for _, image := range r.GetAllImages() {
		if strings.Index(strings.ToLower(image), strings.ToLower(key)) != -1 {
			images = append(images, image)
		}
	}
	return images
}

func (r *registryApiImpl) DeleteImage(image string) {
	tokens := strings.Split(image, ":")
	if len(tokens) != 2 {
		panic("Invalid image, must be in the form image:tag")
	}
	digest, err := r.hub.ManifestDigest(tokens[0], tokens[1])
	if err != nil {
		panic(err)
	}
	err = r.hub.DeleteManifest(tokens[0], digest)
	if err != nil {
		panic(err)
	}
}
