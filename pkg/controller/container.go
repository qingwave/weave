package controller

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/qingwave/weave/pkg/common"
	"github.com/qingwave/weave/pkg/library/docker"
	"github.com/qingwave/weave/pkg/model"
	"github.com/qingwave/weave/pkg/utils/trace"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ContainerController struct {
	client *docker.Client
}

func NewContainerController(client *docker.Client) Controller {
	return &ContainerController{
		client: client,
	}
}

// @Summary Create container
// @Description Create container
// @Accept json
// @Produce json
// @Tags container
// @Security JWT
// @Param container body model.CreatedContainer true "container info"
// @Success 200 {object} common.Response{data=model.Container}
// @Router /api/v1/containers [post]
func (con *ContainerController) Create(c *gin.Context) {
	ccon := new(model.CreatedContainer)
	if err := c.BindJSON(ccon); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.TraceStep(c, "start create container", trace.Field{"name", ccon.Name})
	defer common.TraceStep(c, "create container done", trace.Field{"name", ccon.Name})

	cid, err := con.runContainer(c.Request.Context(), ccon)
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, ccon.GetContainer(cid))
}

// @Summary Get container
// @Description Get container
// @Produce json
// @Tags container
// @Security JWT
// @Param id path string true "container id"
// @Success 200 {object} common.Response{data=model.Container}
// @Router /api/v1/containers/{id} [get]
func (con *ContainerController) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ResponseFailed(c, http.StatusBadRequest, errors.New("empty container id"))
		return
	}

	resp, err := con.client.ContainerInspect(c.Request.Context(), id)
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, model.DockerContainerJSONToContainer(resp))
}

// @Summary List container
// @Description List container
// @Produce json
// @Tags container
// @Security JWT
// @Success 200 {object} common.Response{data=[]model.Container}
// @Router /api/v1/containers [get]
func (con *ContainerController) List(c *gin.Context) {
	common.TraceStep(c, "start list container")
	items, err := con.client.ContainerList(c.Request.Context(), types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("label", model.AppPateformLabel)),
	})
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.TraceStep(c, "list container done")

	containers := make([]model.Container, 0)
	for _, item := range items {
		if item.ID == "" || len(item.Names) == 0 {
			continue
		}
		containers = append(containers, model.DockerContainerToContainer(item))
	}

	common.ResponseSuccess(c, containers)
}

// @Summary Operate container
// @Description Operate container
// @Produce json
// @Tags container
// @Security JWT
// @Param id path string true "container id"
// @Param verb    query     string  true  "verb: start/stop/restart"
// @Success 200 {object} common.Response
// @Router /api/v1/containers/{id} [post]
func (con *ContainerController) Operate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ResponseFailed(c, http.StatusBadRequest, errors.New("empty container id"))
		return
	}

	ctx := c.Request.Context()
	verb := c.Query("verb")
	var err error
	switch verb {
	case "start":
		err = con.client.ContainerStart(ctx, id, types.ContainerStartOptions{})
	case "restart":
		err = con.client.ContainerRestart(ctx, id, nil)
	case "stop":
		err = con.client.ContainerStop(ctx, id, nil)
	default:
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("invaild verb：%s", verb))
		return
	}

	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, fmt.Sprintf("%s container %s success", verb, id))
}

// @Summary Update container
// @Description Update container
// @Accept json
// @Produce json
// @Tags container
// @Security JWT
// @Param container body model.CreatedContainer true "container info"
// @Param id path string true "container id"
// @Success 200 {object} common.Response{data=model.Container}
// @Router /api/v1/containers/:id [put]
func (con *ContainerController) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ResponseFailed(c, http.StatusBadRequest, errors.New("empty container id"))
		return
	}

	ctx := c.Request.Context()
	con.client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		Force: true,
	})

	ccon := new(model.CreatedContainer)
	if err := c.BindJSON(ccon); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	cid, err := con.runContainer(ctx, ccon)
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, ccon.GetContainer(cid))
}

// @Summary Delete container
// @Description Delete container
// @Produce json
// @Tags container
// @Security JWT
// @Param id path string true "container id"
// @Success 200 {object} common.Response
// @Router /api/v1/containers/{id} [delete]
func (con *ContainerController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ResponseFailed(c, http.StatusBadRequest, errors.New("empty container id"))
		return
	}

	err := con.client.ContainerRemove(c.Request.Context(), id, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, fmt.Sprintf("delete container %s success", id))
}

// @Summary Get container log
// @Description Get container log
// @Produce json
// @Tags container
// @Security JWT
// @Param id path string true "container id"
// @Success 200 {string}  string    ""
// @Param follow    query  bool  false  "follow log"
// @Param tail    query  string  false  "tail log all or number"
// @Router /api/v1/containers/{id}/log [get]
func (con *ContainerController) Log(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ResponseFailed(c, http.StatusBadRequest, errors.New("empty container id"))
		return
	}

	reader, err := con.client.ContainerLogs(c.Request.Context(), id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Tail:       c.Query("tail"),
		Follow:     c.Query("follow") == "true",
	})
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		_, err := c.Writer.Write(append(scanner.Bytes(), '\n'))
		c.Writer.Flush()
		if err != nil {
			return
		}
	}
}

// @Summary Proxy container
// @Description proxy container
// @Tags container
// @Security JWT
// @Param id path string true "container id"
// @Param shell query string  false  "shell, sh or bash"
// @Success 200 {string}  string    ""
// @Router /api/v1/containers/{id}/proxy [get]
func (con *ContainerController) Proxy(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ResponseFailed(c, http.StatusBadRequest, errors.New("empty container id"))
		return
	}

	resp, err := con.client.ContainerInspect(c.Request.Context(), id)
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	container := model.DockerContainerJSONToContainer(resp)

	getPath := func(raw string) string {
		paths := strings.SplitAfterN(raw, "proxy", 2)
		if len(paths) != 2 {
			return ""
		}
		return paths[1]
	}

	target := c.Request.URL
	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = container.Address
			req.URL.Host = container.Address
			req.URL.Scheme = "http"
			req.URL.Path = getPath(target.Path)
			req.URL.RawPath = getPath(target.RawPath)
		},
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

// @Summary Exec container
// @Description exec container
// @Tags container
// @Security JWT
// @Param id path string true "container id"
// @Param shell query string  false  "shell, sh or bash"
// @Success 200 {string}  string    ""
// @Router /api/v1/containers/{id}/exec [get]
func (con *ContainerController) Exec(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ResponseFailed(c, http.StatusBadRequest, errors.New("empty container id"))
		return
	}
	c.Done()

	shell := c.Query("shell")
	if shell == "" {
		shell = "sh"
	}

	ctx := c.Request.Context()
	idResp, err := con.client.ContainerExecCreate(ctx, id, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{shell},
		Tty:          true,
	})
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	hijack, err := con.client.ContainerExecAttach(ctx, idResp.ID, types.ExecStartCheck{
		Detach: false,
		Tty:    true,
	})
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	defer hijack.Close()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	defer conn.Close()

	go wsWrite(hijack.Conn, conn)
	wsRead(conn, hijack.Conn)
}

func (con *ContainerController) RegisterRoute(api *gin.RouterGroup) {
	if con.client == nil {
		logrus.Warn("container client is nil, skip register routers")
		return
	}

	api.GET("/containers", con.List)
	api.POST("/containers", con.Create)
	api.GET("/containers/:id", con.Get)
	api.PUT("/containers/:id", con.Update)
	api.POST("/containers/:id", con.Operate)
	api.DELETE("/containers/:id", con.Delete)
	api.GET("/containers/:id/log", con.Log)
	api.GET("/containers/:id/exec", con.Exec)
	api.Any("/containers/:id/proxy/*any", con.Proxy)
	api.GET("/containers/:id/terminal", func(c *gin.Context) {
		c.HTML(200, "terminal.html", nil)
	})
}

func (con *ContainerController) Name() string {
	return "Container"
}

func (con *ContainerController) runContainer(ctx context.Context, ccon *model.CreatedContainer) (string, error) {
	if ccon.Image == "" {
		return "", errors.New("image cannot be empty")
	}

	reader, err := con.client.ImagePull(ctx, ccon.Image, types.ImagePullOptions{})
	if reader != nil {
		defer reader.Close()
	}
	if err != nil {
		return "", err
	}
	io.Copy(io.Discard, reader)

	resp, err := con.client.ContainerCreate(ctx,
		model.ContainerConfig(ccon),
		model.ContainerHostConfig(ccon),
		nil, nil, ccon.Name)

	if err != nil {
		return "", err
	}

	if err := con.client.ContainerStart(ctx, resp.ID,
		types.ContainerStartOptions{}); err != nil {
		return resp.ID, err
	}

	return resp.ID, nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	HandshakeTimeout: 3 * time.Second,
}

func wsWrite(reader io.Reader, writer *websocket.Conn) {
	defer func() {
		writer.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(
			websocket.CloseNormalClosure,
			"closed",
		))
	}()

	for {
		buf := make([]byte, 8192)
		nr, err := reader.Read(buf)
		if nr > 0 {
			writer.WriteMessage(websocket.BinaryMessage, buf[0:nr])
		}
		if err != nil {
			logrus.Infof("ws write err: %v", err)
			return
		}
	}
}

func wsRead(reader *websocket.Conn, writer io.Writer) {
	for {
		messageType, p, err := reader.ReadMessage()
		if err != nil {
			logrus.Infof("ws read err: %v", err)
			return
		}
		switch messageType {
		case websocket.TextMessage, websocket.BinaryMessage:
			writer.Write(p)
		case websocket.CloseMessage:
			return
		}
	}
}
