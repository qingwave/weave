package model

import (
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

const (
	AppNameLabel     = `app`
	AppPateformLabel = `weave.io`
	AppPateformValue = `true`

	timeFormat = `2006-01-02T15:04:05`
)

type CreatedContainer struct {
	Name  string   `json:"name"`
	Image string   `json:"image"`
	Cmd   []string `json:"cmd"`
}

func (c *CreatedContainer) GetContainer(id string) *Container {
	return &Container{
		ID:     id,
		Name:   c.Name,
		Image:  c.Image,
		Cmd:    c.Cmd,
		Status: "",
	}
}

type Container struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Cmd     []string `json:"cmd"`
	Status  string   `json:"status"`
	StartAt string   `json:"startAt"`
}

func ContainerConfig(name, image string, cmd []string) *container.Config {
	labels := map[string]string{
		AppPateformLabel: AppPateformValue,
		AppNameLabel:     name,
	}
	return &container.Config{
		Image:  image,
		Cmd:    cmd,
		Labels: labels,
	}
}

func DockerContainerToContainer(container types.Container) Container {
	var name string
	if len(container.Names) > 0 {
		name = strings.TrimPrefix(container.Names[0], "/")
	}
	return Container{
		ID:      container.ID,
		Name:    name,
		Image:   container.Image,
		Cmd:     []string{container.Command},
		Status:  container.State,
		StartAt: time.Unix(container.Created, 0).Format(timeFormat),
	}
}

func DockerContainerJSONToContainer(container types.ContainerJSON) Container {
	created, _ := time.Parse("2006-01-02T15:04:05.999999999Z", container.Created)
	return Container{
		ID:      container.ID,
		Name:    container.Name,
		Image:   container.Image,
		Cmd:     container.Config.Cmd,
		Status:  container.State.Status,
		StartAt: created.Format(timeFormat),
	}
}
