package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
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
	Port  int      `json:"port"`
}

func (c *CreatedContainer) GetContainer(id string) *Container {
	return &Container{
		ID:      id,
		Name:    c.Name,
		Image:   c.Image,
		Cmd:     c.Cmd,
		Status:  "created",
		StartAt: time.Now().Format(timeFormat),
	}
}

type Container struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Cmd     []string `json:"cmd"`
	Port    int      `json:"port"`
	Status  string   `json:"status"`
	Address string   `json:"address"`
	StartAt string   `json:"startAt"`
}

func ContainerConfig(con *CreatedContainer) *container.Config {
	labels := map[string]string{
		AppPateformLabel: AppPateformValue,
		AppNameLabel:     con.Name,
	}

	config := &container.Config{
		Image:  con.Image,
		Cmd:    con.Cmd,
		Labels: labels,
	}

	if con.Port > 0 {
		config.ExposedPorts = nat.PortSet{
			getContainerPort(con.Port): {},
		}
	}

	return config
}

func ContainerHostConfig(con *CreatedContainer) *container.HostConfig {
	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name:              "on-failure",
			MaximumRetryCount: 5,
		},
	}
	return hostConfig
}

func getContainerPort(port int) nat.Port {
	return nat.Port(fmt.Sprintf("%d/tcp", port))
}

func DockerContainerToContainer(container types.Container) Container {
	var name string
	if len(container.Names) > 0 {
		name = strings.TrimPrefix(container.Names[0], "/")
	}

	var port int
	var address string
	for _, cp := range container.Ports {
		port = int(cp.PrivatePort)
		if len(cp.IP) > 0 {
			address = fmt.Sprintf("%s:%d", cp.IP, port)
		}
	}

	return Container{
		ID:      container.ID,
		Name:    name,
		Image:   container.Image,
		Cmd:     []string{container.Command},
		Port:    port,
		Address: address,
		Status:  container.State,
		StartAt: time.Unix(container.Created, 0).Format(timeFormat),
	}
}

func DockerContainerJSONToContainer(container types.ContainerJSON) Container {
	created, _ := time.Parse("2006-01-02T15:04:05.999999999Z", container.Created)
	ip := container.NetworkSettings.IPAddress
	var port string
	for k := range container.NetworkSettings.Ports {
		port = k.Port()
		break
	}
	var address string
	if len(ip) > 0 && len(port) > 0 {
		address = ip + ":" + port
	}

	cp, _ := strconv.Atoi(port)
	return Container{
		ID:      container.ID,
		Name:    container.Name,
		Image:   container.Image,
		Cmd:     container.Config.Cmd,
		Port:    cp,
		Address: address,
		Status:  container.State.Status,
		StartAt: created.Format(timeFormat),
	}
}
