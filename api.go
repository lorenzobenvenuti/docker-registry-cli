package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lorenzobenvenuti/docker-registry-client/registry"
)

type RegistryApi interface {
	GetAllRepositories() ([]string, error)
	GetAllImages() ([]string, error)
	SearchImages(key string) ([]string, error)
	DeleteImage(image string) error
}

type registryApiImpl struct {
	hub *registry.Registry
}

func NewRegistryApi(hub *registry.Registry) RegistryApi {
	return &registryApiImpl{hub}
}

func (r *registryApiImpl) GetAllRepositories() ([]string, error) {
	return r.hub.Repositories()
}

func (r *registryApiImpl) GetAllImages() ([]string, error) {
	images := []string{}
	repositories, err := r.GetAllRepositories()
	if err != nil {
		return nil, err
	}
	for _, repository := range repositories {
		tags, err := r.getTags(repository)
		if err != nil {
			return nil, err
		}
		for _, tag := range tags {
			images = append(images, fmt.Sprintf("%s:%s", repository, tag))
		}
	}
	return images, nil
}

func (r *registryApiImpl) getTags(repository string) ([]string, error) {
	return r.hub.Tags(repository)
}

func (r *registryApiImpl) SearchImages(key string) ([]string, error) {
	filteredImages := []string{}
	allImages, err := r.GetAllImages()
	if err != nil {
		return nil, err
	}
	for _, image := range allImages {
		if strings.Index(strings.ToLower(image), strings.ToLower(key)) != -1 {
			filteredImages = append(filteredImages, image)
		}
	}
	return filteredImages, nil
}

func (r *registryApiImpl) DeleteImage(image string) error {
	tokens := strings.Split(image, ":")
	if len(tokens) != 2 {
		return errors.New("Invalid image, must be in the form image:tag")
	}
	digest, err := r.hub.ManifestDigest(tokens[0], tokens[1])
	if err != nil {
		return err
	}
	err = r.hub.DeleteManifest(tokens[0], digest)
	if err != nil {
		return err
	}
	return nil
}
