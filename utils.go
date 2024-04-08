package main

import (
	"fmt"
	"strings"
)

type image struct {
	registry string
	path     string
	name     string
	tag      string
}

func (i *image) Values() []string {
	return []string{i.registry, i.path, i.name, i.tag}
}

func splitImageFQN(fqn string) (*image, error) {
	if len(fqn) == 0 {
		return nil, nil
	}

	img := &image{}

	// Image FQNs can only have a single colon (:),
	// so splitting on : yields the registry URL +
	// image name and image tag.
	sp := strings.Split(fqn, ":")

	if len(sp) == 1 {
		img.tag = "latest"
	}

	if len(sp) == 2 && len(sp[1]) > 0 {
		img.tag = sp[1]
	}

	// Split the remainder on forward-slashes (/) to
	// get hold of registry, directory, and image name
	rest := strings.Split(sp[0], "/")

	//
	if len(rest) == 0 {
		return nil, fmt.Errorf("Invalid image FQN: %s", fqn)
	}

	if len(rest) == 1 {
		img.name = rest[0]
		return img, nil
	}

	img.registry = rest[0]
	img.name = rest[len(rest)-1]
	img.path = strings.Join(rest[1:len(rest)-1], "/")

	return img, nil
}
